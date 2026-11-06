package services

import (
	"context"

	"crud-api/models"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (models.User, error)
	Create(ctx context.Context, u models.User) (models.User, error)
	Update(ctx context.Context, u models.User) (models.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

type CreateUserInput struct {
	Name     string
	Email    string
	Document string
}

// UpdateUserInput não tem Document: o documento identifica o usuário e não é alterável.
type UpdateUserInput struct {
	Name string
}

func (s *UserService) GetByID(ctx context.Context, id int64) (models.User, error) {
	if err := validateID(id, "user"); err != nil {
		return models.User{}, err
	}
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, in CreateUserInput) (models.User, error) {
	switch {
	case isBlank(in.Name):
		return models.User{}, &models.ValidationError{Msg: "name is required"}
	case isBlank(in.Email):
		return models.User{}, &models.ValidationError{Msg: "email is required"}
	case isBlank(in.Document):
		return models.User{}, &models.ValidationError{Msg: "document is required"}
	}
	return s.repo.Create(ctx, models.User{Name: in.Name, Email: in.Email, Document: in.Document})
}

func (s *UserService) Update(ctx context.Context, id int64, in UpdateUserInput) (models.User, error) {
	if err := validateID(id, "user"); err != nil {
		return models.User{}, err
	}
	switch {
	case isBlank(in.Name):
		return models.User{}, &models.ValidationError{Msg: "name is required"}
	}
	return s.repo.Update(ctx, models.User{ID: id, Name: in.Name})
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	if err := validateID(id, "user"); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
