package handler

import (
	"embed"
	"html/template"
	"net/http"
	apps "inventorybook/app"
	"inventorybook/auth"
	"inventorybook/db"
	"inventorybook/middleware"

	"github.com/gin-gonic/gin"
)

var templatesFS embed.FS

var app *gin.Engine

func init() {
	conn := db.InitDB()
	db.Migrate(conn)

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	// Load templates dari embed
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/*.html"))
	app.SetHTMLTemplate(tmpl)

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