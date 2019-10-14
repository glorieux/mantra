package mantra

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/mustafaturan/bus"
	"github.com/thejerf/suture"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/internal/structs"
)

// Service is a service
type Service interface {
	// TODO: See if I can get rid of Stop
	Stop() error
}

// Event is an event from bus.Event
type Event *bus.Event

type service struct {
	id   suture.ServiceToken
	ctx  context.Context
	stop context.CancelFunc

	wrappedService        Service
	wrappedServiceName    string
	wrappedServiceMethods []reflect.Method

	Address *Address
	events  chan *bus.Event
}

func newService(wrappedService Service) *service {
	ctx, stop := context.WithCancel(context.Background())
	s := &service{
		ctx:                   ctx,
		stop:                  stop,
		wrappedService:        wrappedService,
		wrappedServiceName:    structs.Name(wrappedService),
		wrappedServiceMethods: structs.Methods(wrappedService),
		events:                make(chan *bus.Event),
	}
	s.Address = newAddress(s.wrappedServiceName)
	bus.RegisterHandler(
		s.wrappedServiceName,
		&bus.Handler{
			Matcher: fmt.Sprintf("^%s.*$", s.Address.String()),
			Handle:  s.handler,
		},
	)
	topics := s.topics()
	log.Debugf("Registering topics: %v", topics)
	bus.RegisterTopics(topics...)
	return s
}

func (s *service) handler(e *bus.Event) {
	s.events <- e
}

// Serve runs the service
func (s *service) Serve() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case e := <-s.events:
			for _, m := range s.wrappedServiceMethods {
				topicMethodName := e.Topic.Name[strings.LastIndex(e.Topic.Name, ".")+1:]
				if topicMethodName != m.Name {
					continue
				}

				numArgs := m.Type.NumIn()
				passedArgs := e.Data.([]interface{})

				log.Debugf("METHOD: %s %d %v", m.Name, m.Type.NumIn(), m.Func)

				if len(passedArgs) != numArgs-1 {
					log.Errorf("Bad argument count! Got %d expected %d.", len(passedArgs), numArgs-1)
				}

				values := []reflect.Value{}
				for _, a := range passedArgs {
					log.Debug("ARG", a)
					values = append(values, reflect.ValueOf(a))
				}

				t := reflect.ValueOf(s.wrappedService)
				meth := t.MethodByName(m.Name)
				meth.Call(values)
			}
		}
	}
}

func (s *service) topics() []string {
	t := []string{}
	for _, m := range s.wrappedServiceMethods {
		t = append(t, newTopic(s.Address, m.Name).String())
	}
	return t
}

// Stop stops the service
func (s *service) Stop() {
	log.Warnf("Stopping %s", s.wrappedServiceName)
	err := s.wrappedService.Stop()
	if err != nil {
		log.Error(err)
	}
	bus.DeregisterTopics(s.topics()...)
	bus.DeregisterHandler(s.wrappedServiceName)
	s.stop()
	close(s.events)
}

func (s *service) String() string {
	return s.wrappedServiceName
}
