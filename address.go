package mantra

import (
	"fmt"
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

	if mailbox == nil {
		log.Fatal("Nil mailbox")
		fmt.Println("NILLLLL")
	}

	return &Address{
		hostname: hostname,
		pid:      os.Getpid(),
		name:     name,
		mailbox:  mailbox,
	}
}

// Send sends a message
func (a *Address) Send(msg interface{}) {
	fmt.Printf("SEND%+v\n", a)
	log.Println("SEND", a.mailbox)
	a.mailbox.send(msg)
}
