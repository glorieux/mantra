package mantra

import (
	"fmt"
	"os"
	"strings"

	"github.com/mustafaturan/bus"
	"pkg.glorieux.io/mantra/internal/log"
)

// Address is a service's address
type Address struct {
	// TODO find a better representation
	// might be better as a node name (string)
	PID         int
	HostName    string
	ServiceName string
}

// newAddress returns a new Address
func newAddress(name string) *Address {
	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err)
	}

	return &Address{
		PID:         os.Getpid(),
		HostName:    hostname,
		ServiceName: name,
	}
}

func parseAddress(a string) *Address {
	pidHostSeparator := strings.Index(a, "@")
	hostServiceSeparator := strings.Index(a, "#")

	return &Address{
		HostName:    a[pidHostSeparator+1 : hostServiceSeparator],
		ServiceName: a[hostServiceSeparator+1:],
	}
}

// Send sends a message
func (a *Address) Send(method string, args ...interface{}) {
	topic := newTopic(a, method)
	log.Debug("TOPIC: ", topic)
	_, err := bus.Emit(topic.String(), args, "")
	if err != nil {
		log.Error(err)
	}
}

func (a *Address) String() string {
	return fmt.Sprintf(
		"%d@%s#%s",
		a.PID,
		a.HostName,
		a.ServiceName,
	)
}
