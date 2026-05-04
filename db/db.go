package db

import (
	"inventorybook/models"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	_ "database/sql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func addNeonEndpoint(connStr string) string {
	u, err := url.Parse(connStr)
	if err != nil {
		return connStr
	}

	host := u.Hostname()
	parts := strings.SplitN(host, ".", 2)
	if len(parts) == 0 || !strings.HasPrefix(parts[0], "ep-") {
		return connStr
	}

	endpointID := strings.TrimSuffix(parts[0], "-pooler")

	q := u.Query()
	if q.Get("options") == "" {
		q.Set("options", fmt.Sprintf("endpoint=%s", endpointID))
		u.RawQuery = q.Encode()
	}

	return u.String()
}

func InitDB() *gorm.DB {
	_ = godotenv.Load(".env")

	conn := os.Getenv("POSTGRES_URL")
	conn = addNeonEndpoint(conn)

	// Ganti "postgres" -> "pgx"
	db, err := gorm.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Books{})

	data := models.Books{}
	if db.Find(&data).RecordNotFound() {
		seederBook(db)
	}
}

func seederBook(db *gorm.DB) {
	data := []models.Books{{
		Title:       "Jojo Bizzare Adventure part 1",
		Author:      "Hirohiko Araki",
		Description: "JOJO!!!!!",
		Stock:       5,
	}, {
		Title:       "Jojo Bizzare Adventure part 2",
		Author:      "Hirohiko Araki",
		Description: "JOJO!!!!!",
		Stock:       5,
	}, {
		Title:       "Jojo Bizzare Adventure part 3",
		Author:      "Hirohiko Araki",
		Description: "JOJO!!!!!",
		Stock:       5,
	}}

	for _, v := range data {
		db.Create(&v)
	}
}