package view

import (
	"time"

	"github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func UserByIDQuery(id, instanceID string, creationDate time.Time) (*es_models.SearchQuery, error) {
	if id == "" {
		return nil, errors.ThrowPreconditionFailed(nil, "EVENT-d8isw", "Errors.User.UserIDMissing")
	}
	return es_models.NewSearchQuery().
		AddQuery().
		AggregateTypeFilter(es_models.AggregateType(user.AggregateType)).
		AggregateIDFilter(id).
		CreationDateNewerFilter(creationDate).
		InstanceIDFilter(instanceID).
		SearchQuery(), nil
}
