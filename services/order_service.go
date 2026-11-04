package services

import (
	"context"
	"fmt"

	"crud-api/models"
)

// As interfaces ficam aqui (em quem consome), e os repositórios Postgres as
// implementam sem conhecê-las. Cada uma pede só os métodos que o pedido usa.
type OrderRepository interface {
	Create(ctx context.Context, order models.Order) error
	GetByID(ctx context.Context, id int64) (models.OrderDetails, error)
}

type OrderUserRepository interface {
	Exists(ctx context.Context, id int64) (bool, error)
}

type OrderProductRepository interface {
	Prices(ctx context.Context, ids []int64) (map[int64]int64, error)
}

type OrderService struct {
	orders   OrderRepository
	users    OrderUserRepository
	products OrderProductRepository
}

func NewOrderService(orders OrderRepository, users OrderUserRepository, products OrderProductRepository) *OrderService {
	return &OrderService{orders: orders, users: users, products: products}
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

	exists, err := s.users.Exists(ctx, in.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: %w (id %d)", models.ErrInvalidReference, models.ErrUserNotFound, in.UserID)
	}

	productIDs := make([]int64, len(in.Items))
	for i, item := range in.Items {
		productIDs[i] = item.ProductID
	}
	prices, err := s.products.Prices(ctx, productIDs)
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
			return fmt.Errorf("%w: %w (id %d)", models.ErrInvalidReference, models.ErrProductNotFound, item.ProductID)
		}
		order.Items = append(order.Items, models.OrderItem{
			ProductID:      item.ProductID,
			Quantity:       item.Quantity,
			UnitPriceCents: price,
		})
		order.TotalCents += price * int64(item.Quantity)
	}

	return s.orders.Create(ctx, order)
}

func (s *OrderService) GetByID(ctx context.Context, id int64) (models.OrderDetails, error) {
	if err := validateID(id, "order"); err != nil {
		return models.OrderDetails{}, err
	}
	return s.orders.GetByID(ctx, id)
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
