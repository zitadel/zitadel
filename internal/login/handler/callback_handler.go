package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	//var callback string
	//if authReq.Ty == model.TYPE_OIDC {
	//	callback = l.oidcAuthCallbackURL + authReq.GetFullID()
	//}
	//http.Redirect(w, r, callback, http.StatusFound)
}
