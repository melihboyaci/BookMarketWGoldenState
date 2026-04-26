package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/bookmarket/golden-state/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAdminRole, Authorization başlığındaki JWT token'ı doğrular
// ve yalnızca "ADMIN" rolüne sahip isteklerin geçmesine izin verir.
// Diğer tüm istekler 401 veya 403 ile reddedilir.
func RequireAdminRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// "Bearer <token>" formatı kontrolü
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Yetkilendirme başlığı eksik veya hatalı formatta. 'Bearer <token>' bekleniyor.",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.DemoClaims{}

		// Token'ı parse et ve imzayı doğrula
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			// Yalnızca HMAC imzalama yöntemine izin ver
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Geçersiz veya süresi dolmuş token.",
			})
			return
		}

		// Rol kontrolü: Yalnızca ADMIN geçebilir
		if claims.Role != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Bu işlem için Admin yetkisi gereklidir.",
			})
			return
		}

		// Doğrulanmış claim bilgilerini context'e aktar (sonraki handler'lar için)
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", string(claims.Role))

		c.Next()
	}
}
