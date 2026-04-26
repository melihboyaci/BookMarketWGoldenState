package handlers

import (
	"net/http"
	"strconv"

	"github.com/bookmarket/golden-state/internal/db"
	"github.com/bookmarket/golden-state/internal/models"
	"github.com/bookmarket/golden-state/internal/repository"
	"github.com/gin-gonic/gin"
)

// GetBooks, demo_active verilerini döner.
func GetBooks(c *gin.Context) {
	books, err := repository.GetAllBooks(db.DB, "demo_active")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitaplar getirilirken hata oluştu"})
		return
	}
	// Return empty array instead of null
	if books == nil {
		books = []models.Book{}
	}
	c.JSON(http.StatusOK, books)
}

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz kitap verisi"})
		return
	}
	
	book.TenantID = "demo_active"
	
	if err := repository.CreateBook(db.DB, &book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap oluşturulamadı"})
		return
	}
	
	c.JSON(http.StatusCreated, book)
}

func UpdateBook(c *gin.Context) {
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
	book.TenantID = "demo_active"

	if err := repository.UpdateBook(db.DB, &book); err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap güncellenemedi"})
		return
	}

	c.JSON(http.StatusOK, book)
}

func DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"})
		return
	}

	if err := repository.DeleteBook(db.DB, id, "demo_active"); err != nil {
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kitap bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kitap silinemedi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitap başarıyla silindi"})
}
