package base

import (
	"fmt"
	"github.com/google/uuid"
)

const aggregateNeverCommitted = 0

var ErrPendingChanges = fmt.Errorf("cannot process command while pending changes exist")

type AggregateRoot struct {
	aggregate Aggregate
	version   int
}

func NewAggregateRoot(aggregate Aggregate, options ...AggregateRootOption) *AggregateRoot {
	r := &AggregateRoot{
		aggregate: aggregate,
		version:   aggregateNeverCommitted,
	}

	for _, option := range options {
		option(r)
	}

	if r.aggregate.ID() == "" {
		r.aggregate.setID(uuid.New().String())
	}

	return r
}

func (a AggregateRoot) ID() string {
	return a.aggregate.ID()
}

func (a AggregateRoot) AggregateID() string {
	return a.aggregate.ID()
}

func (a AggregateRoot) EntityName() string {
	return a.aggregate.EntityName()
}

func (a AggregateRoot) AggregateName() string {
	return a.aggregate.EntityName()
}

func (a AggregateRoot) Aggregate() Aggregate {
	return a.aggregate
}

func (a AggregateRoot) Events() []Event {
	return a.aggregate.Events()
}

func (a AggregateRoot) AddEvents(events ...Event) {
	a.aggregate.AddEvents(events...)
}

func (a AggregateRoot) ClearEvents() {
	a.aggregate.ClearEvents()
}

func (a AggregateRoot) CommitEvents() {
	a.version += len(a.aggregate.Events())
	a.aggregate.ClearEvents()
}

func (a AggregateRoot) LoadEvent(events ...Event) error {
	for _, event := range events {
		err := a.aggregate.ApplyEvent(event)
		if err != nil {
			return err
		}
	}

	a.version += len(events)

	return nil
}

func (a AggregateRoot) LoadSnapshot(snapshot Snapshot, version int) error {
	err := a.aggregate.ApplySnapshot(snapshot)
	if err != nil {
		return err
	}

	a.version = version

	return nil
}

func (a AggregateRoot) PendingVersion() int {
	return a.version + len(a.aggregate.Events())
}

func (a AggregateRoot) Version() int {
	return a.version
}

func (a AggregateRoot) ProcessCommand(command Command) error {
	if len(a.aggregate.Events()) != 0 {
		return ErrPendingChanges
	}

	err := a.aggregate.ProcessCommand(command)
	if err != nil {
		return err
	}

	for _, event := range a.aggregate.Events() {
		aErr := a.aggregate.ApplyEvent(event)
		if aErr != nil {
			return aErr
		}
	}

	return nil
}

func (a AggregateRoot) GetEvent(eventName string) Event {
	return a.aggregate.GetEvent(eventName)
}

func (a AggregateRoot) GetSnapshotType() Snapshot {
	return a.aggregate.GetSnapshot()
}

type AggregateRootOption func(*AggregateRoot)

func WithAggregateRootID(aggregateID string) AggregateRootOption {
	return func(r *AggregateRoot) {
		r.aggregate.setID(aggregateID)
	}
}
