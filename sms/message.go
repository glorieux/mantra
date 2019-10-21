package sms

import "fmt"

// Message is a phone message
type Message struct {
	Sender string `schema:"From"`

	// TODO: enforce phone number
	Recipient string `schema:"To"`

	// TODO: enforce limit
	Content string `schema:"Body"`
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
