package resources

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	scim_schemas "github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/query"
)

type UsersHandler struct {
	command     *command.Commands
	query       *query.Queries
	userCodeAlg crypto.EncryptionAlgorithm
	config      *scim_config.Config
}

type ScimUser struct {
	*Resource
	ID                string                        `json:"id"`
	ExternalID        string                        `json:"externalId,omitempty"`
	UserName          string                        `json:"userName,omitempty"`
	Name              *ScimUserName                 `json:"name,omitempty"`
	DisplayName       string                        `json:"displayName,omitempty"`
	NickName          string                        `json:"nickName,omitempty"`
	ProfileUrl        *scim_schemas.HttpURL         `json:"profileUrl,omitempty"`
	Title             string                        `json:"title,omitempty"`
	PreferredLanguage language.Tag                  `json:"preferredLanguage,omitempty"`
	Locale            string                        `json:"locale,omitempty"`
	Timezone          string                        `json:"timezone,omitempty"`
	Active            *bool                         `json:"active,omitempty"`
	Emails            []*ScimEmail                  `json:"emails,omitempty"`
	PhoneNumbers      []*ScimPhoneNumber            `json:"phoneNumbers,omitempty"`
	Password          *scim_schemas.WriteOnlyString `json:"password,omitempty"`
	Ims               []*ScimIms                    `json:"ims,omitempty"`
	Addresses         []*ScimAddress                `json:"addresses,omitempty"`
	Photos            []*ScimPhoto                  `json:"photos,omitempty"`
	Entitlements      []*ScimEntitlement            `json:"entitlements,omitempty"`
	Roles             []*ScimRole                   `json:"roles,omitempty"`
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
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type ScimPhoneNumber struct {
	Value   string `json:"value"`
	Primary bool   `json:"primary"`
}

type ScimUserName struct {
	Formatted       string `json:"formatted,omitempty"`
	FamilyName      string `json:"familyName,omitempty"`
	GivenName       string `json:"givenName,omitempty"`
	MiddleName      string `json:"middleName,omitempty"`
	HonorificPrefix string `json:"honorificPrefix,omitempty"`
	HonorificSuffix string `json:"honorificSuffix,omitempty"`
}

func NewUsersHandler(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	config *scim_config.Config) ResourceHandler[*ScimUser] {
	return &UsersHandler{command, query, userCodeAlg, config}
}

func (h *UsersHandler) ResourceNameSingular() scim_schemas.ScimResourceTypeSingular {
	return scim_schemas.UserResourceType
}

func (h *UsersHandler) ResourceNamePlural() scim_schemas.ScimResourceTypePlural {
	return scim_schemas.UsersResourceType
}

func (u *ScimUser) GetResource() *Resource {
	return u.Resource
}

func (h *UsersHandler) NewResource() *ScimUser {
	return new(ScimUser)
}

func (h *UsersHandler) SchemaType() scim_schemas.ScimSchemaType {
	return scim_schemas.IdUser
}

func (h *UsersHandler) Create(ctx context.Context, user *ScimUser) (*ScimUser, error) {
	orgID := authz.GetCtxData(ctx).OrgID
	addHuman, err := h.mapToAddHuman(ctx, user)
	if err != nil {
		return nil, err
	}

	err = h.command.AddUserHuman(ctx, orgID, addHuman, true, h.userCodeAlg)
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

func (h *UsersHandler) Delete(ctx context.Context, id string) error {
	memberships, grants, err := h.queryUserDependencies(ctx, id)
	if err != nil {
		return err
	}

	_, err = h.command.RemoveUserV2(ctx, id, memberships, grants...)
	return err
}

func (h *UsersHandler) Get(ctx context.Context, id string) (*ScimUser, error) {
	user, err := h.query.GetUserByID(ctx, false, id)
	if err != nil {
		return nil, err
	}

	metadata, err := h.queryMetadataForUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return h.mapToScimUser(ctx, user, metadata), nil
}

func (h *UsersHandler) queryUserDependencies(ctx context.Context, userID string) ([]*command.CascadingMembership, []string, error) {
	userGrantUserQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}

	grants, err := h.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: []query.SearchQuery{userGrantUserQuery},
	}, true)
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
