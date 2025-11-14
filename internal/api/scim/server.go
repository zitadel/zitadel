package scim

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	zhttp "github.com/zitadel/zitadel/internal/api/http"
	zhttp_middlware "github.com/zitadel/zitadel/internal/api/http/middleware"
	sconfig "github.com/zitadel/zitadel/internal/api/scim/config"
	smiddleware "github.com/zitadel/zitadel/internal/api/scim/middleware"
	sresources "github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
	"github.com/zitadel/zitadel/internal/api/scim/serrors"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
)

func NewServer(
	command *command.Commands,
	query *query.Queries,
	verifier *authz.ApiTokenVerifier,
	userCodeAlg crypto.EncryptionAlgorithm,
	config *sconfig.Config,
	translator *i18n.Translator,
	middlewares ...zhttp_middlware.MiddlewareWithErrorFunc,
) http.Handler {
	verifier.RegisterServer("SCIM-V2", schemas.HandlerPrefix, AuthMapping)
	return buildHandler(command, query, userCodeAlg, config, translator, middlewares...)
}

func buildHandler(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	cfg *sconfig.Config,
	translator *i18n.Translator,
	middlewares ...zhttp_middlware.MiddlewareWithErrorFunc,
) http.Handler {

	router := mux.NewRouter()
	middleware := buildMiddleware(cfg, query, translator, middlewares)

	usersHandler := sresources.NewResourceHandlerAdapter(sresources.NewUsersHandler(command, query, userCodeAlg, cfg))
	mapResource(router, middleware, usersHandler)

	bulkHandler := sresources.NewBulkHandler(cfg.Bulk, translator, usersHandler)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/Bulk", middleware(handleJsonResponse(bulkHandler.BulkFromHttp))).Methods(http.MethodPost)

	serviceProviderHandler := newServiceProviderHandler(cfg, usersHandler)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/ServiceProviderConfig", middleware(handleJsonResponse(serviceProviderHandler.GetConfig))).Methods(http.MethodGet)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/ResourceTypes", middleware(handleJsonResponse(serviceProviderHandler.ListResourceTypes))).Methods(http.MethodGet)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/ResourceTypes/{name}", middleware(handleResourceResponse(serviceProviderHandler.GetResourceType))).Methods(http.MethodGet)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/Schemas", middleware(handleJsonResponse(serviceProviderHandler.ListSchemas))).Methods(http.MethodGet)
	router.Handle("/"+zhttp.OrgIdInPathVariable+"/Schemas/{id}", middleware(handleResourceResponse(serviceProviderHandler.GetSchema))).Methods(http.MethodGet)

	return router
}

func buildMiddleware(
	cfg *sconfig.Config,
	query *query.Queries,
	translator *i18n.Translator,
	middlewares []zhttp_middlware.MiddlewareWithErrorFunc,
) zhttp_middlware.ErrorHandlerFunc {
	// content type middleware needs to run at the very beginning to correctly set content types of errors
	middlewares = append([]zhttp_middlware.MiddlewareWithErrorFunc{smiddleware.ContentTypeMiddleware}, middlewares...)
	middlewares = append(middlewares, smiddleware.ScimContextMiddleware(query))
	scimMiddleware := zhttp_middlware.ChainedWithErrorHandler(serrors.ErrorHandler(translator), middlewares...)
	return func(handler zhttp_middlware.HandlerFuncWithError) http.Handler {
		return http.MaxBytesHandler(scimMiddleware(handler), cfg.MaxRequestBodySize)
	}
}

func mapResource[T sresources.ResourceHolder](router *mux.Router, mw zhttp_middlware.ErrorHandlerFunc, adapter *sresources.ResourceHandlerAdapter[T]) {
	resourceRouter := router.PathPrefix("/" + path.Join(zhttp.OrgIdInPathVariable, string(adapter.Schema().PluralName))).Subrouter()

	resourceRouter.Handle("", mw(handleResourceCreatedResponse(adapter.CreateFromHttp))).Methods(http.MethodPost)
	resourceRouter.Handle("", mw(handleJsonResponse(adapter.ListFromHttp))).Methods(http.MethodGet)
	resourceRouter.Handle("/.search", mw(handleJsonResponse(adapter.ListFromHttp))).Methods(http.MethodPost)
	resourceRouter.Handle("/{id}", mw(handleResourceResponse(adapter.GetFromHttp))).Methods(http.MethodGet)
	resourceRouter.Handle("/{id}", mw(handleResourceResponse(adapter.ReplaceFromHttp))).Methods(http.MethodPut)
	resourceRouter.Handle("/{id}", mw(handleEmptyResponse(adapter.UpdateFromHttp))).Methods(http.MethodPatch)
	resourceRouter.Handle("/{id}", mw(handleEmptyResponse(adapter.DeleteFromHttp))).Methods(http.MethodDelete)
}

func handleJsonResponse[T any](next func(r *http.Request) (T, error)) zhttp_middlware.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		entity, err := next(r)
		if err != nil {
			return err
		}

		err = json.NewEncoder(w).Encode(entity)
		logging.OnError(err).Warn("scim json response encoding failed")
		return nil
	}
}

func handleResourceCreatedResponse[T sresources.ResourceHolder](next func(*http.Request) (T, error)) zhttp_middlware.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		entity, err := next(r)
		if err != nil {
			return err
		}

		resource := entity.GetResource()
		w.Header().Set(zhttp.Location, resource.Meta.Location)
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(entity)
		logging.OnError(err).Warn("scim json response encoding failed")
		return nil
	}
}

func handleResourceResponse[T sresources.ResourceHolder](next func(*http.Request) (T, error)) zhttp_middlware.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		entity, err := next(r)
		if err != nil {
			return err
		}

		resource := entity.GetResource()
		w.Header().Set(zhttp.ContentLocation, resource.Meta.Location)

		err = json.NewEncoder(w).Encode(entity)
		logging.OnError(err).Warn("scim json response encoding failed")
		return nil
	}
}

func handleEmptyResponse(next func(*http.Request) error) zhttp_middlware.HandlerFuncWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		err := next(r)
		if err != nil {
			return err
		}

		w.WriteHeader(http.StatusNoContent)
		return nil
	}
}
