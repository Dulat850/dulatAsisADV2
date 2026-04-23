package repository

import (
    "context"
    "database/sql"
    "time"

    "dulatAsisADV2/payment-service/internal/domain"
)

type PostgresPaymentRepo struct {
    db *sql.DB
}

func NewPostgresPaymentRepo(db *sql.DB) *PostgresPaymentRepo {
    return &PostgresPaymentRepo{db: db}
}

func (r *PostgresPaymentRepo) Create(ctx context.Context, payment *domain.Payment) error {
    query := `
        INSERT INTO payments (id, order_id, transaction_id, amount, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `
    _, err := r.db.ExecContext(ctx, query,
        payment.ID, payment.OrderID, payment.TransactionID, payment.Amount, payment.Status, time.Now())
    return err
}

func (r *PostgresPaymentRepo) GetByOrderID(ctx context.Context, orderID string) (*domain.Payment, error) {
    query := `SELECT id, order_id, transaction_id, amount, status, created_at FROM payments WHERE order_id = $1`

    var payment domain.Payment
    err := r.db.QueryRowContext(ctx, query, orderID).Scan(
        &payment.ID, &payment.OrderID, &payment.TransactionID, &payment.Amount, &payment.Status, &payment.CreatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return &payment, nil
}
