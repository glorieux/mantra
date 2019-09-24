package mantra

import (
	"sync"
)

type message struct {
	msg interface{}
}

// Mailbox receives messages
type Mailbox struct {
	sync.Mutex
	messages    []message
	messageChan <-chan message
}

// newMailbox returns a new Address and Mailbox
func newMailbox(name string) (*Address, *Mailbox) {
	mailbox := &Mailbox{
		messages:    make([]message, 0, 1),
		messageChan: make(chan message),
	}
	return newAddress(name, mailbox), mailbox
}

func (m *Mailbox) send(msg interface{}) {
	m.Lock()
	defer m.Unlock()

	m.messages = append(m.messages, message{msg})
}

// Receive calls handler for a given type
func (m *Mailbox) Receive(handler func(interface{})) {
	for msg := range m.messageChan {
		handler(msg.msg)
	}
}
