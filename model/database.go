package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var db *gorm.DB

func Open() {
	if url := os.Getenv("DATABASE_URL"); url == "" {
		panic("DATABASE_URL must be set")
	} else if dbLocal, err := gorm.Open("sqlite3", url); err != nil {
		panic(err)
	} else {
		db = dbLocal
	}
}

func Close() {
	db.Close()
}

func Debug() {
	db = db.Debug()
}
