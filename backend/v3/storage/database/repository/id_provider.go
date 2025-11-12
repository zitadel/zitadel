package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var _ domain.IDProviderRepository = (*idProvider)(nil)

type idProvider struct{}

func IDProviderRepository() domain.IDProviderRepository {
	return new(idProvider)
}

func (i idProvider) qualifiedTableName() string {
	return "zitadel." + i.unqualifiedTableName()
}

func (idProvider) unqualifiedTableName() string {
	return "identity_providers"
}

func (i idProvider) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IdentityProvider, error) {
	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return get[domain.IdentityProvider](ctx, client, builder)
}

func (i idProvider) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.IdentityProvider, error) {
	builder, err := i.prepareQuery(opts)
	if err != nil {
		return nil, err
	}
	return list[domain.IdentityProvider](ctx, client, builder)
}

func (i idProvider) Create(ctx context.Context, client database.QueryExecutor, idp *domain.IdentityProvider) error {
	builder := database.NewStatementBuilder(`INSERT INTO `)
	builder.WriteString(i.qualifiedTableName())
	builder.WriteString(` (instance_id, organization_id, id, state, name, type, allow_creation, allow_auto_creation, allow_auto_update, allow_linking, auto_linking_field, styling_type, payload, created_at, updated_at) VALUES (`)
	builder.WriteArgs(
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
		defaultValue(idp.AutoLinkingField),
		idp.StylingType,
		string(idp.Payload),
		defaultTimestamp(idp.CreatedAt),
		defaultTimestamp(idp.UpdatedAt),
	)
	builder.WriteString(`) RETURNING created_at, updated_at`)

	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&idp.CreatedAt, &idp.UpdatedAt)
	return err
}

func (i idProvider) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	return update(ctx, client, i, condition, changes...)
}

func (i idProvider) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, i.InstanceIDColumn(), i.OrgIDColumn()); err != nil {
		return 0, err
	}
	return delete(ctx, client, i, condition)
}

func getIDPType[Target any](ctx context.Context, client database.QueryExecutor, i idProvider, t domain.IDPType, opts ...database.QueryOption) (*Target, error) {
	idp, err := i.Get(ctx, client, opts...)
	if err != nil {
		return nil, err
	}

	if idp.Type != nil && *idp.Type != t {
		return nil, domain.NewIDPWrongTypeError(t, *idp.Type)
	}
	var idpType Target
	if err = json.Unmarshal(idp.Payload, &idpType); err != nil {
		return nil, err
	}
	return &idpType, nil
}

func (i idProvider) GetOIDC(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPOIDC, error) {
	return getIDPType[domain.IDPOIDC](ctx, client, i, domain.IDPTypeOIDC, opts...)

}

func (i idProvider) GetJWT(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPJWT, error) {
	return getIDPType[domain.IDPJWT](ctx, client, i, domain.IDPTypeOIDC, opts...)
}

func (i idProvider) GetOAuth(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPOAuth, error) {
	return getIDPType[domain.IDPOAuth](ctx, client, i, domain.IDPTypeOAuth, opts...)
}

func (i idProvider) GetAzureAD(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPAzureAD, error) {
	return getIDPType[domain.IDPAzureAD](ctx, client, i, domain.IDPTypeAzure, opts...)
}

func (i idProvider) GetGoogle(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGoogle, error) {
	return getIDPType[domain.IDPGoogle](ctx, client, i, domain.IDPTypeGoogle, opts...)
}

func (i idProvider) GetGithub(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGithub, error) {
	return getIDPType[domain.IDPGithub](ctx, client, i, domain.IDPTypeGitHub, opts...)
}

func (i idProvider) GetGithubEnterprise(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGithubEnterprise, error) {
	return getIDPType[domain.IDPGithubEnterprise](ctx, client, i, domain.IDPTypeGitHubEnterprise, opts...)
}

func (i idProvider) GetGitlab(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGitlab, error) {
	return getIDPType[domain.IDPGitlab](ctx, client, i, domain.IDPTypeGitLab, opts...)
}

func (i idProvider) GetGitlabSelfHosting(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGitlabSelfHosting, error) {
	return getIDPType[domain.IDPGitlabSelfHosting](ctx, client, i, domain.IDPTypeGitLabSelfHosted, opts...)
}

func (i idProvider) GetLDAP(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPLDAP, error) {
	return getIDPType[domain.IDPLDAP](ctx, client, i, domain.IDPTypeLDAP, opts...)
}

