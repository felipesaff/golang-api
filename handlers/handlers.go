package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func RegisterHandlers(mux *http.ServeMux, db *sql.DB) {
	products := ProductHandler{DB: db}
	users := UserHandler{DB: db}

	// Product routes
	mux.HandleFunc("GET /products", products.GetProducts)
	mux.HandleFunc("GET /products/{id}", products.GetProductByID)
	mux.HandleFunc("POST /products", products.CreateProduct)
	mux.HandleFunc("PUT /products/{id}", products.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", products.DeleteProduct)

	// User routes
	mux.HandleFunc("GET /users/{id}", users.GetUserByID)
	mux.HandleFunc("POST /users", users.CreateUser)
	mux.HandleFunc("PUT /users/{id}", users.UpdateUser)
	mux.HandleFunc("DELETE /users/{id}", users.DeleteUser)
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
