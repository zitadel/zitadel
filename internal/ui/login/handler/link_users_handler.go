package handler

import (
	http_mw "github.com/caos/zitadel/internal/api/http/middleware"
	"net/http"

	"github.com/caos/zitadel/internal/auth_request/model"
)

const (
	tmplLinkUsersDone = "linkusersdone"
)

func (l *Login) linkUsers(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	iam, err := l.authRepo.GetIAM(r.Context())
	if err != nil {
		l.renderLinkUsersDone(w, r, authReq, err)
		return
	}
	resourceOwner := iam.GlobalOrgID
	if authReq.GetScopeOrgID() != "" {
		resourceOwner = authReq.GetScopeOrgID()
	}

	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	err = l.authRepo.LinkExternalUsers(setContext(r.Context(), resourceOwner), authReq.ID, userAgentID)
	l.renderLinkUsersDone(w, r, authReq, err)
}

func (l *Login) renderLinkUsersDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, err error) {
	var errType, errMessage string
	data := l.getUserData(r, authReq, "Linking Users Done", errType, errMessage)
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplLinkUsersDone], data, nil)
}
