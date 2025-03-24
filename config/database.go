package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	log.Println("Database connection trying to connect")
	var err error

	dbUser := os.Getenv("DATABASE_USERNAME")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbName := os.Getenv("DATABASE_NAME")

	dbDns := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	fmt.Println("DB DNS:", dbDns, os.Getenv("APP_ENV"))

	DB, err = gorm.Open(postgres.Open(dbDns), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database. Error: ", err.Error())
	}

	log.Println("Database connection established successfully.")

	return DB
}
