package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bookmarket/golden-state/internal/db"
	"github.com/bookmarket/golden-state/internal/handlers"
	"github.com/bookmarket/golden-state/internal/middleware"
	"github.com/gin-gonic/gin"
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

	// Production modunda Gin'in debug çıktısını kapat
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ── Router Kurulumu ──────────────────────────────────────────────
	r := gin.Default()

	// Sağlık kontrolü (yük dengeleyiciler ve Docker healthcheck için)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 grup
	v1 := r.Group("/api/v1")
	{
		// Kimlik doğrulama (herkese açık)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
		}

		// Yalnızca Admin yetkisi gerektiren korumalı rotalar
		admin := v1.Group("/system")
		admin.Use(middleware.RequireAdminRole())
		{
			// Golden State sıfırlama: demo_active'i sil, blueprint'i kopyala
			admin.POST("/restore", handlers.RestoreGoldenState)
		}
	}

	// ── HTTP Sunucusu ─────────────────────────────────────────────────
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Sunucuyu arka planda başlat
	go func() {
		log.Printf("BİLGİ: Sunucu :%s portunda dinleniyor...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HATA: Sunucu başlatılamadı: %v", err)
		}
	}()

	// Graceful shutdown: CTRL+C veya SIGTERM sinyalini bekle
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("BİLGİ: Kapatma sinyali alındı, sunucu durduruluyor...")

	// Mevcut isteklerin tamamlanması için 5 saniye bekle
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("UYARI: Sunucu zorla kapatıldı: %v", err)
	}

	log.Println("BİLGİ: Uygulama düzgün şekilde kapatıldı.")
}
