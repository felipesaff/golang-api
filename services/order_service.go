package services

import (
	"context"
	"fmt"

	"crud-api/models"
)

// OrderRepository é o que o service precisa do armazenamento. A interface fica
// aqui (em quem consome), e o repositório Postgres a implementa sem conhecê-la.
type OrderRepository interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
	ProductPrices(ctx context.Context, productIDs []int64) (map[int64]int64, error)
	Create(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int64) (models.OrderDetails, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

type CreateOrderInput struct {
	UserID int64
	Items  []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID int64
	Quantity  int
}

// Create valida a entrada, aplica as regras do pedido (preço atual do produto,
// total, status inicial) e grava o pedido.
func (s *OrderService) Create(ctx context.Context, in CreateOrderInput) error {
	if err := in.validate(); err != nil {
		return err
	}

	exists, err := s.repo.UserExists(ctx, in.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: id %d", models.ErrUserNotFound, in.UserID)
	}

	productIDs := make([]int64, len(in.Items))
	for i, item := range in.Items {
		productIDs[i] = item.ProductID
	}
	prices, err := s.repo.ProductPrices(ctx, productIDs)
	if err != nil {
		return err
	}

	order := models.Order{
		UserID: in.UserID,
		Status: models.OrderStatusPending,
		Items:  make([]models.OrderItem, 0, len(in.Items)),
	}
	for _, item := range in.Items {
		price, ok := prices[item.ProductID]
		if !ok {
			return fmt.Errorf("%w: id %d", models.ErrProductNotFound, item.ProductID)
		}
		order.Items = append(order.Items, models.OrderItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPriceCents: price,
		})
		order.TotalCents += price * int64(item.Quantity)
	}

	return s.repo.Create(ctx, order)
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (models.OrderDetails, error) {
	if id < 1 {
		return models.OrderDetails{}, &models.ValidationError{Msg: "invalid order ID"}
	}
	return s.repo.GetByID(ctx, id)
}

func (in CreateOrderInput) validate() error {
	if in.UserID < 1 {
		return &models.ValidationError{Msg: "invalid user ID"}
	}
	if len(in.Items) == 0 {
		return &models.ValidationError{Msg: "items must have at least one item"}
	}

	seen := make(map[int64]bool, len(in.Items))
	for i, item := range in.Items {
		if item.ProductID < 1 {
			return &models.ValidationError{Msg: fmt.Sprintf("items[%d]: product_id must be greater than 0", i)}
		}
		if item.Quantity < 1 {
			return &models.ValidationError{Msg: fmt.Sprintf("items[%d]: quantity must be greater than 0", i)}
		}
		// order_items tem PRIMARY KEY (order_id, product_id): produto repetido
		// estouraria no banco como erro 500.
		if seen[item.ProductID] {
			return &models.ValidationError{Msg: fmt.Sprintf("items[%d]: product_id %d is duplicated", i, item.ProductID)}
		}
		seen[item.ProductID] = true
	}
	return nil
}
