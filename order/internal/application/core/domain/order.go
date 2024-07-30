package domain

import (
	"errors"
	"fmt"
	"order/internal/adapters/base"
)

var (
	ErrOrderUnhandledCommand  = errors.New("unhandled command in order aggregate")
	ErrOrderInvalidState      = errors.New("order state does not allow action")
	ErrOrderUnhandledEvent    = errors.New("unhandled event in order aggregate")
	ErrOrderUnhandledSnapshot = errors.New("unhandled snapshot in order aggregate")
)

type OrderState int

const (
	UnknownOrderState OrderState = iota
	ApprovalPending              // 待确认
	Approved                     // 已确认
	Rejected                     // 已拒绝
	CancelPending                // 待取消
	Cancelled                    // 已取消
	RevisionPending              // 待修改
)

func (s OrderState) String() string {
	switch s {
	case ApprovalPending:
		return "ApprovalPending"
	case Approved:
		return "Approved"
	case Rejected:
		return "Rejected"
	case CancelPending:
		return "CancelPending"
	case Cancelled:
		return "Cancelled"
	case RevisionPending:
		return "RevisionPending"
	default:
		return "Unknown"
	}
}

var _ base.Aggregate = (*Order)(nil)

type Order struct {
	base.AggregateBase
	CustomerID string      `json:"customer_id"`
	State      OrderState  `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ProductId string  `json:"product_id"`
	Price     float32 `json:"price"`
	Number    int32   `json:"number"`
}

func NewOrder() base.Aggregate {
	return &Order{}
}

func (o *Order) EntityName() string {
	return "order"
}

func (o *Order) ProcessCommand(command base.Command) error {
	switch cmd := command.(type) {
	case *CreateOrder:
		if o.State != UnknownOrderState {
			return ErrOrderInvalidState
		}

		var total float32
		orderItems := make([]CreateOrderItem, 0, len(o.OrderItems))
		for _, orderItem := range cmd.OrderItems {
			item := CreateOrderItem{
				ProductId: orderItem.ProductId,
				Price:     orderItem.Price,
				Number:    orderItem.Number,
			}
			total += item.GetTotal()
			orderItems = append(orderItems, item)
		}

		o.AddEvents(&OrderCreated{
			CustomerID: cmd.CustomerID,
			OrderItems: cmd.OrderItems,
			OrderTotal: total,
		})
	default:
		return fmt.Errorf("%w: unhandled command %s", ErrOrderUnhandledCommand, command.CommandName())
	}

	return nil
}

func (o *Order) ApplyEvent(event base.Event) error {
	switch e := event.(type) {
	case *OrderCreated:
		orderItems := make([]OrderItem, 0, len(e.OrderItems))
		for _, orderItem := range e.OrderItems {
			orderItem := OrderItem{
				ProductId: orderItem.ProductId,
				Price:     orderItem.Price,
				Number:    orderItem.Number,
			}
			orderItems = append(orderItems, orderItem)
		}

		o.CustomerID = e.CustomerID
		o.OrderItems = orderItems
		o.State = ApprovalPending
	default:
		return fmt.Errorf("%w: unhandled event %s", ErrOrderUnhandledEvent, event)
	}

	return nil
}

func (o *Order) ApplySnapshot(snapshot base.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OrderSnapshot:
		o.CustomerID = ss.CustomerID
		o.State = ApprovalPending
		o.OrderItems = ss.OrderItems
	default:
		return fmt.Errorf("%w: unhandled snapshot %s", ErrOrderUnhandledSnapshot, snapshot)
	}

	return nil
}

func (o *Order) ToSnapshot() (base.Snapshot, error) {
	return &OrderSnapshot{
		CustomerID: o.CustomerID,
		OrderItems: o.OrderItems,
		State:      o.State,
	}, nil
}

func (o *Order) GetEvent(eventName string) base.Event {
	switch eventName {
	case "OrderCreated":
		return &OrderCreated{}
	}

	return nil
}

func (o *Order) GetSnapshot() base.Snapshot {
	return &OrderSnapshot{}
}
