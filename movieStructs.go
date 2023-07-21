package main

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