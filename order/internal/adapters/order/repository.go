package order

import (
	"context"
	"errors"
	"fmt"
	"order/internal/adapters/base"
	"order/internal/application/core/domain"
	"order/internal/ports"
)

var _ ports.OrderRepository = (*Adapter)(nil)

type Adapter struct {
	store base.AggregateRepository
}

func NewAdapter(store base.Store) *Adapter {
	return &Adapter{store: base.NewAggregateRootRepository(domain.NewOrder, store)}
}

func (a *Adapter) Load(ctx context.Context, aggregateID string) (*domain.Order, error) {
	root, err := a.store.Load(ctx, aggregateID)
	if err != nil {
		if errors.Is(err, base.ErrAggregateNotFound) {
			return nil, fmt.Errorf("order not found: %w", err)
		}
		return nil, err
	}

	return root.Aggregate().(*domain.Order), nil
}

func (a *Adapter) Save(ctx context.Context, command base.Command, options ...base.AggregateRootOption) (*domain.Order, error) {
	root, err := a.store.Save(ctx, command, options...)
	if err != nil {
		return nil, err
	}

	return root.Aggregate().(*domain.Order), nil
}
