package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	productdto "crud-api/dto/product"
	"crud-api/models"
	"crud-api/services"
)

// ProductService é o que o handler precisa da camada de negócio.
type ProductService interface {
	List(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int64) (models.Product, error)
	Create(ctx context.Context, in services.ProductInput) (models.Product, error)
	Update(ctx context.Context, id int64, in services.ProductInput) (models.Product, error)
	Delete(ctx context.Context, id int64) error
}

type ProductHandler struct {
	service ProductService
}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) GetProducts(writer http.ResponseWriter, req *http.Request) {
	products, err := h.service.List(req.Context())
	if err != nil {
		writeServiceError(writer, err)
		return
	}

	res := make([]productdto.ProductResponse, len(products))
	for i, p := range products {
		res[i] = toProductResponse(p)
	}
	writeJSON(writer, http.StatusOK, res)
}

func (h *ProductHandler) GetProductByID(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "product")
	if !ok {
		return
	}

	product, err := h.service.GetByID(req.Context(), id)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toProductResponse(product))
}

func (h *ProductHandler) CreateProduct(writer http.ResponseWriter, req *http.Request) {
	var body productdto.CreateProductDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}

	product, err := h.service.Create(req.Context(), services.ProductInput{
		Name:        body.Name,
		Description: body.Description,
		PriceCents:  body.PriceCents,
		Stock:       int(body.Stock),
	})
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, toProductResponse(product))
}

func (h *ProductHandler) UpdateProduct(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "product")
	if !ok {
		return
	}

	var body productdto.UpdateProductDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}

	product, err := h.service.Update(req.Context(), id, services.ProductInput{
		Name:        body.Name,
		Description: body.Description,
		PriceCents:  body.PriceCents,
		Stock:       int(body.Stock),
	})
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, toProductResponse(product))
}

func (h *ProductHandler) DeleteProduct(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "product")
	if !ok {
		return
	}

	if err := h.service.Delete(req.Context(), id); err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

func toProductResponse(p models.Product) productdto.ProductResponse {
	return productdto.ProductResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		PriceCents:  p.PriceCents,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
	}
}
