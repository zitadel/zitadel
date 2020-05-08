package query

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

type Handler interface {
	ViewModel() string
	EventQuery() (*models.SearchQuery, error)
	Process(*models.Event) error
}
