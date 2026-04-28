package repository_test

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bookmarket/golden-state/internal/repository"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// setupTestDB, testcontainers kullanarak geçici bir PostgreSQL ayağa kaldırır
// ve migrasyon/seed işlemlerini uygular.
func setupTestDB(ctx context.Context, t *testing.T) (*postgres.PostgresContainer, *sql.DB) {
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:16-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	require.NoError(t, err)

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	// Bağlantıyı test et
	require.NoError(t, db.Ping())

	// Migrasyonları ve Seed'i çalıştır
	runSQLFile(t, db, "../db/migrations/001_create_books_table.sql")
	runSQLFile(t, db, "../db/migrations/002_fix_isbn_tenant_unique.sql")
	runSQLFile(t, db, "../db/migrations/003_create_orders_table.sql")
	runSQLFile(t, db, "../db/migrations/004_add_image_url_to_books.sql")
	runSQLFile(t, db, "../db/migrations/005_create_users_table.sql")
	runSQLFile(t, db, "../db/seed.sql")

	return pgContainer, db
}

// runSQLFile, verilen yoldaki SQL dosyasını okur ve veritabanında çalıştırır.
func runSQLFile(t *testing.T, db *sql.DB, path string) {
	content, err := os.ReadFile(filepath.Clean(path))
	require.NoError(t, err, "SQL dosyası okunamadı: %s", path)

	_, err = db.Exec(string(content))
	require.NoError(t, err, "SQL dosyası çalıştırılamadı: %s", path)
}

func TestRestoreGoldenState(t *testing.T) {
	// Test ortamında container ayağa kalkması 5-10 saniye sürebilir.
	if testing.Short() {
		t.Skip("Kısa test modunda Integration Testler atlanıyor.")
	}

	ctx := context.Background()
	pgContainer, db := setupTestDB(ctx, t)
	defer func() {
		db.Close()
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Fatalf("Testcontainer kapatılamadı: %s", err)
		}
	}()

	t.Run("Başarılı: Golden State Restore", func(t *testing.T) {
		// ADIM 1: Veriyi Kirlet (Active ortamındaki tüm kitapların fiyatını 999 yap, birini sil)
		_, err := db.Exec("UPDATE books SET price = 999.99 WHERE tenant_id = 'demo_active'")
		require.NoError(t, err)
		
		_, err = db.Exec("DELETE FROM books WHERE tenant_id = 'demo_active' AND isbn = '978-9750718908'")
		require.NoError(t, err)

		// Kirlenmiş veriyi kontrol et (silinenle birlikte 9 tane kalmış olmalı)
		var activeCount int
		err = db.QueryRow("SELECT count(*) FROM books WHERE tenant_id = 'demo_active'").Scan(&activeCount)
		require.NoError(t, err)
		assert.Equal(t, 9, activeCount)

		// ADIM 2: Restore işlemini çağır
		restoreRepo := repository.NewRestoreRepository(db)
		result, err := restoreRepo.RestoreGoldenState()
		
		// ADIM 3: Sonuçları doğrula
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, int64(12), result.DeletedRows) // 9 kirli kitap + 3 kullanıcı silindi
		assert.Equal(t, int64(13), result.InsertedRows) // 10 kitap + 3 kullanıcı kopyalandı
		assert.GreaterOrEqual(t, result.DurationMs, int64(0)) // Süre 0 veya daha büyük olmalı

		// ADIM 4: Veritabanındaki güncel durumu doğrula (Tekrar 10'a çıkmalı ve orijinal fiyata dönmeli)
		err = db.QueryRow("SELECT count(*) FROM books WHERE tenant_id = 'demo_active'").Scan(&activeCount)
		require.NoError(t, err)
		assert.Equal(t, 10, activeCount)

		var price float64
		err = db.QueryRow("SELECT price FROM books WHERE tenant_id = 'demo_active' AND isbn = '978-9750718908'").Scan(&price)
		require.NoError(t, err)
		assert.Equal(t, 185.00, price) // Orijinal şablondaki (Saatleri Ayarlama Enstitüsü) fiyata dönmüş olmalı
	})
}
