package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/bookmarket/golden-state/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_SECRET", "test-secret")

	// Test edilecek HTTP endpoint'ini kur
	r := gin.Default()
	r.POST("/api/v1/auth/login", handlers.Login)

	tests := []struct {
		name           string
		body           map[string]string
		expectedStatus int
		expectToken    bool
	}{
		{
			name: "Başarılı Login: Admin",
			body: map[string]string{
				"email":    "admin@demo.com",
				"password": "admin123",
			},
			expectedStatus: http.StatusOK,
			expectToken:    true,
		},
		{
			name: "Hata: Yanlış Şifre",
			body: map[string]string{
				"email":    "admin@demo.com",
				"password": "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized, // 401
			expectToken:    false,
		},
		{
			name: "Hata: Olmayan Kullanıcı",
			body: map[string]string{
				"email":    "unknown@demo.com",
				"password": "admin123",
			},
			expectedStatus: http.StatusUnauthorized, // 401
			expectToken:    false,
		},
		{
			name: "Hata: Eksik JSON Body (Email Yok)",
			body: map[string]string{
				"password": "admin123",
			},
			expectedStatus: http.StatusBadRequest, // 400
			expectToken:    false,
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

			if tt.expectToken {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "token")
				assert.Contains(t, response, "role")
			}
		})
	}
}
