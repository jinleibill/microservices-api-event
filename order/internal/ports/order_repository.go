package ports

import (
	"context"
	"order/internal/adapters/base"
	"order/internal/application/core/domain"
)

type OrderRepository interface {
	Load(ctx context.Context, aggregateID string) (*domain.Order, error)
	Save(ctx context.Context, command base.Command, options ...base.AggregateRootOption) (*domain.Order, error)
}
