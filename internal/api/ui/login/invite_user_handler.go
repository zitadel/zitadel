package login

import (
	"net/http"
	"net/url"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryInviteUserCode      = "code"
	queryInviteUserUserID    = "userID"
	queryInviteUserLoginName = "loginname"

	tmplInviteUser     = "inviteuser"
	tmplInviteUserDone = "inviteuserdone"
)

type inviteUserFormData struct {
	Code            string `schema:"code"`
	LoginName       string `schema:"loginname"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	Resend          bool   `schema:"resend"`
}

type inviteUserData struct {
	baseData
	profileData
	Code         string
	LoginName    string
	UserID       string
	MinLength    uint64
	HasUppercase string
	HasLowercase string
	HasNumber    string
	HasSymbol    string
}

func InviteUserLink(origin, userID, loginName, code, orgID string, authRequestID string) string {
	v := url.Values{}
	v.Set(queryInviteUserUserID, userID)
	v.Set(queryInviteUserLoginName, loginName)
	v.Set(queryInviteUserCode, code)
	v.Set(queryOrgID, orgID)
	v.Set(QueryAuthRequestID, authRequestID)
	return externalLink(origin) + EndpointInviteUser + "?" + v.Encode()
}

func (l *Login) handleInviteUser(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryInviteUserUserID)
	code := r.FormValue(queryInviteUserCode)
	loginName := r.FormValue(queryInviteUserLoginName)
	l.renderInviteUser(w, r, authReq, userID, loginName, code, nil)
}

func (l *Login) handleInviteUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(inviteUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInvite(w, r, authReq, data.UserID, data.LoginName)
		return
	}
	l.checkUserInviteCode(w, r, authReq, data, nil)
}

func (l *Login) checkUserInviteCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *inviteUserFormData, err error) {
	if data.Password != data.PasswordConfirm {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-fsdfd", "Errors.User.Password.ConfirmationWrong")
		l.renderInviteUser(w, r, authReq, data.UserID, data.LoginName, data.Code, err)
		return
	}
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	//userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	//err = l.command.VerifyInviteCode(setContext(r.Context(), userOrgID), data.UserID, userOrgID, data.Code, data.Password, userAgentID)
	//if err != nil {
	//	l.renderInviteUser(w, r, authReq, data.UserID, data.LoginName, "", err)
	//	return
	//}
	l.renderInviteUserDone(w, r, authReq, userOrgID)
}

func (l *Login) resendUserInvite(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID string, loginName string) {
	var userOrgID, authRequestID string
	if authReq != nil {
		userOrgID = authReq.UserOrgID
		authRequestID = authReq.ID
	}
	_, err := l.command.ResendInviteCode(setContext(r.Context(), userOrgID), userID, userOrgID, authRequestID)
	l.renderInviteUser(w, r, authReq, userID, loginName, "", err)
}

func (l *Login) renderInviteUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, loginName string, code string, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if authReq != nil {
		userID = authReq.UserID
	}

	translator := l.getTranslator(r.Context(), authReq)
	data := inviteUserData{
		baseData:    l.getBaseData(r, authReq, translator, "InviteUser.Title", "InviteUser.Description", errID, errMessage),
		profileData: l.getProfileData(authReq),
		UserID:      userID,
		Code:        code,
	}
	// if the user clicked on the link in the mail, we need to make sure the loginName is rendered
	if authReq == nil {
		data.LoginName = loginName
		data.UserName = loginName
	}
	policy := l.getPasswordComplexityPolicyByUserID(r, userID)
	if policy != nil {
		data.MinLength = policy.MinLength
		if policy.HasUppercase {
			data.HasUppercase = UpperCaseRegex
		}
		if policy.HasLowercase {
			data.HasLowercase = LowerCaseRegex
		}
		if policy.HasSymbol {
			data.HasSymbol = SymbolRegex
		}
		if policy.HasNumber {
			data.HasNumber = NumberRegex
		}
	}
	if authReq == nil {
		user, err := l.query.GetUserByID(r.Context(), false, userID)
		if err == nil {
			l.customTexts(r.Context(), translator, user.ResourceOwner)
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInviteUser], data, nil)
}

func (l *Login) renderInviteUserDone(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, orgID string) {
	translator := l.getTranslator(r.Context(), authReq)
	data := l.getUserData(r, authReq, translator, "InviteUserDone.Title", "InviteUserDone.Description", "", "")
	if authReq == nil {
		l.customTexts(r.Context(), translator, orgID)
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInviteUserDone], data, nil)
}
