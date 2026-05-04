package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	conn := db.InitDB()
	db.Migrate(conn)

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	_, filename, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(filename))
	templatesPath := filepath.Join(basePath, "templates", "*")

	if _, err := os.Stat(filepath.Join(basePath, "templates")); os.IsNotExist(err) {
		templatesPath = "templates/*"
	}

	app.LoadHTMLGlob(templatesPath)

	h := apps.New(conn)

	app.GET("/", auth.HomeHandler)
	app.GET("/login", auth.LoginGetHandler)
	app.POST("/login", auth.LoginPostHandler)

	app.GET("/books", middleware.AuthValid, h.GetBooks)
	app.GET("/book/:id", middleware.AuthValid, h.GetBookById)
	app.GET("/addBook", middleware.AuthValid, h.AddBook)
	app.POST("/book", middleware.AuthValid, h.PostBook)
	app.GET("/updateBook/:id", middleware.AuthValid, h.UpdateBook)
	app.POST("/updateBook/:id", middleware.AuthValid, h.PutBook)
	app.POST("/deleteBook/:id", middleware.AuthValid, h.DeleteBook)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}