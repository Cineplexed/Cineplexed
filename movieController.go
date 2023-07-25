package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "cineplexed.com/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, err)
		return
	}

	if id != 0 {
		c.IndentedJSON(http.StatusOK, getMovieWithDetail(id))
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
