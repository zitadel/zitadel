package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type AuthRequest struct {
	es_models.ObjectRoot
	BrowserInfo       *BrowserInfo
	ApplicationID     string   //clientID
	CallbackURI       string   //redirectURi
	TransferState     string   //state //oidc only?
	Prompt            Prompt   //name?
	PossibleLOAs      []string //acr_values
	UiLocales         []string //language.Tag?
	LoginHint         string
	PreselectedUserID string
	MaxAuthAge        uint32
	Request           Request

	levelOfAssurance      string   //acr
	projectApplicationIDs []string //aud?
	possibleSteps         []NextStep
	//UserSession   *UserSession

}

func NewAuthRequest(agentID string, info *BrowserInfo, applicationID, callbackURI, transferState string,
	prompt Prompt, possibleLOAs, uiLocales []string, loginHint, preselectedUserID string, maxAuthAge uint32, request Request) *AuthRequest {
	return &AuthRequest{
		ObjectRoot:        es_models.ObjectRoot{AggregateID: agentID},
		BrowserInfo:       info,
		ApplicationID:     applicationID,
		CallbackURI:       callbackURI,
		TransferState:     transferState,
		Prompt:            prompt,
		PossibleLOAs:      possibleLOAs,
		UiLocales:         uiLocales,
		LoginHint:         loginHint,
		PreselectedUserID: preselectedUserID,
		MaxAuthAge:        maxAuthAge,
		Request:           request,
	}
}

func (a *AuthRequest) IsValid() bool {
	return a.AggregateID != "" &&
		a.BrowserInfo != nil && a.BrowserInfo.IsValid() &&
		a.ApplicationID != "" &&
		a.CallbackURI != "" &&
		a.Request != nil && a.Request.IsValid()
}

func (a *AuthRequest) AddPossibleStep(step NextStep) {
	a.possibleSteps = append(a.possibleSteps, step)
}

type Prompt int32

const (
	PromptUnspecified Prompt = iota
	PromptNone
	PromptLogin
	PromptConsent
	PromptSelectAccount
)
