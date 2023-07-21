package main

import (
	"os"
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"io"
	"strings"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type MovieDBResponseArray struct {
	Results []MovieDBResponse `json:"results"`
}

type MovieDBResponse struct {
	Title string `json:"title"`
	ID int `json:"id"`
}

type MovieID struct {
	ID int `json:"id"`
}

type MovieDetails struct {
	Title string `json:"title"`
	Tagline string `json:"tagline"`
	Overview string `json:"overview"`
	Genres []Genre `json:"genres"`
	Revenue int `json:"revenue"`
	Poster string `json:"poster_path"`
	Actors []Actor `json:"actors"` 
}

type Genre struct {
	Genre string `json:"name"`
}

type Actor struct {
	Name string `json:"name"`
}

type Actors struct {
	Actors []Actor `json:"cast"`
}

type Input struct {
	Title string `json:"title"`
	ID int `json:"id"`
}

var key string

func getKey() string {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
		return ""
	} else {
		key = os.Getenv("movieDBKey")
		return key
	}
}

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

func getMovieByName(title string) MovieDBResponseArray {
	req := "https://api.themoviedb.org/3/search/movie?api_key=" + key + "&query=" + strings.ReplaceAll(strings.TrimSpace(title), " ", "+")
	response, err := http.Get(req)
	if err != nil {
		fmt.Println("ERROR")
	} else {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("ERROR")
		} else {
			var entry MovieDBResponseArray
			json.Unmarshal(body, &entry)
			return entry
		}
	}
	var entry MovieDBResponseArray
	return entry
}

func getMovieWithDetail(id int) MovieDetails {
	movieDetailReq := "https://api.themoviedb.org/3/movie/" + fmt.Sprint(id) + "?api_key=" + key
	response, err := http.Get(movieDetailReq)
	if err != nil {
		fmt.Println("ERROR")
	} else {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("ERROR")
		} else {
			var entry MovieDetails
			json.Unmarshal(body, &entry)

			movieActorReq := "https://api.themoviedb.org/3/movie/" + fmt.Sprint(id) + "/credits?api_key=" + key
			response, err := http.Get(movieActorReq)
			if err != nil {
				fmt.Println("ERROR")
			} else {
				body, err := io.ReadAll(response.Body)
				if err != nil {
					fmt.Println("ERROR")
				} else {
					var actors Actors
					json.Unmarshal(body, &actors)
					var arr []Actor
					if len(actors.Actors) < 10 {
						arr = make([]Actor, len(actors.Actors))
					} else {
						arr = make([]Actor, 10)
					}
					for i := 0; i < len(actors.Actors) && i < len(arr); i++ {
						arr[i].Name = actors.Actors[i].Name
					}
					entry.Actors = arr
				}
			}
			return entry
		}
	}
	var entry MovieDetails
	return entry
}

func main() {	
	key = getKey()

	router := gin.Default()
	router.GET("/getMovie", movieByName)
	router.Run("0.0.0.0:5050")
}