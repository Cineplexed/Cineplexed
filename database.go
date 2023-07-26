package main

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"fmt"
	"os"
	"gorm.io/driver/postgres"
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
			fmt.Println("Connected...")
			return db
		}
	}
	return nil
}