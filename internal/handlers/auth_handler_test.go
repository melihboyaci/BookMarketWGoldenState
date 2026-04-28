package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// newTestRouter, verilen db ile login+register route'larını kurar.
func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	cfg := &config.Config{JWTSecret: "test-secret"}
	// Testler için nil repo geçiyoruz (binding testleri için yeterli)
	authHandler := handlers.NewAuthHandler(nil, cfg)

	r.POST("/api/v1/auth/login", authHandler.Login())
	r.POST("/api/v1/auth/register", authHandler.Register())
	return r
}

// TestLoginHandler_InvalidBody, DB bağlantısına gerek duymayan (binding katmanı)
// hata senaryolarını test eder.
func TestLoginHandler_InvalidBody(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	r := newTestRouter() // DB'ye ulaşılmadan 400 dönmeli

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "Hata: Email Alanı Eksik",
			body:           map[string]string{"password": "admin123"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Hata: Geçersiz Email Formatı",
			body:           map[string]string{"email": "notanemail", "password": "admin123"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Hata: Şifre Alanı Eksik",
			body:           map[string]string{"email": "admin@demo.com"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

// TestRegisterHandler_InvalidBody, Register endpoint'inin binding kurallarını test eder.
func TestRegisterHandler_InvalidBody(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret")
	r := newTestRouter() // DB'ye ulaşılmadan 400 dönmeli

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
	}{
		{
			name:           "Hata: Kısa Kullanıcı Adı (1 karakter)",
			body:           map[string]string{"username": "a", "email": "test@demo.com", "password": "pass123"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Hata: Kısa Şifre (5 karakter)",
			body:           map[string]string{"username": "Test User", "email": "test@demo.com", "password": "12345"},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Hata: Email Alanı Eksik",
			body:           map[string]string{"username": "Test User", "password": "pass123"},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonValue, _ := json.Marshal(tt.body)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(jsonValue))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
