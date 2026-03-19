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

func (idProvider) qualifiedTableName() string {
	return "zitadel.identity_providers"
}

const queryIDProviderStmt = `SELECT instance_id, org_id, id, state, name, type, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, auto_linking_field, payload, created_at, updated_at` +
	` FROM zitadel.identity_providers`

func (i *idProvider) Get(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IdentityProvider, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}
	if err := checkRestrictingColumns(options.Condition, i.InstanceIDColumn()); err != nil {
		return nil, err
	}

	builder := database.NewStatementBuilder(queryIDProviderStmt)
	options.Write(builder)

	return scanIDProvider(ctx, client, builder)
}

func (i *idProvider) List(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) ([]*domain.IdentityProvider, error) {
	options := new(database.QueryOpts)
	for _, opt := range opts {
		opt(options)
	}

	if err := checkRestrictingColumns(options.Condition, i.InstanceIDColumn()); err != nil {
		return nil, err
	}

	builder := database.NewStatementBuilder(queryIDProviderStmt)
	options.Write(builder)

	return scanIDProviders(ctx, client, builder)
}

const (
	createIDProviderStmtStart = `INSERT INTO zitadel.identity_providers` +
		` (instance_id, org_id, id, state, name, type, allow_creation, allow_auto_creation,` +
		` allow_auto_update, allow_linking, auto_linking_field, payload, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, `
	createIDProviderStmtEnd = `) RETURNING created_at, updated_at`
)

func (i *idProvider) Create(ctx context.Context, client database.QueryExecutor, idp *domain.IdentityProvider) error {
	builder := database.NewStatementBuilder(createIDProviderStmtStart,
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
		idp.AutoLinkingField,
		idp.Payload,
	)
	var createdAt any = database.NowInstruction
	if !idp.CreatedAt.IsZero() {
		createdAt = idp.CreatedAt
	}
	builder.WriteArgs(createdAt, createdAt)
	builder.WriteString(createIDProviderStmtEnd)

	return client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&idp.CreatedAt, &idp.UpdatedAt)
}

