package handler

import (
	"net/http"
	apps "inventorybook/app"
	"inventorybook/auth"
	"inventorybook/db"
	"inventorybook/middleware"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func init() {
	// Inisialisasi database
	conn := db.InitDB()
	db.Migrate(conn)

	// Inisialisasi router
	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	// Load template HTML
	app.LoadHTMLGlob("./templates/*")
	// Handler app
	h := apps.New(conn)

	// AUTH ROUTES
	app.GET("/", auth.HomeHandler)
	app.GET("/login", auth.LoginGetHandler)
	app.POST("/login", auth.LoginPostHandler)

	// BOOK ROUTES (Protected)
	app.GET("/books", middleware.AuthValid, h.GetBooks)
	app.GET("/book/:id", middleware.AuthValid, h.GetBookById)
	app.GET("/addBook", middleware.AuthValid, h.AddBook)
	app.POST("/book", middleware.AuthValid, h.PostBook)
	app.GET("/updateBook/:id", middleware.AuthValid, h.UpdateBook)
	app.POST("/updateBook/:id", middleware.AuthValid, h.PutBook)
	app.POST("/deleteBook/:id", middleware.AuthValid, h.DeleteBook)
}

// Handler is the entry point for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
