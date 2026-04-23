package domain

import "time"

type PaymentStatus string

const (
    StatusAuthorized PaymentStatus = "Authorized"
    StatusDeclined   PaymentStatus = "Declined"
)

type Payment struct {
    ID            string
    OrderID       string
    TransactionID string
    Amount        float64  // изменили с int64 на float64
    Status        string
    CreatedAt     time.Time
}
