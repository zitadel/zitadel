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
	var createdAt, updatedAt any = database.NowInstruction, database.NowInstruction
	if !link.CreatedAt.IsZero() {
		createdAt = link.CreatedAt
	}
	if !link.UpdatedAt.IsZero() {
		updatedAt = link.UpdatedAt
	}
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("INSERT INTO zitadel.human_identity_provider_links" +
				"(instance_id, user_id, provided_user_id, provided_username, provider_id, created_at, updated_at) SELECT ")
			database.Columns{
				existingHumanUser.instanceIDColumn(),
				existingHumanUser.idColumn(),
			}.WriteQualified(builder)
			builder.WriteArgs(
				link.ProvidedUserID,
				link.ProvidedUsername,
				link.ProviderID,
				createdAt,
				updatedAt,
			)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			builder.WriteString(" RETURNING *")
		}, nil,
	)
}

// RemoveIdentityProviderLink implements [domain.HumanUserRepository.RemoveIdentityProviderLink].
func (u userHuman) RemoveIdentityProviderLink(providerID string, providedUserID string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM zitadel.human_identity_provider_links USING ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(
					existingHumanUser.instanceIDColumn(),
					database.NewColumn("human_identity_provider_links", "instance_id"),
				),
				database.NewColumnCondition(
					existingHumanUser.idColumn(),
					database.NewColumn("human_identity_provider_links", "user_id"),
				),
				database.NewTextCondition(
					database.NewColumn("human_identity_provider_links", "provider_id"),
					database.TextOperationEqual,
					providerID,
				),
				database.NewTextCondition(
					database.NewColumn("human_identity_provider_links", "provided_user_id"),
					database.TextOperationEqual,
					providedUserID,
				),
			))
		}, nil,
	)
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
