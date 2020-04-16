package model

import (
	es_models "github.com/caos/zitadel/internal/eventstore/models"
)

type AuthSession struct {
	es_models.ObjectRoot
	SessionID             string
	Type                  AuthSessionType
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
	OIDC                  *AuthSessionOIDC
	UserSession           UserSession
	PossibleSteps         []*NextStep
}

func NewAuthSession(agentID, sessionID string, sessionType AuthSessionType, info *BrowserInfo,
	applicationID, callbackURI, transferState string, prompt Prompt, requestedPossibleLOAs, requestedUiLocales []string,
	loginHint, preselectedUserID string, maxAuthAge uint32, oidc *AuthSessionOIDC) *AuthSession {
	return &AuthSession{
		ObjectRoot:            es_models.ObjectRoot{ID: agentID},
		SessionID:             sessionID,
		Type:                  sessionType,
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
		OIDC:                  oidc,
	}
}

func (a *AuthSession) IsValid() bool {
	return a.ID != "" &&
		a.SessionID != "" &&
		a.BrowserInfo != nil && a.BrowserInfo.IsValid() &&
		a.ApplicationID != "" &&
		a.CallbackURI != "" &&
		true //todo oidc?
}

type AuthSessionType int32

const (
	AuthSessionTypeOIDC AuthSessionType = iota
	AuthSessionTypeSAML
)

type Prompt int32

const (
	PromptUnspecified Prompt = iota
	PromptNone
	PromptLogin
	PromptConsent
	PromptSelectAccount
)

type AuthSessionOIDC struct {
	Scopes        []string
	ResponseTypes OIDCResponseType
	Nonce         string
	CodeChallenge *OIDCCodeChallenge
}

type OIDCResponseType int32

const (
	CODE OIDCResponseType = iota
	ID_TOKEN
	ID_TOKEN_TOKEN
)

type OIDCCodeChallenge struct {
	Challenge string
	Method    OIDCCodeChallengeMethod
}

type OIDCCodeChallengeMethod int32

const (
	CodeChallengeMethodPlain OIDCCodeChallengeMethod = iota
	CodeChallengeMethodS256
)

type NextStep interface {
	Type() NextStepType
}

type NextStepType int32

const (
	NextStepUnspecified NextStepType = iota
	NextStepLogin
	NextStepUserSelection
	NextStepPassword
	NextStepChangePassword
	NextStepInitPassword
	NextStepVerifyEmail
	NextStepMfaPrompt
	NextStepMfaVerify
	NextStepRedirectToCallback
)

type LoginStep struct {
	ErrMsg string
}

func (s *LoginStep) Type() NextStepType {
	return NextStepLogin
}

type SelectUserStep struct {
	Users []UserSelection
}

func (s *SelectUserStep) Type() NextStepType {
	return NextStepUserSelection
}

type UserSelection struct {
	SessionID        string
	UserID           string
	UserName         string
	UserSessionState UserSessionState
}

type PasswordStep struct {
	ErrMsg       string
	FailureCount uint16
}

func (s *PasswordStep) Type() NextStepType {
	return NextStepPassword
}

type ChangePasswordStep struct {
}

func (s *ChangePasswordStep) Type() NextStepType {
	return NextStepChangePassword
}

type InitPasswordStep struct {
}

func (s *InitPasswordStep) Type() NextStepType {
	return NextStepInitPassword
}

type VerifyEMailStep struct {
}

func (s *VerifyEMailStep) Type() NextStepType {
	return NextStepVerifyEmail
}

type MfaPromptStep struct {
	Required     bool
	MfaProviders []MfaType
}

func (s *MfaPromptStep) Type() NextStepType {
	return NextStepMfaPrompt
}

type MfaVerificationStep struct {
	ErrMsg       string
	FailureCount uint16
	MfaProviders []MfaType
}

func (s *MfaVerificationStep) Type() NextStepType {
	return NextStepMfaVerify
}

type RedirectToCallbackStep struct {
}

func (s *RedirectToCallbackStep) Type() NextStepType {
	return NextStepRedirectToCallback
}
