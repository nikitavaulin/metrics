package filestorage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/nikitavaulin/metrics/internal/validation"
)

type FileStorage struct {
}

func New() *FileStorage {
	return &FileStorage{}
}

func (fs *FileStorage) SaveToJSON(filename string, data any) error {
	if err := validation.ValidateFileExt(filename, ".json"); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(filename, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to save json to file: %w", err)
	}

	return nil
}

func (fs *FileStorage) LoadFromJSON(filename string, dest any) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %s: %w", filename, err)
	}

	if err := json.Unmarshal(jsonData, dest); err != nil {
		return fmt.Errorf("failed to unmarshal json")
	}

	return nil
}
