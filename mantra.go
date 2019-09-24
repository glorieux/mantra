package mantra

import (
	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

type LookupFunc func(string) *Address

type Application interface {
	Lookup(string) *Address
	NewMailbox(string) (*Address, *Mailbox)
}

// New registers a new application
func New(logger *logrus.Logger, services ...Service) (Application, error) {
	supervisor := suture.New("mantra", suture.Spec{
		Log:        func(s string) { logger.Print(s) },
		LogBadStop: badStopLogger(logger),
		LogFailure: failureLogger(logger),
	})
	supervisor.ServeBackground()

	registry := newServiceRegistry(supervisor, logger)
	for _, service := range services {
		registry.addService(service)
	}
	return registry, nil
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
