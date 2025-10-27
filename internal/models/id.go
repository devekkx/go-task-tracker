package models

import (
	"crypto/rand"
	"fmt"
)

// generateID creates a random prefixed identifier.
func generateID(prefix string) string {
	b := make([]byte, 4)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s_%x", prefix, b)
}
