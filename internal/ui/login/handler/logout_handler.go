package handler

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
	data := l.getUserData(r, nil, "Logout Done", "", "")
	l.renderer.RenderTemplate(w, r, l.getTranslator(nil), l.renderer.Templates[tmplLogoutDone], data, nil)
}
