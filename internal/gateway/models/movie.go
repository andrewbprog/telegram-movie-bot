package models

type Movie struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Genre       string   `json:"genre"`
	PosterURL   string   `json:"posterURL"`
	Year        int      `json:"year"`
	Rating      float64  `json:"rating"`
	Country     []string `json:"countries"`
	TrailerURL  string   `json:"trailerURL"`
}
