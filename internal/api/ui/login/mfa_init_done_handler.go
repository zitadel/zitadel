package login

import (
	"net/http"

	"github.com/zitadel/zitadel/internal/domain"
)

const (
	tmplMFAInitDone = "mfainitdone"
)

type mfaInitDoneData struct {
}

func (l *Login) renderMFAInitDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *mfaDoneData) {
	translator := l.getTranslator(r.Context(), authReq)
	data.baseData = l.getBaseData(r, authReq, translator, "InitMFADone.Title", "InitMFADone.Description", nil)
	data.profileData = l.getProfileData(authReq)
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplMFAInitDone], data, nil)
}
