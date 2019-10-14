package sms

import (
	"pkg.glorieux.io/mantra/internal/todo"
)

// ShortMessageService is a service for sending and receiving short messages
type ShortMessageService struct{}

// New returns a new short service message
func New() *ShortMessageService {
	return &ShortMessageService{}
}

// Message is a phone message
type Message struct {
	// TODO: enforce phone number
	Recipient string

	// TODO: enforce limit
	Content string
}

func (*ShortMessageService) Send(msg *Message) {
	todo.NotImplemented("sms.send")
}

// Stop stops the ShortMessageService
func (*ShortMessageService) Stop() error {
	return nil
}
