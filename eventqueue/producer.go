package eventqueue

// Producer the producer interface
type Producer interface {
	Produce(topic string, message string) error
	Close()
}
