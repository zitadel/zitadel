package repository

import (
	"context"
	"encoding/json"
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

const queryIDProviderStmt = `SELECT instance_id, org_id, id, state, name, type, auto_register, allow_creation, allow_auto_creation,` +
	` allow_auto_update, allow_linking, allow_auto_linking, styling_type, payload, created_at, updated_at` +
	` FROM zitadel.identity_providers`

func (i *idProvider) Get(ctx context.Context, id domain.IDPIdentifierCondition, instanceID string, orgID *string) (*domain.IdentityProvider, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(queryIDProviderStmt)

	conditions := []database.Condition{id, i.InstanceIDCondition(instanceID), i.OrgIDCondition(orgID)}

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

func (i *idProvider) Update(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string, changes ...database.Change) (int64, error) {
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

func (i *idProvider) Delete(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (int64, error) {
	builder := database.StatementBuilder{}

	builder.WriteString(`DELETE FROM zitadel.identity_providers`)

	conditions := []database.Condition{
		id,
		i.InstanceIDCondition(instnaceID),
		i.OrgIDCondition(orgID),
	}
	writeCondition(&builder, database.And(conditions...))

	return i.client.Exec(ctx, builder.String(), builder.Args()...)
}

func (i *idProvider) GetOIDC(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPOIDC, error) {
	idpOIDC := &domain.IDPOIDC{}
	var err error

	idpOIDC.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpOIDC.Type != domain.IDPTypeOIDC.String() {
		// TODO
		return nil, nil
	}

	err = json.Unmarshal([]byte(*idpOIDC.Payload), idpOIDC)
	if err != nil {
		return nil, err
	}

	return idpOIDC, nil
}

func (i *idProvider) GetJWT(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPJWT, error) {
	idpJWT := &domain.IDPJWT{}
	var err error

	idpJWT.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpJWT.Type != domain.IDPTypeJWT.String() {
		// TODO
		return nil, nil
	}

	err = json.Unmarshal([]byte(*idpJWT.Payload), idpJWT)
	if err != nil {
		return nil, err
	}

	return idpJWT, nil
}

func (i *idProvider) GetOAuth(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPOAuth, error) {
	idpOAuth := &domain.IDPOAuth{}
	var err error

	idpOAuth.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpOAuth.Type != domain.IDPTypeOAuth.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpOAuth.Payload), idpOAuth)
	if err != nil {
		return nil, err
	}

	return idpOAuth, nil
}

func (i *idProvider) GetOAzureAD(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPOAzureAD, error) {
	idpAzure := &domain.IDPOAzureAD{}
	var err error

	idpAzure.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpAzure.Type != domain.IDPTypeAzure.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpAzure.Payload), idpAzure)
	if err != nil {
		return nil, err
	}

	return idpAzure, nil
}

func (i *idProvider) GetGoogle(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPGoogle, error) {
	idpGoogle := &domain.IDPGoogle{}
	var err error

	idpGoogle.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpGoogle.Type != domain.IDPTypeGoogle.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpGoogle.Payload), idpGoogle)
	if err != nil {
		return nil, err
	}

	return idpGoogle, nil
}

func (i *idProvider) GetGithub(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPGithub, error) {
	idpGithub := &domain.IDPGithub{}
	var err error

	idpGithub.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpGithub.Type != domain.IDPTypeGitHub.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpGithub.Payload), idpGithub)
	if err != nil {
		return nil, err
	}

	return idpGithub, nil
}

func (i *idProvider) GetGithubEnterprise(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPGithubEnterprise, error) {
	idpGithubEnterprise := &domain.IDPGithubEnterprise{}
	var err error

	idpGithubEnterprise.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpGithubEnterprise.Type != domain.IDPTypeGitHubEnterprise.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpGithubEnterprise.Payload), idpGithubEnterprise)
	if err != nil {
		return nil, err
	}

	return idpGithubEnterprise, nil
}

func (i *idProvider) GetGitlab(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPGitlab, error) {
	idpGitlab := &domain.IDPGitlab{}
	var err error

	idpGitlab.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpGitlab.Type != domain.IDPTypeGitLab.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpGitlab.Payload), idpGitlab)
	if err != nil {
		return nil, err
	}

	return idpGitlab, nil
}

func (i *idProvider) GetGitlabSelfHosting(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPGitlabSelfHosting, error) {
	idpGitlabSelfHosting := &domain.IDPGitlabSelfHosting{}
	var err error

	idpGitlabSelfHosting.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if idpGitlabSelfHosting.Type != domain.IDPTypeGitLabSelfHosted.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*idpGitlabSelfHosting.Payload), idpGitlabSelfHosting)
	if err != nil {
		return nil, err
	}

	return idpGitlabSelfHosting, nil
}

func (i *idProvider) GetLDAP(ctx context.Context, id domain.IDPIdentifierCondition, instnaceID string, orgID *string) (*domain.IDPLDAP, error) {
	ldap := &domain.IDPLDAP{}
	var err error

	ldap.IdentityProvider, err = i.Get(ctx, id, instnaceID, orgID)
	if err != nil {
		return nil, err
	}

	if ldap.Type != domain.IDPTypeLDAP.String() {
		// TODO
		return nil, errors.New("WRONG TYPE")
	}

	err = json.Unmarshal([]byte(*ldap.Payload), ldap)
	if err != nil {
		return nil, err
	}

	return ldap, nil
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

func (idProvider) AutoRegisterColumn() database.Column {
	return database.NewColumn("auto_register")
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

func (idProvider) AllowAutoLinkingColumn() database.Column {
	return database.NewColumn("allow_auto_linking")
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

func (i idProvider) TypeCondition(typee domain.IDPType) database.Condition {
	return database.NewTextCondition(i.TypeColumn(), database.TextOperationEqual, typee.String())
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

func (i idProvider) AllowAutoLinkingCondition(linkingType domain.IDPAutoLinkingOption) database.Condition {
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
