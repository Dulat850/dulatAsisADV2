package repository

import (
    "context"
    "dulatAsisADV2/payment-service/internal/domain"
)

type PaymentRepository interface {
    Create(ctx context.Context, payment *domain.Payment) error
    GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
}
