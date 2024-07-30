package base

type Aggregate interface {
	Entity
	setID(id string)
	ProcessCommand(command Command) error
	ApplyEvent(event Event) error
	ApplySnapshot(snapshot Snapshot) error
	ToSnapshot() (Snapshot, error)
	GetSnapshot() Snapshot
}

type AggregateBase struct {
	EntityBase
	id string
}

func (a *AggregateBase) ID() string {
	return a.id
}

func (a *AggregateBase) setID(id string) {
	a.id = id
}
