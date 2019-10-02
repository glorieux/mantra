package mantra

import (
	"context"
	"fmt"

	"github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

// Service is a service
type Service interface {
	fmt.Stringer

	Serve(ServeMux)
	Stop() error
}

// Event is an event from bus.Event
type Event *bus.Event

// Handler is function handling the Event
type Handler func(e Event)

// ServeMux is a Handler multiplexer
type ServeMux interface {
	Handle(string, Handler)
}

type service struct {
	id             suture.ServiceToken
	ctx            context.Context
	stop           context.CancelFunc
	log            *logrus.Logger
	wrappedService Service
	events         chan *bus.Event

	eventHandlers map[string]Handler
}

func newService(wrappedService Service, logger *logrus.Logger) *service {
	ctx, stop := context.WithCancel(context.Background())
	s := &service{
		ctx:            ctx,
		stop:           stop,
		log:            logger,
		wrappedService: wrappedService,
		events:         make(chan *bus.Event),
		eventHandlers:  make(map[string]Handler),
	}
	bus.RegisterHandler(
		wrappedService.String(),
		&bus.Handler{
			Matcher: fmt.Sprintf("^%s.*$", wrappedService.String()),
			Handle:  s.handler,
		},
	)
	s.wrappedService.Serve(s)
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
			handler := s.eventHandlers[e.Topic.Name]
			if handler != nil {
				handler(e)
				continue
			}
			s.log.Warnf("Unknown topic [%s]", e.Topic.Name)
		}
	}
}

func (s *service) Handle(name string, handler Handler) {
	topic := fmt.Sprintf("%s.%s", s.wrappedService.String(), name)
	s.log.Debug("Registering topic ", topic)
	s.eventHandlers[topic] = handler
	bus.RegisterTopics(topic)
}

func (s *service) topics() []string {
	t := []string{}
	for topic := range s.eventHandlers {
		t = append(t, topic)
	}
	return t
}

// Stop stops the service
func (s *service) Stop() {
	err := s.wrappedService.Stop()
	if err != nil {
		s.log.Error(err)
	}
	bus.DeregisterTopics(s.topics()...)
	bus.DeregisterHandler(s.wrappedService.String())
	close(s.events)
	s.stop()
}
