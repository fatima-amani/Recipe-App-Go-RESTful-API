package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recipe-crud/pkg/recipes"

	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
)

func main() {

	store := recipes.NewMemStore()
	router := mux.NewRouter()
	router.Handle("/", &homeHandler{})
	router.Use(loggingMiddleware)

	s := router.PathPrefix("/recipes").Subrouter()

	NewRecipesHandler(store, s)

	http.ListenAndServe(":8080", router)
}

type RecipesHandler struct {
	store recipeStore
}

func NewRecipesHandler(s recipeStore, router *mux.Router) *RecipesHandler {

	handler := &RecipesHandler{
		store: s,
	}

	router.HandleFunc("/", handler.CreateRecipe).Methods("POST")
	router.HandleFunc("/", handler.GetRecipes).Methods("GET")
	router.HandleFunc("/{id}", handler.GetRecipe).Methods("GET")
	router.HandleFunc("/{id}", handler.UpdateRecipe).Methods("PUT")
	router.HandleFunc("/{id}", handler.DeleteRecipe).Methods("DELETE")

	return handler

}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Request received:", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type recipeStore interface {
	Add(name string, recipe recipes.Recipe) error
	Get(name string) (recipes.Recipe, error)
	List() (map[string]recipes.Recipe, error)
	Update(name string, recipe recipes.Recipe) error
	Remove(name string) error
}

func (h *RecipesHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	resourceID := slug.Make(recipe.Name)
	if err := h.store.Add(resourceID, recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Recipe created successfully"))
}
func (h *RecipesHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := h.store.List()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(recipes)

	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
func (h *RecipesHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	recipe, err := h.store.Get(id)
	if err != nil {
		if err == recipes.ErrNotFound {
			NotFoundErrorHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h *RecipesHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var recipe recipes.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(id, recipe); err != nil {
		if err == recipes.ErrNotFound {
			NotFoundErrorHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}
	jsonBytes, err := json.Marshal(recipe)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)

}
func (h *RecipesHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.store.Remove(id); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Internal Server Error Occured"))
}

func NotFoundErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Page Not found !!"))
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is home Page"))
}
