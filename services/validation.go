package services

import (
	"fmt"
	"strings"

	"crud-api/models"
)

func validateID(id int64, entity string) error {
	if id < 1 {
		return &models.ValidationError{Msg: fmt.Sprintf("invalid %s ID", entity)}
	}
	return nil
}

// isBlank trata como vazio também um texto só com espaços.
func isBlank(s string) bool {
	return strings.TrimSpace(s) == ""
}
