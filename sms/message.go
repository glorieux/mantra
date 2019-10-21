package sms

import "fmt"

// Message is a phone message
type Message struct {
	Sender string

	// TODO: enforce phone number
	Recipient string

	// TODO: enforce limit
	Content string
}

// Response returns a response to the message
func (m *Message) Response(content string) *Message {
	return &Message{
		Sender:    m.Recipient,
		Recipient: m.Sender,
		Content:   content,
	}
}

func (m *Message) String() string {
	return fmt.Sprintf("Message from: %s, to: %s\n%s", m.Sender, m.Recipient, m.Content)
}
