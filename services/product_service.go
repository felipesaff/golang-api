package services

import (
	"context"

	"crud-api/models"
)

type ProductRepository interface {
	List(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int64) (models.Product, error)
	Create(ctx context.Context, p models.Product) (models.Product, error)
	Update(ctx context.Context, p models.Product) (models.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// ProductInput é usado tanto na criação quanto na atualização.
type ProductInput struct {
	Name        string
	Description string
	PriceCents  int64
	Stock       int
}

func (in ProductInput) validate() error {
	if isBlank(in.Name) || isBlank(in.Description) || in.PriceCents < 1 || in.Stock < 1 {
		return &models.ValidationError{Msg: "name, description are required; price_cents and stock must be >= 1"}
	}
	return nil
}

func (in ProductInput) toModel(id int64) models.Product {
	return models.Product{
		ID:          id,
		Name:        in.Name,
		Description: in.Description,
		PriceCents:  in.PriceCents,
		Stock:       in.Stock,
	}
}

func (s *ProductService) List(ctx context.Context) ([]models.Product, error) {
	return s.repo.List(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (models.Product, error) {
	if err := validateID(id, "product"); err != nil {
		return models.Product{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) Create(ctx context.Context, in ProductInput) (models.Product, error) {
	if err := in.validate(); err != nil {
		return models.Product{}, err
	}
	return s.repo.Create(ctx, in.toModel(0))
}

func (s *ProductService) Update(ctx context.Context, id int64, in ProductInput) (models.Product, error) {
	if err := validateID(id, "product"); err != nil {
		return models.Product{}, err
	}
	if err := in.validate(); err != nil {
		return models.Product{}, err
	}
	return s.repo.Update(ctx, in.toModel(id))
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	if err := validateID(id, "product"); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
