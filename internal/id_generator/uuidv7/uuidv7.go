package uuidv7

import (
	"github.com/google/uuid"
)

type Generator struct{}

func (s *Generator) Next() (string, error) {
	_id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return _id.String(), nil
}
