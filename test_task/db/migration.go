package db

import (
	"gorm.io/gorm"
	"log"
	Song "test_task/models"
)

func Migrate(db *gorm.DB) {

	err := db.AutoMigrate(&Song.Song{})
	if err != nil {
		log.Fatal("Failed to migrate database")
	}
}
