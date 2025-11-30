package event

import "context"

// EventPublisher defines the interface for any event publisher
type EventPublisher interface {
	Publish(ctx context.Context, evt Event) error
}

// Publisher wraps any EventPublisher implementation (noop or real)
type Publisher struct {
	impl EventPublisher
}

// NewPublisher returns a wrapped publisher
func NewPublisher(impl EventPublisher) *Publisher {
	return &Publisher{impl: impl}
}

// PublishEvent publishes an event using the underlying implementation
func (p *Publisher) PublishEvent(ctx context.Context, evt Event) error {
	return p.impl.Publish(ctx, evt)
}
