package base

import "context"

type Store interface {
	Load(ctx context.Context, root *AggregateRoot) error
	Save(ctx context.Context, root *AggregateRoot) error
}

type StoreMiddleware func(store Store) Store
