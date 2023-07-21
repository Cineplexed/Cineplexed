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
var baseUrl string
var searchUrl string

func getEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
	} else {
		key = os.Getenv("movieDBKey")
		baseUrl = os.Getenv("movieDBRootUrl")
		searchUrl = os.Getenv("movieDBSearchUrl")
	}
}

func getMovieByName(title string) MovieDBResponseArray {
	req := searchUrl + "?api_key=" + key + "&query=" + strings.ReplaceAll(strings.TrimSpace(title), " ", "+")
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

			var limitedEntry MovieDBResponseArray
			if (len(entry.Results) < 20) {
				limitedEntry.Results = make([]MovieDBResponse, len(entry.Results))
			} else {
				limitedEntry.Results = make([]MovieDBResponse, 20)
			}
			for i := 0; i < len(limitedEntry.Results); i++ {
				limitedEntry.Results[i] = entry.Results[i]
			}
			return limitedEntry
		}
	}
	var entry MovieDBResponseArray
	return entry
}

func getMovieWithDetail(id int) MovieDetails {
	movieDetailReq := baseUrl + "/" + fmt.Sprint(id) + "?api_key=" + key
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
			
			var producers Producers
			json.Unmarshal(body, &producers)

			entry.Producer = producers.Companies[0].Name
			entry.ReleaseYear = entry.ReleaseYear[0:4]

			movieActorReq := baseUrl + "/" + fmt.Sprint(id) + "/credits?api_key=" + key
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

					var crew Crew
					json.Unmarshal(body, &crew)

					for i := 0; i < len(crew.EntireCrew); i++ {
						fmt.Println(crew.EntireCrew[i].Name + " " + crew.EntireCrew[i].Job)
						if crew.EntireCrew[i].Job == "Producer" || crew.EntireCrew[i].Job == "Executive Producer" {
							entry.Director = crew.EntireCrew[i].Name
							break
						}
					}

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