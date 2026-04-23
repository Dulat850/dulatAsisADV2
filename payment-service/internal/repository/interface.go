package repository

import (
    "context"
    "dulatAsisADV2/payment-service/internal/domain"
)

type PaymentRepository interface {
    Create(ctx context.Context, payment *domain.Payment) error
    GetByID(ctx context.Context, id string) (*domain.Payment, error)
    GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error)
    Update(ctx context.Context, payment *domain.Payment) error
    ListByStatus(ctx context.Context, status string) ([]*domain.Payment, error) // новый метод
}
