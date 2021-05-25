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
	router          *mux.Router
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
	h.router = mux.NewRouter()
	RegisterRoutes(h.router, h)
	//h.router.HandleFunc(defaultLabelPolicyLogoURL, h.UploadHandleFunc(&labelPolicyLogoUploader{idGenerator, false, true})).Methods("POST")
	//h.router.HandleFunc(defaultLabelPolicyLogoURL, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, true, false})).Methods("GET")
	//h.router.HandleFunc(defaultLabelPolicyLogoURL+preview, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, true, true})).Methods("GET")
	//
	//h.router.HandleFunc(defaultLabelPolicyLogoDarkURL, h.UploadHandleFunc(&labelPolicyLogoUploader{idGenerator, true, true})).Methods("POST")
	//h.router.HandleFunc(defaultLabelPolicyLogoDarkURL, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, true, false})).Methods("GET")
	//h.router.HandleFunc(defaultLabelPolicyLogoDarkURL+preview, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, true, true})).Methods("GET")
	//
	//h.router.HandleFunc(defaultLabelPolicyIconURL, h.UploadHandleFunc(&labelPolicyIconUploader{idGenerator, false, true})).Methods("POST")
	//h.router.HandleFunc(defaultLabelPolicyIconURL, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, false, true, false})).Methods("GET")
	//h.router.HandleFunc(defaultLabelPolicyIconURL+preview, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, false, true, true})).Methods("GET")
	//
	//h.router.HandleFunc(defaultLabelPolicyIconDarkURL, h.UploadHandleFunc(&labelPolicyIconUploader{idGenerator, true, true})).Methods("POST")
	//h.router.HandleFunc(defaultLabelPolicyIconDarkURL, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, true, true, false})).Methods("GET")
	//h.router.HandleFunc(defaultLabelPolicyIconDarkURL+preview, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, true, true, true})).Methods("GET")
	//
	//h.router.HandleFunc(defaultLabelPolicyFontURL, h.UploadHandleFunc(&labelPolicyFontUploader{idGenerator, true})).Methods("POST")
	//h.router.HandleFunc(defaultLabelPolicyFontURL, h.DownloadHandleFunc(&labelPolicyFontDownloader{orgRepo, true, false})).Methods("GET")
	//h.router.HandleFunc(defaultLabelPolicyFontURL+preview, h.DownloadHandleFunc(&labelPolicyFontDownloader{orgRepo, true, true})).Methods("GET")
	//
	//h.router.HandleFunc(orgLabelPolicyLogoURL, h.UploadHandleFunc(&labelPolicyLogoUploader{idGenerator, false, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyLogoURL, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, false, false})).Methods("GET")
	//h.router.HandleFunc(orgLabelPolicyLogoURL+preview, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, false, false, true})).Methods("GET")
	//
	//h.router.HandleFunc(orgLabelPolicyLogoDarkURL, h.UploadHandleFunc(&labelPolicyLogoUploader{idGenerator, true, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyLogoDarkURL, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, true, false, false})).Methods("GET")
	//h.router.HandleFunc(orgLabelPolicyLogoDarkURL+preview, h.DownloadHandleFunc(&labelPolicyLogoDownloader{orgRepo, true, false, true})).Methods("GET")
	//
	//h.router.HandleFunc(orgLabelPolicyIconDarkURL, h.UploadHandleFunc(&labelPolicyIconUploader{idGenerator, false, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyIconDarkURL, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, false, false, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyIconDarkURL+preview, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, false, false, true})).Methods("POST")
	//
	//h.router.HandleFunc(orgLabelPolicyIconURL, h.UploadHandleFunc(&labelPolicyIconUploader{idGenerator, true, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyIconURL, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, true, false, false})).Methods("GET")
	//h.router.HandleFunc(orgLabelPolicyIconURL+preview, h.DownloadHandleFunc(&labelPolicyIconDownloader{orgRepo, true, false, true})).Methods("GET")
	//
	//h.router.HandleFunc(orgLabelPolicyFontURL, h.UploadHandleFunc(&labelPolicyFontUploader{idGenerator, false})).Methods("POST")
	//h.router.HandleFunc(orgLabelPolicyFontURL, h.DownloadHandleFunc(&labelPolicyFontDownloader{orgRepo, false, false})).Methods("GET")
	//h.router.HandleFunc(orgLabelPolicyFontURL+preview, h.DownloadHandleFunc(&labelPolicyFontDownloader{orgRepo, false, true})).Methods("GET")
	//
	//h.router.HandleFunc(myUserAvatarURL, h.UploadHandleFunc(&myHumanAvatarUploader{})).Methods("POST")
	//h.router.HandleFunc(myUserAvatarURL, h.DownloadHandleFunc(&myHumanAvatarDownloader{})).Methods("GET")
	return h.router
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
