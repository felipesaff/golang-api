package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crud-api/models"

	"github.com/lib/pq"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// UserExists informa se o usuário existe.
// TODO: mover para um UserRepository quando users for migrado para esta arquitetura.
func (r *OrderRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", userID,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user %d: %w", userID, err)
	}
	return exists, nil
}

// ProductPrices devolve o preço atual de cada produto encontrado, indexado pelo ID.
// IDs inexistentes simplesmente não aparecem no mapa.
// TODO: mover para um ProductRepository quando products for migrado.
func (r *OrderRepository) ProductPrices(ctx context.Context, productIDs []int64) (map[int64]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, price_cents FROM products WHERE id = ANY($1)", pq.Array(productIDs),
	)
	if err != nil {
		return nil, fmt.Errorf("select product prices: %w", err)
	}
	defer rows.Close()

	prices := make(map[int64]int64, len(productIDs))
	for rows.Next() {
		var id, price int64
		if err := rows.Scan(&id, &price); err != nil {
			return nil, fmt.Errorf("scan product price: %w", err)
		}
		prices[id] = price
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product prices: %w", err)
	}
	return prices, nil
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
