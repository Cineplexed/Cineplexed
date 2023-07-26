package main

import pq "github.com/lib/pq"

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
	ReleaseYear string `json:"release_date"`
	Director string `json:"director"`
	Producer string `json:"producer"`
}

type selections struct {
	Date string `gorm:"column:date"`
	Movie string `gorm:"column:movie"`
	NumCorrect int `gorm:"column:num_correct"`
	NumIncorrect int `gorm:"column:num_incorrect"`
	Tagline string `gorm:"column:tagline"`
	Overview string `gorm:"column:overview"`
	Genres pq.StringArray `gorm:"type:text[]; column:genres"`
	Actors pq.StringArray `gorm:"type:text[]; column:actors"` 
	Revenue int `gorm:"column:revenue"`
	Poster string `gorm:"column:poster"`
	ReleaseYear string `gorm:"column:year"`
	Director string `gorm:"column:director"`
	Producer string `gorm:"column:producer"`
}

type Genre struct {
	GenreVal string `json:"name"`
}

type Actor struct {
	Name string `json:"name"`
}

type Actors struct {
	Actors []Actor `json:"cast"`
}

type CrewMember struct {
	Name string `json:"name"`
	Job string `json:"job"`
}

type Crew struct {
	EntireCrew []CrewMember `json:"crew"`
}

type Producer struct {
	Name string `json:"name"`
}

type Producers struct {
	Companies []Producer `json:"production_companies"`
}

type Input struct {
	Title string `json:"title"`
	ID int `json:"id"`
}

//lint:ignore U1000 Used for Swagger
type docs_ID struct {
	ID int `json:"id"`
}

//lint:ignore U1000 Used for Swagger
type docs_Title struct {
	Title string `json:"title"`
}