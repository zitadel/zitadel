package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type AuthSession struct {
	es_models.ObjectRoot
	SessionID string
	//Type                  AuthSessionType
	BrowserInfo           *BrowserInfo
	ApplicationID         string   //clientID
	CallbackURI           string   //redirectURi
	TransferState         string   //state //oidc only?
	Prompt                Prompt   //name?
	LevelOfAssurance      string   //acr
	RequestedPossibleLOAs []string //acr_values
	RequestedUiLocales    []string //language.Tag?
	LoginHint             string
	PreselectedUserID     string
	MaxAuthAge            uint32
	ProjectApplicationIDs []string //aud?
	//OIDC                  *AuthSessionOIDC
	Request       Request
	UserSession   *UserSession
	PossibleSteps []*NextStep
}

func NewAuthSession(agentID, sessionID string, info *BrowserInfo,
	applicationID, callbackURI, transferState string, prompt Prompt, requestedPossibleLOAs, requestedUiLocales []string,
	loginHint, preselectedUserID string, maxAuthAge uint32, request Request) *AuthSession {
	return &AuthSession{
		ObjectRoot: es_models.ObjectRoot{AggregateID: agentID},
		SessionID:  sessionID,
		//Type:                  sessionType,
		BrowserInfo:           info,
		ApplicationID:         applicationID,
		CallbackURI:           callbackURI,
		TransferState:         transferState,
		Prompt:                prompt,
		RequestedPossibleLOAs: requestedPossibleLOAs,
		RequestedUiLocales:    requestedUiLocales,
		LoginHint:             loginHint,
		PreselectedUserID:     preselectedUserID,
		MaxAuthAge:            maxAuthAge,
		Request:               request,
	}
}

func (a *AuthSession) IsValid() bool {
	return a.AggregateID != "" &&
		a.SessionID != "" &&
		a.BrowserInfo != nil && a.BrowserInfo.IsValid() &&
		a.ApplicationID != "" &&
		a.CallbackURI != "" &&
		a.Request != nil && a.Request.IsValid()
}

type Prompt int32

const (
	PromptUnspecified Prompt = iota
	PromptNone
	PromptLogin
	PromptConsent
	PromptSelectAccount
)

type OIDCResponseType int32

const (
	CODE OIDCResponseType = iota
	ID_TOKEN
	ID_TOKEN_TOKEN
)
