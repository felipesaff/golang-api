package services

import (
	"context"
	"errors"
	"testing"

	"crud-api/models"
)

// fakeUserRepo registra se o repositório foi chamado, para garantir que a
// validação barra a entrada antes de chegar ao banco.
type fakeUserRepo struct {
	called bool
}

func (f *fakeUserRepo) GetByID(_ context.Context, id int64) (models.User, error) {
	f.called = true
	return models.User{ID: id}, nil
}

func (f *fakeUserRepo) Create(_ context.Context, u models.User) (models.User, error) {
	f.called = true
	u.ID = 1
	return u, nil
}

func (f *fakeUserRepo) Update(_ context.Context, u models.User) (models.User, error) {
	f.called = true
	return u, nil
}

func (f *fakeUserRepo) Delete(_ context.Context, _ int64) error {
	f.called = true
	return nil
}

func TestUserService_Validation(t *testing.T) {
	tests := []struct {
		name string
		call func(s *UserService) error
	}{
		{"create sem nome", func(s *UserService) error {
			_, err := s.Create(context.Background(), CreateUserInput{Name: "  ", Email: "a@b.com", Document: "123"})
			return err
		}},
		{"create sem documento", func(s *UserService) error {
			_, err := s.Create(context.Background(), CreateUserInput{Name: "Ana", Email: "a@b.com"})
			return err
		}},
		{"update sem email", func(s *UserService) error {
			_, err := s.Update(context.Background(), 1, UpdateUserInput{Name: "Ana"})
			return err
		}},
		{"delete com id inválido", func(s *UserService) error {
			return s.Delete(context.Background(), 0)
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeUserRepo{}
			err := tt.call(NewUserService(repo))

			var vErr *models.ValidationError
			if !errors.As(err, &vErr) {
				t.Fatalf("err = %v, want ValidationError", err)
			}
			if repo.called {
				t.Error("repository should not have been called")
			}
		})
	}
}

func TestUserService_CreateValid(t *testing.T) {
	repo := &fakeUserRepo{}
	u, err := NewUserService(repo).Create(context.Background(), CreateUserInput{Name: "Ana", Email: "a@b.com", Document: "123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if u.ID != 1 || u.Name != "Ana" {
		t.Errorf("user = %+v, want ID 1 and name Ana", u)
	}
}
