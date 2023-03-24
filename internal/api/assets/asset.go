package assets

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
)

const (
	HandlerPrefix = "/assets/v1"
)

type Handler struct {
	errorHandler    ErrorHandler
	storage         static.Storage
	commands        *command.Commands
	authInterceptor *http_mw.AuthInterceptor
	idGenerator     id.Generator
	query           *query.Queries
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

func AssetAPI(externalSecure bool) func(context.Context) string {
	return func(ctx context.Context) string {
		return http_util.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix
	}
}

func AssetAPIFromDomain(externalSecure bool, externalPort uint16) func(context.Context) string {
	return func(ctx context.Context) string {
		return http_util.BuildHTTP(authz.GetInstance(ctx).RequestedDomain(), externalPort, externalSecure) + HandlerPrefix
	}
}

type Uploader interface {
	UploadAsset(ctx context.Context, info string, asset *command.AssetUpload, commands *command.Commands) error
	ObjectName(data authz.CtxData) (string, error)
	ResourceOwner(instance authz.Instance, data authz.CtxData) string
	ContentTypeAllowed(contentType string) bool
	MaxFileSize() int64
	ObjectType() static.ObjectType
}

type Downloader interface {
	ObjectName(ctx context.Context, path string) (string, error)
	ResourceOwner(ctx context.Context, ownerPath string) string
}

type ErrorHandler func(http.ResponseWriter, *http.Request, error, int)

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error, code int) {
	logging.WithFields("uri", r.RequestURI).WithError(err).Warn("error occurred on asset api")
	http.Error(w, err.Error(), code)
}

func NewHandler(commands *command.Commands, verifier *authz.TokenVerifier, authConfig authz.Config, idGenerator id.Generator, storage static.Storage, queries *query.Queries, callDurationInterceptor, instanceInterceptor, assetCacheInterceptor, accessInterceptor func(handler http.Handler) http.Handler) http.Handler {
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler,
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, authConfig),
		idGenerator:     idGenerator,
		storage:         storage,
		query:           queries,
	}

	verifier.RegisterServer("Assets-API", "assets", AssetsService_AuthMethods)
	router := mux.NewRouter()
	router.Use(callDurationInterceptor, instanceInterceptor, assetCacheInterceptor, accessInterceptor)
	RegisterRoutes(router, h)
	router.PathPrefix("/{owner}").Methods("GET").HandlerFunc(DownloadHandleFunc(h, h.GetFile()))
	return http_util.CopyHeadersToContext(http_mw.CORSInterceptor(router))
}

func (h *Handler) GetFile() Downloader {
	return &publicFileDownloader{}
}

type publicFileDownloader struct{}

func (l *publicFileDownloader) ObjectName(_ context.Context, path string) (string, error) {
	return path, nil
}

func (l *publicFileDownloader) ResourceOwner(_ context.Context, ownerPath string) string {
	return ownerPath
}

const maxMemory = 2 << 20
const paramFile = "file"

func UploadHandleFunc(s AssetsService, uploader Uploader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxData := authz.GetCtxData(ctx)
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		file, handler, err := r.FormFile(paramFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			err = file.Close()
			logging.OnError(err).Warn("could not close file")
		}()
		contentType := handler.Header.Get("content-type")
		size := handler.Size
		if !uploader.ContentTypeAllowed(contentType) {
			s.ErrorHandler()(w, r, fmt.Errorf("invalid content-type: %s", contentType), http.StatusBadRequest)
			return
		}
		if size > uploader.MaxFileSize() {
			s.ErrorHandler()(w, r, fmt.Errorf("file to big, max file size is %vKB", uploader.MaxFileSize()/1024), http.StatusBadRequest)
			return
		}

		resourceOwner := uploader.ResourceOwner(authz.GetInstance(ctx), ctxData)
		objectName, err := uploader.ObjectName(ctxData)
		if err != nil {
			s.ErrorHandler()(w, r, fmt.Errorf("upload failed: %v", err), http.StatusInternalServerError)
			return
		}
		uploadInfo := &command.AssetUpload{
			ResourceOwner: resourceOwner,
			ObjectName:    objectName,
			ContentType:   contentType,
			ObjectType:    uploader.ObjectType(),
			File:          file,
			Size:          size,
		}
		err = uploader.UploadAsset(ctx, ctxData.OrgID, uploadInfo, s.Commands())
		if err != nil {
			s.ErrorHandler()(w, r, fmt.Errorf("upload failed: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func DownloadHandleFunc(s AssetsService, downloader Downloader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.Storage() == nil {
			return
		}
		ctx := r.Context()
		ownerPath := mux.Vars(r)["owner"]
		resourceOwner := downloader.ResourceOwner(ctx, ownerPath)
		path := ""
		if ownerPath != "" {
			path = strings.Split(r.RequestURI, ownerPath+"/")[1]
		}
		objectName, err := downloader.ObjectName(ctx, path)
		if err != nil {
			s.ErrorHandler()(w, r, fmt.Errorf("download failed: %v", err), http.StatusInternalServerError)
			return
		}
		if objectName == "" {
			s.ErrorHandler()(w, r, fmt.Errorf("file not found: %v", path), http.StatusNotFound)
			return
		}
		if err = GetAsset(w, r, resourceOwner, objectName, s.Storage()); err != nil {
			s.ErrorHandler()(w, r, err, http.StatusInternalServerError)
		}
	}
}

func GetAsset(w http.ResponseWriter, r *http.Request, resourceOwner, objectName string, storage static.Storage) error {
	split := strings.Split(objectName, "?v=")
	if len(split) == 2 {
		objectName = split[0]
	}
	data, getInfo, err := storage.GetObject(r.Context(), authz.GetInstance(r.Context()).InstanceID(), resourceOwner, objectName)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	info, err := getInfo()
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	if info.Hash == strings.Trim(r.Header.Get(http_util.IfNoneMatch), "\"") {
		w.Header().Set(http_util.LastModified, info.LastModified.Format(time.RFC1123))
		w.Header().Set(http_util.Etag, "\""+info.Hash+"\"")
		w.WriteHeader(304)
		return nil
	}
	w.Header().Set(http_util.ContentLength, strconv.FormatInt(info.Size, 10))
	w.Header().Set(http_util.ContentType, info.ContentType)
	w.Header().Set(http_util.LastModified, info.LastModified.Format(time.RFC1123))
	w.Header().Set(http_util.Etag, "\""+info.Hash+"\"")
	_, err = w.Write(data)
	logging.New().OnError(err).Error("error writing response for asset")
	return nil
}
