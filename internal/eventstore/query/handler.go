package query

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Handler interface {
	ViewModel() string
	EventQuery() (*models.SearchQuery, error)
	Reduce(*models.Event) error
	OnError(event *models.Event, err error) error
}
