package dto

type CreateOrderDTO struct {
	CustomerID string
	OrderItems []CreateOrderItemDTO
}

type CreateOrderItemDTO struct {
	ProductId string
	Price     float32
	Number    int32
}
