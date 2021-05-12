package upload

import (
	"context"
	"net/http"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/static"
)

type Handler struct {
	router          *http.ServeMux
	errorHandler    ErrorHandler
	storage         static.Storage
	commands        *command.Commands
	authInterceptor *http_mw.AuthInterceptor
	idGenerator     id.Generator
}

type Uploader interface {
	Callback(ctx context.Context, info *domain.AssetInfo, orgID string, commands *command.Commands) error
	ObjectName(data authz.CtxData) (string, error)
	BucketName(data authz.CtxData) string
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
) http.Handler {
	h := &Handler{
		commands:        commands,
		errorHandler:    DefaultErrorHandler,
		authInterceptor: http_mw.AuthorizationInterceptor(verifier, authConfig),
		idGenerator:     idGenerator,
	}
	h.router = http.NewServeMux()
	h.router.HandleFunc(defaultLabelPolicyLogoURL, h.UploadHandleFunc(&labelPolicyLogo{idGenerator, false, true}))
	h.router.HandleFunc(defaultLabelPolicyLogoDarkURL, h.UploadHandleFunc(&labelPolicyLogo{idGenerator, true, true}))
	h.router.HandleFunc(defaultLabelPolicyIconURL, h.UploadHandleFunc(&labelPolicyIcon{idGenerator, false, true}))
	h.router.HandleFunc(defaultLabelPolicyIconDarkURL, h.UploadHandleFunc(&labelPolicyIcon{idGenerator, true, true}))
	h.router.HandleFunc(defaultLabelPolicyFontURL, h.UploadHandleFunc(&labelPolicyFont{idGenerator, true}))

	h.router.HandleFunc(orgLabelPolicyLogoURL, h.UploadHandleFunc(&labelPolicyLogo{idGenerator, false, false}))
	h.router.HandleFunc(orgLabelPolicyLogoDarkURL, h.UploadHandleFunc(&labelPolicyLogo{idGenerator, true, false}))
	h.router.HandleFunc(orgLabelPolicyIconDarkURL, h.UploadHandleFunc(&labelPolicyIcon{idGenerator, false, false}))
	h.router.HandleFunc(orgLabelPolicyIconURL, h.UploadHandleFunc(&labelPolicyIcon{idGenerator, true, false}))
	h.router.HandleFunc(orgLabelPolicyFontURL, h.UploadHandleFunc(&labelPolicyFont{idGenerator, false}))
	h.router.HandleFunc(userAvatarURL, h.UploadHandleFunc(&humanAvatar{}))
	return h.router
}

const maxMemory = 10 << 20
const paramFile = "file"

func (h *Handler) UploadHandleFunc(uploader Uploader) func(http.ResponseWriter, *http.Request) {
	return h.authInterceptor.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctxData := authz.GetCtxData(ctx)
			err := r.ParseMultipartForm(maxMemory)
			file, handler, err := r.FormFile(paramFile)
			if err != nil {
				h.errorHandler(w, r, err)
				return
			}
			defer func() {
				err = file.Close()
				logging.Log("UPLOAD-GDg34").OnError(err).Warn("could not close file")
			}()

			bucketName := uploader.BucketName(ctxData)
			objectName, err := uploader.ObjectName(ctxData)
			if err != nil {
				h.errorHandler(w, r, err)
				return
			}
			info, err := h.commands.UploadAsset(ctx, bucketName, objectName, handler.Header.Get("content-type"), file, handler.Size)
			if err != nil {
				h.errorHandler(w, r, err)
				return
			}
			err = uploader.Callback(ctx, info, ctxData.OrgID, h.commands)
			if err != nil {
				h.errorHandler(w, r, err)
				return
			}
		})
}
