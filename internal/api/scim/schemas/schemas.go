package schemas

import (
	"context"
	"path"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
)

type ScimSchemaType string
type ScimResourceTypeSingular string
type ScimResourceTypePlural string

const (
	idPrefixMessages        = "urn:ietf:params:scim:api:messages:2.0:"
	idPrefixCore            = "urn:ietf:params:scim:schemas:core:2.0:"
	idPrefixZitadelMessages = "urn:ietf:params:scim:api:zitadel:messages:2.0:"

	IdUser                  ScimSchemaType = idPrefixCore + "User"
	IdServiceProviderConfig ScimSchemaType = idPrefixCore + "ServiceProviderConfig"
	IdResourceType          ScimSchemaType = idPrefixCore + "ResourceType"
	IdSchema                ScimSchemaType = idPrefixCore + "Schema"
	IdListResponse          ScimSchemaType = idPrefixMessages + "ListResponse"
	IdPatchOperation        ScimSchemaType = idPrefixMessages + "PatchOp"
	IdSearchRequest         ScimSchemaType = idPrefixMessages + "SearchRequest"
	IdBulkRequest           ScimSchemaType = idPrefixMessages + "BulkRequest"
	IdBulkResponse          ScimSchemaType = idPrefixMessages + "BulkResponse"
	IdError                 ScimSchemaType = idPrefixMessages + "Error"
	IdZitadelErrorDetail    ScimSchemaType = idPrefixZitadelMessages + "ErrorDetail"

	UserResourceType  ScimResourceTypeSingular = "User"
	UsersResourceType ScimResourceTypePlural   = "Users"

	ServiceProviderConfigResourceType  ScimResourceTypeSingular = "ServiceProviderConfig"
	ServiceProviderConfigsResourceType ScimResourceTypePlural   = "ServiceProviderConfig"

	SchemaResourceType  ScimResourceTypeSingular = "Schema"
	SchemasResourceType ScimResourceTypePlural   = "Schemas"

	ResourceTypesResourceType ScimResourceTypePlural = "ResourceTypes"

	HandlerPrefix = "/scim/v2"
)

type Resource struct {
	ID      string           `json:"-"`
	Schemas []ScimSchemaType `json:"schemas"`
	Meta    *ResourceMeta    `json:"meta"`
}

type ResourceMeta struct {
	ResourceType ScimResourceTypeSingular `json:"resourceType"`
	Created      *time.Time               `json:"created,omitempty"`
	LastModified *time.Time               `json:"lastModified,omitempty"`
	Version      string                   `json:"version,omitempty"`
	Location     string                   `json:"location,omitempty"`
}

type ResourceType struct {
	*Resource
	ID          ScimResourceTypeSingular `json:"id"`
	Name        ScimResourceTypeSingular `json:"name"`
	Endpoint    ScimResourceTypePlural   `json:"endpoint"`
	Schema      ScimSchemaType           `json:"schema"`
	Description string                   `json:"description"`
}

type ResourceSchema struct {
	*Resource
	ID          ScimSchemaType           `json:"id"`
	Name        ScimResourceTypeSingular `json:"name"`
	PluralName  ScimResourceTypePlural   `json:"-"`
	Description string                   `json:"description,omitempty"`
	Attributes  []*SchemaAttribute       `json:"attributes"`
}

type SchemaAttribute struct {
	Name          string                    `json:"name"`
	Description   string                    `json:"description"`
	Type          SchemaAttributeType       `json:"type"`
	SubAttributes []*SchemaAttribute        `json:"subAttributes,omitempty"`
	MultiValued   bool                      `json:"multiValued"`
	Required      bool                      `json:"required"`
	CaseExact     bool                      `json:"caseExact"`
	Mutability    SchemaAttributeMutability `json:"mutability"`
	Returned      SchemaAttributeReturned   `json:"returned"`
	Uniqueness    SchemaAttributeUniqueness `json:"uniqueness"`
}

type SchemaAttributeType string

const (
	SchemaAttributeTypeString   SchemaAttributeType = "string"
	SchemaAttributeTypeBoolean  SchemaAttributeType = "boolean"
	SchemaAttributeTypeDecimal  SchemaAttributeType = "decimal"
	SchemaAttributeTypeInteger  SchemaAttributeType = "integer"
	SchemaAttributeTypeDateTime SchemaAttributeType = "dateTime"
	SchemaAttributeTypeComplex  SchemaAttributeType = "complex"
)

type SchemaAttributeMutability string

const (
	SchemaAttributeMutabilityReadWrite SchemaAttributeMutability = "readWrite"
	SchemaAttributeMutabilityWriteOnly SchemaAttributeMutability = "writeOnly"
)

type SchemaAttributeReturned string

const (
	SchemaAttributeReturnedAlways SchemaAttributeReturned = "always"
	SchemaAttributeReturnedNever  SchemaAttributeReturned = "never"
)

type SchemaAttributeUniqueness string

const (
	SchemaAttributeUniquenessNone   SchemaAttributeUniqueness = "none"
	SchemaAttributeUniquenessServer SchemaAttributeUniqueness = "server"
)

func (s *ResourceType) GetSchemas() []ScimSchemaType {
	return s.Resource.Schemas
}

func (s *ResourceType) GetResource() *Resource {
	return s.Resource
}

func (s *ResourceSchema) GetSchemas() []ScimSchemaType {
	return s.Resource.Schemas
}

func (s *ResourceSchema) GetResource() *Resource {
	return s.Resource
}

func (s *ResourceSchema) ToResourceType(ctx context.Context, orgID string) *ResourceType {
	return &ResourceType{
		Resource: &Resource{
			Schemas: []ScimSchemaType{IdResourceType},
			ID:      string(s.Name),
			Meta: &ResourceMeta{
				ResourceType: s.Name,
				Location:     BuildLocationWithOrg(ctx, orgID, ResourceTypesResourceType, string(s.Name)),
			},
		},
		ID:          s.Name,
		Name:        s.Name,
		Endpoint:    s.PluralName,
		Schema:      s.ID,
		Description: s.Description,
	}
}

func BuildLocationForResource(ctx context.Context, resourceName ScimResourceTypePlural, id string) string {
	return BuildLocationWithOrg(ctx, authz.GetCtxData(ctx).OrgID, resourceName, id)
}

func BuildLocationWithOrg(ctx context.Context, orgID string, resourceName ScimResourceTypePlural, id string) string {
	return http.DomainContext(ctx).Origin() + path.Join(HandlerPrefix, orgID, string(resourceName), id)
}
