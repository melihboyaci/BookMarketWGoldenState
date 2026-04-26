package db

import (
	_ "embed"
	"log"
)

// go:embed direktifi, seed.sql dosyasını derleme zamanında binary'ye gömer.
// Bu sayede go run, go build ve production binary'de aynı şekilde çalışır.
//
//go:embed seed.sql
var seedSQL string

// Seed, uygulama her açıldığında çağrılır.
// seed.sql dosyasındaki INSERT'leri çalıştırarak demo_blueprint ve demo_active
// verilerinin veritabanında hazır olmasını sağlar.
// ON CONFLICT DO NOTHING sayesinde idempotent'tir: mevcut kayıtları bozmaz.
func Seed() {
	// Tüm seed betiğini tek bir Transaction içinde çalıştır
	tx, err := DB.Begin()
	if err != nil {
		log.Fatalf("HATA: Seed için transaction başlatılamadı: %v", err)
	}

	if _, err = tx.Exec(seedSQL); err != nil {
		// Hata durumunda değişiklikleri geri al
		_ = tx.Rollback()
		log.Fatalf("HATA: Seed SQL çalıştırılamadı: %v", err)
	}

	if err = tx.Commit(); err != nil {
		log.Fatalf("HATA: Seed transaction commit edilemedi: %v", err)
	}

	log.Println("BİLGİ: Seed verisi başarıyla yüklendi (demo_blueprint ve demo_active hazır).")
}
