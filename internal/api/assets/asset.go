package assets

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/caos/logging"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/superseriousbusiness/exifremove/pkg/exifremove"

	"github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/static"
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

type Uploader interface {
	UploadAsset(ctx context.Context, info string, asset *command.AssetUpload, commands *command.Commands) error
	ObjectName(data authz.CtxData) (string, error)
	ResourceOwner(data authz.CtxData) string
	ContentTypeAllowed(contentType string) bool
	MaxFileSize() int64
	ObjectType() static.ObjectType
}

type Downloader interface {
	ObjectName(ctx context.Context, path string) (string, error)
	BucketName(ctx context.Context, id string) string
}

type ErrorHandler func(http.ResponseWriter, *http.Request, error, int)

func DefaultErrorHandler(w http.ResponseWriter, r *http.Request, err error, code int) {
	logging.Log("ASSET-g5ef1").WithError(err).WithField("uri", r.RequestURI).Error("error occurred on asset api")
	http.Error(w, err.Error(), code)
}

func NewHandler(commands *command.Commands, verifier *authz.TokenVerifier, authConfig authz.Config, idGenerator id.Generator, storage static.Storage, queries *query.Queries) http.Handler {
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler,
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, authConfig),
		idGenerator:     idGenerator,
		storage:         storage,
		query:           queries,
	}

	verifier.RegisterServer("Management-API", "assets", AssetsService_AuthMethods) //TODO: separate api?
	router := mux.NewRouter()
	router.Use(sentryhttp.New(sentryhttp.Options{}).Handle)
	RegisterRoutes(router, h)
	router.PathPrefix("/{id}").Methods("GET").HandlerFunc(DownloadHandleFunc(h, h.GetFile()))
	return router
}

func (h *Handler) GetFile() Downloader {
	return &publicFileDownloader{}
}

type publicFileDownloader struct{}

func (l *publicFileDownloader) ObjectName(_ context.Context, path string) (string, error) {
	return path, nil
}

func (l *publicFileDownloader) BucketName(_ context.Context, id string) string {
	return id
}

const maxMemory = 2 << 20
const paramFile = "file"

func UploadHandleFunc(s AssetsService, uploader Uploader) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctxData := authz.GetCtxData(ctx)
		err := r.ParseMultipartForm(maxMemory)
		file, handler, err := r.FormFile(paramFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			err = file.Close()
			logging.Log("UPLOAD-GDg34").OnError(err).Warn("could not close file")
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

		resourceOwner := uploader.ResourceOwner(ctxData)
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
		id := mux.Vars(r)["id"]
		resourceOwner := downloader.BucketName(ctx, id)
		path := ""
		if id != "" {
			path = strings.Split(r.RequestURI, id+"/")[1]
		}
		objectName, err := downloader.ObjectName(ctx, path)
		if err != nil {
			s.ErrorHandler()(w, r, fmt.Errorf("download failed: %v", err), http.StatusInternalServerError)
			return
		}
		if objectName == "" {
			s.ErrorHandler()(w, r, fmt.Errorf("file not found: %v", objectName), http.StatusNotFound)
			return
		}
		if err = GetAsset(w, r, resourceOwner, objectName, s.Storage()); err != nil {
			s.ErrorHandler()(w, r, err, http.StatusInternalServerError)
		}
	}
}

func GetAsset(w http.ResponseWriter, r *http.Request, resourceOwner, objectName string, storage static.Storage) error {
	data, getInfo, err := storage.GetObject(r.Context(), "0", resourceOwner, objectName)
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	info, err := getInfo()
	if err != nil {
		return fmt.Errorf("download failed: %v", err)
	}
	if info.Hash == r.Header.Get(http_util.IfNoneMatch) {
		w.WriteHeader(304)
		return nil
	}
	w.Header().Set(http_util.ContentLength, strconv.FormatInt(info.Size, 10))
	w.Header().Set(http_util.ContentType, info.ContentType)
	w.Header().Set(http_util.LastModified, info.LastModified.Format(time.RFC1123))
	w.Header().Set(http_util.Etag, info.Hash)
	_, err = w.Write(data)
	logging.New().OnError(err).Error("error writing response for asset")
	return nil
}

func removeExif(file io.Reader, size int64, contentType string) (io.Reader, int64, error) {
	if !isAllowedContentType(contentType) {
		return file, size, nil
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		return file, 0, err
	}
	data, err := exifremove.Remove(buf.Bytes())
	if err != nil {
		return nil, 0, err
	}
	return bytes.NewReader(data), int64(len(data)), nil
}

func isAllowedContentType(contentType string) bool {
	return strings.HasSuffix(contentType, "png") ||
		strings.HasSuffix(contentType, "jpg") ||
		strings.HasSuffix(contentType, "jpeg")
}
