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

const queryIDProviderStmt = `SELECT instance_id, org_id, id, state, name, type, auto_register, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, auto_linking_field, styling_type, payload, created_at, updated_at` +
	` FROM zitadel.identity_providers`

func (i *idProvider) Get(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IdentityProvider, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryIDProviderStmt)

	conditions := []database.Condition{id, i.InstanceIDCondition(instanceID), i.OrgIDCondition(orgID)}

	writeCondition(&builder, database.And(conditions...))

	return scanIDProvider(ctx, client, &builder)
}

func (i *idProvider) List(ctx context.Context, client database.QueryExecutor, conditions ...database.Condition) ([]*domain.IdentityProvider, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryIDProviderStmt)

	if conditions != nil {
		writeCondition(&builder, database.And(conditions...))
	}

	orderBy := database.OrderBy(i.CreatedAtColumn())
	orderBy.Write(&builder)

	return scanIDProviders(ctx, client, &builder)
}

const createIDProviderStmtStart = `INSERT INTO zitadel.identity_providers` +
	` (instance_id, org_id, id, state, name, type, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, styling_type, payload) VALUES (`

const createIDProviderStmtEnd = `) RETURNING created_at, updated_at`

func (i *idProvider) Create(ctx context.Context, client database.QueryExecutor, idp *domain.IdentityProvider) error {
	builder := database.StatementBuilder{}

	builder.WriteString(createIDProviderStmtStart)

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
		idp.StylingType,
		string(idp.Payload))

	builder.WriteString(createIDProviderStmtEnd)

	err := client.QueryRow(ctx, builder.String(), builder.Args()...).Scan(&idp.CreatedAt, &idp.UpdatedAt)
	return err
}

func (i *idProvider) Update(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string, changes ...database.Change) (int64, error) {
	if changes == nil {
		return 0, database.ErrNoChanges
	}
	changes = append(changes, i.SetUpdatedAt(nil))
	builder := database.StatementBuilder{}
	builder.WriteString(`UPDATE zitadel.identity_providers SET `)

	conditions := []database.Condition{
		id,
		i.InstanceIDCondition(instanceID),
		i.OrgIDCondition(orgID),
	}
	database.Changes(changes).Write(&builder)
	writeCondition(&builder, database.And(conditions...))

	stmt := builder.String()

	return client.Exec(ctx, stmt, builder.Args()...)
}

