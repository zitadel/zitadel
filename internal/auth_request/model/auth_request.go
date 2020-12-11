package model

import (
	"github.com/caos/zitadel/internal/iam/model"
	"golang.org/x/text/language"
	"strings"
	"time"

	"github.com/caos/zitadel/internal/errors"
)

type AuthRequest struct {
	ID            string
	AgentID       string
	CreationDate  time.Time
	ChangeDate    time.Time
	BrowserInfo   *BrowserInfo
	ApplicationID string
	CallbackURI   string
	TransferState string
	Prompt        Prompt
	PossibleLOAs  []LevelOfAssurance
	UiLocales     []string
	LoginHint     string
	MaxAuthAge    uint32
	Request       Request

	levelOfAssurance    LevelOfAssurance
	UserID              string
	LoginName           string
	DisplayName         string
	UserOrgID           string
	RequestedOrgID      string
	SelectedIDPConfigID string
	LinkingUsers        []*ExternalUser
	PossibleSteps       []NextStep
	PasswordVerified    bool
	MFAsVerified        []MFAType
	Audience            []string
	AuthTime            time.Time
	Code                string
	LoginPolicy         *model.LoginPolicyView
	AllowedExternalIDPs []*model.IDPProviderView
}

type ExternalUser struct {
	IDPConfigID       string
	ExternalUserID    string
	DisplayName       string
	PreferredUsername string
	FirstName         string
	LastName          string
	NickName          string
	Email             string
	IsEmailVerified   bool
	PreferredLanguage language.Tag
	Phone             string
	IsPhoneVerified   bool
}

type Prompt int32

const (
	PromptUnspecified Prompt = iota
	PromptNone
	PromptLogin
	PromptConsent
	PromptSelectAccount
)

type LevelOfAssurance int

const (
	LevelOfAssuranceNone LevelOfAssurance = iota
)

func NewAuthRequest(id, agentID string, info *BrowserInfo, applicationID, callbackURI, transferState string,
	prompt Prompt, possibleLOAs []LevelOfAssurance, uiLocales []string, loginHint, preselectedUserID string, maxAuthAge uint32, request Request) *AuthRequest {
	return &AuthRequest{
		ID:            id,
		AgentID:       agentID,
		BrowserInfo:   info,
		ApplicationID: applicationID,
		CallbackURI:   callbackURI,
		TransferState: transferState,
		Prompt:        prompt,
		PossibleLOAs:  possibleLOAs,
		UiLocales:     uiLocales,
		LoginHint:     loginHint,
		UserID:        preselectedUserID,
		MaxAuthAge:    maxAuthAge,
		Request:       request,
	}
}

func NewAuthRequestFromType(requestType AuthRequestType) (*AuthRequest, error) {
	request, ok := authRequestTypeMapping[requestType]
	if !ok {
		return nil, errors.ThrowInvalidArgument(nil, "MODEL-ds2kl", "invalid request type")
	}
	return &AuthRequest{Request: request}, nil
}

func (a *AuthRequest) IsValid() bool {
	return a.ID != "" &&
		a.AgentID != "" &&
		a.BrowserInfo != nil && a.BrowserInfo.IsValid() &&
		a.ApplicationID != "" &&
		a.CallbackURI != "" &&
		a.Request != nil && a.Request.IsValid()
}

func (a *AuthRequest) MFALevel() MFALevel {
	return -1
	//PLANNED: check a.PossibleLOAs (and Prompt Login?)
}

func (a *AuthRequest) WithCurrentInfo(info *BrowserInfo) *AuthRequest {
	a.BrowserInfo = info
	return a
}

func (a *AuthRequest) SetUserInfo(userID, loginName, displayName, userOrgID string) {
	a.UserID = userID
	a.LoginName = loginName
	a.DisplayName = displayName
	a.UserOrgID = userOrgID
}

func (a *AuthRequest) GetScopeOrgPrimaryDomain() string {
	switch request := a.Request.(type) {
	case *AuthRequestOIDC:
		for _, scope := range request.Scopes {
			if strings.HasPrefix(scope, OrgDomainPrimaryScope) {
				return strings.TrimPrefix(scope, OrgDomainPrimaryScope)
			}
		}
	}
	return ""
}

func (a *AuthRequest) GetScopeProjectIDsForAud() []string {
	projectIDs := make([]string, 0)
	switch request := a.Request.(type) {
	case *AuthRequestOIDC:
		for _, scope := range request.Scopes {
			if strings.HasPrefix(scope, ProjectIDScope) && strings.HasSuffix(scope, AudSuffix) {
				projectIDs = append(projectIDs, strings.TrimSuffix(strings.TrimPrefix(scope, ProjectIDScope), AudSuffix))
			}
		}
	}
	return projectIDs
}

func (a *AuthRequest) AppendAudIfNotExisting(aud string) {
	for _, a := range a.Audience {
		if a == aud {
			return
		}
	}
	a.Audience = append(a.Audience, aud)
}
