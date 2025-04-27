package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"recipe-crud/pkg/recipes"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func readTestData(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile("../../testdata/" + name)
	if err != nil {
		t.Errorf("Could not read %v", name)
	}

	return content
}

func TestRecipesHandlerCRUD_Integration(t *testing.T) {
	store := recipes.NewMemStore()
	recipesHandler := NewRecipesHandler(store)

	router := mux.NewRouter()
	router.HandleFunc("/recipes", recipesHandler.CreateRecipe).Methods(http.MethodPost)
	router.HandleFunc("/recipes", recipesHandler.GetRecipes).Methods(http.MethodGet)
	router.HandleFunc("/recipes/{id}", recipesHandler.GetRecipe).Methods(http.MethodGet)
	router.HandleFunc("/recipes/{id}", recipesHandler.UpdateRecipe).Methods(http.MethodPut)
	router.HandleFunc("/recipes/{id}", recipesHandler.DeleteRecipe).Methods(http.MethodDelete)

	chickenBiryani := readTestData(t, "chicken_biryani.json")
	chickenBiryaniReader := bytes.NewReader(chickenBiryani)

	chickenBiryaniWithRaita := readTestData(t, "chicken_biryani_with_raita.json")
	chickenBiryaniWithRaitaReader := bytes.NewReader(chickenBiryaniWithRaita)

	// create
	req := httptest.NewRequest(http.MethodPost, "/recipes", chickenBiryaniReader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ := store.List()
	assert.Len(t, saved, 1)

	// Get one
	req = httptest.NewRequest(http.MethodGet, "/recipes/chicken-biryani", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	data, err := io.ReadAll(res.Body)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	assert.JSONEq(t, string(chickenBiryani), string(data))

	// Update
	req = httptest.NewRequest(http.MethodPut, "/recipes/chicken-biryani", chickenBiryaniWithRaitaReader)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	updatedChickenBiryani, err := store.Get("chicken-biryani")
	assert.NoError(t, err)

	assert.Contains(t, updatedChickenBiryani.Ingredients, recipes.Ingredient{Name: "dahi"})

	// Remove
	req = httptest.NewRequest(http.MethodDelete, "/recipes/chicken-biryani", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res = w.Result()
	defer res.Body.Close()
	assert.Equal(t, 200, res.StatusCode)

	saved, _ = store.List()
	assert.Len(t, saved, 0)
}
