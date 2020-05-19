package handler

import (
	"net/http"

)

func (l *Login) redirectToCallback(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession) {
	var callback string
	if authSession.Type == model.TYPE_OIDC {
		callback = l.oidcAuthCallbackURL + authSession.GetFullID()
	}
	http.Redirect(w, r, callback, http.StatusFound)
}
