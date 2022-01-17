package models

import "time"

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" bson:"name"`
	Tags         []string  `json:"tags" bson:"tags"`
	Ingredients  []string  `json:"ingredients" bson:"ingredients"`
	Instructions []string  `json:"instructions" bson:"instructions"`
	PublishAt    time.Time `json:"publishAt"`
}
