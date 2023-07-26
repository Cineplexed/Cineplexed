package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"os"
	"github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
	"github.com/gin-contrib/cors"
	_ "cineplexed.com/docs"
	"strconv"
)

// moviesByName godoc
// @Summary moviesByName
// @Description Get a list of possible movies by it's title
// @Tags movie
// @Param Item body docs_Title true "Title"
// @Accept */*
// @Produce json
// @Router /getMovieOptions [GET]
func moviesByName(c *gin.Context) {
	title := c.Query("title")
	if title != "" {
		c.IndentedJSON(http.StatusOK, getMovieByName(title))
	}
}

// movieWithDetails godoc
// @Summary movieWithDetails
// @Description Get a movie with extensive details using it's ID
// @Tags movie
// @Param Item body docs_ID true "ID"
// @Accept */*
// @Produce json
// @Router /getMovieDetails [GET]
func movieWithDetails(c *gin.Context) {
	id := c.Query("id")
	numId, err := strconv.Atoi(id)
	if err != nil {	
		fmt.Println("ERROR")
	} else {
		if numId != 0 {
			c.IndentedJSON(http.StatusOK, getMovieWithDetail(numId))
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

// @title Gin Swagger Cineplexed
// @version 1.0
// @description This is Cineplexed.
// @host localhost:5050
// @BasePath /
// @schemes http
func main() {
	getEnv()
	getTargetTime()
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}

	router.Use(cors.New(config))

	router.GET("/getMovieOptions", moviesByName)
	router.GET("/getMovieDetails", movieWithDetails)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run(getHost())
}