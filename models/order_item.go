package models

type OrderItem struct {
	OrderID        int   `json:"order_id"`
	ProductID      int   `json:"product_id"`
	Quantity       int   `json:"quantity"`
	UnitPriceCents int64 `json:"unit_price_cents"`
}
