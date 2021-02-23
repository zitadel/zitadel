package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
)

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest) {
	callback := l.oidcAuthCallbackURL + authReq.ID
	http.Redirect(w, r, callback, http.StatusFound)
}
