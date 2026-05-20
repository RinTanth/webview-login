package util

import "github.com/google/uuid"

// NewUUID generates a new time-ordered UUID v7.
func NewUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate UUID v7: " + err.Error())
	}
	return id
}
