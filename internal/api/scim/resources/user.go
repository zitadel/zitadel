package resources

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	"github.com/zitadel/zitadel/internal/api/scim/resources/filter"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	scim_schemas "github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type UsersHandler struct {
	command         *command.Commands
	query           *query.Queries
	userCodeAlg     crypto.EncryptionAlgorithm
	config          *scim_config.Config
	filterEvaluator *filter.Evaluator
	schema          *scim_schemas.ResourceSchema
}

type ScimUser struct {
	*scim_schemas.Resource `scim:"ignoreInSchema"`
	ID                     string                        `json:"id" scim:"ignoreInSchema"`
	ExternalID             string                        `json:"externalId,omitempty"`
	UserName               string                        `json:"userName,omitempty" scim:"required,unique,caseInsensitive"`
	Name                   *ScimUserName                 `json:"name,omitempty" scim:"required"`
	DisplayName            string                        `json:"displayName,omitempty"`
	NickName               string                        `json:"nickName,omitempty"`
	ProfileUrl             *scim_schemas.HttpURL         `json:"profileUrl,omitempty"`
	Title                  string                        `json:"title,omitempty"`
	PreferredLanguage      language.Tag                  `json:"preferredLanguage,omitempty"`
	Locale                 string                        `json:"locale,omitempty"`
	Timezone               string                        `json:"timezone,omitempty"`
	Active                 *scim_schemas.RelaxedBool     `json:"active,omitempty"`
	Emails                 []*ScimEmail                  `json:"emails,omitempty" scim:"required"`
	PhoneNumbers           []*ScimPhoneNumber            `json:"phoneNumbers,omitempty"`
	Password               *scim_schemas.WriteOnlyString `json:"password,omitempty"`
	Ims                    []*ScimIms                    `json:"ims,omitempty"`
	Addresses              []*ScimAddress                `json:"addresses,omitempty"`
	Photos                 []*ScimPhoto                  `json:"photos,omitempty"`
	Entitlements           []*ScimEntitlement            `json:"entitlements,omitempty"`
	Roles                  []*ScimRole                   `json:"roles,omitempty"`
}

