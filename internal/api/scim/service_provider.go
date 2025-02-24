package scim

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	scim_config "github.com/zitadel/zitadel/internal/api/scim/config"
	sresources "github.com/zitadel/zitadel/internal/api/scim/resources"
	sschemas "github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type serviceProviderHandler struct {
	config                *scim_config.Config
	schemas               []*sschemas.ResourceSchema
	schemasByID           map[sschemas.ScimSchemaType]*sschemas.ResourceSchema
	schemasByResourceName map[sschemas.ScimResourceTypeSingular]*sschemas.ResourceSchema
}

type serviceProviderConfig struct {
	*sschemas.Resource
	DocumentationUri      string                                                   `json:"documentationUri"`
	Patch                 serviceProviderConfigSupported                           `json:"patch"`
	Bulk                  serviceProviderConfigBulk                                `json:"bulk"`
	Filter                serviceProviderFilterSupported                           `json:"filter"`
	ChangePassword        serviceProviderConfigSupported                           `json:"changePassword"`
	Sort                  serviceProviderConfigSupported                           `json:"sort"`
	ETag                  serviceProviderConfigSupported                           `json:"etag"`
	AuthenticationSchemes []*scim_config.ServiceProviderConfigAuthenticationScheme `json:"authenticationSchemes,omitempty"`
}

type serviceProviderConfigSupported struct {
	Supported bool `json:"supported"`
}

type serviceProviderFilterSupported struct {
	Supported  bool `json:"supported"`
	MaxResults int  `json:"maxResults"`
}

type serviceProviderConfigBulk struct {
	Supported      bool  `json:"supported"`
	MaxOperations  int   `json:"maxOperations"`
	MaxPayloadSize int64 `json:"maxPayloadSize"`
}

var (
	defaultConfigSearchRequest = query.SearchRequest{
		Offset: 0,
		Limit:  100,
	}
)

func newServiceProviderHandler(cfg *scim_config.Config, handlers ...sresources.RawResourceHandlerAdapter) *serviceProviderHandler {
	schemas := make([]*sschemas.ResourceSchema, len(handlers))
	schemasByID := make(map[sschemas.ScimSchemaType]*sschemas.ResourceSchema, len(handlers))
	schemasByResourceName := make(map[sschemas.ScimResourceTypeSingular]*sschemas.ResourceSchema, len(handlers))
	for i, handler := range handlers {
		schema := handler.Schema()
		schemas[i] = schema
		schemasByID[schema.ID] = schema
		schemasByResourceName[schema.Name] = schema
	}

	return &serviceProviderHandler{
		config:                cfg,
		schemas:               schemas,
		schemasByID:           schemasByID,
		schemasByResourceName: schemasByResourceName,
	}
}

func (h *serviceProviderHandler) GetConfig(r *http.Request) (*serviceProviderConfig, error) {
	// the request is unauthenticated, read the orgID from the url instead of the context
	orgID := mux.Vars(r)[zhttp.OrgIdInPathVariableName]
	return &serviceProviderConfig{
		Resource: &sschemas.Resource{
			Schemas: []sschemas.ScimSchemaType{sschemas.IdServiceProviderConfig},
			Meta: &sschemas.ResourceMeta{
				ResourceType: sschemas.ServiceProviderConfigResourceType,
				Location:     sschemas.BuildLocationWithOrg(r.Context(), orgID, sschemas.ServiceProviderConfigsResourceType, ""),
			},
		},
		DocumentationUri: h.config.DocumentationUrl,
		Patch: serviceProviderConfigSupported{
			Supported: true,
		},
		Bulk: serviceProviderConfigBulk{
			Supported:      true,
			MaxOperations:  h.config.Bulk.MaxOperationsCount,
			MaxPayloadSize: h.config.MaxRequestBodySize,
		},
		Filter: serviceProviderFilterSupported{
			Supported:  true,
			MaxResults: sresources.MaxListCount,
		},
		ChangePassword: serviceProviderConfigSupported{
			Supported: true,
		},
		Sort: serviceProviderConfigSupported{
			Supported: true,
		},
		ETag: serviceProviderConfigSupported{
			Supported: false,
		},
		AuthenticationSchemes: h.config.AuthenticationSchemes,
	}, nil
}

func (h *serviceProviderHandler) ListResourceTypes(r *http.Request) (*sresources.ListResponse[*sschemas.ResourceType], error) {
	// the request is unauthenticated, read the orgID from the url instead of the context
	ctx := r.Context()
	orgID := mux.Vars(r)[zhttp.OrgIdInPathVariableName]

	resourceTypes := make([]*sschemas.ResourceType, len(h.schemas))
	for i, schema := range h.schemas {
		resourceTypes[i] = schema.ToResourceType(ctx, orgID)
	}

	return sresources.NewListResponse(uint64(len(resourceTypes)), defaultConfigSearchRequest, resourceTypes), nil
}

func (h *serviceProviderHandler) GetResourceType(r *http.Request) (*sschemas.ResourceType, error) {
	// the request is unauthenticated, read the orgID from the url instead of the context
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars[zhttp.OrgIdInPathVariableName]
	name := sschemas.ScimResourceTypeSingular(vars["name"])

	schema, ok := h.schemasByResourceName[name]
	if !ok {
		return nil, zerrors.ThrowNotFoundf(nil, "SCIMSP-148z", "Scim resource type %s not found", name)
	}

	return schema.ToResourceType(ctx, orgID), nil
}

func (h *serviceProviderHandler) ListSchemas(r *http.Request) (*sresources.ListResponse[*sschemas.ResourceSchema], error) {
	// the request is unauthenticated, read the orgID from the url instead of the context
	ctx := r.Context()
	orgID := mux.Vars(r)[zhttp.OrgIdInPathVariableName]

	schemas := make([]*sschemas.ResourceSchema, len(h.schemas))
	for i, schema := range h.schemas {
		schemas[i] = buildSchema(ctx, orgID, schema)
	}

	return sresources.NewListResponse(uint64(len(h.schemas)), defaultConfigSearchRequest, schemas), nil
}

func (h *serviceProviderHandler) GetSchema(r *http.Request) (*sschemas.ResourceSchema, error) {
	// the request is unauthenticated, read the orgID from the url instead of the context
	ctx := r.Context()
	vars := mux.Vars(r)
	orgID := vars[zhttp.OrgIdInPathVariableName]
	id := sschemas.ScimSchemaType(vars["id"])

	schema, ok := h.schemasByID[id]
	if !ok {
		return nil, zerrors.ThrowNotFoundf(nil, "SCIMSP-148y", "Scim schema %s not found", id)
	}

	return buildSchema(ctx, orgID, schema), nil
}

// buildSchema shallow copies the provided schema and sets the correct location based on the provided context information.
func buildSchema(ctx context.Context, orgID string, schema *sschemas.ResourceSchema) *sschemas.ResourceSchema {
	newSchema := *schema
	newSchema.Resource = &sschemas.Resource{
		ID:      schema.Resource.ID,
		Schemas: schema.Resource.Schemas,
		Meta: &sschemas.ResourceMeta{
			ResourceType: schema.Resource.Meta.ResourceType,
			Location:     sschemas.BuildLocationWithOrg(ctx, orgID, sschemas.SchemasResourceType, string(schema.ID)),
		},
	}
	return &newSchema
}
