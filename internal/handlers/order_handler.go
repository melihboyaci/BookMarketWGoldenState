package handlers

import (
	"net/http"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo repository.OrderRepository
}

func NewOrderHandler(orderRepo repository.OrderRepository) *OrderHandler {
	return &OrderHandler{orderRepo: orderRepo}
}

func (h *OrderHandler) Checkout() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CheckoutRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz sipariş verisi"})
			return
		}

		if err := h.orderRepo.Checkout(config.TenantActive, &req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Sipariş başarıyla tamamlandı"})
	}
}

func (h *OrderHandler) GetSales() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.orderRepo.GetSalesStats(config.TenantActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Satış istatistikleri alınamadı"})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}
