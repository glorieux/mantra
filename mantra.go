package mantra

import (
	"time"

	"github.com/mustafaturan/bus"
	"github.com/mustafaturan/monoton"
	"github.com/mustafaturan/monoton/sequencer"
	"github.com/thejerf/suture"
	"pkg.glorieux.io/mantra/internal/log"
	"pkg.glorieux.io/mantra/internal/structs"
)

// VERSION is mantra's version
const VERSION = "0.2.0"

var rootSupervisor *suture.Supervisor

// New registers a new application
// TODO: Rename to Start
func New(services ...Service) error {
	setFlags()
	rootSupervisor = suture.New("mantra", suture.Spec{
		Timeout:    10 * time.Second,
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

	registry := newServiceRegistry()
	for _, service := range services {
		registry.AddService(service)
	}

	rootSupervisor.ServeBackground()
	return nil
}

// Stop stops the application
func Stop() {
	log.Warn("Stopping all services...")
	rootSupervisor.Stop()
}

// SendMessage supperseeds Send
func SendMessage(address *Address, method interface{}, args ...interface{}) {
	topic := newTopic(address, structs.FuncName(method))
	log.Debug("TOPIC: ", topic)
	_, err := bus.Emit(topic.String(), args, "")
	if err != nil {
		log.Error(err)
	}
}

// Lookup returns a services address
func Lookup(name string) *Address {
	// TODO verify is actually exists
	return newAddress(name)
}

func badStopLogger() suture.BadStopLogger {
	return func(sup *suture.Supervisor, service suture.Service, msg string) {
		log.Errorf("Bad stop: %v %s", service, msg)
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
