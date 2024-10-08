package main

import (
	"log"
	"net/http"
	"os"
	"test_task/controllers"
	Song "test_task/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Подключение к базе данных
	dsn := os.Getenv("DB_DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	// Миграции
	db.AutoMigrate(&Song.Song{})

	// Настройка роутера
	router := mux.NewRouter()

	// Маршруты
	router.HandleFunc("/songs", controllers.GetSongs).Methods("GET")
	router.HandleFunc("/songs", controllers.CreateSong).Methods("POST")
	router.HandleFunc("/songs/{id}", controllers.GetSongLyrics).Methods("GET")
	router.HandleFunc("/songs/{id}", controllers.UpdateSong).Methods("PUT")
	router.HandleFunc("/songs/{id}", controllers.DeleteSong).Methods("DELETE")

	// Запуск сервера
	log.Fatal(http.ListenAndServe(":8080", router))
}