func (i idProvider) GetApple(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPApple, error) {
	return getIDPType[domain.IDPApple](ctx, client, i, domain.IDPTypeApple, opts...)
}

func (i idProvider) GetSAML(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPSAML, error) {
	return getIDPType[domain.IDPSAML](ctx, client, i, domain.IDPTypeSAML, opts...)
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

// PrimaryKeyColumns implements [domain.Repository].
func (i idProvider) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		i.InstanceIDColumn(),
		i.IDColumn(),
	}
}

func (i idProvider) InstanceIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "instance_id")
}

func (i idProvider) OrgIDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "organization_id")
}

func (i idProvider) IDColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "id")
}

func (i idProvider) StateColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "state")
}

func (i idProvider) NameColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "name")
}

func (i idProvider) TypeColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "type")
}

func (i idProvider) AutoRegisterColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "auto_register")
}

func (i idProvider) AllowCreationColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "allow_creation")
}

func (i idProvider) AllowAutoCreationColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "allow_auto_creation")
}

func (i idProvider) AllowAutoUpdateColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "allow_auto_update")
}

func (i idProvider) AllowLinkingColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "allow_linking")
}

func (i idProvider) AllowAutoLinkingColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "auto_linking_field")
}

func (i idProvider) StylingTypeColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "styling_type")
}

func (i idProvider) PayloadColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "payload")
}

func (i idProvider) CreatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "created_at")
}

func (i idProvider) UpdatedAtColumn() database.Column {
	return database.NewColumn(i.unqualifiedTableName(), "updated_at")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (i idProvider) PrimaryKeyCondition(instanceID, id string) database.Condition {
	return database.And(
		i.InstanceIDCondition(instanceID),
		i.IDCondition(id),
	)
}

func (i idProvider) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) OrgIDCondition(id *string) database.Condition {
	if id == nil {
		return database.IsNull(i.OrgIDColumn())
	}
	return database.NewTextCondition(i.OrgIDColumn(), database.TextOperationEqual, *id)
}

func (i idProvider) IDCondition(id string) database.Condition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) StateCondition(state domain.IDPState) database.Condition {
	return database.NewTextCondition(i.StateColumn(), database.TextOperationEqual, state.String())
}

func (i idProvider) NameCondition(name string) database.Condition {
	return database.NewTextCondition(i.NameColumn(), database.TextOperationEqual, name)
}

func (i idProvider) TypeCondition(typ domain.IDPType) database.Condition {
	return database.NewNumberCondition(i.TypeColumn(), database.NumberOperationEqual, typ)
}

func (i idProvider) AutoRegisterCondition(allow bool) database.Condition {
	return database.NewBooleanCondition(i.AutoRegisterColumn(), allow)
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

func (i idProvider) AllowAutoLinkingCondition(linkingType domain.IDPAutoLinkingField) database.Condition {
	return database.NewNumberCondition(i.AllowAutoLinkingColumn(), database.NumberOperationEqual, linkingType)
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

func (i idProvider) SetAutoRegister(allow bool) database.Change {
	return database.NewChange(i.AutoRegisterColumn(), allow)
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

func (i idProvider) SetAutoAllowLinking(allow bool) database.Change {
	return database.NewChange(i.AllowAutoLinkingColumn(), allow)
}

func (i idProvider) SetStylingType(stylingType int16) database.Change {
	return database.NewChange(i.StylingTypeColumn(), stylingType)
}

func (i idProvider) SetPayload(payload string) database.Change {
	return database.NewChange(i.PayloadColumn(), payload)
}

func (i idProvider) SetUpdatedAt(updatedAt *time.Time) database.Change {
	return database.NewChangePtr(i.UpdatedAtColumn(), updatedAt)
}

// -------------------------------------------------------------
// helpers
// -------------------------------------------------------------

const queryIDProviderStmt = `SELECT instance_id, organization_id, id, state, name, type, auto_register, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, auto_linking_field, styling_type, payload, created_at, updated_at` +
	` FROM `

func (i idProvider) prepareQuery(opts []database.QueryOption) (*database.StatementBuilder, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, i.InstanceIDColumn()); err != nil {
		return nil, err
	}
	builder := database.NewStatementBuilder(queryIDProviderStmt + i.qualifiedTableName())
	options.Write(builder)

	return builder, nil
}
