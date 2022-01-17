package handlers

import (
	"github.com/HiChen85/RedditRecipesWithGin/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"time"
)

var recipes []*models.Recipe

func init() {
	recipes = make([]*models.Recipe, 0)
	
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
