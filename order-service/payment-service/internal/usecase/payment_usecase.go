package usecase

import (
    "context"
    "crypto/rand"
    "encoding/hex"
    "fmt"

    "dulatAsisADV2/payment-service/internal/domain"
    "dulatAsisADV2/payment-service/internal/repository"
)

type PaymentUseCase struct {
    paymentRepo repository.PaymentRepository
}

func NewPaymentUseCase(paymentRepo repository.PaymentRepository) *PaymentUseCase {
    return &PaymentUseCase{
        paymentRepo: paymentRepo,
    }
}

func (uc *PaymentUseCase) ProcessPayment(ctx context.Context, orderID string, amount int64) (*domain.Payment, error) {
    if amount <= 0 {
        return nil, fmt.Errorf("invalid amount: %d", amount)
    }

    payment := &domain.Payment{
        ID:            generateID(),
        OrderID:       orderID,
        TransactionID: generateTransactionID(),
        Amount:        amount,
    }

    if domain.CanAuthorize(amount) {
        payment.Status = "authorized"
    } else {
        payment.Status = "declined"
    }

    if err := uc.paymentRepo.Create(ctx, payment); err != nil {
        return nil, err
    }

    return payment, nil
}

func generateID() string {
    bytes := make([]byte, 16)
    rand.Read(bytes)
    return hex.EncodeToString(bytes)
}

func generateTransactionID() string {
    bytes := make([]byte, 32)
    rand.Read(bytes)
    return "txn_" + hex.EncodeToString(bytes)
}
