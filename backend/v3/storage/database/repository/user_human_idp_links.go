package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

func (u userHuman) unqualifiedIDPLinksTableName() string {
	return "user_identity_provider_links"
}

func (u userHuman) qualifiedIDPLinksTableName() string {
	return "zitadel." + u.unqualifiedIDPLinksTableName()
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
			builder.WriteString("INSERT INTO ")
			builder.WriteString(u.qualifiedIDPLinksTableName())
			builder.WriteString("(instance_id, user_id, provided_user_id, provided_username, provider_id, created_at, updated_at) SELECT ")
			database.Columns{
				existingHumanUser.InstanceIDColumn(),
				existingHumanUser.IDColumn(),
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
			builder.WriteString("DELETE FROM ")
			builder.WriteString(u.qualifiedIDPLinksTableName())
			builder.WriteString(" USING ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(
					existingHumanUser.InstanceIDColumn(),
					u.linkedIdentityProviderInstanceIDColumn(),
				),
				database.NewColumnCondition(
					existingHumanUser.IDColumn(),
					u.linkedIdentityProviderInstanceIDColumn(),
				),
				database.NewTextCondition(
					u.linkedIdentityProviderIDColumn(),
					database.TextOperationEqual,
					providerID,
				),
				database.NewTextCondition(
					u.providedUserIDColumn(),
					database.TextOperationEqual,
					providedUserID,
				),
			))
		}, nil,
	)
}

// SetIdentityProviderLinkProvidedID implements [domain.HumanUserRepository.SetIdentityProviderLinkProvidedID].
func (u userHuman) SetIdentityProviderLinkProvidedID(providerID string, currentProvidedUserID string, newProvidedUserID string) database.Change {
	return database.NewChange(u.providedUserIDColumn(), newProvidedUserID)
}

// SetIdentityProviderLinkUsername implements [domain.HumanUserRepository.SetIdentityProviderLinkUsername].
func (u userHuman) SetIdentityProviderLinkUsername(providerID string, providedUserID string, username string) database.Change {
	return database.NewChange(u.providedUsernameColumn(), username)
}

// UpdateIdentityProviderLink implements [domain.HumanUserRepository.UpdateIdentityProviderLink].
func (u userHuman) UpdateIdentityProviderLink(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE zitadel.human_identity_provider_links SET ")
			database.Changes(changes).Write(builder)
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.linkedIdentityProviderInstanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.IDColumn(), u.linkedIdentityProviderInstanceIDColumn()),
				condition,
			))
		}, nil,
	)
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

func (u userHuman) linkedIdentityProviderInstanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "instance_id")
}

func (u userHuman) linkedIdentityProviderUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "user_id")
}

func (u userHuman) linkedIdentityProviderIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "identity_provider_id")
}

func (u userHuman) providedUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "provided_user_id")
}

func (u userHuman) providedUsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedIDPLinksTableName(), "provided_username")
}
