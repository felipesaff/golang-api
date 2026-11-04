package handlers

import (
	"crud-api/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/lib/pq"
)

type UserHandler struct {
	DB *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
	return &UserHandler{DB: db}
}

func (userHandler *UserHandler) GetUserByID(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	row := userHandler.DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
	var user models.User
	err = row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(writer, http.StatusOK, user)
}

func (userHandler *UserHandler) CreateUser(writer http.ResponseWriter, req *http.Request) {
	var user models.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Name == "" {
		writeError(writer, http.StatusBadRequest, "Invalid user name")
		return
	}
	if user.Email == "" {
		writeError(writer, http.StatusBadRequest, "Invalid user email")
		return
	}
	if user.Document == "" {
		writeError(writer, http.StatusBadRequest, "Invalid user document")
		return
	}

	tx, err := userHandler.DB.Begin()
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err := tx.Exec(`
		INSERT INTO users
		(name, email, document)
		VALUES ($1, $2, $3)`,
		user.Name,
		user.Email,
		user.Document); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			tx.Rollback()
			writeError(writer, http.StatusConflict, "Email or document already registered")
			return
		}
		tx.Rollback()
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	tx.Commit()
	writeJSON(writer, http.StatusCreated, user)
}

func (userHandler *UserHandler) UpdateUser(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var user models.User
	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	_, err = userHandler.DB.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", user.Name, user.Email, id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Failed to update user")
		return
	}

	writeJSON(writer, http.StatusOK, user)
}

func (userHandler *UserHandler) DeleteUser(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid user ID")
		return
	}

	_, err = userHandler.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	writeJSON(writer, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
