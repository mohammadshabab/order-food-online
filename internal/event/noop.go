package event

import "context"

// NoOpPublisher is a placeholder for future event-driven system
type NoOpPublisher struct{}

func NewNoOpPublisher() EventPublisher {
	return &NoOpPublisher{}
}

// Publish implements EventPublisher interface
func (n *NoOpPublisher) Publish(ctx context.Context, evt Event) error {
	return nil
}
