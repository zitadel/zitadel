package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userIdentityProviderLinkRepo struct{}

func (u userIdentityProviderLinkRepo) unqualifiedTableName() string {
	return "identity_provider_links"
}

func (u userIdentityProviderLinkRepo) qualifiedTableName() string {
	return "zitadel.identity_provider_links"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddIdentityProviderLink implements [domain.HumanUserRepository.AddIdentityProviderLink].
func (u userIdentityProviderLinkRepo) AddIdentityProviderLink(link *domain.IdentityProviderLink) database.Change {
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
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString("(instance_id, user_id, provided_user_id, provided_username, identity_provider_id, created_at, updated_at) SELECT ")
			database.Columns{
				existingHumanUser.InstanceIDColumn(),
				existingHumanUser.IDColumn(),
			}.WriteQualified(builder)
			builder.WriteString(", ")
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

// UpdateIdentityProviderLink implements [domain.HumanUserRepository.UpdateIdentityProviderLink].
func (u userIdentityProviderLinkRepo) UpdateIdentityProviderLink(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE zitadel.identity_provider_links SET ")
			database.Changes(changes).Write(builder)
			builder.WriteString(" FROM ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(existingHumanUser.InstanceIDColumn(), u.instanceIDColumn()),
				database.NewColumnCondition(existingHumanUser.IDColumn(), u.userIDColumn()),
				condition,
			))
		}, nil,
	)
}

// RemoveIdentityProviderLink implements [domain.HumanUserRepository.RemoveIdentityProviderLink].
func (u userIdentityProviderLinkRepo) RemoveIdentityProviderLink(providerID string, providedUserID string) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("DELETE FROM ")
			builder.WriteString(u.qualifiedTableName())
			builder.WriteString(" USING ")
			builder.WriteString(existingHumanUser.unqualifiedTableName())
			writeCondition(builder, database.And(
				database.NewColumnCondition(
					existingHumanUser.InstanceIDColumn(),
					u.instanceIDColumn(),
				),
				database.NewColumnCondition(
					existingHumanUser.IDColumn(),
					u.userIDColumn(),
				),
				database.NewTextCondition(
					u.providerIDColumn(),
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

// SetIdentityProviderLinkUsername implements [domain.HumanUserRepository.SetIdentityProviderLinkUsername].
func (u userIdentityProviderLinkRepo) SetIdentityProviderLinkUsername(username string) database.Change {
	return database.NewChange(u.providedUsernameColumn(), username)
}

// SetIdentityProviderLinkProvidedID implements [domain.HumanUserRepository.SetIdentityProviderLinkProvidedID].
func (u userIdentityProviderLinkRepo) SetIdentityProviderLinkProvidedID(providerID string, currentProvidedUserID string, newProvidedUserID string) database.Change {
	return database.NewChange(u.providedUserIDColumn(), newProvidedUserID)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// ExistsIdentityProviderLink implements [domain.HumanUserRepository.ExistsIdentityProviderLink].
func (u userIdentityProviderLinkRepo) ExistsIdentityProviderLink(condition database.Condition) database.Condition {
	panic("unimplemented")
}

// IdentityProviderLinkConditions implements [domain.HumanUserRepository.IdentityProviderLinkConditions].
func (u userIdentityProviderLinkRepo) IdentityProviderLinkConditions() domain.HumanIdentityProviderLinkConditions {
	return u
}

// ProviderIDCondition implements [domain.HumanUserRepository.ProviderIDCondition].
func (u userIdentityProviderLinkRepo) ProviderIDCondition(idpID string) database.Condition {
	return database.NewTextCondition(u.providerIDColumn(), database.TextOperationEqual, idpID)
}

// ProvidedUserIDCondition implements [domain.HumanUserRepository.ProvidedUserIDCondition].
func (u userIdentityProviderLinkRepo) ProvidedUserIDCondition(providedUserID string) database.Condition {
	return database.NewTextCondition(u.providedUserIDColumn(), database.TextOperationEqual, providedUserID)
}

// ProvidedUsernameCondition implements [domain.HumanUserRepository.ProvidedUsernameCondition].
func (u userIdentityProviderLinkRepo) ProvidedUsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.providedUsernameColumn(), op, username)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userIdentityProviderLinkRepo) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}
func (u userIdentityProviderLinkRepo) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userIdentityProviderLinkRepo) providerIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "identity_provider_id")
}

func (u userIdentityProviderLinkRepo) providedUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "provided_user_id")
}

func (u userIdentityProviderLinkRepo) providedUsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "provided_username")
}
