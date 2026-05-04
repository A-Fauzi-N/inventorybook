package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	apps "inventorybook/app"
	"inventorybook/auth"
	"inventorybook/db"
	"inventorybook/middleware"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func findTemplates() string {
	candidates := []string{
		"/vercel/path0/templates/*",
		"templates/*",
		"../templates/*",
		"/var/task/templates/*",
	}

	for _, path := range candidates {
		matches, err := filepath.Glob(path)
		fmt.Printf("[DEBUG] trying path: %s -> matches: %v, err: %v\n", path, matches, err)
		if err == nil && len(matches) > 0 {
			return path
		}
	}

	// Log isi direktori untuk debug
	wd, _ := os.Getwd()
	fmt.Printf("[DEBUG] working dir: %s\n", wd)
	entries, _ := os.ReadDir(wd)
	for _, e := range entries {
		fmt.Printf("[DEBUG] found in wd: %s\n", e.Name())
	}

	return "templates/*" // fallback
}

func init() {
	conn := db.InitDB()
	db.Migrate(conn)

	app = gin.New()
	app.Use(gin.Logger(), gin.Recovery())

	templatesPath := findTemplates()
	fmt.Printf("[DEBUG] using templates path: %s\n", templatesPath)
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