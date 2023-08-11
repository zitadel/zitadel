package login

//
//import (
//	"net/http"
//
//	"github.com/zitadel/zitadel/internal/domain"
//
//	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
//)
//
//const (
//	tmplSMSVerification = "smsverification"
//)
//
//type mfaSMSData struct {
//	userData
//	MFAProviders     []domain.MFAType
//	SelectedProvider domain.MFAType
//}
//
//type mfaSMSFormData struct {
//	Resend           bool           `schema:"resend"`
//	Code             string         `schema:"code"`
//	SelectedProvider domain.MFAType `schema:"provider"`
//}
//
//func (l *Login) renderOTPSMSVerification(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, providers []domain.MFAType, err error) {
//	var errID, errMessage string
//	if err == nil {
//		userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//		err = l.authRepo.SendMFAOTPSMS(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, authReq.ID, userAgentID)
//	}
//	if err != nil {
//		errID, errMessage = l.getErrorMessage(r, err)
//	}
//	data := &mfaSMSData{
//		userData:         l.getUserData(r, authReq, "VerifyMFAU2F.Title", "VerifyMFAU2F.Description", errID, errMessage),
//		MFAProviders:     removeSelectedProviderFromList(providers, domain.MFATypeOTPSMS),
//		SelectedProvider: domain.MFATypeOTPSMS,
//	}
//	l.renderer.RenderTemplate(w, r, l.getTranslator(r.Context(), authReq), l.renderer.Templates[tmplSMSVerification], data, nil)
//}
//
//func (l *Login) handleOTPSMSVerification(w http.ResponseWriter, r *http.Request) {
//	formData := new(mfaSMSFormData)
//	authReq, err := l.getAuthRequestAndParseData(r, formData)
//	if err != nil {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	step, ok := authReq.PossibleSteps[0].(*domain.MFAVerificationStep)
//	if !ok {
//		l.renderError(w, r, authReq, err)
//		return
//	}
//	if formData.Resend {
//		l.renderOTPSMSVerification(w, r, authReq, step.MFAProviders, nil)
//		return
//	}
//	if formData.Code == "" {
//		l.renderMFAVerifySelected(w, r, authReq, step, formData.SelectedProvider, nil)
//		return
//	}
//	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
//	err = l.authRepo.VerifyMFAOTPSMS(setContext(r.Context(), authReq.UserOrgID), authReq.UserID, authReq.UserOrgID, formData.Code, authReq.ID, userAgentID, domain.BrowserInfoFromRequest(r))
//
//	metadata, actionErr := l.runPostInternalAuthenticationActions(authReq, r, authMethodOTPSMS, err)
//	if err == nil && actionErr == nil && len(metadata) > 0 {
//		_, err = l.command.BulkSetUserMetadata(r.Context(), authReq.UserID, authReq.UserOrgID, metadata...)
//	} else if actionErr != nil && err == nil {
//		err = actionErr
//	}
//
//	if err != nil {
//		l.renderOTPSMSVerification(w, r, authReq, step.MFAProviders, err)
//		return
//	}
//	l.renderNextStep(w, r, authReq)
//}
