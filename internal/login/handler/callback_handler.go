package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest) {
	callback := l.oidcAuthCallbackURL + authReq.ID
	http.Redirect(w, r, callback, http.StatusFound)
}
