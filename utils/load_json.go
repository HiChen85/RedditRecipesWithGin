package utils

import (
	"encoding/json"
	"github.com/HiChen85/RedditRecipesWithGin/recipes_service/models"
	"io/ioutil"
	"log"
)

func LoadRecipesJson(fileName string) []*models.Recipe {
	recipes := make([]*models.Recipe, 0)
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &recipes)
	if err != nil {
		log.Fatal(err)
	}
	return recipes
}
