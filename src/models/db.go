package models

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {
	sqlite_file := os.Getenv("SQLITE_PATH")
	db, err := gorm.Open(sqlite.Open(sqlite_file), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Record{})

	DB = db
}
