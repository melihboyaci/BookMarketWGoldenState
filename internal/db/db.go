package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	// PostgreSQL sürücüsünü yan etki (side-effect) olarak import ediyoruz.
	// Bu, database/sql'in "postgres" sürücüsünü tanımasını sağlar.
	_ "github.com/lib/pq"
)

// DB, uygulama genelinde kullanılacak tek veritabanı bağlantı havuzudur.
var DB *sql.DB

// Connect, .env dosyasındaki değişkenleri okuyarak PostgreSQL'e bağlanır
// ve bağlantı havuzunu yapılandırır. Uygulama başlangıcında bir kez çağrılmalıdır.
func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		// Bağlantı dizesi hatalıysa uygulamayı başlatmak anlamsızdır.
		log.Fatalf("HATA: Veritabanı bağlantısı açılamadı: %v", err)
	}

	// Bağlantı havuzu ayarları
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)

	// Gerçek bir TCP bağlantısı kuruluyor mu diye kontrol ediyoruz.
	if err = DB.Ping(); err != nil {
		log.Fatalf("HATA: Veritabanına ping atılamadı. Sunucu çalışıyor mu? Hata: %v", err)
	}

	log.Println("BİLGİ: PostgreSQL bağlantısı başarıyla kuruldu.")
}
