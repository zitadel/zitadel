package assets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	"github.com/zitadel/logging"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	HandlerPrefix = "/assets/v1"
)

type Handler struct {
	errorHandler    ErrorHandler
	storage         static.Storage
	commands        *command.Commands
	authInterceptor *http_mw.AuthInterceptor
	query           *query.Queries
}

func (h *Handler) AuthInterceptor() *http_mw.AuthInterceptor {
	return h.authInterceptor
}

func (h *Handler) Commands() *command.Commands {
	return h.commands
}

func (h *Handler) ErrorHandler() ErrorHandler {
	return h.errorHandler
}

func (h *Handler) Storage() static.Storage {
	return h.storage
}

func AssetAPI(externalSecure bool) func(context.Context) string {
	return func(ctx context.Context) string {
		return http_util.BuildOrigin(authz.GetInstance(ctx).RequestedHost(), externalSecure) + HandlerPrefix
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

type ErrorHandler func(w http.ResponseWriter, r *http.Request, err error, defaultCode int)

func DefaultErrorHandler(translator *i18n.Translator) func(w http.ResponseWriter, r *http.Request, err error, defaultCode int) {
	return func(w http.ResponseWriter, r *http.Request, err error, defaultCode int) {
		logging.WithFields("uri", r.RequestURI).WithError(err).Warn("error occurred on asset api")
		code, ok := http_util.ZitadelErrorToHTTPStatusCode(err)
		if !ok {
			code = defaultCode
		}
		zErr := new(zerrors.ZitadelError)
		if errors.As(err, &zErr) {
			zErr.SetMessage(translator.LocalizeFromCtx(r.Context(), zErr.GetMessage(), nil))
			zErr.Parent = nil // ensuring we don't leak any unwanted information
			err = zErr
		}
		http.Error(w, err.Error(), code)
	}
}

func NewHandler(commands *command.Commands, verifier authz.APITokenVerifier, authConfig authz.Config, storage static.Storage, queries *query.Queries, callDurationInterceptor, instanceInterceptor, assetCacheInterceptor, accessInterceptor func(handler http.Handler) http.Handler) http.Handler {
	translator, err := i18n.NewZitadelTranslator(language.English)
	logging.OnError(err).Panic("unable to get translator")
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler(translator),
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, authConfig),
		storage:         storage,
		query:           queries,
	}

	verifier.RegisterServer("Assets-API", "assets", AssetsService_AuthMethods)
	router := mux.NewRouter()
	csp := http_mw.SecurityHeaders(&http_mw.DefaultSCP, nil)
	router.Use(callDurationInterceptor, instanceInterceptor, assetCacheInterceptor, accessInterceptor, csp)
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

		mimeType, err := mimetype.DetectReader(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		size := handler.Size
		if !uploader.ContentTypeAllowed(mimeType.String()) {
			s.ErrorHandler()(w, r, fmt.Errorf("invalid content-type: %s", mimeType), http.StatusBadRequest)
			return
		}
		if size > uploader.MaxFileSize() {
			s.ErrorHandler()(w, r, fmt.Errorf("file too big, max file size is %vKB", uploader.MaxFileSize()/1024), http.StatusBadRequest)
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
			ContentType:   mimeType.String(),
			ObjectType:    uploader.ObjectType(),
			File:          file,
			Size:          size,
		}
		err = uploader.UploadAsset(ctx, ctxData.OrgID, uploadInfo, s.Commands())
		if err != nil {
			s.ErrorHandler()(w, r, fmt.Errorf("upload failed: %w", err), http.StatusInternalServerError)
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
		return fmt.Errorf("download failed: %w", err)
	}
	info, err := getInfo()
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
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
