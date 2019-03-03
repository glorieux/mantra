package mantra

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

// Service is a service
type Service interface {
	fmt.Stringer
	Serve(context.Context, <-chan Message, SendFunc) error
	Stop() error
}

type service struct {
	id             suture.ServiceToken
	ctx            context.Context
	stop           context.CancelFunc
	log            *logrus.Logger
	messageChan    chan Message
	send           SendFunc
	wrappedService Service
}

func newService(s Service, logger *logrus.Logger, send SendFunc) *service {
	ctx, stop := context.WithCancel(context.Background())
	return &service{
		ctx:            ctx,
		stop:           stop,
		log:            logger,
		messageChan:    make(chan Message),
		send:           send,
		wrappedService: s,
	}
}

// Serve runs the service
func (s *service) Serve() {
	err := s.wrappedService.Serve(s.ctx, s.messageChan, s.send)
	if err != nil {
		s.log.Error(err)
		return
	}
	for {
		select {
		case <-s.ctx.Done():
			return
		}
	}
}

// Stop stops the service
func (s *service) Stop() {
	err := s.wrappedService.Stop()
	if err != nil {
		s.log.Error(err)
	}
	s.stop()
}
