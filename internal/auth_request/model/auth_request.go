package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type AuthRequest struct {
	es_models.ObjectRoot
	BrowserInfo       *BrowserInfo
	ApplicationID     string             //clientID
	CallbackURI       string             //redirectURi
	TransferState     string             //state //oidc only?
	Prompt            Prompt             //name?
	PossibleLOAs      []LevelOfAssurance //acr_values
	UiLocales         []string           //language.Tag?
	LoginHint         string
	PreselectedUserID string
	MaxAuthAge        uint32
	Request           Request

	levelOfAssurance      LevelOfAssurance //acr
	projectApplicationIDs []string         //aud?
	PossibleSteps         []NextStep
	//UserSession   *UserSession

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

func NewAuthRequest(agentID string, info *BrowserInfo, applicationID, callbackURI, transferState string,
	prompt Prompt, possibleLOAs []LevelOfAssurance, uiLocales []string, loginHint, preselectedUserID string, maxAuthAge uint32, request Request) *AuthRequest {
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
	a.PossibleSteps = append(a.PossibleSteps, step)
}

func (a *AuthRequest) MfaLevel() model.MfaLevel {
	return -1
	//TODO: check a.PossibleLOAs
}

func (a *AuthRequest) WithCurrentInfo(info *BrowserInfo) *AuthRequest {
	a.BrowserInfo = info
	return a
}
