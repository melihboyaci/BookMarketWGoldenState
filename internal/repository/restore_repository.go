package repository

import (
	"database/sql"
	"fmt"
	"time"
)

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
//  1. Tüm demo_active satırlarını sil (DELETE)
//  2. demo_blueprint satırlarını 'demo_active' tenant_id ile kopyala (INSERT ... SELECT)
//  3. Commit — başarılı ise sonucu döndür, hata varsa Rollback ile geri al
func RestoreGoldenState(db *sql.DB) (*RestoreResult, error) {
	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("transaction başlatılamadı: %w", err)
	}

	// Adım 1: demo_active kayıtlarını sil
	delRes, err := tx.Exec(`DELETE FROM books WHERE tenant_id = 'demo_active'`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_active kayıtları silinemedi: %w", err)
	}
	deletedRows, _ := delRes.RowsAffected()

	// Adım 2: demo_blueprint'i demo_active olarak kopyala.
	// id sütunu BIGSERIAL olduğu için SELECT listesine dahil edilmez;
	// her kopyalanan kitap yeni ve benzersiz bir id alır.
	insRes, err := tx.Exec(`
		INSERT INTO books (tenant_id, title, author, isbn, price, stock, created_at, updated_at)
		SELECT 'demo_active', title, author, isbn, price, stock, NOW(), NOW()
		FROM books
		WHERE tenant_id = 'demo_blueprint'
	`)
	if err != nil {
		_ = tx.Rollback()
		return nil, fmt.Errorf("demo_blueprint kopyalanamadı: %w", err)
	}
	insertedRows, _ := insRes.RowsAffected()

	// Adım 3: Değişiklikleri kalıcı hale getir
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("transaction commit edilemedi: %w", err)
	}

	return &RestoreResult{
		DeletedRows:  deletedRows,
		InsertedRows: insertedRows,
		DurationMs:   time.Since(start).Milliseconds(),
	}, nil
}
