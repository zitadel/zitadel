package resources

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/gorilla/mux"

	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ResourceHandlerAdapter[T ResourceHolder] struct {
	handler ResourceHandler[T]
}

type ListRequest struct {
	// Count An integer indicating the desired maximum number of query results per page. OPTIONAL.
	Count uint64 `json:"count" schema:"count"`

	// StartIndex An integer indicating the 1-based index of the first query result. Optional.
	StartIndex uint64 `json:"startIndex" schema:"startIndex"`
}

type ListResponse[T any] struct {
	Schemas      []schemas.ScimSchemaType `json:"schemas"`
	ItemsPerPage uint64                   `json:"itemsPerPage"`
	TotalResults uint64                   `json:"totalResults"`
	StartIndex   uint64                   `json:"startIndex"`
	Resources    []T                      `json:"Resources"` // according to the rfc this is the only field in PascalCase...
}

func NewResourceHandlerAdapter[T ResourceHolder](handler ResourceHandler[T]) *ResourceHandlerAdapter[T] {
	return &ResourceHandlerAdapter[T]{
		handler,
	}
}

func (adapter *ResourceHandlerAdapter[T]) Create(r *http.Request) (T, error) {
	entity, err := adapter.readEntityFromBody(r)
	if err != nil {
		return entity, err
	}

	return adapter.handler.Create(r.Context(), entity)
}

func (adapter *ResourceHandlerAdapter[T]) Replace(r *http.Request) (T, error) {
	entity, err := adapter.readEntityFromBody(r)
	if err != nil {
		return entity, err
	}

	id := mux.Vars(r)["id"]
	return adapter.handler.Replace(r.Context(), id, entity)
}

func (adapter *ResourceHandlerAdapter[T]) Delete(r *http.Request) error {
	id := mux.Vars(r)["id"]
	return adapter.handler.Delete(r.Context(), id)
}

func (adapter *ResourceHandlerAdapter[T]) Get(r *http.Request) (T, error) {
	id := mux.Vars(r)["id"]
	return adapter.handler.Get(r.Context(), id)
}

func (adapter *ResourceHandlerAdapter[T]) readEntityFromBody(r *http.Request) (T, error) {
	entity := adapter.handler.NewResource()
	err := json.NewDecoder(r.Body).Decode(entity)
	if err != nil {
		if zerrors.IsZitadelError(err) {
			return entity, err
		}

		return entity, serrors.ThrowInvalidSyntax(zerrors.ThrowInvalidArgumentf(nil, "SCIM-ucrjson", "Could not deserialize json: %v", err.Error()))
	}

	resource := entity.GetResource()
	if resource == nil {
		return entity, serrors.ThrowInvalidSyntax(zerrors.ThrowInvalidArgument(nil, "SCIM-xxrjson", "Could not get resource, is the schema correct?"))
	}

	if !slices.Contains(resource.Schemas, adapter.handler.SchemaType()) {
		return entity, serrors.ThrowInvalidSyntax(zerrors.ThrowInvalidArgumentf(nil, "SCIM-xxrschema", "Expected schema %v is not provided", adapter.handler.SchemaType()))
	}

	return entity, nil
}
