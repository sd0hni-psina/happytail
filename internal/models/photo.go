package models

import "time"

type AnimalPhoto struct {
	ID        int       `json:"id"`
	AnimalID  int       `json:"animal_id"`
	URL       string    `json:"url"`
	IsMain    bool      `json:"is_main"`
	CreatedAt time.Time `json:"created_at"`
}

type AnimalPhotoInput struct {
	AnimalID int    `json:"animal_id"`
	URL      string `json:"url"`
	IsMain   bool   `json:"is_main"`
}
