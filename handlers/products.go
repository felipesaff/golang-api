package handlers

import (
	"crud-api/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	DB *sql.DB
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{DB: db}
}

func (productHandler *ProductHandler) GetProducts(writer http.ResponseWriter, req *http.Request) {
	rows, err := productHandler.DB.Query("SELECT * FROM products")
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.PriceCents, &product.Stock, &product.CreatedAt)
		if err != nil {
			writeError(writer, http.StatusInternalServerError, err.Error())
			return
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, products)
}

func (productHandler *ProductHandler) GetProductByID(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product
	row := productHandler.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)
	err = row.Scan(&product.ID, &product.Name, &product.Description, &product.PriceCents, &product.Stock, &product.CreatedAt)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, product)
}

func (productHandler *ProductHandler) CreateProduct(writer http.ResponseWriter, req *http.Request) {
	var product models.Product
	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}
	if product.Name == "" || product.Description == "" || product.PriceCents < 1 || product.Stock < 1 {
		writeError(writer, http.StatusBadRequest, "name, description are required; price_cents and stock must be >= 1")
		return
	}

	err := productHandler.DB.QueryRowContext(req.Context(),
		`INSERT INTO products (name, description, price_cents, stock)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		product.Name, product.Description, product.PriceCents, product.Stock).Scan(&product.ID)
	if err != nil {
		log.Printf("create product: %v", err)
		writeError(writer, http.StatusInternalServerError, "internal error")
		return
	}

	writeJSON(writer, http.StatusCreated, product)
}

func (productHandler *ProductHandler) UpdateProduct(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var product models.Product

	if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if product.Name == "" || product.Description == "" || product.PriceCents < 1 || product.Stock < 1 {
		writeError(writer, http.StatusBadRequest, "name, description are required; price_cents and stock must be >= 1")
		return
	}

	_, err = productHandler.DB.Exec("UPDATE products SET name = $1, description = $2, price_cents = $3, stock = $4 WHERE id = $5", product.Name, product.Description, product.PriceCents, product.Stock, id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	row := productHandler.DB.QueryRow("SELECT * FROM products WHERE id = $1", id)
	var updatedProduct models.Product
	err = row.Scan(&updatedProduct.ID, &updatedProduct.Name, &updatedProduct.Description, &updatedProduct.PriceCents, &updatedProduct.Stock, &updatedProduct.CreatedAt)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, updatedProduct)
}

func (productHandler *ProductHandler) DeleteProduct(writer http.ResponseWriter, req *http.Request) {
	id, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "Invalid product ID")
		return
	}

	_, err = productHandler.DB.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		writeError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(writer, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
