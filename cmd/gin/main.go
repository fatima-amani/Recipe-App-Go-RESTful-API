package main

import (
	"net/http"
	"recipe-crud/pkg/recipes"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func main() {
	router := gin.Default()

	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	router.GET("/", homePage)
	router.GET("/recipes", recipesHandler.GetRecipes)
	router.GET("/recipes/:id", recipesHandler.GetRecipe)
	router.POST("recipes", recipesHandler.CreateRecipe)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipe)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipe)

	router.Run()
}

func homePage(c *gin.Context) {
	c.String(http.StatusOK, "This is Home Route")
}

type RecipesHandler struct {
	store recipeStore
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

func (h RecipesHandler) CreateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := slug.Make(recipe.Name)

	h.store.Add(id, recipe)

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) GetRecipes(c *gin.Context) {
	r, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, r)
}
func (h RecipesHandler) GetRecipe(c *gin.Context) {
	id := c.Param("id")

	recipe, err := h.store.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipe)
}
func (h RecipesHandler) UpdateRecipe(c *gin.Context) {
	var recipe recipes.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	err := h.store.Update(id, recipe)
	if err != nil {
		if err == recipes.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
func (h RecipesHandler) DeleteRecipe(c *gin.Context) {
	id := c.Param("id")

	err := h.store.Remove(id)
	if err != nil {
		if err == recipes.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted successfully"})
}

func NewRecipesHandler(s recipeStore) *RecipesHandler {
	return &RecipesHandler{
		store: s,
	}
}
