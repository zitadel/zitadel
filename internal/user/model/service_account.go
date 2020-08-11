package model

import es_models "github.com/caos/zitadel/internal/eventstore/models"

type ServiceAccount struct {
	es_models.ObjectRoot

	Name        string
	Email       string
	Description string
}

type ServiceAccountState int32

const (
	ServiceAccountStateUnspecified = iota
	ServiceAccountStateActive
	ServiceAccountStateInactive
	ServiceAccountStateDeleted
	ServiceAccountStateLocked
)

func (sa *ServiceAccount) IsValid() bool {
	return sa.Name != "" && sa.Email != ""
}

type ServiceAccountSearchRequest struct{}

type ServiceAccountSearchResult struct{}

type ServiceAccountChanges struct{}

type ServiceAccountView struct{}
