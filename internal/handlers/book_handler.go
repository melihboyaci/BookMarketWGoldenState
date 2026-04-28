package handlers

import (
	"net/http"
	"strconv"

	"github.com/bookmarket/golden-state/internal/config"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookRepo repository.BookRepository
}

func NewBookHandler(bookRepo repository.BookRepository) *BookHandler {
	return &BookHandler{bookRepo: bookRepo}
}

// GetBooks, demo_active verilerini döner.
func (h *BookHandler) GetBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		books, err := h.bookRepo.GetAll(config.TenantActive)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitaplar getirilirken hata oluştu"})
			return
		}
		if books == nil {
			books = []models.Book{}
		}
		c.JSON(http.StatusOK, books)
	}
}

func (h *BookHandler) CreateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		var book models.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kitap verisi"})
			return
		}

		book.TenantID = config.TenantActive

		if err := h.bookRepo.Create(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap oluşturulamadı"})
			return
		}

		c.JSON(http.StatusCreated, book)
	}
}

func (h *BookHandler) UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
			return
		}

		var book models.Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kitap verisi"})
			return
		}

		book.ID = id
		book.TenantID = config.TenantActive

		if err := h.bookRepo.Update(&book); err != nil {
			if err.Error() == "book not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap güncellenemedi"})
			return
		}

		c.JSON(http.StatusOK, book)
	}
}

func (h *BookHandler) DeleteBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
			return
		}

		if err := h.bookRepo.Delete(id, config.TenantActive); err != nil {
			if err.Error() == "book not found" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap silinemedi"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Kitap başarıyla silindi"})
	}
}
