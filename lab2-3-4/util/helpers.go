package util

import "strings"

func IsDuplicateKeyError(err error) bool {
	// This will depend on your specific database and error handling
	// For SQLite, you might check for a unique constraint violation
	return strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "unique constraint") ||
		strings.Contains(err.Error(), "duplicate key")
}
