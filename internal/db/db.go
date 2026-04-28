package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// PostgreSQL sürücüsünü yan etki (side-effect) olarak import ediyoruz.
	// Bu, database/sql'in "postgres" sürücüsünü tanımasını sağlar.
	_ "github.com/lib/pq"
)

// Connect, .env dosyasındaki değişkenleri okuyarak PostgreSQL'e bağlanır
// ve bağlantı havuzunu yapılandırır. Uygulama başlangıcında bir kez çağrılmalıdır.
func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("veritabanı bağlantısı açılamadı: %w", err)
	}

	// Bağlantı havuzu ayarları
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	// Stale connection hatasını önlemek için bağlantı ömrünü sınırla
	db.SetConnMaxLifetime(30 * time.Minute)

	// Gerçek bir TCP bağlantısı kuruluyor mu diye kontrol ediyoruz.
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("veritabanına ping atılamadı: %w", err)
	}

	log.Println("BİLGİ: PostgreSQL bağlantısı başarıyla kuruldu.")
	return db, nil
}
