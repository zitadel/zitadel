package handler

import (
	"github.com/caos/zitadel/internal/auth_request/model"
	"net/http"
)

const (
	tmplMfaInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMfaInitDone(w http.ResponseWriter, r *http.Request, authReq *model.AuthRequest, data *mfaDoneData) {
	var errType, errMessage string
	data.baseData = l.getBaseData(r, authReq, "Mfa Init Done", errType, errMessage)
	data.UserName = authReq.UserName
	l.renderer.RenderTemplate(w, r, l.renderer.Templates[tmplMfaInitDone], data, nil)
}
