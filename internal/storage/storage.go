package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/devekkx/go-task-tracker/internal/models"
)

const dataFileName = "data.json"

// Store holds all application data.
type Store struct {
	Tasks     []models.Task     `json:"tasks"`
	TodoLists []models.TodoList `json:"todo_lists"`
	path      string
}

// New creates a Store backed by the user's home directory.
func New() (*Store, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine data directory: %w", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	s := &Store{
		Tasks:     []models.Task{},
		TodoLists: []models.TodoList{},
		path:      filepath.Join(dir, dataFileName),
	}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".task-tracker"), nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to read data file: %w", err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("failed to parse data file: %w", err)
	}
	return nil
}

// save flushes the store to disk
func (s *Store) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}
	if err := os.WriteFile(s.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}
	return nil
}

// DataPath returns the path to the backing JSON file.
func (s *Store) DataPath() string { return s.path }
