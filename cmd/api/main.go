package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bookmarket/golden-state/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle (dosya yoksa ortam değişkenlerini olduğu gibi kullan)
	if err := godotenv.Load(); err != nil {
		log.Println("BİLGİ: .env dosyası bulunamadı, sistem ortam değişkenleri kullanılıyor.")
	}

	// Veritabanı bağlantısını kur
	db.Connect()
	defer db.DB.Close()

	// Uygulama her açıldığında seed verisini yükle (idempotent)
	db.Seed()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("BİLGİ: Uygulama %s portunda başlatılıyor... (Phase 3'te router eklenecek)", port)
	log.Println("BİLGİ: Çıkmak için CTRL+C tuşlarına basın.")

	// Sinyal yakalayıcı: Uygulamanın deadlock hatası vermeden beklemesini sağlar
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("BİLGİ: Uygulama kapatılıyor...")
}
