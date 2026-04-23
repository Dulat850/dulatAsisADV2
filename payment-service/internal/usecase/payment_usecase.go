package usecase

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
    "dulatAsisADV2/payment-service/internal/domain"
    "dulatAsisADV2/payment-service/internal/repository"
)

type PaymentUseCase struct {
    repo repository.PaymentRepository
}

func NewPaymentUseCase(repo repository.PaymentRepository) *PaymentUseCase {
    return &PaymentUseCase{repo: repo}
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount float64) (*domain.Payment, error) {
    payment := &domain.Payment{
        ID:            uuid.New().String(),
        OrderID:       orderID,
        TransactionID: fmt.Sprintf("TXN-%d", time.Now().Unix()),
        Amount:        amount,  // теперь float64
        Status:        string(domain.StatusAuthorized),
        CreatedAt:     time.Now(),
    }
    
    err := uc.repo.Create(ctx, payment)
    if err != nil {
        return nil, err
    }
    
    return payment, nil
}

func (uc *PaymentUseCase) GetPaymentStatus(ctx context.Context, paymentID string) (*domain.Payment, error) {
    return uc.repo.GetByID(ctx, paymentID)
}

func (uc *PaymentUseCase) ListByStatus(ctx context.Context, status string) ([]*domain.Payment, error) {
    return uc.repo.ListByStatus(ctx, status)
}
