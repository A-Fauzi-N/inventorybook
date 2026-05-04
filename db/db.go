package db

import (
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

	endpointID := parts[0]

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

	db, err := gorm.Open("pgx", conn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	// Buat table manual pakai raw SQL agar tidak pakai AUTO_INCREMENT
	db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			description TEXT NOT NULL,
			stock INTEGER NOT NULL
		)
	`)

	var count int
	db.Raw("SELECT COUNT(*) FROM books").Scan(&count)
	if count == 0 {
		seederBook(db)
	}
}

func seederBook(db *gorm.DB) {
	db.Exec(`INSERT INTO books (title, author, description, stock) VALUES
		('Jojo Bizzare Adventure part 1', 'Hirohiko Araki', 'JOJO!!!!!', 5),
		('Jojo Bizzare Adventure part 2', 'Hirohiko Araki', 'JOJO!!!!!', 5),
		('Jojo Bizzare Adventure part 3', 'Hirohiko Araki', 'JOJO!!!!!', 5)
	`)
}