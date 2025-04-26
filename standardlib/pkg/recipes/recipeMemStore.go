package recipes

import "sync"

type MemStore struct {
	mu      sync.RWMutex
	recipes map[string]Recipe
}

func NewMemStore() *MemStore {
	return &MemStore{
		recipes: make(map[string]Recipe),
	}
}

func (m *MemStore) Add(name string, recipe Recipe) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recipes[name] = recipe
	return nil
}

func (m *MemStore) Get(name string) (Recipe, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	recipe, ok := m.recipes[name]
	if !ok {
		return Recipe{}, ErrNotFound
	}
	return recipe, nil
}

func (m *MemStore) Update(name string, recipe Recipe) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.recipes[name] = recipe
	return nil
}

func (m *MemStore) List() (map[string]Recipe, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	recipesCopy := make(map[string]Recipe)
	for k, v := range m.recipes {
		recipesCopy[k] = v
	}
	return recipesCopy, nil
}

func (m *MemStore) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.recipes, name)
	return nil
}

// ErrNotFound error if recipe doesn't exist
var ErrNotFound = errorString("recipe not found")

type errorString string

func (e errorString) Error() string {
	return string(e)
}
