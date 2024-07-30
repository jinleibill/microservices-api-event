package ports

import (
	"context"
	"order/internal/application/core/domain"
	"order/internal/application/core/dto"
)

type Application interface {
	GetOrder(ctx context.Context, aggregateID string) (*domain.Order, error)
	CreateOrder(ctx context.Context, dto dto.CreateOrderDTO) (domain.Order, error)
}
