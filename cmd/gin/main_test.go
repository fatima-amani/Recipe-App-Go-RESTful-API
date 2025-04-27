package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"recipe-crud/pkg/recipes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Fatalf("Could not read test data file %v: %v", name, err)
	}
	return content
}

func setupTestRouter() (*gin.Engine, recipeStore) {
	store := recipes.NewMemStore()
	handler := NewRecipesHandler(store)

	router := gin.Default()
	router.GET("/", homePage)
	router.GET("/recipes", handler.GetRecipes)
	router.GET("/recipes/:id", handler.GetRecipe)
	router.POST("/recipes", handler.CreateRecipe)
	router.PUT("/recipes/:id", handler.UpdateRecipe)
	router.DELETE("/recipes/:id", handler.DeleteRecipe)

	return router, store
}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {
	router, store := setupTestRouter()

	chickenBiryani := readTestData(t, "chicken_biryani.json")
	chickenBiryaniWithRaita := readTestData(t, "chicken_biryani_with_raita.json")

	// CREATE
	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(chickenBiryani))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)

	// GET ONE
	req = httptest.NewRequest(http.MethodGet, "/recipes/chicken-biryani", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	assert.NoError(t, err)

	assert.JSONEq(t, string(chickenBiryani), string(data))

	// UPDATE
	req = httptest.NewRequest(http.MethodPut, "/recipes/chicken-biryani", bytes.NewReader(chickenBiryaniWithRaita))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	updatedChickenBiryani, err := store.Get("chicken-biryani")
	assert.NoError(t, err)
	assert.Contains(t, updatedChickenBiryani.Ingredients, recipes.Ingredient{Name: "dahi"})

	// DELETE
	req = httptest.NewRequest(http.MethodDelete, "/recipes/chicken-biryani", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	saved, _ = store.List()
	assert.Len(t, saved, 0)
}
