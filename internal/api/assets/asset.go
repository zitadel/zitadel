package assets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gorilla/mux"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_util "github.com/zitadel/zitadel/internal/api/http"
	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/id"
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
	return h.errorHandler
}

func (h *Handler) Storage() static.Storage {
	return h.storage
}

func AssetAPI() func(context.Context) string {
	return func(ctx context.Context) string {
		return http_util.DomainContext(ctx).Origin() + HandlerPrefix
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
		code, ok := http_util.ZitadelErrorToHTTPStatusCode(r.Context(), err)
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

func NewHandler(
	commands *command.Commands,
	verifier authz.APITokenVerifier,
	systemAuthCOnfig authz.Config,
	authConfig authz.Config,
	idGenerator id.Generator,
	storage static.Storage,
	queries *query.Queries,
	callDurationInterceptor, instanceInterceptor, assetCacheInterceptor, accessInterceptor func(handler http.Handler) http.Handler,
	translator *i18n.Translator,
) http.Handler {
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler(translator),
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, systemAuthCOnfig, authConfig),
		idGenerator:     idGenerator,
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

const (
	maxMemory = 2 << 20
	paramFile = "file"
)

// fontExtensionContentTypes maps well-known font file extensions to their
// correct content type. Some valid font files, in particular TrueType and
// OpenType fonts whose table directory does not start with one of the
// commonly recognized SFNT table tags, are misdetected as
// application/octet-stream by content-based mimetype sniffing.
var fontExtensionContentTypes = map[string]string{
	".ttf":   "font/ttf",
	".otf":   "font/otf",
	".woff":  "font/woff",
	".woff2": "font/woff2",
	".eot":   "application/vnd.ms-fontobject",
	".ttc":   "font/collection",
}

// refineContentType corrects a generic application/octet-stream detection
// based on the uploaded file's extension, for extensions whose content-based
// sniffing is known to be unreliable. Any other detected content type is
// returned unchanged.
func refineContentType(contentType, filename string) string {
	if contentType != "application/octet-stream" {
		return contentType
	}
	if refined, ok := fontExtensionContentTypes[strings.ToLower(filepath.Ext(filename))]; ok {
		return refined
	}
	return contentType
}

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

		contentType := refineContentType(mimeType.String(), handler.Filename)

		size := handler.Size
		if !uploader.ContentTypeAllowed(contentType) {
			s.ErrorHandler()(w, r, fmt.Errorf("invalid content-type: %s", contentType), http.StatusBadRequest)
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
			ContentType:   contentType,
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
			splitPath := strings.Split(r.RequestURI, ownerPath+"/")
			if len(splitPath) < 2 {
				s.ErrorHandler()(w, r, fmt.Errorf("invalid request URI format: %v", r.RequestURI), http.StatusNotFound)
				return
			}
			path = splitPath[1]
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
