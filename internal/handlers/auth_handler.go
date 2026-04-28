package handlers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthHandler(userRepo repository.UserRepository, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

// loginRequest, POST /auth/login için beklenen istek gövdesidir.
type loginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// loginResponse, başarılı giriş sonucunda dönen yanıt yapısıdır.
type loginResponse struct {
	Token string      `json:"token"`
	Role  models.Role `json:"role"`
}

// registerRequest, POST /auth/register için beklenen istek gövdesidir.
type registerRequest struct {
	Username string `json:"username" binding:"required,min=2,max=100"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Login, POST /api/v1/auth/login endpoint'ini karşılar.
func (h *AuthHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek formatı. 'email' ve 'password' alanları zorunludur."})
			return
		}

		// Kullanıcıyı yalnızca demo_active tenant'ından ara
		user, err := h.userRepo.FindByEmail(strings.ToLower(req.Email), config.TenantActive)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kimlik doğrulama sırasında bir hata oluştu."})
			return
		}

		// Bcrypt karşılaştırması
		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "E-posta veya şifre hatalı."})
			return
		}

		signedToken, err := h.generateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulurken bir hata meydana geldi."})
			return
		}

		c.JSON(http.StatusOK, loginResponse{
			Token: signedToken,
			Role:  user.Role,
		})
	}
}

// Register, POST /api/v1/auth/register endpoint'ini karşılar.
func (h *AuthHandler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz istek. username (min 2), email ve password (min 6) zorunludur."})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Şifre işlenirken bir hata oluştu."})
			return
		}

		newUser := &models.User{
			Username:     req.Username,
			Email:        strings.ToLower(req.Email),
			PasswordHash: string(hash),
			Role:         models.RoleBuyer,
			TenantID:     config.TenantActive,
		}

		if err := h.userRepo.Create(newUser); err != nil {
			if strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate") {
				c.JSON(http.StatusConflict, gin.H{"error": "Bu e-posta adresi zaten kayıtlı."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı oluşturulurken bir hata oluştu."})
			return
		}

		signedToken, err := h.generateJWT(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token oluşturulurken bir hata meydana geldi."})
			return
		}

		c.JSON(http.StatusCreated, loginResponse{
			Token: signedToken,
			Role:  newUser.Role,
		})
	}
}

// GetUsers, GET /api/v1/admin/users endpoint'ini karşılar.
func (h *AuthHandler) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := h.userRepo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcılar getirilirken bir hata oluştu."})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (h *AuthHandler) generateJWT(user *models.User) (string, error) {
	claims := models.DemoClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "golden-state-demo",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.cfg.JWTSecret))
}
