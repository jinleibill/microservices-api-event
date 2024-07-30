package base

type Entity interface {
	ID() string
	EntityName() string
	Events() []Event
	AddEvents(events ...Event)
	ClearEvents()
	GetEvent(eventName string) Event
}

type EntityBase struct {
	events []Event
}

func (e *EntityBase) Events() []Event {
	return e.events
}

func (e *EntityBase) AddEvents(events ...Event) {
	e.events = append(e.events, events...)
}

func (e *EntityBase) ClearEvents() {
	e.events = []Event{}
}
