package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crud-api/models"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create grava o pedido e seus itens numa única transação.
func (r *OrderRepository) Create(ctx context.Context, order models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	// Rollback após um Commit bem-sucedido é no-op, então é seguro sempre adiar.
	defer tx.Rollback()

	var orderID int64
	err = tx.QueryRowContext(ctx,
		"INSERT INTO orders (user_id, total_cents, status) VALUES ($1, $2, $3) RETURNING id",
		order.UserID, order.TotalCents, order.Status,
	).Scan(&orderID)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO order_items (order_id, product_id, quantity, unit_price_cents) VALUES ($1, $2, $3, $4)",
			orderID, item.ProductID, item.Quantity, item.UnitPriceCents,
		)
		if err != nil {
			return fmt.Errorf("insert order_item (product_id=%d): %w", item.ProductID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// GetByID busca o pedido com dados do usuário e dos itens.
// Devolve models.ErrOrderNotFound se o pedido não existir.
func (r *OrderRepository) GetByID(ctx context.Context, id int64) (models.OrderDetails, error) {
	var details models.OrderDetails

	err := r.db.QueryRowContext(ctx, `
		SELECT o.id, o.total_cents, u.name, u.email, u.document
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE o.id = $1`, id,
	).Scan(&details.ID, &details.TotalCents, &details.Customer.Name, &details.Customer.Email, &details.Customer.Document)
	if errors.Is(err, sql.ErrNoRows) {
		return details, models.ErrOrderNotFound
	}
	if err != nil {
		return details, fmt.Errorf("select order %d: %w", id, err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT p.name, p.description, oi.quantity, oi.unit_price_cents
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
		ORDER BY p.id`, id,
	)
	if err != nil {
		return details, fmt.Errorf("select order items %d: %w", id, err)
	}
	defer rows.Close()

	details.Items = []models.OrderDetailsItem{}
	for rows.Next() {
		var item models.OrderDetailsItem
		if err := rows.Scan(&item.ProductName, &item.ProductDescription, &item.Quantity, &item.UnitPriceCents); err != nil {
			return details, fmt.Errorf("scan order item: %w", err)
		}
		details.Items = append(details.Items, item)
	}
	if err := rows.Err(); err != nil {
		return details, fmt.Errorf("iterate order items: %w", err)
	}

	return details, nil
}
