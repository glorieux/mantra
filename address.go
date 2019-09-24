package mantra

import (
	"log"
	"os"
)

// Address is a service's address
type Address struct {
	hostname string
	pid      int
	name     string
	mailbox  *Mailbox
}

func newAddress(name string, mailbox *Mailbox) *Address {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	return &Address{
		hostname: hostname,
		pid:      os.Getpid(),
		name:     name,
		mailbox:  mailbox,
	}
}

func (a *Address) Send(msg interface{}) {
	a.mailbox.send(msg)
}
