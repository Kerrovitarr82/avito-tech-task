package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func ConnectToDB() *gorm.DB {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	for count := 0; count < 5; count++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db // успешное подключение
		}
		log.Println("Failed to connect to database. Retrying...")
		time.Sleep(5 * time.Second)
	}
	log.Fatal("Could not connect to the database after 5 attempts:", err)
	return nil
}
