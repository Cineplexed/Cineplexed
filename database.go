package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"fmt"
	"os"
	"gorm.io/driver/postgres"
	"github.com/google/uuid"
	"time"
)

var db = connectDB()

func connectDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
	} else {
		dsn := os.Getenv("conString")
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(err.Error())
		} else {
			db.AutoMigrate(&selections{})
			fmt.Println("Connected to Database...")
			return db
		}
	}
	return nil
}

func log(severity string, content string) {
	var entry Log
	entry.ID = uuid.New().String()
	entry.Severity = severity
	entry.Content = content
	entry.Timestamp = time.Now().String()
	result := db.Create(&entry)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
	}
}