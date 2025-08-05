package login

import (
	"net/http"
)

const (
	tmplLogoutDone = "logoutdone"
)

func (l *Login) handleLogoutDone(w http.ResponseWriter, r *http.Request) {
	l.renderLogoutDone(w, r)
}

func (l *Login) renderLogoutDone(w http.ResponseWriter, r *http.Request) {
	translator := l.getTranslator(r.Context(), nil)
	data := l.getUserData(r, nil, translator, "LogoutDone.Title", "LogoutDone.Description", nil)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplLogoutDone], data, nil)
}
