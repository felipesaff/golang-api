package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	orderdto "crud-api/dto/order"
	"crud-api/models"
	"crud-api/services"
)

// OrderService é o que o handler precisa da camada de negócio.
type OrderService interface {
	Create(ctx context.Context, in services.CreateOrderInput) error
	GetByID(ctx context.Context, id int64) (models.OrderDetails, error)
}

type OrderHandler struct {
	service OrderService
}

func NewOrderHandler(service OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetOrderByID(writer http.ResponseWriter, req *http.Request) {
	id, ok := parseID(writer, req, "order")
	if !ok {
		return
	}

	details, err := h.service.GetByID(req.Context(), id)
	if err != nil {
		writeServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, toOrderResponse(details))
}

func (h *OrderHandler) CreateOrder(writer http.ResponseWriter, req *http.Request) {
	var body orderdto.CreateOrderRequestDto
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid json")
		return
	}

	if err := h.service.Create(req.Context(), toCreateOrderInput(body)); err != nil {
		writeServiceError(writer, err)
		return
	}

	writeJSON(writer, http.StatusCreated, map[string]string{"status": "success"})
}

func toCreateOrderInput(body orderdto.CreateOrderRequestDto) services.CreateOrderInput {
	in := services.CreateOrderInput{
		UserID: body.UserID,
		Items:  make([]services.CreateOrderItemInput, len(body.Items)),
	}
	for i, item := range body.Items {
		in.Items[i] = services.CreateOrderItemInput{ProductID: item.ProductID, Quantity: item.Quantity}
	}
	return in
}

func toOrderResponse(d models.OrderDetails) orderdto.OrderResponse {
	res := orderdto.OrderResponse{
		OrderID:         d.ID,
		TotalPriceCents: d.TotalCents,
		Items:           make([]orderdto.OrderItem, len(d.Items)),
	}
	res.User.Name = d.Customer.Name
	res.User.Email = d.Customer.Email
	res.User.Document = d.Customer.Document

	for i, item := range d.Items {
		res.Items[i] = orderdto.OrderItem{
			Name:        item.ProductName,
			Description: item.ProductDescription,
			Quantity:    item.Quantity,
			PriceCents:  item.UnitPriceCents,
		}
	}
	return res
}
