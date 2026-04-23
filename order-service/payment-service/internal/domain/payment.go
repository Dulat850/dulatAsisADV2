package domain

import "time"

type Payment struct {
    ID            string
    OrderID       string
    TransactionID string
    Amount        int64
    Status        string
    CreatedAt     time.Time
}

const MaxAmount int64 = 100000

type PaymentStatus string

const (
    StatusAuthorized PaymentStatus = "authorized"
    StatusDeclined   PaymentStatus = "declined"
)

func CanAuthorize(amount int64) bool {
    return amount <= MaxAmount && amount > 0
}
