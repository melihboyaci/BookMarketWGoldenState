package handlers

import (
	"database/sql"
	"net/http"

	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

// GetUsers, GET /api/v1/admin/users endpoint'ini karşılar.
// Yalnızca ADMIN rolüne sahip kullanıcılar erişebilir (RequireAdminRole middleware'i ile korunur).
func GetUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := repository.GetAllUsers(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcılar alınırken hata oluştu."})
			return
		}

		// nil slice yerine boş dizi döndür
		if users == nil {
			users = []models.User{}
		}

		c.JSON(http.StatusOK, users)
	}
}
