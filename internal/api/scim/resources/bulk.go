package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/zitadel/logging"

	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	"github.com/zitadel/zitadel/internal/api/scim/metadata"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type BulkHandler struct {
	cfg                          *scim_config.BulkConfig
	handlersByPluralResourceName map[schemas.ScimResourceTypePlural]RawResourceHandlerAdapter
	translator                   *i18n.Translator
}

type BulkRequest struct {
	Schemas      []schemas.ScimSchemaType `json:"schemas"`
	FailOnErrors *int                     `json:"failOnErrors"`
	Operations   []*BulkRequestOperation  `json:"Operations"`
}

type BulkRequestOperation struct {
	Method string          `json:"method"`
	BulkID string          `json:"bulkId"`
	Path   string          `json:"path"`
	Data   json.RawMessage `json:"data"`
}

type BulkResponse struct {
	Schemas    []schemas.ScimSchemaType `json:"schemas"`
	Operations []*BulkResponseOperation `json:"Operations"`
}

type BulkResponseOperation struct {
	Method   string             `json:"method"`
	BulkID   string             `json:"bulkId,omitempty"`
	Location string             `json:"location,omitempty"`
	Error    *serrors.ScimError `json:"response,omitempty"`
	Status   string             `json:"status"`
}

func (r *BulkRequest) GetSchemas() []schemas.ScimSchemaType {
	return r.Schemas
}

func NewBulkHandler(
	cfg scim_config.BulkConfig,
	translator *i18n.Translator,
	handlers ...RawResourceHandlerAdapter,
) *BulkHandler {
	handlersByPluralResourceName := make(map[schemas.ScimResourceTypePlural]RawResourceHandlerAdapter, len(handlers))
	for _, handler := range handlers {
		handlersByPluralResourceName[handler.Schema().PluralName] = handler
	}

	return &BulkHandler{
		&cfg,
		handlersByPluralResourceName,
		translator,
	}
}

func (h *BulkHandler) BulkFromHttp(r *http.Request) (*BulkResponse, error) {
	req, err := h.readBulkRequest(r)
	if err != nil {
		return nil, err
	}

	return h.processRequest(r.Context(), req)
}

