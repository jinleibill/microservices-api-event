package base

import (
	"context"
	"errors"
)

var ErrAggregateNotFound = errors.New("aggregate not found")

type AggregateRepository interface {
	Load(ctx context.Context, aggregateID string) (*AggregateRoot, error)
	Save(ctx context.Context, command Command, options ...AggregateRootOption) (*AggregateRoot, error)
}

type AggregateRootRepository struct {
	constructor func() Aggregate
	store       Store
}

func NewAggregateRootRepository(constructor func() Aggregate, store Store) *AggregateRootRepository {
	r := &AggregateRootRepository{
		constructor: constructor,
		store:       store,
	}

	return r
}

func (a *AggregateRootRepository) Load(ctx context.Context, aggregateID string) (*AggregateRoot, error) {
	root := a.root(WithAggregateRootID(aggregateID))

	err := a.store.Load(ctx, root)
	if err != nil {
		return nil, err
	}

	if root.version == aggregateNeverCommitted {
		return nil, ErrAggregateNotFound
	}

	return root, a.store.Load(ctx, root)
}

func (a *AggregateRootRepository) Save(ctx context.Context, command Command, options ...AggregateRootOption) (*AggregateRoot, error) {
	root := a.root(options...)

	return root, a.save(ctx, command, root)
}

func (a *AggregateRootRepository) root(options ...AggregateRootOption) *AggregateRoot {
	return NewAggregateRoot(a.constructor(), options...)
}

func (a *AggregateRootRepository) save(ctx context.Context, command Command, root *AggregateRoot) error {
	err := root.ProcessCommand(command)
	if err != nil {
		return err
	}

	if root.PendingVersion() == root.Version() {
		return nil
	}

	err = a.store.Save(ctx, root)
	if err != nil {
		return err
	}

	return nil
}
