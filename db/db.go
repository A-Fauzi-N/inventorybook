package db

import (
	"inventorybook/models"
	"log"
	"os"

	_ "database/sql"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitDB() *gorm.DB {
	_ = godotenv.Load(".env")

	conn := os.Getenv("POSTGRES_URL")
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Books{})
		
		data := models.Books{}
		if db.Find(&data).RecordNotFound(){
			seederBook(db)
		}
}

func seederBook(db *gorm.DB) {
	data := []models.Books{{
	Title		:	"Jojo Bizzare Adventure part 1",
	Author		:	"Hirohiko Araki",
	Description :	"JOJO!!!!!",
	Stock		:	5,
	}, {
	Title		:	"Jojo Bizzare Adventure part 2",
	Author		:	"Hirohiko Araki",
	Description :	"JOJO!!!!!",
	Stock		:	5,
	}, {
	Title		:	"Jojo Bizzare Adventure part 3",
	Author		:	"Hirohiko Araki",
	Description :	"JOJO!!!!!",
	Stock		:	5,
	}}

	for _, v := range data {
		db.Create(&v)
	}
}
