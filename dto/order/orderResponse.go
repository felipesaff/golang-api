package order

type OrderResponse struct {
	OrderID         int64       `json:"order_id"`
	TotalPriceCents int64       `json:"total_price_cents"`
	User            orderUser   `json:"user"`
	Items           []OrderItem `json:"items"`
}

type orderUser struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
}

type OrderItem struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Quantity    int     `json:"quantity"`
	PriceCents  int64   `json:"price_cents"`
}
