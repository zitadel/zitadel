package oidc

import (
	"context"
	"net"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"golang.org/x/text/language"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/auth_request/model"
	"github.com/caos/zitadel/internal/errors"
)

const (
	amrPassword     = "password"
	amrMFA          = "mfa"
	amrOTP          = "otp"
	amrUserPresence = "user"
)

type AuthRequest struct {
	*model.AuthRequest
}

func (a *AuthRequest) GetID() string {
	return a.ID
}

func (a *AuthRequest) GetACR() string {
	// return a.
	return "" //PLANNED: impl
}

func (a *AuthRequest) GetAMR() []string {
	amr := make([]string, 0)
	if a.PasswordVerified {
		amr = append(amr, amrPassword)
	}
	if len(a.MFAsVerified) > 0 {
		amr = append(amr, amrMFA)
		for _, mfa := range a.MFAsVerified {
			if amrMFA := AMRFromMFAType(mfa); amrMFA != "" {
				amr = append(amr, amrMFA)
			}
		}
	}
	return amr
}

func (a *AuthRequest) GetAudience() []string {
	return a.Audience
}

func (a *AuthRequest) GetAuthTime() time.Time {
	return a.AuthTime
}

func (a *AuthRequest) GetClientID() string {
	return a.ApplicationID
}

func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return CodeChallengeToOIDC(a.oidc().CodeChallenge)
}

func (a *AuthRequest) GetNonce() string {
	return a.oidc().Nonce
}

func (a *AuthRequest) GetRedirectURI() string {
	return a.CallbackURI
}

func (a *AuthRequest) GetResponseType() oidc.ResponseType {
	return ResponseTypeToOIDC(a.oidc().ResponseType)
}

func (a *AuthRequest) GetScopes() []string {
	return a.oidc().Scopes
}

func (a *AuthRequest) GetState() string {
	return a.TransferState
}

func (a *AuthRequest) GetSubject() string {
	return a.UserID
}

func (a *AuthRequest) Done() bool {
	for _, step := range a.PossibleSteps {
		if step.Type() == model.NextStepRedirectToCallback {
			return true
		}
	}
	return false
}

func (a *AuthRequest) oidc() *model.AuthRequestOIDC {
	return a.Request.(*model.AuthRequestOIDC)
}

func AuthRequestFromBusiness(authReq *model.AuthRequest) (_ op.AuthRequest, err error) {
	if _, ok := authReq.Request.(*model.AuthRequestOIDC); !ok {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-Haz7A", "auth request is not of type oidc")
	}
	return &AuthRequest{authReq}, nil
}

func CreateAuthRequestToBusiness(ctx context.Context, authReq *oidc.AuthRequest, userAgentID, userID string) *model.AuthRequest {
	return &model.AuthRequest{
		AgentID:       userAgentID,
		BrowserInfo:   ParseBrowserInfoFromContext(ctx),
		ApplicationID: authReq.ClientID,
		CallbackURI:   authReq.RedirectURI,
		TransferState: authReq.State,
		Prompt:        PromptToBusiness(authReq.Prompt),
		PossibleLOAs:  ACRValuesToBusiness(authReq.ACRValues),
		UiLocales:     UILocalesToBusiness(authReq.UILocales),
		LoginHint:     authReq.LoginHint,
		MaxAuthAge:    authReq.MaxAge,
		UserID:        userID,
		Request: &model.AuthRequestOIDC{
			Scopes:        authReq.Scopes,
			ResponseType:  ResponseTypeToBusiness(authReq.ResponseType),
			Nonce:         authReq.Nonce,
			CodeChallenge: CodeChallengeToBusiness(authReq.CodeChallenge, authReq.CodeChallengeMethod),
		},
	}
}

func ParseBrowserInfoFromContext(ctx context.Context) *model.BrowserInfo {
	userAgent, acceptLang := HttpHeadersFromContext(ctx)
	ip := IpFromContext(ctx)
	return &model.BrowserInfo{RemoteIP: ip, UserAgent: userAgent, AcceptLanguage: acceptLang}
}

func HttpHeadersFromContext(ctx context.Context) (userAgent, acceptLang string) {
	ctxHeaders, ok := http_utils.HeadersFromCtx(ctx)
	if !ok {
		return
	}
	if agents, ok := ctxHeaders[http_utils.UserAgentHeader]; ok {
		userAgent = agents[0]
	}
	if langs, ok := ctxHeaders[http_utils.AcceptLanguage]; ok {
		acceptLang = langs[0]
	}
	return userAgent, acceptLang
}

func IpFromContext(ctx context.Context) net.IP {
	ipString := http_utils.RemoteIPFromCtx(ctx)
	if ipString == "" {
		return nil
	}
	return net.ParseIP(ipString)
}

func PromptToBusiness(prompt oidc.Prompt) model.Prompt {
	switch prompt {
	case oidc.PromptNone:
		return model.PromptNone
	case oidc.PromptLogin:
		return model.PromptLogin
	case oidc.PromptConsent:
		return model.PromptConsent
	case oidc.PromptSelectAccount:
		return model.PromptSelectAccount
	default:
		return model.PromptUnspecified
	}
}

func ACRValuesToBusiness(values []string) []model.LevelOfAssurance {
	return nil
}

func UILocalesToBusiness(tags []language.Tag) []string {
	if tags == nil {
		return nil
	}
	locales := make([]string, len(tags))
	for i, tag := range tags {
		locales[i] = tag.String()
	}
	return locales
}

func ResponseTypeToBusiness(responseType oidc.ResponseType) model.OIDCResponseType {
	switch responseType {
	case oidc.ResponseTypeCode:
		return model.OIDCResponseTypeCode
	case oidc.ResponseTypeIDTokenOnly:
		return model.OIDCResponseTypeIdToken
	case oidc.ResponseTypeIDToken:
		return model.OIDCResponseTypeIdTokenToken
	default:
		return model.OIDCResponseTypeCode
	}
}

func ResponseTypeToOIDC(responseType model.OIDCResponseType) oidc.ResponseType {
	switch responseType {
	case model.OIDCResponseTypeCode:
		return oidc.ResponseTypeCode
	case model.OIDCResponseTypeIdTokenToken:
		return oidc.ResponseTypeIDToken
	case model.OIDCResponseTypeIdToken:
		return oidc.ResponseTypeIDTokenOnly
	default:
		return oidc.ResponseTypeCode
	}
}

func CodeChallengeToBusiness(challenge string, method oidc.CodeChallengeMethod) *model.OIDCCodeChallenge {
	if challenge == "" {
		return nil
	}
	challengeMethod := model.CodeChallengeMethodPlain
	if method == oidc.CodeChallengeMethodS256 {
		challengeMethod = model.CodeChallengeMethodS256
	}
	return &model.OIDCCodeChallenge{
		Challenge: challenge,
		Method:    challengeMethod,
	}
}

func CodeChallengeToOIDC(challenge *model.OIDCCodeChallenge) *oidc.CodeChallenge {
	if challenge == nil {
		return nil
	}
	challengeMethod := oidc.CodeChallengeMethodPlain
	if challenge.Method == model.CodeChallengeMethodS256 {
		challengeMethod = oidc.CodeChallengeMethodS256
	}
	return &oidc.CodeChallenge{
		Challenge: challenge.Challenge,
		Method:    challengeMethod,
	}
}

func AMRFromMFAType(mfaType model.MFAType) string {
	switch mfaType {
	case model.MFATypeOTP:
		return amrOTP
	case model.MFATypeU2F,
		model.MFATypeU2FUserVerification:
		return amrUserPresence
	default:
		return ""
	}
}
