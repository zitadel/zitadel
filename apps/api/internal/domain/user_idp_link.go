package domain

import es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"

type UserIDPLink struct {
	es_models.ObjectRoot

	IDPConfigID    string
	ExternalUserID string
	DisplayName    string
}

func (idp *UserIDPLink) IsValid() bool {
	return idp.IDPConfigID != "" && idp.ExternalUserID != ""
}

type UserIDPLinkState int32

const (
	UserIDPLinkStateUnspecified UserIDPLinkState = iota
	UserIDPLinkStateActive
	UserIDPLinkStateRemoved

	userIDPLinkStateCount
)

func (s UserIDPLinkState) Valid() bool {
	return s >= 0 && s < userIDPLinkStateCount
}
