package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/bookmarket/golden-state/internal/config"
)

// Seed, uygulama her açıldığında blueprint tenant'ından active tenant'a
// veri kopyalar. Bu işlem idempotent'tir.
func Seed(db *sql.DB) {
	// seed.sql dosyasını oku
	content, err := os.ReadFile("internal/db/seed.sql")
	if err != nil {
		log.Printf("UYARI: Seed dosyası okunamadı: %v", err)
		return
	}

	// SQL komutlarını çalıştır (blueprint verisini hazırla)
	_, err = db.Exec(string(content))
	if err != nil {
		log.Printf("UYARI: Seed işlemi sırasında hata: %v", err)
		return
	}

	// 'demo_active' tenant'ı boşsa, 'demo_blueprint'ten kopyala
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM books WHERE tenant_id = $1", config.TenantActive).Scan(&count)
	if err == nil && count == 0 {
		log.Println("BİLGİ: demo_active boş, blueprint verileri aktarılıyor...")
		copyQuery := fmt.Sprintf(`
			INSERT INTO books (title, author, price, stock, image_url, tenant_id)
			SELECT title, author, price, stock, image_url, '%s'
			FROM books WHERE tenant_id = '%s'
		`, config.TenantActive, config.TenantBlueprint)

		_, err = db.Exec(copyQuery)
		if err != nil {
			log.Printf("HATA: Veri aktarımı başarısız: %v", err)
		}
	}

	log.Println("BİLGİ: Veritabanı seed işlemi tamamlandı.")
}
