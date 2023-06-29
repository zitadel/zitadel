package oidc

import (
	"strings"
	"time"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/errors"
)

const IDPrefix = "V2_"

func StripIDPrefix(authRequestID string) (string, error) {
	after, found := strings.CutPrefix(authRequestID, IDPrefix)
	if !found {
		return "", errors.ThrowInvalidArgumentf(nil, "OIDC-Aumu8", "auth_request_id wrong version, missing %s prefix", IDPrefix)
	}
	return after, nil
}

type AuthRequestV2 struct {
	*command.AuthRequest
}

func (a *AuthRequestV2) GetID() string {
	return IDPrefix + a.ID
}

func (a *AuthRequestV2) GetACR() string {
	return "" //PLANNED: impl
}

func (a *AuthRequestV2) GetAMR() []string {
	//TODO: get from linked session?
	return nil
}

func (a *AuthRequestV2) GetAudience() []string {
	return a.Audience
}

func (a *AuthRequestV2) GetAuthTime() time.Time {
	return time.Time{} //TODO: get from linked session?
}

func (a *AuthRequestV2) GetClientID() string {
	return a.ClientID
}

func (a *AuthRequestV2) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.CodeChallenge)
}

func (a *AuthRequestV2) GetNonce() string {
	return a.Nonce
}

func (a *AuthRequestV2) GetRedirectURI() string {
	return a.RedirectURI
}

func (a *AuthRequestV2) GetResponseType() oidc.ResponseType {
	return ResponseTypeToOIDC(a.ResponseType)
}

func (a *AuthRequestV2) GetResponseMode() oidc.ResponseMode {
	return ""
}

func (a *AuthRequestV2) GetScopes() []string {
	return a.Scope
}

func (a *AuthRequestV2) GetState() string {
	return a.State
}

func (a *AuthRequestV2) GetSubject() string {
	return "" //TODO: get from linked session?
}

func (a *AuthRequestV2) Done() bool {
	return false //TODO: get from linked session?
}
