package login

import (
	"fmt"
	"net/http"
	"net/url"

	http_mw "github.com/zitadel/zitadel/internal/api/http/middleware"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	queryInviteUserCode      = "code"
	queryInviteUserUserID    = "userID"
	queryInviteUserLoginName = "loginname"

	tmplInviteUser = "inviteuser"
)

type inviteUserFormData struct {
	Code            string `schema:"code"`
	LoginName       string `schema:"loginname"`
	Password        string `schema:"password"`
	PasswordConfirm string `schema:"passwordconfirm"`
	UserID          string `schema:"userID"`
	OrgID           string `schema:"orgID"`
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

func InviteUserLinkTemplate(origin, userID, orgID string, authRequestID string) string {
	return fmt.Sprintf("%s%s?%s=%s&%s=%s&%s=%s&%s=%s&%s=%s",
		externalLink(origin), EndpointInviteUser,
		queryInviteUserUserID, userID,
		queryInviteUserLoginName, "{{.LoginName}}",
		queryInviteUserCode, "{{.Code}}",
		queryOrgID, orgID,
		QueryAuthRequestID, authRequestID)
}

func (l *Login) handleInviteUser(w http.ResponseWriter, r *http.Request) {
	authReq := l.checkOptionalAuthRequestOfEmailLinks(r)
	userID := r.FormValue(queryInviteUserUserID)
	orgID := r.FormValue(queryOrgID)
	code := r.FormValue(queryInviteUserCode)
	loginName := r.FormValue(queryInviteUserLoginName)
	l.renderInviteUser(w, r, authReq, userID, orgID, loginName, code, nil)
}

func (l *Login) handleInviteUserCheck(w http.ResponseWriter, r *http.Request) {
	data := new(inviteUserFormData)
	authReq, err := l.getAuthRequestAndParseData(r, data)
	if err != nil {
		l.renderError(w, r, nil, err)
		return
	}

	if data.Resend {
		l.resendUserInvite(w, r, authReq, data.UserID, data.OrgID, data.LoginName)
		return
	}
	l.checkUserInviteCode(w, r, authReq, data)
}

func (l *Login) checkUserInviteCode(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, data *inviteUserFormData) {
	if data.Password != data.PasswordConfirm {
		err := zerrors.ThrowInvalidArgument(nil, "VIEW-KJS3h", "Errors.User.Password.ConfirmationWrong")
		l.renderInviteUser(w, r, authReq, data.UserID, data.OrgID, data.LoginName, data.Code, err)
		return
	}
	userOrgID := ""
	if authReq != nil {
		userOrgID = authReq.UserOrgID
	}
	userAgentID, _ := http_mw.UserAgentIDFromCtx(r.Context())
	_, err := l.command.VerifyInviteCodeSetPassword(setUserContext(r.Context(), data.UserID, userOrgID), data.UserID, data.Code, data.Password, userAgentID)
	if err != nil {
		l.renderInviteUser(w, r, authReq, data.UserID, data.OrgID, data.LoginName, "", err)
		return
	}
	if authReq == nil {
		l.defaultRedirect(w, r)
		return
	}
	l.renderNextStep(w, r, authReq)
}

func (l *Login) resendUserInvite(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, orgID, loginName string) {
	var userOrgID, authRequestID string
	if authReq != nil {
		userOrgID = authReq.UserOrgID
		authRequestID = authReq.ID
	}
	_, err := l.command.ResendInviteCode(setUserContext(r.Context(), userID, userOrgID), userID, userOrgID, authRequestID)
	l.renderInviteUser(w, r, authReq, userID, orgID, loginName, "", err)
}

func (l *Login) renderInviteUser(w http.ResponseWriter, r *http.Request, authReq *domain.AuthRequest, userID, orgID, loginName string, code string, err error) {
	var errID, errMessage string
	if err != nil {
		errID, errMessage = l.getErrorMessage(r, err)
	}
	if authReq != nil {
		userID = authReq.UserID
		orgID = authReq.UserOrgID
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
		if err == nil {
			l.customTexts(r.Context(), translator, orgID)
		}
	}
	l.renderer.RenderTemplate(w, r, translator, l.renderer.Templates[tmplInviteUser], data, nil)
}
