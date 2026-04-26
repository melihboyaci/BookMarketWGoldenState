package models

import "time"

type Order struct {
	ID          int64       `json:"id"`
	TenantID    string      `json:"tenant_id"`
	TotalAmount float64     `json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	Items       []OrderItem `json:"items,omitempty"`
}

type OrderItem struct {
	ID       int64   `json:"id"`
	OrderID  int64   `json:"order_id"`
	BookID   int64   `json:"book_id"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

type CartItemRequest struct {
	BookID   int64 `json:"book_id" binding:"required"`
	Quantity int   `json:"quantity" binding:"required,min=1"`
}

type CheckoutRequest struct {
	Items []CartItemRequest `json:"items" binding:"required,min=1"`
}

type SalesStat struct {
	Date        string  `json:"date"`
	TotalSales  float64 `json:"total_sales"`
	TotalOrders int     `json:"total_orders"`
}
