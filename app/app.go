package app

import (
	"fmt"
	"inventorybook/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Handler {
	return Handler{DB: db}
}

// GET /books
func (h *Handler) GetBooks(c *gin.Context) {
	var books []models.Books

	if err := h.DB.Find(&books).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Gagal mengambil data buku",
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":   "Home Page",
		"payload": books,
		"auth":    c.Query("auth"),
	})
}

// GET /book/:id
func (h *Handler) GetBookById(c *gin.Context) {
	bookId := c.Param("id")
	var book models.Books

	if h.DB.First(&book, "id = ?", bookId).RecordNotFound() {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Book not found",
		})
		return
	}

	c.HTML(http.StatusOK, "book.html", gin.H{
		"title":   book.Title,
		"payload": book,
		"auth":    c.Query("auth"),
	})
}

// GET /addBook
func (h *Handler) AddBook(c *gin.Context) {
	c.HTML(http.StatusOK, "formbook.html", gin.H{
		"title": "Add Book",
		"auth":  c.Query("auth"),
	})
}

// POST /book
func (h *Handler) PostBook(c *gin.Context) {
	var book models.Books

	if err := c.ShouldBind(&book); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Input book tidak valid",
		})
		return
	}

	if err := h.DB.Create(&book).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Gagal menambahkan buku",
		})
		return
	}

	c.Redirect(http.StatusFound,
		fmt.Sprintf("/books?auth=%s", c.PostForm("auth")))
}

// GET /updateBook/:id
func (h *Handler) UpdateBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.First(&book, "id = ?", bookId).RecordNotFound() {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Book not found",
		})
		return
	}

	c.HTML(http.StatusOK, "formbook.html", gin.H{
		"title":   "Update Book",
		"payload": book,
		"auth":    c.Query("auth"),
	})
}

// POST /updateBook/:id
func (h *Handler) PutBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	// Cek apakah buku ada
	if h.DB.First(&book, "id = ?", bookId).RecordNotFound() {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Book not found",
		})
		return
	}

	// Ambil data baru dari form
	var reqBook models.Books
	if err := c.ShouldBind(&reqBook); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error": "Input update tidak valid",
		})
		return
	}

	// Update field
	if err := h.DB.Model(&book).Updates(reqBook).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Gagal update buku",
		})
		return
	}

	c.Redirect(http.StatusFound,
		fmt.Sprintf("/book/%s?auth=%s", bookId, c.PostForm("auth")))
}

// POST /deleteBook/:id
func (h *Handler) DeleteBook(c *gin.Context) {
	var book models.Books
	bookId := c.Param("id")

	if h.DB.First(&book, "id = ?", bookId).RecordNotFound() {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Book not found",
		})
		return
	}

	if err := h.DB.Delete(&book).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"error": "Gagal menghapus buku",
		})
		return
	}

	c.Redirect(http.StatusFound,
		fmt.Sprintf("/books?auth=%s", c.PostForm("auth")))
}