func (i *idProvider) Update(ctx context.Context, client database.QueryExecutor, condition database.Condition, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, database.ErrNoChanges
	}
	if err := checkRestrictingColumns(condition, i.InstanceIDColumn(), i.IDColumn()); err != nil {
		return 0, err
	}
	dbChanges := database.Changes(changes)
	if !dbChanges.IsOnColumn(i.UpdatedAtColumn()) {
		dbChanges = append(dbChanges, i.SetUpdatedAt(nil))
	}

	builder := database.NewStatementBuilder(`UPDATE zitadel.identity_providers SET `)
	err := dbChanges.Write(builder)
	if err != nil {
		return 0, err
	}
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (i *idProvider) Delete(ctx context.Context, client database.QueryExecutor, condition database.Condition) (int64, error) {
	if err := checkRestrictingColumns(condition, i.InstanceIDColumn(), i.IDColumn()); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder(`DELETE FROM `)
	builder.WriteString(i.qualifiedTableName())
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (i *idProvider) GetOIDC(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPOIDC, error) {
	idp, typ, err := getIDP[domain.OIDC](ctx, client, i, domain.IDPTypeOIDC, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPOIDC{
		IdentityProvider: idp,
		OIDC:             typ,
	}, nil
}

func (i *idProvider) GetJWT(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPJWT, error) {
	idp, typ, err := getIDP[domain.JWT](ctx, client, i, domain.IDPTypeJWT, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPJWT{
		IdentityProvider: idp,
		JWT:              typ,
	}, nil
}

func (i *idProvider) GetOAuth(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPOAuth, error) {
	idp, typ, err := getIDP[domain.OAuth](ctx, client, i, domain.IDPTypeOAuth, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPOAuth{
		IdentityProvider: idp,
		OAuth:            typ,
	}, nil
}

func (i *idProvider) GetAzureAD(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPAzureAD, error) {
	idp, typ, err := getIDP[domain.Azure](ctx, client, i, domain.IDPTypeAzure, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPAzureAD{
		IdentityProvider: idp,
		Azure:            typ,
	}, nil
}

func (i *idProvider) GetGoogle(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGoogle, error) {
	idp, typ, err := getIDP[domain.Google](ctx, client, i, domain.IDPTypeGoogle, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPGoogle{
		IdentityProvider: idp,
		Google:           typ,
	}, nil
}

func (i *idProvider) GetGithub(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGithub, error) {
	idp, typ, err := getIDP[domain.Github](ctx, client, i, domain.IDPTypeGitHub, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPGithub{
		IdentityProvider: idp,
		Github:           typ,
	}, nil
}

func (i *idProvider) GetGithubEnterprise(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGithubEnterprise, error) {
	idp, typ, err := getIDP[domain.GithubEnterprise](ctx, client, i, domain.IDPTypeGitHubEnterprise, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPGithubEnterprise{
		IdentityProvider: idp,
		GithubEnterprise: typ,
	}, nil
}

func (i *idProvider) GetGitlab(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGitlab, error) {
	idp, typ, err := getIDP[domain.Gitlab](ctx, client, i, domain.IDPTypeGitLab, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPGitlab{
		IdentityProvider: idp,
		Gitlab:           typ,
	}, nil
}

func (i *idProvider) GetGitlabSelfHosted(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPGitlabSelfHosted, error) {
	idp, typ, err := getIDP[domain.GitlabSelfHosted](ctx, client, i, domain.IDPTypeGitLabSelfHosted, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPGitlabSelfHosted{
		IdentityProvider: idp,
		GitlabSelfHosted: typ,
	}, nil
}

func (i *idProvider) GetLDAP(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPLDAP, error) {
	idp, typ, err := getIDP[domain.LDAP](ctx, client, i, domain.IDPTypeLDAP, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPLDAP{
		IdentityProvider: idp,
		LDAP:             typ,
	}, nil
}

func (i *idProvider) GetApple(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPApple, error) {
	idp, typ, err := getIDP[domain.Apple](ctx, client, i, domain.IDPTypeApple, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPApple{
		IdentityProvider: idp,
		Apple:            typ,
	}, nil
}

func (i *idProvider) GetSAML(ctx context.Context, client database.QueryExecutor, opts ...database.QueryOption) (*domain.IDPSAML, error) {
	idp, typ, err := getIDP[domain.SAML](ctx, client, i, domain.IDPTypeSAML, opts...)
	if err != nil {
		return nil, err
	}

	return &domain.IDPSAML{
		IdentityProvider: idp,
		SAML:             typ,
	}, nil
}

type idpType interface {
	domain.OIDC | domain.JWT | domain.OAuth | domain.Azure | domain.Google |
		domain.Github | domain.GithubEnterprise | domain.Gitlab | domain.GitlabSelfHosted |
		domain.LDAP | domain.Apple | domain.SAML
}

func getIDP[T idpType](ctx context.Context, client database.QueryExecutor, repo *idProvider, expectedType domain.IDPType, opts ...database.QueryOption) (*domain.IdentityProvider, T, error) {
	var zero T
	idp, err := repo.Get(ctx, client, opts...)
	if err != nil {

		return nil, zero, err
	}

	var idpType domain.IDPType
	if idp.Type != nil {
		idpType = *idp.Type
	}

	if idpType != expectedType {
		return nil, zero, domain.NewIDPWrongTypeError(expectedType, idpType)
	}

	var specificIDP T
	err = json.Unmarshal(idp.Payload, &specificIDP)
	if err != nil {
		return nil, zero, err
	}

	return idp, specificIDP, nil
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

func (idProvider) InstanceIDColumn() database.Column {
	return database.NewColumn("identity_providers", "instance_id")
}

func (idProvider) OrgIDColumn() database.Column {
	return database.NewColumn("identity_providers", "org_id")
}

func (idProvider) IDColumn() database.Column {
	return database.NewColumn("identity_providers", "id")
}

func (idProvider) StateColumn() database.Column {
	return database.NewColumn("identity_providers", "state")
}

func (idProvider) NameColumn() database.Column {
	return database.NewColumn("identity_providers", "name")
}

func (idProvider) TypeColumn() database.Column {
	return database.NewColumn("identity_providers", "type")
}

func (idProvider) AllowCreationColumn() database.Column {
	return database.NewColumn("identity_providers", "allow_creation")
}

func (idProvider) AllowAutoCreationColumn() database.Column {
	return database.NewColumn("identity_providers", "allow_auto_creation")
}

func (idProvider) AllowAutoUpdateColumn() database.Column {
	return database.NewColumn("identity_providers", "allow_auto_update")
}

func (idProvider) AllowLinkingColumn() database.Column {
	return database.NewColumn("identity_providers", "allow_linking")
}

func (idProvider) AutoLinkingFieldColumn() database.Column {
	return database.NewColumn("identity_providers", "auto_linking_field")
}

func (idProvider) PayloadColumn() database.Column {
	return database.NewColumn("identity_providers", "payload")
}

func (idProvider) CreatedAtColumn() database.Column {
	return database.NewColumn("identity_providers", "created_at")
}

func (idProvider) UpdatedAtColumn() database.Column {
	return database.NewColumn("identity_providers", "updated_at")
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

func (i idProvider) NameCondition(op database.TextOperation, name string) database.Condition {
	return database.NewTextCondition(i.NameColumn(), op, name)
}

func (i idProvider) TypeCondition(typ domain.IDPType) database.Condition {
	return database.NewNumberCondition(i.TypeColumn(), database.NumberOperationEqual, typ)
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

func (i idProvider) SetAutoLinkingField(field *domain.IDPAutoLinkingField) database.Change {
	return database.NewChangePtr(i.AutoLinkingFieldColumn(), field)
}

func (i idProvider) SetPayload(payload string) database.Change {
	return database.NewChange(i.PayloadColumn(), payload)
}

func (i idProvider) SetUpdatedAt(updatedAt *time.Time) database.Change {
	return database.NewChangePtr(i.UpdatedAtColumn(), updatedAt)
}

func (i idProvider) SetType(typee domain.IDPType) database.Change {
	return database.NewChange(i.TypeColumn(), typee)
}

func scanIDProvider(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*domain.IdentityProvider, error) {
	idp := &domain.IdentityProvider{}
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).CollectExactlyOneRow(idp)
	if err != nil {
		return nil, err
	}

	return idp, err
}

func scanIDProviders(ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*domain.IdentityProvider, error) {
	idps := []*domain.IdentityProvider{}

	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}

	err = rows.(database.CollectableRows).Collect(&idps)
	if err != nil {
		return nil, err
	}

	return idps, nil
}
