package domain

import "time"

type OrderStatus string

const (
    StatusPending   OrderStatus = "Pending"
    StatusPaid      OrderStatus = "Paid"
    StatusFailed    OrderStatus = "Failed"
    StatusCancelled OrderStatus = "Cancelled"
)

type OrderItem struct {
    ProductID string
    Quantity  int32
    Price     float64
}

type Order struct {
    ID          string
    UserID      string
    Items       []OrderItem
    TotalAmount float64
    Status      string
    CreatedAt   time.Time
}
