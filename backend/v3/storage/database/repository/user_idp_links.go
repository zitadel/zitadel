package repository

import (
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type userIdentityProviderLink struct{}

func (u userIdentityProviderLink) unqualifiedTableName() string {
	return "user_identity_provider_links"
}

func (u userIdentityProviderLink) qualifiedTableName() string {
	return "zitadel.user_identity_provider_links"
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

// AddIdentityProviderLink implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) AddIdentityProviderLink(link *domain.IdentityProviderLink) database.Change {
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

// UpdateIdentityProviderLink implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) UpdateIdentityProviderLink(condition database.Condition, changes ...database.Change) database.Change {
	return database.NewCTEChange(
		func(builder *database.StatementBuilder) {
			builder.WriteString("UPDATE zitadel.user_identity_provider_links SET ")
			err := database.Changes(changes).Write(builder)
			logging.New(logging.StreamRuntime).Debug("write changes in cte failed", "error", err)
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

// RemoveIdentityProviderLink implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) RemoveIdentityProviderLink(providerID string, providedUserID string) database.Change {
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

// SetIdentityProviderLinkUsername implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) SetIdentityProviderLinkUsername(username string) database.Change {
	return database.NewChange(u.providedUsernameColumn(), username)
}

// SetIdentityProviderLinkProvidedID implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) SetIdentityProviderLinkProvidedID(providedUserID string) database.Change {
	return database.NewChange(u.providedUserIDColumn(), providedUserID)
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

// IdentityProviderLinkConditions implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) IdentityProviderLinkConditions() domain.HumanIdentityProviderLinkConditions {
	return u
}

// ProviderIDCondition implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) ProviderIDCondition(idpID string) database.Condition {
	return database.NewTextCondition(u.providerIDColumn(), database.TextOperationEqual, idpID)
}

// ProvidedUserIDCondition implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) ProvidedUserIDCondition(providedUserID string) database.Condition {
	return database.NewTextCondition(u.providedUserIDColumn(), database.TextOperationEqual, providedUserID)
}

// ProvidedUsernameCondition implements [domain.HumanUserRepository].
func (u userIdentityProviderLink) ProvidedUsernameCondition(op database.TextOperation, username string) database.Condition {
	return database.NewTextCondition(u.providedUsernameColumn(), op, username)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (u userIdentityProviderLink) instanceIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "instance_id")
}
func (u userIdentityProviderLink) userIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "user_id")
}

func (u userIdentityProviderLink) providerIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "identity_provider_id")
}

func (u userIdentityProviderLink) providedUserIDColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "provided_user_id")
}

func (u userIdentityProviderLink) providedUsernameColumn() database.Column {
	return database.NewColumn(u.unqualifiedTableName(), "provided_username")
}
