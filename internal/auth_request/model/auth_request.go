package model

import (
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

	levelOfAssurance LevelOfAssurance
	UserID           string
	LoginName        string
	UserOrgID        string
	PossibleSteps    []NextStep
	PasswordVerified bool
	MfasVerified     []MfaType
	Audience         []string
	AuthTime         time.Time
	Code             string
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

func (a *AuthRequest) MfaLevel() MfaLevel {
	return -1
	//PLANNED: check a.PossibleLOAs (and Prompt Login?)
}

func (a *AuthRequest) WithCurrentInfo(info *BrowserInfo) *AuthRequest {
	a.BrowserInfo = info
	return a
}

func (a *AuthRequest) SetUserInfo(userID string, loginName string, userOrgID string) {
	a.UserID = userID
	a.LoginName = loginName
	a.UserOrgID = userOrgID
}
