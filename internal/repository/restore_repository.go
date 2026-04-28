package repository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// restoreMu, aynı anda yalnızca bir restore işleminin çalışmasını garanti eder.
// TryLock kullanıyoruz: eş zamanlı istek gelirse bekletmek yerine anında reddediyoruz.
var restoreMu sync.Mutex

// RestoreResult, Golden State sıfırlama işleminin sonucunu tutar.
// Bu bilgiler API yanıtında demo kullanıcısına gösterilmek üzere döner.
type RestoreResult struct {
	DeletedRows  int64
	InsertedRows int64
	DurationMs   int64
}

// RestoreGoldenState, demo_active verilerini silerek demo_blueprint'in
// birebir kopyasını demo_active olarak yeniden oluşturur.
//
// İş adımları (tümü tek bir Transaction içinde):
//
//  1. DELETE (Foreign Key sırasına göre): orders → users → books
//     (Sadece demo_active olanlar silinir)
//
//  2. INSERT ... SELECT (Bağımlılık sırasına göre): books → users
//     (demo_blueprint satırları 'demo_active' tenant_id ile kopyalanır)
//
//  3. Commit — başarılı ise sonucu döndür, hata varsa Rollback ile geri al
func RestoreGoldenState(db *sql.DB) (*RestoreResult, error) {
	// Eş zamanlı restore isteklerini önle: kilit alınamazsa işlemi hemen reddet
	if !restoreMu.TryLock() {
		return nil, fmt.Errorf("başka bir restore işlemi zaten devam ediyor")
	}
	defer restoreMu.Unlock()

	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction başlatılamadı: %w", err)
	}

	// ── Adım 1: DELETE (FK sırasına göre) ────────────────────────────────────

	// 1a. Önce orders
	resOrders, err := tx.Exec(`DELETE FROM orders WHERE tenant_id = 'demo_active'`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active orders silinemedi: %w", err)
	}
	delOrders, _ := resOrders.RowsAffected()

	// 1b. Sonra users
	resUsers, err := tx.Exec(`DELETE FROM users WHERE tenant_id = 'demo_active'`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active users silinemedi: %w", err)
	}
	delUsers, _ := resUsers.RowsAffected()

	// 1c. Son olarak books
	resBooks, err := tx.Exec(`DELETE FROM books WHERE tenant_id = 'demo_active'`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active books silinemedi: %w", err)
	}
	delBooks, _ := resBooks.RowsAffected()
	deletedRows := delOrders + delUsers + delBooks

	// ── Adım 2: INSERT ... SELECT (Bağımlılık sırasına göre) ─────────────────

	// 2a. Önce books: users'dan bağımsız, önce oluşturulmalı
	// id BIGSERIAL olduğu için SELECT listesine dahil edilmez; yeni id alır
	insBooks, err := tx.Exec(`
		INSERT INTO books (tenant_id, title, author, isbn, image_url, price, stock, created_at, updated_at)
		SELECT 'demo_active', title, author, isbn, image_url, price, stock, NOW(), NOW()
		FROM books
		WHERE tenant_id = 'demo_blueprint'
	`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_blueprint books kopyalanamadı: %w", err)
	}
	insertedBooksRows, _ := insBooks.RowsAffected()

	// 2b. Sonra users: id UUID olduğu için SELECT listesine dahil edilmez; yeni UUID alır
	insUsers, err := tx.Exec(`
		INSERT INTO users (username, email, password_hash, role, tenant_id, created_at)
		SELECT username, email, password_hash, role, 'demo_active', NOW()
		FROM users
		WHERE tenant_id = 'demo_blueprint'
	`)
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
