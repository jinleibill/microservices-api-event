package domain

type OrderEvent struct{}

func (OrderEvent) DestinationChannel() string { return "Order" }

type OrderCreated struct {
	OrderEvent
	CustomerID string            `json:"customer_id"`
	OrderItems []CreateOrderItem `json:"order_items"`
	OrderTotal float32           `json:"order_total"`
}

func (OrderCreated) EventName() string { return "OrderCreated" }
