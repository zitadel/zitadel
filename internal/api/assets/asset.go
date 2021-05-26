package assets

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/caos/logging"
	"github.com/gorilla/mux"

	"github.com/caos/zitadel/internal/api/authz"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/management/repository"
	"github.com/caos/zitadel/internal/static"
)

type Handler struct {
	errorHandler    ErrorHandler
	storage         static.Storage
	commands        *command.Commands
	authInterceptor *http_mw.AuthInterceptor
	idGenerator     id.Generator
	orgRepo         repository.OrgRepository
}

func (h *Handler) AuthInterceptor() *http_mw.AuthInterceptor {
	return h.authInterceptor
}

func (h *Handler) Commands() *command.Commands {
	return h.commands
}

func (h *Handler) ErrorHandler() ErrorHandler {
	return DefaultErrorHandler
}

func (h *Handler) Storage() static.Storage {
	return h.storage
}

type Uploader interface {
	Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error
	ObjectName(data authz.CtxData) (string, error)
	BucketName(data authz.CtxData) string
}

type Downloader interface {
	ObjectName(ctx context.Context) (string, error)
	BucketName(ctx context.Context) string
}

type ErrorHandler func(http.ResponseWriter, *http.Request, error)

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func NewHandler(
	commands *command.Commands,
	verifier *authz.TokenVerifier,
	authConfig authz.Config,
	idGenerator id.Generator,
	storage static.Storage,
	orgRepo repository.OrgRepository,
) http.Handler {
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler,
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, authConfig),
		idGenerator:     idGenerator,
		storage:         storage,
		orgRepo:         orgRepo,
	}

	verifier.RegisterServer("Management-API", "assets", AssetsService_AuthMethods)
	router := mux.NewRouter()
	RegisterRoutes(router, h)
	return router
}

const maxMemory = 10 << 20
const paramFile = "file"

func UploadHandleFunc(s AssetsService, uploader Uploader) func(http.ResponseWriter, *http.Request) {
	return s.AuthInterceptor().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctxData := authz.GetCtxData(ctx)
			err := r.ParseMultipartForm(maxMemory)
			file, handler, err := r.FormFile(paramFile)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			defer func() {
				err = file.Close()
				logging.Log("UPLOAD-GDg34").OnError(err).Warn("could not close file")
			}()

			bucketName := uploader.BucketName(ctxData)
			objectName, err := uploader.ObjectName(ctxData)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			info, err := s.Commands().UploadAsset(ctx, bucketName, objectName, handler.Header.Get("content-type"), file, handler.Size)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			err = uploader.Callback(ctx, info, ctxData.OrgID, s.Commands())
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
		})
}

func DownloadHandleFunc(s AssetsService, downloader Downloader) func(http.ResponseWriter, *http.Request) {
	return s.AuthInterceptor().HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if s.Storage() == nil {
				return
			}
			ctx := r.Context()

			bucketName := downloader.BucketName(ctx)
			objectName, err := downloader.ObjectName(ctx)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			reader, getInfo, err := s.Storage().GetObject(ctx, bucketName, objectName)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			data, err := ioutil.ReadAll(reader)
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			info, err := getInfo()
			if err != nil {
				s.ErrorHandler()(w, r, err)
				return
			}
			w.Header().Set("content-length", strconv.FormatInt(info.Size, 16))
			w.Header().Set("content-type", info.ContentType)
			w.Header().Set("ETag", info.ETag)
			w.Write(data)
		})
}
