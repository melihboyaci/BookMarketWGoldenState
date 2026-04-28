package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/db"
	"github.com/bookmarket/golden-state/internal/handlers"
	"github.com/bookmarket/golden-state/internal/middleware"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env dosyasını yükle
	if err := godotenv.Load(); err != nil {
		log.Println("BİLGİ: .env dosyası bulunamadı, sistem ortam değişkenleri kullanılıyor.")
	}

	// Konfigürasyonu yükle
	cfg := config.LoadConfig()

	// Veritabanı bağlantısını kur
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("HATA: Veritabanı bağlantısı kurulamadı: %v", err)
	}
	defer database.Close()

	// Uygulama her açıldığında seed verisini yükle (idempotent)
	db.Seed(database)

	// Bağımlılıkları Başlat (Dependency Injection)
	// Repositories
	userRepo := repository.NewUserRepository(database)
	bookRepo := repository.NewBookRepository(database)
	orderRepo := repository.NewOrderRepository(database)
	restoreRepo := repository.NewRestoreRepository(database)

	// Handlers
	authHandler := handlers.NewAuthHandler(userRepo, cfg)
	bookHandler := handlers.NewBookHandler(bookRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo)
	systemHandler := handlers.NewSystemHandler(restoreRepo)

	// Production modunda Gin'in debug çıktısını kapat
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// ── Router Kurulumu ──────────────────────────────────────────────
	r := gin.Default()

	// Sağlık kontrolü
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API v1 grup
	v1 := r.Group("/api/v1")
	{
		// Kimlik doğrulama
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login())
			auth.POST("/register", authHandler.Register())
		}

		// Kitaplar
		books := v1.Group("/books")
		{
			books.GET("", bookHandler.GetBooks())
			
			protectedBooks := books.Group("")
			protectedBooks.Use(middleware.RequireSellerRole())
			{
				protectedBooks.POST("", bookHandler.CreateBook())
				protectedBooks.PUT("/:id", bookHandler.UpdateBook())
				protectedBooks.DELETE("/:id", bookHandler.DeleteBook())
			}
		}

		// Checkout ve Satışlar
		orders := v1.Group("")
		orders.Use(middleware.RequireAuth())
		{
			orders.POST("/checkout", orderHandler.Checkout())
			orders.GET("/sales", orderHandler.GetSales())
		}

		// Admin rotaları
		admin := v1.Group("/admin")
		admin.Use(middleware.RequireAdminRole())
		{
			admin.POST("/system/restore", systemHandler.RestoreGoldenState())
			admin.GET("/users", authHandler.GetUsers())
		}
	}

	// ── HTTP Sunucusu ─────────────────────────────────────────────────
	srv := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	go func() {
		log.Printf("BİLGİ: Sunucu :%s portunda dinleniyor...", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HATA: Sunucu başlatılamadı: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("BİLGİ: Kapatma sinyali alındı, sunucu durduruluyor...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("UYARI: Sunucu zorla kapatıldı: %v", err)
	}

	log.Println("BİLGİ: Uygulama düzgün şekilde kapatıldı.")
}
