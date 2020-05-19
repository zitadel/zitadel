package handler

import (
	"net/http"

	"github.com/caos/citadel/login/internal/model"
)

const (
	tmplMfaInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMfaInitDone(w http.ResponseWriter, r *http.Request, authSession *model.AuthSession, data *mfaDoneData) {
	var errType, errMessage string
	data.baseData = l.getBaseData(r, authSession, "Mfa Init Done", errType, errMessage)
	data.UserName = authSession.UserSession.User.UserName
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaInitDone], data, nil)
}
