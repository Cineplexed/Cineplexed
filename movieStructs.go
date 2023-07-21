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
	ReleaseYear string `json:"release_date"`
	Director string `json:"director"`
	Producer string `json:"producer"`
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