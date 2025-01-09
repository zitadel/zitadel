package resources

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	schemas2 "github.com/zitadel/zitadel/internal/api/scim/schemas"
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
	ID                string                    `json:"id"`
	ExternalID        string                    `json:"externalId,omitempty"`
	UserName          string                    `json:"userName,omitempty"`
	Name              *ScimUserName             `json:"name,omitempty"`
	DisplayName       string                    `json:"displayName,omitempty"`
	NickName          string                    `json:"nickName,omitempty"`
	ProfileUrl        *schemas2.HttpURL         `json:"profileUrl,omitempty"`
	Title             string                    `json:"title,omitempty"`
	PreferredLanguage language.Tag              `json:"preferredLanguage,omitempty"`
	Locale            string                    `json:"locale,omitempty"`
	Timezone          string                    `json:"timezone,omitempty"`
	Active            bool                      `json:"active,omitempty"`
	Emails            []*ScimEmail              `json:"emails,omitempty"`
	PhoneNumbers      []*ScimPhoneNumber        `json:"phoneNumbers,omitempty"`
	Password          *schemas2.WriteOnlyString `json:"password,omitempty"`
	Ims               []*ScimIms                `json:"ims,omitempty"`
	Addresses         []*ScimAddress            `json:"addresses,omitempty"`
	Photos            []*ScimPhoto              `json:"photos,omitempty"`
	Entitlements      []*ScimEntitlement        `json:"entitlements,omitempty"`
	Roles             []*ScimRole               `json:"roles,omitempty"`
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
	Value   schemas2.HttpURL `json:"value"`
	Display string           `json:"display,omitempty"`
	Type    string           `json:"type"`
	Primary bool             `json:"primary,omitempty"`
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

func (h *UsersHandler) ResourceNameSingular() schemas2.ScimResourceTypeSingular {
	return schemas2.UserResourceType
}

func (h *UsersHandler) ResourceNamePlural() schemas2.ScimResourceTypePlural {
	return schemas2.UsersResourceType
}

func (u *ScimUser) GetResource() *Resource {
	return u.Resource
}

func (h *UsersHandler) NewResource() *ScimUser {
	return new(ScimUser)
}

func (h *UsersHandler) SchemaType() schemas2.ScimSchemaType {
	return schemas2.IdUser
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

	user.ID = addHuman.Details.ID
	user.Resource = buildResource(ctx, h, addHuman.Details)
	return user, err
}