func (i *idProvider) Delete(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.identity_providers`)

	conditions := []database.Condition{
		id,
		i.InstanceIDCondition(instanceID),
		i.OrgIDCondition(orgID),
	}
	writeCondition(&builder, database.And(conditions...))

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

func (i *idProvider) GetOIDC(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPOIDC, error) {
	idpOIDC := &domain.IDPOIDC{}
	var err error

	idpOIDC.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpOIDC.Type != nil {
		idpType = *idpOIDC.Type
	}

	if idpType != domain.IDPTypeOIDC {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeOIDC, idpType)
	}

	err = json.Unmarshal(idpOIDC.Payload, idpOIDC)
	if err != nil {
		return nil, err
	}

	return idpOIDC, nil
}

func (i *idProvider) GetJWT(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPJWT, error) {
	idpJWT := &domain.IDPJWT{}
	var err error

	idpJWT.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpJWT.Type != nil {
		idpType = *idpJWT.Type
	}

	if idpType != domain.IDPTypeJWT {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeJWT, idpType)
	}

	err = json.Unmarshal(idpJWT.Payload, idpJWT)
	if err != nil {
		return nil, err
	}

	return idpJWT, nil
}

func (i *idProvider) GetOAuth(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPOAuth, error) {
	idpOAuth := &domain.IDPOAuth{}
	var err error

	idpOAuth.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpOAuth.Type != nil {
		idpType = *idpOAuth.Type
	}

	if idpType != domain.IDPTypeOAuth {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeOAuth, idpType)
	}

	err = json.Unmarshal(idpOAuth.Payload, idpOAuth)
	if err != nil {
		return nil, err
	}

	return idpOAuth, nil
}

func (i *idProvider) GetAzureAD(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPAzureAD, error) {
	idpAzure := &domain.IDPAzureAD{}
	var err error

	idpAzure.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpAzure.Type != nil {
		idpType = *idpAzure.Type
	}

	if idpType != domain.IDPTypeAzure {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeAzure, idpType)
	}

	err = json.Unmarshal(idpAzure.Payload, idpAzure)
	if err != nil {
		return nil, err
	}

	return idpAzure, nil
}

func (i *idProvider) GetGoogle(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGoogle, error) {
	idpGoogle := &domain.IDPGoogle{}
	var err error

	idpGoogle.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpGoogle.Type != nil {
		idpType = *idpGoogle.Type
	}

	if idpType != domain.IDPTypeGoogle {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeGoogle, idpType)
	}

	err = json.Unmarshal(idpGoogle.Payload, idpGoogle)
	if err != nil {
		return nil, err
	}

	return idpGoogle, nil
}

func (i *idProvider) GetGithub(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGithub, error) {
	idpGithub := &domain.IDPGithub{}
	var err error

	idpGithub.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpGithub.Type != nil {
		idpType = *idpGithub.Type
	}

	if idpType != domain.IDPTypeGitHub {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeGitHub, idpType)
	}

	err = json.Unmarshal(idpGithub.Payload, idpGithub)
	if err != nil {
		return nil, err
	}

	return idpGithub, nil
}

func (i *idProvider) GetGithubEnterprise(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGithubEnterprise, error) {
	idpGithubEnterprise := &domain.IDPGithubEnterprise{}
	var err error

	idpGithubEnterprise.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpGithubEnterprise.Type != nil {
		idpType = *idpGithubEnterprise.Type
	}

	if idpType != domain.IDPTypeGitHubEnterprise {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeGitHubEnterprise, idpType)
	}

	err = json.Unmarshal(idpGithubEnterprise.Payload, idpGithubEnterprise)
	if err != nil {
		return nil, err
	}

	return idpGithubEnterprise, nil
}

func (i *idProvider) GetGitlab(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGitlab, error) {
	idpGitlab := &domain.IDPGitlab{}
	var err error

	idpGitlab.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpGitlab.Type != nil {
		idpType = *idpGitlab.Type
	}

	if idpType != domain.IDPTypeGitLab {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeGitLab, idpType)
	}

	err = json.Unmarshal(idpGitlab.Payload, idpGitlab)
	if err != nil {
		return nil, err
	}

	return idpGitlab, nil
}

func (i *idProvider) GetGitlabSelfHosting(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPGitlabSelfHosting, error) {
	idpGitlabSelfHosting := &domain.IDPGitlabSelfHosting{}
	var err error

	idpGitlabSelfHosting.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if idpGitlabSelfHosting.Type != nil {
		idpType = *idpGitlabSelfHosting.Type
	}

	if idpType != domain.IDPTypeGitLabSelfHosted {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeGitLabSelfHosted, idpType)
	}

	err = json.Unmarshal(idpGitlabSelfHosting.Payload, idpGitlabSelfHosting)
	if err != nil {
		return nil, err
	}

	return idpGitlabSelfHosting, nil
}

func (i *idProvider) GetLDAP(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPLDAP, error) {
	ldap := &domain.IDPLDAP{}
	var err error

	ldap.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if ldap.Type != nil {
		idpType = *ldap.Type
	}

	if idpType != domain.IDPTypeLDAP {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeLDAP, idpType)
	}

	err = json.Unmarshal(ldap.Payload, ldap)
	if err != nil {
		return nil, err
	}

	return ldap, nil
}

func (i *idProvider) GetApple(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPApple, error) {
	apple := &domain.IDPApple{}
	var err error

	apple.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if apple.Type != nil {
		idpType = *apple.Type
	}

	if idpType != domain.IDPTypeApple {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeApple, idpType)
	}

	err = json.Unmarshal(apple.Payload, apple)
	if err != nil {
		return nil, err
	}

	return apple, nil
}

func (i *idProvider) GetSAML(ctx context.Context, client database.QueryExecutor, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IDPSAML, error) {
	saml := &domain.IDPSAML{}
	var err error

	saml.IdentityProvider, err = i.Get(ctx, client, id, instanceID, orgID)
	if err != nil {
		return nil, err
	}

	var idpType domain.IDPType
	if saml.Type != nil {
		idpType = *saml.Type
	}

	if idpType != domain.IDPTypeSAML {
		return nil, domain.NewIDPWrongTypeError(domain.IDPTypeSAML, idpType)
	}

	err = json.Unmarshal(saml.Payload, saml)
	if err != nil {
		return nil, err
	}

	return saml, nil
}

// -------------------------------------------------------------
// columns
// -------------------------------------------------------------

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

func (idProvider) AutoRegisterColumn() database.Column {
	return database.NewColumn("identity_providers", "auto_register")
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

func (idProvider) AllowAutoLinkingColumn() database.Column {
	return database.NewColumn("identity_providers", "auto_linking_field")
}

func (idProvider) StylingTypeColumn() database.Column {
	return database.NewColumn("identity_providers", "styling_type")
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

func (i idProvider) InstanceIDCondition(id string) database.Condition {
	return database.NewTextCondition(i.InstanceIDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) OrgIDCondition(id *string) database.Condition {
	if id == nil {
		return database.IsNull(i.OrgIDColumn())
	}
	return database.NewTextCondition(i.OrgIDColumn(), database.TextOperationEqual, *id)
}

func (i idProvider) IDCondition(id string) domain.IDPIdentifierCondition {
	return database.NewTextCondition(i.IDColumn(), database.TextOperationEqual, id)
}

func (i idProvider) StateCondition(state domain.IDPState) database.Condition {
	return database.NewTextCondition(i.StateColumn(), database.TextOperationEqual, state.String())
}

func (i idProvider) NameCondition(name string) domain.IDPIdentifierCondition {
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
	return database.NewTextCondition(i.AllowAutoLinkingColumn(), database.TextOperationEqual, linkingType.String())
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