type ScimEntitlement struct {
	Value   string `json:"value,omitempty"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type ScimRole struct {
	Value   string `json:"value,omitempty"`
	Display string `json:"display,omitempty"`
	Type    string `json:"type,omitempty"`
	Primary bool   `json:"primary,omitempty"`
}

type ScimPhoto struct {
	Value   scim_schemas.HttpURL `json:"value"`
	Display string               `json:"display,omitempty"`
	Type    string               `json:"type"`
	Primary bool                 `json:"primary,omitempty"`
}

type ScimAddress struct {
	Type          string `json:"type,omitempty"`
	StreetAddress string `json:"streetAddress,omitempty"`
	Locality      string `json:"locality,omitempty"`
	Region        string `json:"region,omitempty"`
	PostalCode    string `json:"postalCode,omitempty"`
	Country       string `json:"country,omitempty"`
	Formatted     string `json:"formatted,omitempty"`
	Primary       bool   `json:"primary,omitempty"`
}

type ScimIms struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ScimEmail struct {
	Value   string `json:"value" scim:"required"`
	Primary bool   `json:"primary"`
	Type    string `json:"type,omitempty"`
}

type ScimPhoneNumber struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type ScimUserName struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty" scim:"required"`
	GivenName       string `json:"givenName,omitempty" scim:"required"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

func NewUsersHandler(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	config *scim_config.Config) ResourceHandler[*ScimUser] {
	return &UsersHandler{
		command,
		query,
		userCodeAlg,
		config,
		filter.NewEvaluator(scim_schemas.IdUser),
		scim_schemas.BuildSchema(scim_schemas.SchemaBuilderArgs{
			ID:           scim_schemas.IdUser,
			Name:         scim_schemas.UserResourceType,
			EndpointName: scim_schemas.UsersResourceType,
			Description:  "User Account",
			Resource:     new(ScimUser),
		}),
	}
}

func (u *ScimUser) GetResource() *scim_schemas.Resource {
	return u.Resource
}

func (u *ScimUser) GetSchemas() []scim_schemas.ScimSchemaType {
	if u.Resource == nil {
		return nil
	}

	return u.Resource.Schemas
}

func (h *UsersHandler) Schema() *scim_schemas.ResourceSchema {
	return h.schema
}

func (h *UsersHandler) NewResource() *ScimUser {
	return new(ScimUser)
}

func (h *UsersHandler) Create(ctx context.Context, user *ScimUser) (*ScimUser, error) {
	orgID := authz.GetCtxData(ctx).OrgID
	addHuman, err := h.mapToAddHuman(ctx, user)
	if err != nil {
		return nil, err
	}

	err = h.command.AddUserHuman(ctx, orgID, addHuman, false, h.userCodeAlg)
	if err != nil {
		return nil, err
	}

	h.mapAddCommandToScimUser(ctx, user, addHuman)
	return user, nil
}

func (h *UsersHandler) Replace(ctx context.Context, id string, user *ScimUser) (*ScimUser, error) {
	user.ID = id
	changeHuman, err := h.mapToChangeHuman(ctx, user)
	if err != nil {
		return nil, err
	}

	err = h.command.ChangeUserHuman(ctx, changeHuman, h.userCodeAlg)
	if err != nil {
		return nil, err
	}

	h.mapChangeCommandToScimUser(ctx, user, changeHuman)
	return user, nil
}

func (h *UsersHandler) Update(ctx context.Context, id string, operations patch.OperationCollection) error {
	orgID := authz.GetCtxData(ctx).OrgID
	userWM, err := h.command.UserHumanWriteModel(ctx, id, orgID, true, true, true, true, false, false, true)
	if err != nil {
		return err
	}

	user := h.mapWriteModelToScimUser(ctx, userWM)
	changeHuman, err := h.applyPatchesToChangeHuman(ctx, user, operations)
	if err != nil {
		return err
	}

	// ensure the identity of the user is not modified
	changeHuman.ID = id
	changeHuman.ResourceOwner = orgID
	return h.command.ChangeUserHuman(ctx, changeHuman, h.userCodeAlg)
}

func (h *UsersHandler) Delete(ctx context.Context, id string) error {
	memberships, grants, err := h.queryUserDependencies(ctx, id)
	if err != nil {
		return err
	}
	_, err = h.command.RemoveUserV2(ctx, id, authz.GetCtxData(ctx).OrgID, memberships, grants...)
	return err
}

func (h *UsersHandler) Get(ctx context.Context, id string) (*ScimUser, error) {
	user, err := h.query.GetUserByIDWithResourceOwner(ctx, false, id, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}

	if user.Type != domain.UserTypeHuman {
		return nil, zerrors.ThrowNotFound(nil, "SCIM-USRT1", "Errors.Users.NotFound")
	}

	metadata, err := h.queryMetadataForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return h.mapToScimUser(ctx, user, metadata), nil
}

func (h *UsersHandler) List(ctx context.Context, request *ListRequest) (*ListResponse[*ScimUser], error) {
	q, err := h.buildListQuery(ctx, request)
	if err != nil {
		return nil, err
	}

	if request.Count == 0 {
		count, err := h.query.CountUsers(ctx, q)
		if err != nil {
			return nil, err
		}

		return NewListResponse(count, q.SearchRequest, make([]*ScimUser, 0)), nil
	}

	users, err := h.query.SearchUsers(ctx, q, nil)
	if err != nil {
		return nil, err
	}

	metadata, err := h.queryMetadataForUsers(ctx, usersToIDs(users.Users))
	if err != nil {
		return nil, err
	}

	scimUsers := h.mapToScimUsers(ctx, users.Users, metadata)
	return NewListResponse(users.SearchResponse.Count, q.SearchRequest, scimUsers), nil
}

func (h *UsersHandler) queryUserDependencies(ctx context.Context, userID string) ([]*command.CascadingMembership, []string, error) {
	userGrantUserQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}

	grants, err := h.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{userGrantUserQuery},
	}, true, nil)
	if err != nil {
		return nil, nil, err
	}

	membershipsUserQuery, err := query.NewMembershipUserIDQuery(userID)
	if err != nil {
		return nil, nil, err
	}

	memberships, err := h.query.Memberships(ctx, &query.MembershipSearchQuery{
		Queries: []query.SearchQuery{membershipsUserQuery},
	}, false)

	if err != nil {
		return nil, nil, err
	}
	return cascadingMemberships(memberships.Memberships), userGrantsToIDs(grants.UserGrants), nil
}
