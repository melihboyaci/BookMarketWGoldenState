package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/bookmarket/golden-state/internal/db"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

// RestoreGoldenState, POST /api/v1/system/restore endpoint'ini karşılar.
// RequireAdminRole middleware'i tarafından korunur; yalnızca ADMIN rolü erişebilir.
// Başarılı restore sonrası kaç satır silindiğini, kaç satır eklendiğini
// ve işlemin kaç millisaniyede tamamlandığını döner.
func RestoreGoldenState(c *gin.Context) {
	result, err := repository.RestoreGoldenState(db.DB)
	if err != nil {
		// Güvenlik: iç hata detayını (tablo adı, SQL mesajı vb.) logla ama dışarıya sızdırma
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
