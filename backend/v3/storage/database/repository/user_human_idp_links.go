package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func (u userHuman) unqualifiedIDPLinksTableName() string {
	return "user_identity_provider_links"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddIdentityProviderLink implements [domain.HumanUserRepository.AddIdentityProviderLink].
func (u userHuman) AddIdentityProviderLink(link *domain.IdentityProviderLink) database.Change {
	panic("unimplemented")
}

// RemoveIdentityProviderLink implements [domain.HumanUserRepository.RemoveIdentityProviderLink].
func (u userHuman) RemoveIdentityProviderLink(providerID string, providedUserID string) database.Change {
	panic("unimplemented")
}

// SetIdentityProviderLinkProvidedID implements [domain.HumanUserRepository.SetIdentityProviderLinkProvidedID].
func (u userHuman) SetIdentityProviderLinkProvidedID(providerID string, currentProvidedUserID string, newProvidedUserID string) database.Change {
	panic("unimplemented")
}

// SetIdentityProviderLinkUsername implements [domain.HumanUserRepository.SetIdentityProviderLinkUsername].
func (u userHuman) SetIdentityProviderLinkUsername(providerID string, providedUserID string, username string) database.Change {
	panic("unimplemented")
}

// UpdateIdentityProviderLink implements [domain.HumanUserRepository.UpdateIdentityProviderLink].
func (u userHuman) UpdateIdentityProviderLink(changes ...database.Change) database.Change {
	panic("unimplemented")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// LinkedIdentityProviderIDCondition implements [domain.HumanUserRepository.LinkedIdentityProviderIDCondition].
func (u userHuman) LinkedIdentityProviderIDCondition(idpID string) database.Condition {
	return database.NewTextCondition(u.linkedIdentityProviderIDColumn(), database.TextOperationEqual, idpID)
}

// ProvidedUserIDCondition implements [domain.HumanUserRepository.ProvidedUserIDCondition].
func (u userHuman) ProvidedUserIDCondition(providedUserID string) database.Condition {
	return database.NewTextCondition(u.providedUserIDColumn(), database.TextOperationEqual, providedUserID)
}

// ProvidedUsernameCondition implements [domain.HumanUserRepository.ProvidedUsernameCondition].
func (u userHuman) ProvidedUsernameCondition(username string) database.Condition {
	return database.NewTextCondition(u.providedUsernameColumn(), database.TextOperationEqual, username)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userHuman) linkedIdentityProviderIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "identity_provider_id")
}

func (u userHuman) providedUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "provided_user_id")
}

func (u userHuman) providedUsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "provided_username")
}
