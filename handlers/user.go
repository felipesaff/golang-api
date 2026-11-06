package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	userdto "crud-api/dto/user"
	"crud-api/models"
	"crud-api/services"
)

// UserService é o que o handler precisa da camada de negócio.
type UserService interface {
	GetByID(ctx context.Context, id int64) (models.User, error)
	Create(ctx context.Context, in services.CreateUserInput) (models.User, error)
	Update(ctx context.Context, id int64, in services.UpdateUserInput) (models.User, error)
	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "user")
	if !ok {
		return
	}

	user, err := h.service.GetByID(req.Context(), id)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) CreateUser(writer http.ResponseWriter, req *http.Request) {
	var body userdto.CreateUserRequestDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}

	user, err := h.service.Create(req.Context(), services.CreateUserInput{
		Name:     body.Name,
		Email:    body.Email,
		Document: body.Document,
	})
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, toUserResponse(user))
}

func (h *UserHandler) UpdateUser(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "user")
	if !ok {
		return
	}

	var body userdto.UpdateUserRequestDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}

	user, err := h.service.Update(req.Context(), id, services.UpdateUserInput{
		Name: body.Name,
	})
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toUserResponse(user))
}

func (h *UserHandler) DeleteUser(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "user")
	if !ok {
		return
	}

	if err := h.service.Delete(req.Context(), id); err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}

func toUserResponse(u models.User) userdto.UserResponse {
	return userdto.UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Document:  u.Document,
		CreatedAt: u.CreatedAt,
	}
}
