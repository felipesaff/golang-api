package services

import (
	"context"
	"errors"
	"testing"

	"crud-api/models"
)

type fakeProductRepo struct {
	updated *models.Product
}

func (f *fakeProductRepo) List(_ context.Context) ([]models.Product, error) { return nil, nil }

func (f *fakeProductRepo) GetByID(_ context.Context, id int64) (models.Product, error) {
	return models.Product{ID: id}, nil
}

func (f *fakeProductRepo) Create(_ context.Context, p models.Product) (models.Product, error) {
	return p, nil
}

func (f *fakeProductRepo) Update(_ context.Context, p models.Product) (models.Product, error) {
	f.updated = &p
	return p, nil
}

func (f *fakeProductRepo) Delete(_ context.Context, _ int64) error { return nil }

func validProductInput() ProductInput {
	return ProductInput{Name: "Caneta", Description: "Azul", PriceCents: 350, Stock: 10}
}

func TestProductService_CreateValidation(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ProductInput)
	}{
		{"sem nome", func(in *ProductInput) { in.Name = "" }},
		{"sem descrição", func(in *ProductInput) { in.Description = " " }},
		{"preço zero", func(in *ProductInput) { in.PriceCents = 0 }},
		{"estoque zero", func(in *ProductInput) { in.Stock = 0 }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := validProductInput()
			tt.mutate(&in)

			_, err := NewProductService(&fakeProductRepo{}).Create(context.Background(), in)
			var vErr *models.ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("err = %v, want ValidationError", err)
			}
		})
	}
}

func TestProductService_UpdateUsesPathID(t *testing.T) {
	repo := &fakeProductRepo{}
	if _, err := NewProductService(repo).Update(context.Background(), 7, validProductInput()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.updated.ID != 7 {
		t.Errorf("updated ID = %d, want 7", repo.updated.ID)
	}
}
