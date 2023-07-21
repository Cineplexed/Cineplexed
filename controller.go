package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func movieByName(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		fmt.Println("ERROR")
	} else {
		var entry Input
		json.Unmarshal(body, &entry)
		if entry.ID != 0 {
			c.IndentedJSON(http.StatusOK, getMovieWithDetail(entry.ID))
		} else if entry.Title != "" {
			c.IndentedJSON(http.StatusOK, getMovieByName(entry.Title))
		}
	}
}

func main() {	
	key = getKey()

	router := gin.Default()
	router.GET("/getMovie", movieByName)
	router.Run("0.0.0.0:5050")
}