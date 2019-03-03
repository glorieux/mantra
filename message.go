package mantra

// Message is a command exchanged between services
type Message interface {
	To() string
}
