package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/bookmarket/golden-state/internal/config"
)

// RestoreResult, Golden State sıfırlama işleminin sonucunu tutar.
type RestoreResult struct {
	DeletedRows  int64
	InsertedRows int64
	DurationMs   int64
}

type RestoreRepository interface {
	RestoreGoldenState() (*RestoreResult, error)
}

type PostgresRestoreRepository struct {
	db        *sql.DB
	restoreMu sync.Mutex
}

func NewRestoreRepository(db *sql.DB) RestoreRepository {
	return &PostgresRestoreRepository{db: db}
}

// RestoreGoldenState, demo_active verilerini silerek demo_blueprint'in
// birebir kopyasını demo_active olarak yeniden oluşturur.
func (r *PostgresRestoreRepository) RestoreGoldenState() (*RestoreResult, error) {
	if !r.restoreMu.TryLock() {
		return nil, fmt.Errorf("başka bir restore işlemi zaten devam ediyor")
	}
	defer r.restoreMu.Unlock()

	start := time.Now()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction başlatılamadı: %w", err)
	}

	// ── Adım 1: DELETE (FK sırasına göre) ────────────────────────────────────

	// 1a. Önce orders
	resOrders, err := tx.Exec(`DELETE FROM orders WHERE tenant_id = $1`, config.TenantActive)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active orders silinemedi: %w", err)
	}
	delOrders, _ := resOrders.RowsAffected()

	// 1b. Sonra users
	resUsers, err := tx.Exec(`DELETE FROM users WHERE tenant_id = $1`, config.TenantActive)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active users silinemedi: %w", err)
	}
	delUsers, _ := resUsers.RowsAffected()

	// 1c. Son olarak books
	resBooks, err := tx.Exec(`DELETE FROM books WHERE tenant_id = $1`, config.TenantActive)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active books silinemedi: %w", err)
	}
	delBooks, _ := resBooks.RowsAffected()
	deletedRows := delOrders + delUsers + delBooks

	// ── Adım 2: INSERT ... SELECT (Bağımlılık sırasına göre) ─────────────────

	// 2a. Önce books
	insBooks, err := tx.Exec(`
		INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock, created_at, updated_at)
		SELECT $1, title, author, isbn, image_url, price, stock, NOW(), NOW()
		FROM books
		WHERE tenant_id = $2
	`, config.TenantActive, config.TenantBlueprint)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_blueprint books kopyalanamadı: %w", err)
	}
	insertedBooksRows, _ := insBooks.RowsAffected()

	// 2b. Sonra users
	insUsers, err := tx.Exec(`
		INSERT INTO users (username, email, password_hash, role, tenant_id, created_at)
		SELECT username, email, password_hash, role, $1, NOW()
		FROM users
		WHERE tenant_id = $2
	`, config.TenantActive, config.TenantBlueprint)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_blueprint users kopyalanamadı: %w", err)
	}
	insertedUsersRows, _ := insUsers.RowsAffected()

	// ── Adım 3: Değişiklikleri kalıcı hale getir ─────────────────────────────
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit edilemedi: %w", err)
	}

	return &RestoreResult{
		DeletedRows:  deletedRows,
		InsertedRows: insertedBooksRows + insertedUsersRows,
		DurationMs:   time.Since(start).Milliseconds(),
	}, nil
}
