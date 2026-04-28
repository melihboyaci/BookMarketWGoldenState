package repository

import (
	"database/sql"
	"fmt"

	"github.com/bookmarket/golden-state/internal/models"
)

type OrderRepository interface {
	Checkout(tenantID string, req *models.CheckoutRequest) error
	GetSalesStats(tenantID string) ([]models.SalesStat, error)
}

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Checkout(tenantID string, req *models.CheckoutRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	// Calculate total and prepare items
	for _, item := range req.Items {
		var price float64
		var stock int
		err := tx.QueryRow(`SELECT price, stock FROM books WHERE id = $1 AND tenant_id = $2 FOR UPDATE`, item.BookID, tenantID).Scan(&price, &stock)
		if err != nil {
			_ = tx.Rollback()
			if err == sql.ErrNoRows {
				return fmt.Errorf("book id %d not found", item.BookID)
			}
			return err
		}

		if stock < item.Quantity {
			_ = tx.Rollback()
			return fmt.Errorf("not enough stock for book id %d", item.BookID)
		}

		// Reduce stock
		_, err = tx.Exec(`UPDATE books SET stock = stock - $1 WHERE id = $2 AND tenant_id = $3`, item.Quantity, item.BookID, tenantID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		totalAmount += price * float64(item.Quantity)
		orderItems = append(orderItems, models.OrderItem{
			BookID:   item.BookID,
			Quantity: item.Quantity,
			Price:    price,
		})
	}

	// Create order
	var orderID int64
	err = tx.QueryRow(`INSERT INTO orders (tenant_id, total_amount) VALUES ($1, $2) RETURNING id`, tenantID, totalAmount).Scan(&orderID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert items
	for _, item := range orderItems {
		_, err = tx.Exec(`INSERT INTO order_items (order_id, book_id, quantity, price) VALUES ($1, $2, $3, $4)`, orderID, item.BookID, item.Quantity, item.Price)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresOrderRepository) GetSalesStats(tenantID string) ([]models.SalesStat, error) {
	query := `
		SELECT 
			TO_CHAR(created_at, 'YYYY-MM-DD') as date,
			SUM(total_amount) as total_sales,
			COUNT(id) as total_orders
		FROM orders
		WHERE tenant_id = $1
		GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		ORDER BY date DESC
		LIMIT 30
	`
	rows, err := r.db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.SalesStat
	for rows.Next() {
		var s models.SalesStat
		if err := rows.Scan(&s.Date, &s.TotalSales, &s.TotalOrders); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	
	if stats == nil {
		stats = []models.SalesStat{}
	}
	
	return stats, nil
}
