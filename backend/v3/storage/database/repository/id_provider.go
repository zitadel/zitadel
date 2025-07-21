package repository

import (
	"context"
	"errors"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.IDProviderRepository = (*idProvider)(nil)

type idProvider struct {
	repository
}

func IDProviderRepository(client database.QueryExecutor) domain.IDProviderRepository {
	return &idProvider{
		repository: repository{
			client: client,
		},
	}
}

const queryIDProviderStmt = `SELECT instance_id, org_id, id, state, name, type, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, styling_type, payload, created_at, updated_at` +
	` FROM zitadel.identity_providers`

func (i *idProvider) Get(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID string) (*domain.IdentityProvider, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryIDProviderStmt)

	conditions := []database.Condition{id, i.InstanceIDCondition(instnaceID), i.OrgIDCondition(orgID)}

	writeCondition(&builder, database.And(conditions...))

	return scanIDProvider(ctx, i.client, &builder)
}

func (i *idProvider) List(ctx context.Context, conditions ...database.Condition) ([]*domain.IdentityProvider, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryIDProviderStmt)

	if conditions != nil {
		writeCondition(&builder, database.And(conditions...))
	}

	orderBy := database.OrderBy(i.CreatedAtColumn())
	orderBy.Write(&builder)

	return scanIDProviders(ctx, i.client, &builder)
}

const createIDProviderStmt = `INSERT INTO zitadel.identity_providers` +
	` (instance_id, org_id, id, state, name, type, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, styling_type, payload)` +
	` VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)` +
	` RETURNING created_at, updated_at`

func (i *idProvider) Create(ctx context.Context, idp *domain.IdentityProvider) error {
	builder := database.StatementBuilder{}
	builder.AppendArgs(
		idp.InstanceID,
		idp.OrgID,
		idp.ID,
		idp.State,
		idp.Name,
		idp.Type,
		idp.AllowCreation,
		idp.AllowAutoCreation,
		idp.AllowAutoUpdate,
		idp.AllowLinking,
		idp.StylingType,
		idp.Payload)
	builder.WriteString(createIDProviderStmt)

	err := i.client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&idp.CreatedAt, &idp.UpdatedAt)
	if err != nil {
		return checkCreateOrgErr(err)
	}
	return nil
}

func (i *idProvider) Update(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID string, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, errors.New("Update must contain at least one change")
	}
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.identity_providers SET `)

	conditions := []database.Condition{
		id,
		i.InstanceIDCondition(instnaceID),
		i.OrgIDCondition(orgID),
	}
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return i.client.Exec(ctx, stmt, builder.Args()...)
}

func (i *idProvider) Delete(ctx context.Context, id domain.IDPIdentifierCondition) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.identity_providers`)

	// conditions := []database.Condition{i.IDCondition(id)}
	// writeCondition(&builder, database.And(conditions...))

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

func (idProvider) InstanceIDColumn() database.Column {
	return database.NewColumn("instance_id")
}

func (idProvider) OrgIDColumn() database.Column {
	return database.NewColumn("org_id")
}

func (idProvider) IDColumn() database.Column {
	return database.NewColumn("id")
}

func (idProvider) StateColumn() database.Column {
	return database.NewColumn("state")
}

func (idProvider) NameColumn() database.Column {
	return database.NewColumn("name")
}

func (idProvider) TypeColumn() database.Column {
	return database.NewColumn("type")
}

func (idProvider) AllowCreationColumn() database.Column {
	return database.NewColumn("allow_creation")
}

func (idProvider) AllowAutoCreationColumn() database.Column {
	return database.NewColumn("allow_auto_creation")
}

func (idProvider) AllowAutoUpdateColumn() database.Column {
	return database.NewColumn("allow_auto_update")
}

func (idProvider) AllowLinkingColumn() database.Column {
	return database.NewColumn("allow_linking")
}

func (idProvider) StylingTypeColumn() database.Column {
	return database.NewColumn("styling_type")
}

func (idProvider) PayloadColumn() database.Column {
	return database.NewColumn("payload")
}

func (idProvider) CreatedAtColumn() database.Column {
	return database.NewColumn("created_at")
}

func (idProvider) UpdatedAtColumn() database.Column {
	return database.NewColumn("updated_at")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (i idProvider) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) OrgIDCondition(id string) database.Condition {
	return database.NewTextCondition(i.OrgIDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) IDCondition(id string) domain.IDPIdentifierCondition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) StateCondition(state domain.IDPState) database.Condition {
	return database.NewTextCondition(i.OrgIDColumn(), database.TextOperationEqual, state.String())
}

func (i idProvider) NameCondition(name string) domain.IDPIdentifierCondition {
	return database.NewTextCondition(i.NameColumn(), database.TextOperationEqual, name)
}

func (i idProvider) TypeCondition(typee domain.IDPType) database.Condition {
	return database.NewTextCondition(i.TypeColumn(), database.TextOperationEqual, typee.String())
}

func (i idProvider) AllowCreationCondition(allow bool) database.Condition {
	return database.NewBooleanCondition(i.AllowCreationColumn(), allow)
}

func (i idProvider) AllowAutoCreationCondition(allow bool) database.Condition {
	return database.NewBooleanCondition(i.AllowAutoCreationColumn(), allow)
}

func (i idProvider) AllowAutoUpdateCondition(allow bool) database.Condition {
	return database.NewBooleanCondition(i.AllowAutoUpdateColumn(), allow)
}

func (i idProvider) AllowLinkingCondition(allow bool) database.Condition {
	return database.NewBooleanCondition(i.AllowLinkingColumn(), allow)
}

func (i idProvider) StylingTypeCondition(style int16) database.Condition {
	return database.NewNumberCondition(i.StylingTypeColumn(), database.NumberOperationEqual, style)
}

func (i idProvider) PayloadCondition(payload string) database.Condition {
	return database.NewTextCondition(i.PayloadColumn(), database.TextOperationEqual, payload)
}

// -------------------------------------------------------------
// changes
// -------------------------------------------------------------

func (i idProvider) SetName(name string) database.Change {
	return database.NewChange(i.NameColumn(), name)
}

func (i idProvider) SetState(state domain.IDPState) database.Change {
	return database.NewChange(i.StateColumn(), state)
}

func (i idProvider) SetAllowCreation(allow bool) database.Change {
	return database.NewChange(i.AllowCreationColumn(), allow)
}

func (i idProvider) SetAllowAutoCreation(allow bool) database.Change {
	return database.NewChange(i.AllowAutoCreationColumn(), allow)
}

func (i idProvider) SetAllowAutoUpdate(allow bool) database.Change {
	return database.NewChange(i.AllowAutoUpdateColumn(), allow)
}

func (i idProvider) SetAllowLinking(allow bool) database.Change {
	return database.NewChange(i.AllowLinkingColumn(), allow)
}

func (i idProvider) SetStylingType(stylingType int16) database.Change {
	return database.NewChange(i.StylingTypeColumn(), stylingType)
}

func (i idProvider) SetPayload(payload string) database.Change {
	return database.NewChange(i.PayloadColumn(), payload)
}

func scanIDProvider(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.IdentityProvider, error) {
	idp := &domain.IdentityProvider{}
	err := scan(ctx, querier, builder, idp)
	if err != nil {
		return nil, err
	}
	return idp, err
}

func scanIDProviders(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.IdentityProvider, error) {
	idps := []*domain.IdentityProvider{}
	err := scanMultiple(ctx, querier, builder, &idps)
	if err != nil {
		return nil, err
	}
	return idps, nil
}
