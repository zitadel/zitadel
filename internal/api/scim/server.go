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
	"github.com/zitadel/zitadel/internal/query"
)

func NewServer(
	command *command.Commands,
	query *query.Queries,
	verifier *authz.ApiTokenVerifier,
	userCodeAlg crypto.EncryptionAlgorithm,
	config *sconfig.Config,
	middlewares ...zhttp_middlware.MiddlewareWithErrorFunc) http.Handler {
	verifier.RegisterServer("SCIM-V2", schemas.HandlerPrefix, AuthMapping)
	return buildHandler(command, query, userCodeAlg, config, middlewares...)
}

func buildHandler(
	command *command.Commands,
	query *query.Queries,
	userCodeAlg crypto.EncryptionAlgorithm,
	cfg *sconfig.Config,
	middlewares ...zhttp_middlware.MiddlewareWithErrorFunc) http.Handler {

	router := mux.NewRouter()

	// content type middleware needs to run at the very beginning to correctly set content types of errors
	middlewares = append([]zhttp_middlware.MiddlewareWithErrorFunc{smiddleware.ContentTypeMiddleware}, middlewares...)
	middlewares = append(middlewares, smiddleware.ScimContextMiddleware(query))
	scimMiddleware := zhttp_middlware.ChainedWithErrorHandler(serrors.ErrorHandler, middlewares...)
	mapResource(router, scimMiddleware, sresources.NewUsersHandler(command, query, userCodeAlg, cfg))
	return router
}

func mapResource[T sresources.ResourceHolder](router *mux.Router, mw zhttp_middlware.ErrorHandlerFunc, handler sresources.ResourceHandler[T]) {
	adapter := sresources.NewResourceHandlerAdapter[T](handler)
	resourceRouter := router.PathPrefix("/" + path.Join(zhttp.OrgIdInPathVariable, string(handler.ResourceNamePlural()))).Subrouter()

	resourceRouter.Handle("", mw(handleResourceCreatedResponse(adapter.Create))).Methods(http.MethodPost)
	resourceRouter.Handle("/{id}", mw(handleResourceResponse(adapter.Get))).Methods(http.MethodGet)
	resourceRouter.Handle("/{id}", mw(handleResourceResponse(adapter.Replace))).Methods(http.MethodPut)
	resourceRouter.Handle("/{id}", mw(handleEmptyResponse(adapter.Delete))).Methods(http.MethodDelete)
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
