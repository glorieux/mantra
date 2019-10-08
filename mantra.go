package mantra

import (
	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
	"github.com/thejerf/suture"
	"pkg.glorieux.io/mantra/internal/log"
)

// VERSION is mantra's version
const VERSION = "0.2.0"

var rootSupervisor *suture.Supervisor

// New registers a new application
// TODO: Rename to Start
func New(services ...Service) error {
	rootSupervisor = suture.New("mantra", suture.Spec{
		Log:        func(s string) { log.Info(s) },
		LogBadStop: badStopLogger(),
		LogFailure: failureLogger(),
	})

	node := uint(1)
	initialTime := uint(0)
	err := monoton.Configure(sequencer.NewMillisecond(), node, initialTime)
	if err != nil {
		log.Fatal(err)
	}
	err = bus.Configure(bus.Config{Next: monoton.Next})
	if err != nil {
		log.Fatal(err)
	}
	bus.RegisterHandler("eventLogger", &bus.Handler{
		Matcher: ".*",
		Handle: func(e *bus.Event) {
			log.Debugf("Event [%s] <- %+v", e.Topic.Name, e.Data)
		},
	})

	registry := newServiceRegistry(rootSupervisor)
	for _, service := range services {
		registry.addService(service)
	}

	rootSupervisor.ServeBackground()
	return nil
}

// Stop stops the application
func Stop() {
	rootSupervisor.Stop()
}

// Send message to a given topic
func Send(topic string, data interface{}) error {
	// Leave transaction ID blank to let bus package auto assigns an ID using the
	// provided generator
	_, err := bus.Emit(topic, data, "")
	return err
}

func badStopLogger() suture.BadStopLogger {
	return func(sup *suture.Supervisor, service suture.Service, msg string) {
		log.Error(service, msg)
	}
}

func failureLogger() suture.FailureLogger {
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
