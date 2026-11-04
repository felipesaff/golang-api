package config

import (
	"crud-api/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func LoadDB() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err.Error())
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")

	dbConn, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME))
	if err != nil {
		log.Fatal("Error connecting to the database: ", err.Error())
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal("Error pinging the database: ", err.Error())
	}

	_, err = dbConn.Exec(models.CREATE_TABLES)
	if err != nil {
		log.Fatal("Error creating tables: ", err.Error())
	}

	log.Println("Successfully connected to the database")
	return dbConn
}
