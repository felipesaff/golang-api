package handlers

import (
	"crud-api/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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

	_, err := userHandler.DB.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, "Failed to create user")
		return
	}

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
