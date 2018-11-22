package mantra

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

// Application represents a mantra application
type Application interface {
	fmt.Stringer

	// Init initializes the application
	Init(SendFunc) error
}

// New registers a new application
func New(app Application, logger *logrus.Logger) error {
	supervisor := suture.New(app.String(), suture.Spec{
		Log:        func(s string) { logger.Print(s) },
		LogBadStop: badStopLogger(logger),
		LogFailure: failureLogger(logger),
	})
	supervisor.ServeBackground()
	registry := newServiceRegistry(supervisor, logger)

	// Make the init call syncronous
	err := app.Init(registry.send)
	if err != nil {
		return err
	}
	return nil
}

// Service is a service
type Service interface {
	fmt.Stringer
	Serve(context.Context, <-chan Message, SendFunc) error
	Stop() error
}

// Message is a command exchanged between services
type Message interface {
	To() string
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
	s.wrappedService.Serve(s.ctx, s.messageChan, s.send)
	for {
		select {
		case <-s.ctx.Done():
			return
		}
	}
}

// Stop stops the service
func (s *service) Stop() {
	s.wrappedService.Stop()
	s.stop()
}

func badStopLogger(log *logrus.Logger) suture.BadStopLogger {
	return func(sup *suture.Supervisor, service suture.Service, msg string) {
		log.Error(service, msg)
	}
}

func failureLogger(log *logrus.Logger) suture.FailureLogger {
	return func(
		supervisor *suture.Supervisor,
		service suture.Service,
		serviceName string,
		currentFailures float64,
		failureThreshold float64,
		restarting bool,
		error interface{},
		stacktrace []byte,
	) {
		log.Errorf("Service: %s\nError: %s\nStacktrace: %s\n", serviceName, error, stacktrace)
	}
}
