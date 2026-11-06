package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"crud-api/models"
)

const userColumns = "id, name, email, document, created_at"

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func scanUser(row rowScanner) (models.User, error) {
	var u models.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Document, &u.CreatedAt)
	return u, err
}

// GetByID devolve models.ErrUserNotFound se o usuário não existir.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (models.User, error) {
	u, err := scanUser(r.db.QueryRowContext(ctx,
		"SELECT "+userColumns+" FROM users WHERE id = $1", id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return u, models.ErrUserNotFound
	}
	if err != nil {
		return u, fmt.Errorf("select user %d: %w", id, err)
	}
	return u, nil
}

// Exists informa se o usuário existe.
func (r *UserRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx,
		"SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", id,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check user %d: %w", id, err)
	}
	return exists, nil
}

// Create devolve models.ErrUserAlreadyExists se email ou documento já estiverem cadastrados.
func (r *UserRepository) Create(ctx context.Context, u models.User) (models.User, error) {
	created, err := scanUser(r.db.QueryRowContext(ctx,
		"INSERT INTO users (name, email, document) VALUES ($1, $2, $3) RETURNING "+userColumns,
		u.Name, u.Email, u.Document,
	))
	if hasPgCode(err, pgUniqueViolation) {
		return created, models.ErrUserAlreadyExists
	}
	if err != nil {
		return created, fmt.Errorf("insert user: %w", err)
	}
	return created, nil
}

// Update altera nome e email. Devolve models.ErrUserNotFound ou
// models.ErrUserAlreadyExists (email já usado por outro usuário).
func (r *UserRepository) Update(ctx context.Context, u models.User) (models.User, error) {
	updated, err := scanUser(r.db.QueryRowContext(ctx,
		"UPDATE users SET name = $1 WHERE id = $2 RETURNING "+userColumns,
		u.Name, u.ID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return updated, models.ErrUserNotFound
	}
	if hasPgCode(err, pgUniqueViolation) {
		return updated, models.ErrUserAlreadyExists
	}
	if err != nil {
		return updated, fmt.Errorf("update user %d: %w", u.ID, err)
	}
	return updated, nil
}

// Delete devolve models.ErrUserNotFound, ou models.ErrUserHasOrders se houver
// pedidos ligados ao usuário (orders.user_id é uma foreign key).
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if hasPgCode(err, pgForeignKeyViolation) {
		return models.ErrUserHasOrders
	}
	if err != nil {
		return fmt.Errorf("delete user %d: %w", id, err)
	}
	return ensureAffected(res, models.ErrUserNotFound)
}

// ensureAffected devolve notFound quando o comando não alterou nenhuma linha.
func ensureAffected(res sql.Result, notFound error) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}
	if n == 0 {
		return notFound
	}
	return nil
}
