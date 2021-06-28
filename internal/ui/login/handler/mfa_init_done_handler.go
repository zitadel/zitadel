package handler

import (
	"github.com/caos/zitadel/internal/domain"
	"net/http"
)

const (
	tmplMFAInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMFAInitDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaDoneData) {
	var errType, errMessage string
	translator := l.getTranslator(authReq)
	data.baseData = l.getBaseData(r, authReq, translator, "MFA Init Done", errType, errMessage)
	data.profileData = l.getProfileData(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFAInitDone], data, nil)
}
