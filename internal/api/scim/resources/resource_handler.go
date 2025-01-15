package resources

import (
	"context"
	"path"
	"strconv"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/domain"
)

type ResourceHandler[T ResourceHolder] interface {
	ResourceNameSingular() schemas.ScimResourceTypeSingular
	ResourceNamePlural() schemas.ScimResourceTypePlural
	SchemaType() schemas.ScimSchemaType
	NewResource() T

	Create(ctx context.Context, resource T) (T, error)
	Replace(ctx context.Context, id string, resource T) (T, error)
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (T, error)
}

type Resource struct {
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
	GetResource() *Resource
}

func buildResource[T ResourceHolder](ctx context.Context, handler ResourceHandler[T], details *domain.ObjectDetails) *Resource {
	created := details.CreationDate.UTC()
	if created.IsZero() {
		created = details.EventDate.UTC()
	}

	return &Resource{
		Schemas: []schemas.ScimSchemaType{handler.SchemaType()},
		Meta: &ResourceMeta{
			ResourceType: handler.ResourceNameSingular(),
			Created:      created,
			LastModified: details.EventDate.UTC(),
			Version:      strconv.FormatUint(details.Sequence, 10),
			Location:     buildLocation(ctx, handler, details.ID),
		},
	}
}

func buildLocation[T ResourceHolder](ctx context.Context, handler ResourceHandler[T], id string) string {
	return http.DomainContext(ctx).Origin() + path.Join(schemas.HandlerPrefix, authz.GetCtxData(ctx).OrgID, string(handler.ResourceNamePlural()), id)
}
