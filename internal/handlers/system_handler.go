package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	restoreRepo repository.RestoreRepository
}

func NewSystemHandler(restoreRepo repository.RestoreRepository) *SystemHandler {
	return &SystemHandler{restoreRepo: restoreRepo}
}

// RestoreGoldenState, POST /api/v1/system/restore endpoint'ini karşılar.
func (h *SystemHandler) RestoreGoldenState() gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := h.restoreRepo.RestoreGoldenState()
		if err != nil {
			log.Printf("HATA: Golden State restore başarısız: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Sistem sıfırlama işlemi sırasında bir hata oluştu.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "Golden State başarıyla geri yüklendi.",
			"deleted_rows":  result.DeletedRows,
			"inserted_rows": result.InsertedRows,
			"duration_ms":   result.DurationMs,
			"restored_at":   time.Now().UTC().Format(time.RFC3339),
		})
	}
}
