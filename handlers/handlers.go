package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"crud-api/models"
)

func RegisterHandlers(mux *http.ServeMux, users *UserHandler, products *ProductHandler, orders *OrderHandler) {
	// User routes
	mux.HandleFunc("GET /users/{id}", users.GetUserByID)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("PUT /users/{id}", users.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", users.DeleteUser)

	// Product routes
	mux.HandleFunc("GET /products", products.GetProducts)
	mux.HandleFunc("GET /products/{id}", products.GetProductByID)
	mux.HandleFunc("POST /products", products.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", products.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", products.DeleteProduct)

	// Order routes
	mux.HandleFunc("GET /orders/{id}", orders.GetOrderByID)
	mux.HandleFunc("POST /orders", orders.CreateOrder)
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(writer http.ResponseWriter, status int, msg string) {
	writeJSON(writer, status, map[string]string{"error": msg})
}

// writeServiceError é o único lugar que traduz erros de domínio em status HTTP.
// Erros desconhecidos vão para o log e o cliente recebe uma mensagem genérica.
func writeServiceError(writer http.ResponseWriter, err error) {
	var vErr *models.ValidationError
	switch {
	case errors.As(err, &vErr):
		writeError(writer, http.StatusBadRequest, vErr.Error())
	case errors.Is(err, models.ErrOrderNotFound):
		writeError(writer, http.StatusNotFound, err.Error())
	case errors.Is(err, models.ErrUserNotFound), errors.Is(err, models.ErrProductNotFound):
		writeError(writer, http.StatusUnprocessableEntity, err.Error())
	default:
		log.Printf("internal error: %v", err)
		writeError(writer, http.StatusInternalServerError, "internal error")
	}
}
