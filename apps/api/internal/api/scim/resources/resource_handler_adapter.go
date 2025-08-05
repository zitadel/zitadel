package resources

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"slices"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/scim/resources/patch"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// RawResourceHandlerAdapter adapts the ResourceHandler[T] without any generics
type RawResourceHandlerAdapter interface {
	Schema() *schemas.ResourceSchema

	Create(ctx context.Context, data io.ReadCloser) (ResourceHolder, error)
	Replace(ctx context.Context, resourceID string, data io.ReadCloser) (ResourceHolder, error)
	Update(ctx context.Context, resourceID string, data io.ReadCloser) error
	Delete(ctx context.Context, resourceID string) error
}

type ResourceHandlerAdapter[T ResourceHolder] struct {
	handler ResourceHandler[T]
}

func NewResourceHandlerAdapter[T ResourceHolder](handler ResourceHandler[T]) *ResourceHandlerAdapter[T] {
	return &ResourceHandlerAdapter[T]{
		handler,
	}
}

func (adapter *ResourceHandlerAdapter[T]) Schema() *schemas.ResourceSchema {
	return adapter.handler.Schema()
}

func (adapter *ResourceHandlerAdapter[T]) CreateFromHttp(r *http.Request) (ResourceHolder, error) {
	return adapter.Create(r.Context(), r.Body)
}

func (adapter *ResourceHandlerAdapter[T]) Create(ctx context.Context, data io.ReadCloser) (ResourceHolder, error) {
	entity, err := adapter.readEntity(data)
	if err != nil {
		return entity, err
	}

	return adapter.handler.Create(ctx, entity)
}

func (adapter *ResourceHandlerAdapter[T]) ReplaceFromHttp(r *http.Request) (ResourceHolder, error) {
	return adapter.Replace(r.Context(), mux.Vars(r)["id"], r.Body)
}

func (adapter *ResourceHandlerAdapter[T]) Replace(ctx context.Context, resourceID string, data io.ReadCloser) (ResourceHolder, error) {
	entity, err := adapter.readEntity(data)
	if err != nil {
		return entity, err
	}

	return adapter.handler.Replace(ctx, resourceID, entity)
}

func (adapter *ResourceHandlerAdapter[T]) UpdateFromHttp(r *http.Request) error {
	return adapter.Update(r.Context(), mux.Vars(r)["id"], r.Body)
}

func (adapter *ResourceHandlerAdapter[T]) Update(ctx context.Context, resourceID string, data io.ReadCloser) error {
	request := new(patch.OperationRequest)
	if err := readSchema(data, request, schemas.IdPatchOperation); err != nil {
		return err
	}

	if err := request.Validate(); err != nil {
		return err
	}

	if len(request.Operations) == 0 {
		return nil
	}

	return adapter.handler.Update(ctx, resourceID, request.Operations)
}

func (adapter *ResourceHandlerAdapter[T]) DeleteFromHttp(r *http.Request) error {
	return adapter.Delete(r.Context(), mux.Vars(r)["id"])
}

func (adapter *ResourceHandlerAdapter[T]) Delete(ctx context.Context, resourceID string) error {
	return adapter.handler.Delete(ctx, resourceID)
}

func (adapter *ResourceHandlerAdapter[T]) ListFromHttp(r *http.Request) (*ListResponse[T], error) {
	request, err := adapter.readListRequest(r)
	if err != nil {
		return nil, err
	}

	return adapter.handler.List(r.Context(), request)
}

func (adapter *ResourceHandlerAdapter[T]) GetFromHttp(r *http.Request) (T, error) {
	id := mux.Vars(r)["id"]
	return adapter.handler.Get(r.Context(), id)
}

func (adapter *ResourceHandlerAdapter[T]) readEntity(data io.ReadCloser) (T, error) {
	entity := adapter.handler.NewResource()
	return entity, readSchema(data, entity, adapter.handler.Schema().ID)
}

func readSchema(data io.ReadCloser, entity SchemasHolder, schema schemas.ScimSchemaType) error {
	defer func() {
		err := data.Close()
		logging.OnError(err).Warn("Failed to close http request body")
	}()

	err := json.NewDecoder(data).Decode(&entity)
	if err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			return serrors.ThrowPayloadTooLarge(zerrors.ThrowInvalidArgumentf(err, "SCIM-hmaxb1", "Request payload too large, max %d bytes allowed.", maxBytesErr.Limit))
		}

		if serrors.IsScimOrZitadelError(err) {
			return err
		}

		return serrors.ThrowInvalidSyntax(zerrors.ThrowInvalidArgumentf(err, "SCIM-ucrjson", "Could not deserialize json"))
	}

	providedSchemas := entity.GetSchemas()
	if !slices.Contains(providedSchemas, schema) {
		return serrors.ThrowInvalidSyntax(zerrors.ThrowInvalidArgumentf(nil, "SCIM-xxrschema", "Expected schema %v is not provided", schema))
	}

	return nil
}
