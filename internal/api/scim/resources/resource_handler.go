package resources

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/domain"
)

type ResourceHandler[T ResourceHolder] interface {
	SchemaType() schemas.ScimSchemaType
	ResourceNameSingular() schemas.ScimResourceTypeSingular
	ResourceNamePlural() schemas.ScimResourceTypePlural
	NewResource() T

	Create(ctx context.Context, resource T) (T, error)
	Replace(ctx context.Context, id string, resource T) (T, error)
	Update(ctx context.Context, id string, operations patch.OperationCollection) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (T, error)
	List(ctx context.Context, request *ListRequest) (*ListResponse[T], error)
}

type Resource struct {
	ID      string                   `json:"-"`
	Schemas []schemas.ScimSchemaType `json:"schemas"`
	Meta    *ResourceMeta            `json:"meta"`
}

type ResourceMeta struct {
	ResourceType schemas.ScimResourceTypeSingular `json:"resourceType"`
	Created      time.Time                        `json:"created"`
	LastModified time.Time                        `json:"lastModified"`
	Version      string                           `json:"version"`
	Location     string                           `json:"location"`
}

type ResourceHolder interface {
	SchemasHolder
	GetResource() *Resource
}

type SchemasHolder interface {
	GetSchemas() []schemas.ScimSchemaType
}

func buildResource[T ResourceHolder](ctx context.Context, handler ResourceHandler[T], details *domain.ObjectDetails) *Resource {
	created := details.CreationDate.UTC()
	if created.IsZero() {
		created = details.EventDate.UTC()
	}

	return &Resource{
		ID:      details.ID,
		Schemas: []schemas.ScimSchemaType{handler.SchemaType()},
		Meta: &ResourceMeta{
			ResourceType: handler.ResourceNameSingular(),
			Created:      created,
			LastModified: details.EventDate.UTC(),
			Version:      strconv.FormatUint(details.Sequence, 10),
			Location:     buildLocation(ctx, handler.ResourceNamePlural(), details.ID),
		},
	}
}

func buildLocation(ctx context.Context, resourceName schemas.ScimResourceTypePlural, id string) string {
	return http.DomainContext(ctx).Origin() + path.Join(schemas.HandlerPrefix, authz.GetCtxData(ctx).OrgID, string(resourceName), id)
}
