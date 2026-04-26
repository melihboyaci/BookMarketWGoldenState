package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bookmarket/golden-state/internal/middleware"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Testler için yardımcı fonksiyon: İstenilen rolde geçerli veya süresi geçmiş JWT üretir.
func generateTestToken(role models.Role, secret string, expired bool) string {
	expiration := time.Now().Add(1 * time.Hour)
	if expired {
		expiration = time.Now().Add(-1 * time.Hour)
	}

	claims := models.DemoClaims{
		UserID: 1,
		Email:  "test@demo.com",
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestRequireAdminRole(t *testing.T) {
	// Gin debug modunu kapat (konsolun kirlenmemesi için)
	gin.SetMode(gin.TestMode)

	// Test ortamı için JWT_SECRET ayarla
	os.Setenv("JWT_SECRET", "test-secret-key")

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		expectedStatus int
	}{
		{
			name: "Başarılı: Admin Rolü",
			setupRequest: func(req *http.Request) {
				token := generateTestToken(models.RoleAdmin, "test-secret-key", false)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Hata: Yetkisiz Rol (Seller)",
			setupRequest: func(req *http.Request) {
				token := generateTestToken(models.RoleSeller, "test-secret-key", false)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusForbidden, // 403
		},
		{
			name: "Hata: Authorization Başlığı Eksik",
			setupRequest: func(req *http.Request) {
				// Header eklenmiyor
			},
			expectedStatus: http.StatusUnauthorized, // 401
		},
		{
			name: "Hata: Hatalı Format (Bearer Yok)",
			setupRequest: func(req *http.Request) {
				token := generateTestToken(models.RoleAdmin, "test-secret-key", false)
				req.Header.Set("Authorization", token) // "Bearer " prefixi eksik
			},
			expectedStatus: http.StatusUnauthorized, // 401
		},
		{
			name: "Hata: Süresi Dolmuş Token",
			setupRequest: func(req *http.Request) {
				token := generateTestToken(models.RoleAdmin, "test-secret-key", true)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusUnauthorized, // 401
		},
		{
			name: "Hata: Yanlış Secret Key ile İmzalanmış",
			setupRequest: func(req *http.Request) {
				token := generateTestToken(models.RoleAdmin, "wrong-secret", false)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusUnauthorized, // 401
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Sahte (Mock) Response Recorder
			w := httptest.NewRecorder()
			// Sahte Gin Context oluştur
			_, r := gin.CreateTestContext(w)

			// Test edilecek route ve middleware tanımı
			r.GET("/protected", middleware.RequireAdminRole(), func(c *gin.Context) {
				// Middleware başarılı olursa bu handler çalışır ve 200 döner
				c.String(http.StatusOK, "success")
			})

			// HTTP isteği oluştur
			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			tt.setupRequest(req)

			// İsteği router'a gönder
			r.ServeHTTP(w, req)

			// Sonucu doğrula
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
