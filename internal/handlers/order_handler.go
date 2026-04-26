package handlers

import (
	"net/http"

	"github.com/bookmarket/golden-state/internal/db"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

func Checkout(c *gin.Context) {
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sipariş verisi"})
		return
	}

	if err := repository.Checkout(db.DB, "demo_active", &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sipariş başarıyla tamamlandı"})
}

func GetSales(c *gin.Context) {
	stats, err := repository.GetSalesStats(db.DB, "demo_active")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Satış istatistikleri alınamadı"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
