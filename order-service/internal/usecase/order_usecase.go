package usecase

import (
    "context"
    "time"

    "github.com/google/uuid"
    "dulatAsisADV2/order-service/internal/domain"
    "dulatAsisADV2/order-service/internal/repository"
)

type OrderUseCase struct {
    repo repository.OrderRepository
}

func NewOrderUseCase(repo repository.OrderRepository) *OrderUseCase {
    return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, userID string, items []*domain.OrderItem, totalAmount float64) (*domain.Order, error) {
    // Convert items
    var domainItems []domain.OrderItem
    for _, item := range items {
        domainItems = append(domainItems, domain.OrderItem{
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
            Price:     item.Price,
        })
    }
    
    order := &domain.Order{
        ID:          uuid.New().String(),
        UserID:      userID,
        Items:       domainItems,
        TotalAmount: totalAmount,
        Status:      string(domain.StatusPending),
        CreatedAt:   time.Now(),
    }
    
    err := uc.repo.Create(ctx, order)
    if err != nil {
        return nil, err
    }
    
    return order, nil
}

func (uc *OrderUseCase) GetOrder(ctx context.Context, orderID string) (*domain.Order, error) {
    return uc.repo.GetByID(ctx, orderID)
}

func (uc *OrderUseCase) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
    return uc.repo.UpdateStatus(ctx, orderID, status)
}
