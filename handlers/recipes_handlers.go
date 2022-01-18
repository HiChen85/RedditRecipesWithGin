package handlers

import (
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/HiChen85/RedditRecipesWithGin/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"os"
	"time"
)

var recipes []*models.Recipe

func init() {
	//recipes = make([]*models.Recipe, 0)
	path, _ := os.Getwd()
	recipes = utils.LoadRecipesJson(path + "/recipes.json")
}

func NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishAt = time.Now()
	recipes = append(recipes, &recipe)
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	for i := range recipes {
		if recipes[i].ID == id {
			recipe.ID = recipes[i].ID
			recipe.PublishAt = recipes[i].PublishAt
			recipes[i] = &recipe
			c.JSON(http.StatusOK, gin.H{
				"message": "OK",
				"updated": recipe,
			})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Recipe not found",
	})
}
