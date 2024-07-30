package application

import (
	"context"
	"order/internal/application/core/domain"
	"order/internal/application/core/dto"
	"order/internal/ports"
)

var _ ports.Application = (*Application)(nil)

type Application struct {
	orderRepo ports.OrderRepository
}

func NewApplication(orderRepo ports.OrderRepository) *Application {
	return &Application{
		orderRepo: orderRepo,
	}
}

func (app *Application) CreateOrder(ctx context.Context, dto dto.CreateOrderDTO) (domain.Order, error) {
	var orderItems []domain.CreateOrderItem
	for _, item := range dto.OrderItems {
		orderItems = append(orderItems, domain.CreateOrderItem{
			ProductId: item.ProductId,
			Price:     item.Price,
			Number:    item.Number,
		})
	}

	order, err := app.orderRepo.Save(ctx, &domain.CreateOrder{
		CustomerID: dto.CustomerID,
		OrderItems: orderItems,
	})
	if err != nil {
		return domain.Order{}, err
	}

	return *order, nil
}

func (app *Application) GetOrder(ctx context.Context, aggregateID string) (*domain.Order, error) {
	return app.orderRepo.Load(ctx, aggregateID)
}
