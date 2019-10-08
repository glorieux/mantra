package sms

import (
	"pkg.glorieux.io/mantra"
	"pkg.glorieux.io/mantra/internal/todo"
)

// ShortMessageService is a service for sending and receiving short messages
type ShortMessageService struct{}

// New returns a new short service message
func New() mantra.Service {
	return &ShortMessageService{}
}

// Serve handles calls to the ShortMessageService
func (s *ShortMessageService) Serve(mux mantra.ServeMux) {
	mux.Handle("send", func(e mantra.Event) {
		s.send(e.Data.(*Message))
	})
}

// Stop stops the ShortMessageService
func (*ShortMessageService) Stop() error {
	return nil
}

func (*ShortMessageService) String() string {
	return "sms"
}

// Message is a phone message
type Message struct {
	// TODO: enforce phone number
	Recipient string

	// TODO: enforce limit
	Content string
}

// send sends a short message
func (*ShortMessageService) send(message *Message) {
	todo.NotImplemented("sms.send")
}
