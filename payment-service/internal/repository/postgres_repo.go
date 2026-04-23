package repository

import (
    "context"
    "database/sql"
    
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
    _, err := r.db.ExecContext(ctx, query, payment.ID, payment.OrderID, payment.TransactionID, payment.Amount, payment.Status, payment.CreatedAt)
    return err
}

func (r *PostgresPaymentRepo) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
    query := `SELECT id, order_id, transaction_id, amount, status, created_at FROM payments WHERE id = $1`
    
    var payment domain.Payment
    err := r.db.QueryRowContext(ctx, query, id).Scan(
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

func (r *PostgresPaymentRepo) Update(ctx context.Context, payment *domain.Payment) error {
    query := `UPDATE payments SET status = $1, transaction_id = $2 WHERE id = $3`
    _, err := r.db.ExecContext(ctx, query, payment.Status, payment.TransactionID, payment.ID)
    return err
}

func (r *PostgresPaymentRepo) ListByStatus(ctx context.Context, status string) ([]*domain.Payment, error) {
    query := `SELECT id, order_id, transaction_id, amount, status, created_at FROM payments WHERE status = $1 ORDER BY created_at DESC`
    
    rows, err := r.db.QueryContext(ctx, query, status)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var payments []*domain.Payment
    for rows.Next() {
        var payment domain.Payment
        err := rows.Scan(&payment.ID, &payment.OrderID, &payment.TransactionID, &payment.Amount, &payment.Status, &payment.CreatedAt)
        if err != nil {
            return nil, err
        }
        payments = append(payments, &payment)
    }
    
    return payments, nil
}
