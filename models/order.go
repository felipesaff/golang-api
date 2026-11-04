package models

const OrderStatusPending = "pending"

// Order é o pedido como o domínio o enxerga: o que será gravado no banco.
type Order struct {
	ID         int64
	UserID     int64
	TotalCents int64
	Status     string
	Items      []OrderItem
}

// OrderDetails é a visão de leitura de um pedido, com dados do usuário e dos produtos.
type OrderDetails struct {
	ID         int64
	TotalCents int64
	Customer   OrderCustomer
	Items      []OrderDetailsItem
}

type OrderCustomer struct {
	Name     string
	Email    string
	Document string
}

type OrderDetailsItem struct {
	ProductName        string
	ProductDescription *string
	Quantity           int
	UnitPriceCents     int64
}
