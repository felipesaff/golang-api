package order

type CreateOrderRequestDto struct {
	UserID int64                       `json:"user_id"`
	Items  []createOrderItemRequestDto `json:"items"`
}

type createOrderItemRequestDto struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
