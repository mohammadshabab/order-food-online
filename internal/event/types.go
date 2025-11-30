package event

type EventType string

const (
	EventOrderCreated   EventType = "order.created"
	EventOrderPaid      EventType = "order.paid"
	EventOrderCancelled EventType = "order.cancelled"
)

type Event struct {
	Type EventType
	Data any
}
