package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crud-api/models"

	"github.com/lib/pq"
)

// description é anulável no banco; COALESCE evita falha no scan para string.
const productColumns = "id, name, COALESCE(description, ''), price_cents, stock, created_at"

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func scanProduct(row rowScanner) (models.Product, error) {
	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.PriceCents, &p.Stock, &p.CreatedAt)
	return p, err
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT "+productColumns+" FROM products ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("select products: %w", err)
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return products, nil
}

// GetByID devolve models.ErrProductNotFound se o produto não existir.
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (models.Product, error) {
	p, err := scanProduct(r.db.QueryRowContext(ctx,
		"SELECT "+productColumns+" FROM products WHERE id = $1", id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return p, models.ErrProductNotFound
	}
	if err != nil {
		return p, fmt.Errorf("select product %d: %w", id, err)
	}
	return p, nil
}

// Prices devolve o preço atual de cada produto encontrado, indexado pelo ID.
// IDs inexistentes simplesmente não aparecem no mapa.
func (r *ProductRepository) Prices(ctx context.Context, ids []int64) (map[int64]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, price_cents FROM products WHERE id = ANY($1)", pq.Array(ids),
	)
	if err != nil {
		return nil, fmt.Errorf("select product prices: %w", err)
	}
	defer rows.Close()

	prices := make(map[int64]int64, len(ids))
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

func (r *ProductRepository) Create(ctx context.Context, p models.Product) (models.Product, error) {
	created, err := scanProduct(r.db.QueryRowContext(ctx,
		"INSERT INTO products (name, description, price_cents, stock) VALUES ($1, $2, $3, $4) RETURNING "+productColumns,
		p.Name, p.Description, p.PriceCents, p.Stock,
	))
	if err != nil {
		return created, fmt.Errorf("insert product: %w", err)
	}
	return created, nil
}

// Update devolve models.ErrProductNotFound se o produto não existir.
func (r *ProductRepository) Update(ctx context.Context, p models.Product) (models.Product, error) {
	updated, err := scanProduct(r.db.QueryRowContext(ctx,
		"UPDATE products SET name = $1, description = $2, price_cents = $3, stock = $4 WHERE id = $5 RETURNING "+productColumns,
		p.Name, p.Description, p.PriceCents, p.Stock, p.ID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return updated, models.ErrProductNotFound
	}
	if err != nil {
		return updated, fmt.Errorf("update product %d: %w", p.ID, err)
	}
	return updated, nil
}

// Delete devolve models.ErrProductNotFound, ou models.ErrProductInUse se o
// produto aparecer em algum pedido (order_items.product_id é uma foreign key).
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
	if hasPgCode(err, pgForeignKeyViolation) {
		return models.ErrProductInUse
	}
	if err != nil {
		return fmt.Errorf("delete product %d: %w", id, err)
	}
	return ensureAffected(res, models.ErrProductNotFound)
}
