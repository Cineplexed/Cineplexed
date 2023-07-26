package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strings"
	"github.com/joho/godotenv"
	"os"
	"math/rand"
	"time"
)

var key string
var baseUrl string
var searchUrl string
var randUrl string

var tomorrow time.Time
var nextTime time.Time
var updatingDaily = false

func getEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("ERROR")
	} else {
		key = os.Getenv("movieDBKey")
		baseUrl = os.Getenv("movieDBRootUrl")
		searchUrl = os.Getenv("movieDBSearchUrl")
		randUrl = os.Getenv("movieDBRandUrl")
	}
}

func checkTime() {
	if time.Now().After(nextTime) && !updatingDaily {
		updatingDaily = true
		getDailyMovie()
	}
}

func getTargetTime() {
	var entry selections
	db.Last(&entry)
	tomorrow, _ = time.Parse("2006-01-02", strings.ReplaceAll(entry.Date, "/", "-"))
	nextTime = time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day() + 1, 0, 0, 0, 0, time.Now().Location())
}

func getMovieByName(title string) MovieDBResponseArray {
	checkTime()
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
	checkTime()
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

			if len(producers.Companies) > 0 {
				entry.Producer = producers.Companies[0].Name
			} 
			if len(entry.ReleaseYear) >= 4 { 
				entry.ReleaseYear = entry.ReleaseYear[0:4]
			}

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
						if crew.EntireCrew[i].Job == "Director" {
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

func getDailyMovie() {
	page := rand.Intn(25) + 1
	req := randUrl + "?api_key=" + key + "&page=" + fmt.Sprint(page) 
	response, err := http.Get(req)
	if err != nil {
		fmt.Println("ERROR")
	} else {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("ERROR")
		} else {
			var collection MovieDBResponseArray
			json.Unmarshal(body, &collection)
			index := rand.Intn(20)
			entry := collection.Results[index]
			detailedEntry := getMovieWithDetail(entry.ID)

			var arrGenres []string = make([]string, len(detailedEntry.Genres))
			for i := 0; i < len(detailedEntry.Genres); i++ {
				arrGenres[i] = string(detailedEntry.Genres[i].GenreVal)
			}

			var arrActors []string = make([]string, len(detailedEntry.Actors))
			for i := 0; i < len(detailedEntry.Actors); i++ {
				arrActors[i] = string(detailedEntry.Actors[i].Name)
			}

			var complete selections = selections{
				Date: nextTime.Format("2006") + "/" + nextTime.Format("01") + "/" + nextTime.Format("02"), 
				Movie: detailedEntry.Title, 
				NumCorrect: 0,
				NumIncorrect: 0,
				Tagline: detailedEntry.Tagline,
				Overview: detailedEntry.Overview,
				Genres: arrGenres,
				Actors: arrActors,
				Revenue: detailedEntry.Revenue,
				Poster: detailedEntry.Poster,
				ReleaseYear: detailedEntry.ReleaseYear,
				Director: detailedEntry.Director,
				Producer: detailedEntry.Producer}
			db.Table("selections")
			result := db.Create(&complete)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
			} else {
				fmt.Println("Daily movie updated to " + complete.Movie)
			}
		}
	}
	getTargetTime()
	updatingDaily = false
}