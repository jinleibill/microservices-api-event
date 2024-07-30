package grpc

import (
	"context"
	"github.com/jinleibill/microservices-proto/golang/order"
	"google.golang.org/grpc"
	"order/internal/adapters/base"
	"order/internal/application/core/dto"
	"order/internal/ports"
)

type Adapter struct {
	app ports.Application
	order.UnimplementedOrderServer
	client base.Client
}

func NewAdapter(api ports.Application, client base.Client) *Adapter {
	return &Adapter{app: api, client: client}
}

func (a *Adapter) Mount(registrar grpc.ServiceRegistrar) {
	order.RegisterOrderServer(registrar, a)
}

func (a *Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []dto.CreateOrderItemDTO
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, dto.CreateOrderItemDTO{
			ProductId: orderItem.ProductCode,
			Price:     orderItem.UnitPrice,
			Number:    orderItem.Quantity,
		})
	}

	result, err := a.app.CreateOrder(ctx, dto.CreateOrderDTO{
		CustomerID: request.UserId,
		OrderItems: orderItems,
	})
	if err != nil {
		return nil, err
	}

	return &order.CreateOrderResponse{OrderId: result.ID()}, nil
}

func (a *Adapter) Get(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	o, err := a.app.GetOrder(ctx, request.OrderId)
	if err != nil {
		return nil, err
	}

	orderItems := make([]*order.OrderItem, 0, len(o.OrderItems))
	for _, orderItem := range o.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			ProductCode: orderItem.ProductId,
			UnitPrice:   orderItem.Price,
			Quantity:    orderItem.Number,
		})
	}

	return &order.GetOrderResponse{
		UserId:     o.CustomerID,
		OrderItems: orderItems,
	}, nil
}
