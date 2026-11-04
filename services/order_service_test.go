package services

import (
	"context"
	"errors"
	"testing"

	"crud-api/models"
)

// fakeOrderRepo substitui o banco nos testes: é isso que as interfaces no service permitem.
// Implementa as três interfaces de que o OrderService depende.
type fakeOrderRepo struct {
	users   map[int64]bool
	prices  map[int64]int64
	created *models.Order
}

func (f *fakeOrderRepo) Exists(_ context.Context, id int64) (bool, error) {
	return f.users[id], nil
}

func (f *fakeOrderRepo) Prices(_ context.Context, ids []int64) (map[int64]int64, error) {
	out := map[int64]int64{}
	for _, id := range ids {
		if p, ok := f.prices[id]; ok {
			out[id] = p
		}
	}
	return out, nil
}

func (f *fakeOrderRepo) Create(_ context.Context, o models.Order) error {
	f.created = &o
	return nil
}

func (f *fakeOrderRepo) GetByID(_ context.Context, _ int64) (models.OrderDetails, error) {
	return models.OrderDetails{}, models.ErrOrderNotFound
}

func newFakeOrderRepo() *fakeOrderRepo {
	return &fakeOrderRepo{
		users:  map[int64]bool{1: true},
		prices: map[int64]int64{10: 1500, 20: 250},
	}
}

func newOrderServiceWith(repo *fakeOrderRepo) *OrderService {
	return NewOrderService(repo, repo, repo)
}

func TestCreateOrder_CalculatesTotalAndStatus(t *testing.T) {
	repo := newFakeOrderRepo()

	err := newOrderServiceWith(repo).Create(context.Background(), CreateOrderInput{
		UserID: 1,
		Items: []CreateOrderItemInput{
			{ProductID: 10, Quantity: 2},
			{ProductID: 20, Quantity: 3},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.created.TotalCents != 2*1500+3*250 {
		t.Errorf("total = %d, want %d", repo.created.TotalCents, 2*1500+3*250)
	}
	if repo.created.Status != models.OrderStatusPending {
		t.Errorf("status = %q, want %q", repo.created.Status, models.OrderStatusPending)
	}
	if repo.created.Items[0].UnitPriceCents != 1500 {
		t.Errorf("unit price = %d, want 1500", repo.created.Items[0].UnitPriceCents)
	}
}

func TestCreateOrder_Errors(t *testing.T) {
	tests := []struct {
		name string
		in   CreateOrderInput
		want error // nil => espera ValidationError
	}{
		{"sem itens", CreateOrderInput{UserID: 1}, nil},
		{"quantidade zero", CreateOrderInput{UserID: 1, Items: []CreateOrderItemInput{{ProductID: 10, Quantity: 0}}}, nil},
		{"produto duplicado", CreateOrderInput{UserID: 1, Items: []CreateOrderItemInput{{10, 1}, {10, 2}}}, nil},
		{"usuário inexistente", CreateOrderInput{UserID: 99, Items: []CreateOrderItemInput{{10, 1}}}, models.ErrUserNotFound},
		{"produto inexistente", CreateOrderInput{UserID: 1, Items: []CreateOrderItemInput{{10, 1}, {99, 1}}}, models.ErrProductNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeOrderRepo()
			err := newOrderServiceWith(repo).Create(context.Background(), tt.in)

			if tt.want == nil {
				var vErr *models.ValidationError
				if !errors.As(err, &vErr) {
					t.Fatalf("err = %v, want ValidationError", err)
				}
			} else {
				if !errors.Is(err, tt.want) {
					t.Fatalf("err = %v, want %v", err, tt.want)
				}
				// Referência inválida num pedido é 422, não 404.
				if !errors.Is(err, models.ErrInvalidReference) {
					t.Errorf("err = %v, want it to wrap ErrInvalidReference", err)
				}
			}
			if repo.created != nil {
				t.Error("order should not have been created")
			}
		})
	}
}
