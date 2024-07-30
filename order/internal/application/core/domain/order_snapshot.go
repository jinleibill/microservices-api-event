package domain

type OrderSnapshot struct {
	CustomerID string      `json:"customer_id"`
	State      OrderState  `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

func (OrderSnapshot) SnapshotName() string { return "OrderSnapshot" }
