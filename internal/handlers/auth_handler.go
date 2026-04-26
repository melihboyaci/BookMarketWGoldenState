package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/bookmarket/golden-state/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// demoUsers: Demo ortamı için sabit kullanıcı listesi.
// Bu uygulama bir satış demosu olduğu için gerçek bir users tablosu gerekmez.
var demoUsers = map[string]struct {
	Password string
	Role     models.Role
	UserID   int64
}{
	"admin@demo.com":  {Password: "admin123", Role: models.RoleAdmin, UserID: 1},
	"seller@demo.com": {Password: "seller123", Role: models.RoleSeller, UserID: 2},
	"buyer@demo.com":  {Password: "buyer123", Role: models.RoleBuyer, UserID: 3},
}

type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Token string      `json:"token"`
	Role  models.Role `json:"role"`
}

// Login, POST /api/v1/auth/login endpoint'ini karşılar.
// Başarılı kimlik doğrulama sonucunda 24 saatlik JWT token döner.
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı. 'email' ve 'password' alanları zorunludur."})
		return
	}

	user, ok := demoUsers[req.Email]
	if !ok || user.Password != req.Password {
		// Güvenlik: hangi alanın hatalı olduğunu ifşa etme
		c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı."})
		return
	}

	claims := models.DemoClaims{
		UserID: user.UserID,
		Email:  req.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "golden-state-demo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulurken bir hata meydana geldi."})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		Token: signedToken,
		Role:  user.Role,
	})
}
