package mantra

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
	"github.com/sirupsen/logrus"
	"github.com/thejerf/suture"
)

// New registers a new application
func New(logger *logrus.Logger, services ...Service) error {
	supervisor := suture.New("mantra", suture.Spec{
		Log:        func(s string) { logger.Print(s) },
		LogBadStop: badStopLogger(logger),
		LogFailure: failureLogger(logger),
	})

	node := uint(1)
	initialTime := uint(0)
	err := monoton.Configure(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		logger.Fatal(err)
	}
	err = bus.Configure(bus.Config{Next: monoton.Next})
	if err != nil {
		logger.Fatal(err)
	}
	bus.RegisterHandler("eventLogger", &bus.Handler{
		Matcher: ".*",
		Handle: func(e *bus.Event) {
			logger.Debugf("Event [%s] <- %+v", e.Topic.Name, e.Data)
		},
	})

	registry := newServiceRegistry(supervisor, logger)
	for _, service := range services {
		registry.addService(service)
	}

	supervisor.ServeBackground()
	return nil
}

// Send message to a given topic
func Send(topic string, data interface{}) error {
	// Leave transaction ID blank to let bus package auto assigns an ID using the provided gen
	_, err := bus.Emit(topic, data, "")
	return err
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
