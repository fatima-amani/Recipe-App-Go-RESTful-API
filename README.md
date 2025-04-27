# Recipe-App-Go-RESTful-API

A simple **CRUD** (Create, Read, Update, Delete) **Recipe API** built in **Go** with three different HTTP libraries:
* `net/http` (stdlib)
* `Gorilla Mux`
* `Gin`

Data is stored **in memory** (no external DB). This project is a great starter for building full REST APIs in Go!

## ✨ Features

* Full RESTful API for managing recipes
* Different implementations using `net/http`, `Gorilla Mux`, and `Gin`
* In-memory storage for recipes (no database needed)
* Organized project structure
* Sample JSON files for testing recipes

## 📂 Folder Structure

| Folder/File | Description |
|-------------|-------------|
| `/cmd/gin` | Gin framework based REST API (`main.go`, `main_test.go`) |
| `/cmd/gorilla` | Gorilla Mux based REST API (`main.go`, `main_test.go`) |
| `/cmd/stdlib` | Standard Library (`net/http`) based REST API (`main.go`, `main_test.go`) |
| `/pkg/recipes/models.go` | Defines the `Recipe` model (data structure) |
| `/pkg/recipes/recipeMemStore.go` | In-memory storage for recipes (Add, Get, List, Update, Remove methods) |
| `/testdata/chicken_biryani.json` | Sample recipe data for testing create endpoints |
| `/testdata/chicken_biryani_with_raita.json` | Sample updated recipe data for testing update endpoints |
| `go.mod`, `go.sum` | Go modules and dependencies |

## 🛠️ Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/fatima-amani/Recipe-App-Go-RESTful-API.git
cd Recipe-App-Go-RESTful-API
```

### 2. Install Go (if you haven't)

Make sure you have Go installed:

```bash
go version
```

If not, download from: https://golang.org/dl/

### 3. Run the Server

You can choose any of the implementations:

#### ➡️ Gorilla Mux Version

```bash
cd cmd/gorilla
go run main.go
```

#### ➡️ Gin Gonic Version

```bash
cd cmd/gin
go run main.go
```

#### ➡️ Stdlib Version

```bash
cd cmd/stdlib
go run main.go
```

## 🔥 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/recipes` | Create a new recipe |
| GET | `/recipes` | Get all recipes |
| GET | `/recipes/{id}` | Get a single recipe by ID |
| PUT | `/recipes/{id}` | Update a recipe |
| DELETE | `/recipes/{id}` | Delete a recipe |

## 📬 Sample API Requests

You can use `curl`, `Postman`, or any API testing tool.

### 🥘 Create a New Recipe (POST)

```bash
curl -X POST http://localhost:8080/recipes \
  -H "Content-Type: application/json" \
  -d @testdata/chicken_biryani.json
```

or manually:

```bash
curl -X POST http://localhost:8080/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "Name": "Chicken Biryani",
    "Ingredients": [
      {"Name": "Chicken"},
      {"Name": "Rice"},
      {"Name": "Spices"}
    ],
    "Steps": ["Marinate chicken", "Cook rice", "Layer and steam"]
  }'
```

### 📖 Get All Recipes (GET)

```bash
curl http://localhost:8080/recipes
```

### 📋 Get a Single Recipe (GET)

```bash
curl http://localhost:8080/recipes/chicken-biryani
```
(`chicken-biryani` is automatically generated from the name using `slug.Make()`)

### ✏️ Update a Recipe (PUT)

```bash
curl -X PUT http://localhost:8080/recipes/chicken-biryani \
  -H "Content-Type: application/json" \
  -d @testdata/chicken_biryani_with_raita.json
```

### ❌ Delete a Recipe (DELETE)

```bash
curl -X DELETE http://localhost:8080/recipes/chicken-biryani
```

## 🧪 Testing

Each implementation also comes with a basic **integration test** (`main_test.go`) under each folder (`cmd/gorilla/main_test.go`, etc).

To run tests:

```bash
go test ./...
```

It will:
* Create a recipe
* Retrieve it
* Update it
* Delete it
* Check all expected behaviors

## 👩‍💻 Author

Made with ❤️ by Fatima Amani
GitHub: [@fatima-amani](https://github.com/fatima-amani)
