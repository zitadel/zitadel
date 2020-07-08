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
	data := userData{
		baseData: l.getBaseData(r, nil, "Logout Done", "", ""),
	}
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplLogoutDone], data, nil)
}
