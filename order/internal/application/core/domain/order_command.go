package domain

type CreateOrder struct {
	CustomerID string
	OrderItems []CreateOrderItem
}

type CreateOrderItem struct {
	ProductId string
	Price     float32
	Number    int32
}

func (CreateOrder) CommandName() string {
	return "CreateOrder"
}

func (i CreateOrderItem) GetTotal() float32 {
	return i.Price * float32(i.Number)
}
