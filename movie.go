package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strings"
	"github.com/joho/godotenv"
	"os"
)

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