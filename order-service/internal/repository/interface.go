package repository

import (
    "context"
    "dulatAsisADV2/order-service/internal/domain"
)

type OrderRepository interface {
    Create(ctx context.Context, order *domain.Order) error
    GetByID(ctx context.Context, id string) (*domain.Order, error)
    UpdateStatus(ctx context.Context, id string, status string) error
}
