package repository

import (
	"database/sql"
	"fmt"

	"github.com/bookmarket/golden-state/internal/models"
)

func GetAllBooks(db *sql.DB, tenantID string) ([]models.Book, error) {
	query := `
		SELECT id, tenant_id, title, author, isbn, image_url, price, stock, created_at, updated_at
		FROM books
		WHERE tenant_id = $1
		ORDER BY id DESC
	`
	rows, err := db.Query(query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		var imgURL sql.NullString
		if err := rows.Scan(&b.ID, &b.TenantID, &b.Title, &b.Author, &b.ISBN, &imgURL, &b.Price, &b.Stock, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		b.ImageURL = imgURL.String
		books = append(books, b)
	}
	return books, nil
}

func CreateBook(db *sql.DB, book *models.Book) error {
	query := `
		INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	return db.QueryRow(
		query, book.TenantID, book.Title, book.Author, book.ISBN, book.ImageURL, book.Price, book.Stock,
	).Scan(&book.ID, &book.CreatedAt, &book.UpdatedAt)
}

func UpdateBook(db *sql.DB, book *models.Book) error {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, image_url = $4, price = $5, stock = $6, updated_at = NOW()
		WHERE id = $7 AND tenant_id = $8
		RETURNING updated_at
	`
	err := db.QueryRow(
		query, book.Title, book.Author, book.ISBN, book.ImageURL, book.Price, book.Stock, book.ID, book.TenantID,
	).Scan(&book.UpdatedAt)
	if err == sql.ErrNoRows {
		return fmt.Errorf("book not found")
	}
	return err
}

func DeleteBook(db *sql.DB, id int64, tenantID string) error {
	query := `DELETE FROM books WHERE id = $1 AND tenant_id = $2`
	res, err := db.Exec(query, id, tenantID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("book not found")
	}
	return nil
}
