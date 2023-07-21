package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
)

func moviesByName(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		fmt.Println("ERROR")
	} else {
		var entry Input
		json.Unmarshal(body, &entry)
		if entry.Title != "" {
			c.IndentedJSON(http.StatusOK, getMovieByName(entry.Title))
		}
	}
}

func movieWithDetails(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		fmt.Println("ERROR")
	} else {
		var entry Input
		json.Unmarshal(body, &entry)
		if entry.ID != 0 {
			c.IndentedJSON(http.StatusOK, getMovieWithDetail(entry.ID))
		}
	}
}

func getHost() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
		return ""
	} else {
		return os.Getenv("host")
	}
}

func main() {	
	getEnv()

	router := gin.Default()
	router.GET("/getMovieOptions", moviesByName)
	router.GET("/getMovieDetails", movieWithDetails)
	router.Run(getHost())
}