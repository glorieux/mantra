package mantra

import (
	"sync"
)

type message struct {
	msg interface{}
}

type ReceiveFunc func(interface{})

// Mailbox receives messages
type Mailbox struct {
	sync.Mutex
	messages    []message
	messageChan chan message
}

// newMailbox returns a new Mailbox
func newMailbox(name string) *Mailbox {
	return &Mailbox{
		messages:    make([]message, 0, 1),
		messageChan: make(chan message, 1),
	}
}

func (m *Mailbox) send(msg interface{}) {
	m.Lock()
	defer m.Unlock()

	m.messages = append(m.messages, message{msg})
	m.messageChan <- message{msg}
}

// Receive calls handler for a given type
func (m *Mailbox) Receive(handler func(interface{})) {
	go func() {
		for msg := range m.messageChan {
			handler(msg.msg)
		}
	}()
}