func (h *BulkHandler) readBulkRequest(r *http.Request) (*BulkRequest, error) {
	request := new(BulkRequest)
	if err := readSchema(r.Body, request, schemas.IdBulkRequest); err != nil {
		return nil, err
	}

	if len(request.Operations) > h.cfg.MaxOperationsCount {
		return nil, serrors.ThrowPayloadTooLarge(zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK19", "Too many bulk operations in one request, max %d allowed.", h.cfg.MaxOperationsCount))
	}
	return request, nil
}

func (h *BulkHandler) processRequest(ctx context.Context, req *BulkRequest) (*BulkResponse, error) {
	errorBudget := math.MaxInt32
	if req.FailOnErrors != nil {
		errorBudget = *req.FailOnErrors
	}

	resp := &BulkResponse{
		Schemas:    []schemas.ScimSchemaType{schemas.IdBulkResponse},
		Operations: make([]*BulkResponseOperation, 0, len(req.Operations)),
	}

	for _, operation := range req.Operations {
		opResp := h.processOperation(ctx, operation)
		resp.Operations = append(resp.Operations, opResp)

		if opResp.Error == nil {
			continue
		}

		errorBudget--
		if errorBudget <= 0 {
			return resp, nil
		}
	}

	return resp, nil
}

func (h *BulkHandler) processOperation(ctx context.Context, op *BulkRequestOperation) (opResp *BulkResponseOperation) {
	var statusCode int
	var resourceNamePlural schemas.ScimResourceTypePlural
	var resourceID string
	var err error
	opResp = &BulkResponseOperation{
		Method: op.Method,
		BulkID: op.BulkID,
	}

	defer func() {
		if r := recover(); r != nil {
			logging.WithFields("panic", r).Error("Bulk operation panic")
			err = zerrors.ThrowInternal(nil, "SCIM-BLK12", "Internal error while processing bulk operation")
		}

		if resourceNamePlural != "" && resourceID != "" {
			opResp.Location = schemas.BuildLocationForResource(ctx, resourceNamePlural, resourceID)
		}

		opResp.Status = strconv.Itoa(statusCode)
		if err != nil {
			opResp.Error = serrors.MapToScimError(ctx, h.translator, err)
			opResp.Status = opResp.Error.Status
		}
	}()

	resourceNamePlural, resourceID, err = h.parsePath(op.Path)
	if err != nil {
		return opResp
	}

	resourceID, err = metadata.ResolveScimBulkIDIfNeeded(ctx, resourceID)
	if err != nil {
		return opResp
	}

	resourceHandler, ok := h.handlersByPluralResourceName[resourceNamePlural]
	if !ok {
		err = zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK13", "Unknown resource %s", resourceNamePlural)
		return opResp
	}

	switch op.Method {
	case http.MethodPatch:
		statusCode = http.StatusNoContent
		err = h.processPatchOperation(ctx, resourceHandler, resourceID, op.Data)
	case http.MethodPut:
		statusCode = http.StatusOK
		err = h.processPutOperation(ctx, resourceHandler, resourceID, op)
	case http.MethodPost:
		statusCode = http.StatusCreated
		resourceID, err = h.processPostOperation(ctx, resourceHandler, resourceID, op)
	case http.MethodDelete:
		statusCode = http.StatusNoContent
		err = h.processDeleteOperation(ctx, resourceHandler, resourceID)
	default:
		err = zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK14", "Unsupported operation %s", op.Method)
	}

	return opResp
}

func (h *BulkHandler) processPutOperation(ctx context.Context, resourceHandler RawResourceHandlerAdapter, resourceID string, op *BulkRequestOperation) error {
	data := io.NopCloser(bytes.NewReader(op.Data))
	_, err := resourceHandler.Replace(ctx, resourceID, data)
	return err
}

func (h *BulkHandler) processPostOperation(ctx context.Context, resourceHandler RawResourceHandlerAdapter, resourceID string, op *BulkRequestOperation) (string, error) {
	if resourceID != "" {
		return "", zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK56", "Cannot post with a resourceID")
	}

	data := io.NopCloser(bytes.NewReader(op.Data))
	createdResource, err := resourceHandler.Create(ctx, data)
	if err != nil {
		return "", err
	}

	id := createdResource.GetResource().ID
	if op.BulkID != "" {
		metadata.SetScimBulkIDMapping(ctx, op.BulkID, id)
	}
	return id, nil
}

func (h *BulkHandler) processPatchOperation(ctx context.Context, resourceHandler RawResourceHandlerAdapter, resourceID string, data json.RawMessage) error {
	if resourceID == "" {
		return zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK16", "To patch a resource, a resourceID is required")
	}

	return resourceHandler.Update(ctx, resourceID, io.NopCloser(bytes.NewReader(data)))
}

func (h *BulkHandler) processDeleteOperation(ctx context.Context, resourceHandler RawResourceHandlerAdapter, resourceID string) error {
	if resourceID == "" {
		return zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK17", "To delete a resource, a resourceID is required")
	}

	return resourceHandler.Delete(ctx, resourceID)
}

func (h *BulkHandler) parsePath(path string) (schemas.ScimResourceTypePlural, string, error) {
	if !strings.HasPrefix(path, "/") {
		return "", "", zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK19", "Invalid path: has to start with a /")
	}

	// part 0 is always an empty string due to the leading /
	pathParts := strings.Split(path, "/")
	switch len(pathParts) {
	case 2:
		return schemas.ScimResourceTypePlural(pathParts[1]), "", nil
	case 3:
		return schemas.ScimResourceTypePlural(pathParts[1]), pathParts[2], nil
	default:
		return "", "", zerrors.ThrowInvalidArgumentf(nil, "SCIM-BLK20", "Invalid resource path")
	}
}
