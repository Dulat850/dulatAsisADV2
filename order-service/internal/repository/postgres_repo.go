package repository

import (
    "context"
    "database/sql"
    "encoding/json"
    
    "dulatAsisADV2/order-service/internal/domain"
)

type PostgresOrderRepo struct {
    db *sql.DB
}

func NewPostgresOrderRepo(db *sql.DB) *PostgresOrderRepo {
    return &PostgresOrderRepo{db: db}
}

func (r *PostgresOrderRepo) Create(ctx context.Context, order *domain.Order) error {
    itemsJSON, err := json.Marshal(order.Items)
    if err != nil {
        return err
    }
    
    query := `INSERT INTO orders (id, user_id, items, total_amount, status, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
    _, err = r.db.ExecContext(ctx, query, order.ID, order.UserID, itemsJSON, order.TotalAmount, order.Status, order.CreatedAt)
    return err
}

func (r *PostgresOrderRepo) GetByID(ctx context.Context, id string) (*domain.Order, error) {
    query := `SELECT id, user_id, items, total_amount, status, created_at FROM orders WHERE id = $1`
    
    var order domain.Order
    var itemsJSON []byte
    
    err := r.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &itemsJSON, &order.TotalAmount, &order.Status, &order.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    
    err = json.Unmarshal(itemsJSON, &order.Items)
    if err != nil {
        return nil, err
    }
    
    return &order, nil
}

func (r *PostgresOrderRepo) UpdateStatus(ctx context.Context, id string, status string) error {
    query := `UPDATE orders SET status = $1 WHERE id = $2`
    _, err := r.db.ExecContext(ctx, query, status, id)
    return err
}
