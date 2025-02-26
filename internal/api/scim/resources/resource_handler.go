package resources

import (
	"context"
	"strconv"

	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/domain"
)

type ResourceHandler[T ResourceHolder] interface {
	Schema() *schemas.ResourceSchema
	NewResource() T

	Create(ctx context.Context, resource T) (T, error)
	Replace(ctx context.Context, id string, resource T) (T, error)
	Update(ctx context.Context, id string, operations patch.OperationCollection) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (T, error)
	List(ctx context.Context, request *ListRequest) (*ListResponse[T], error)
}

type ResourceHolder interface {
	SchemasHolder
	GetResource() *schemas.Resource
}

type SchemasHolder interface {
	GetSchemas() []schemas.ScimSchemaType
}

func buildResource[T ResourceHolder](ctx context.Context, handler ResourceHandler[T], details *domain.ObjectDetails) *schemas.Resource {
	created := details.CreationDate.UTC()
	if created.IsZero() {
		created = details.EventDate.UTC()
	}

	schema := handler.Schema()
	return &schemas.Resource{
		ID:      details.ID,
		Schemas: []schemas.ScimSchemaType{schema.ID},
		Meta: &schemas.ResourceMeta{
			ResourceType: schema.Name,
			Created:      &created,
			LastModified: gu.Ptr(details.EventDate.UTC()),
			Version:      strconv.FormatUint(details.Sequence, 10),
			Location:     schemas.BuildLocationForResource(ctx, schema.PluralName, details.ID),
		},
	}
}
