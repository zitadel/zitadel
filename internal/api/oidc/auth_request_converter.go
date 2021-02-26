package oidc

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"net"
	"time"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/oidc/pkg/op"
	"golang.org/x/text/language"

	http_utils "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/errors"
)

const (
	amrPassword     = "password"
	amrMFA          = "mfa"
	amrOTP          = "otp"
	amrUserPresence = "user"
)

type AuthRequest struct {
	*domain.AuthRequest
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
		if step.Type() == domain.NextStepRedirectToCallback {
			return true
		}
	}
	return false
}

func (a *AuthRequest) oidc() *domain.AuthRequestOIDC {
	return a.Request.(*domain.AuthRequestOIDC)
}

func AuthRequestFromBusiness(authReq *domain.AuthRequest) (_ op.AuthRequest, err error) {
	if _, ok := authReq.Request.(*domain.AuthRequestOIDC); !ok {
		return nil, errors.ThrowInvalidArgument(nil, "OIDC-Haz7A", "auth request is not of type oidc")
	}
	return &AuthRequest{authReq}, nil
}

func CreateAuthRequestToBusiness(ctx context.Context, authReq *oidc.AuthRequest, userAgentID, userID string) *domain.AuthRequest {
	return &domain.AuthRequest{
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
		Request: &domain.AuthRequestOIDC{
			Scopes:        authReq.Scopes,
			ResponseType:  ResponseTypeToBusiness(authReq.ResponseType),
			Nonce:         authReq.Nonce,
			CodeChallenge: CodeChallengeToBusiness(authReq.CodeChallenge, authReq.CodeChallengeMethod),
		},
	}
}

func ParseBrowserInfoFromContext(ctx context.Context) *domain.BrowserInfo {
	userAgent, acceptLang := HttpHeadersFromContext(ctx)
	ip := IpFromContext(ctx)
	return &domain.BrowserInfo{RemoteIP: ip, UserAgent: userAgent, AcceptLanguage: acceptLang}
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

func PromptToBusiness(prompt oidc.Prompt) domain.Prompt {
	switch prompt {
	case oidc.PromptNone:
		return domain.PromptNone
	case oidc.PromptLogin:
		return domain.PromptLogin
	case oidc.PromptConsent:
		return domain.PromptConsent
	case oidc.PromptSelectAccount:
		return domain.PromptSelectAccount
	default:
		return domain.PromptUnspecified
	}
}

func ACRValuesToBusiness(values []string) []domain.LevelOfAssurance {
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

func ResponseTypeToBusiness(responseType oidc.ResponseType) domain.OIDCResponseType {
	switch responseType {
	case oidc.ResponseTypeCode:
		return domain.OIDCResponseTypeCode
	case oidc.ResponseTypeIDTokenOnly:
		return domain.OIDCResponseTypeIDToken
	case oidc.ResponseTypeIDToken:
		return domain.OIDCResponseTypeIDTokenToken
	default:
		return domain.OIDCResponseTypeCode
	}
}

func ResponseTypeToOIDC(responseType domain.OIDCResponseType) oidc.ResponseType {
	switch responseType {
	case domain.OIDCResponseTypeCode:
		return oidc.ResponseTypeCode
	case domain.OIDCResponseTypeIDTokenToken:
		return oidc.ResponseTypeIDToken
	case domain.OIDCResponseTypeIDToken:
		return oidc.ResponseTypeIDTokenOnly
	default:
		return oidc.ResponseTypeCode
	}
}

func CodeChallengeToBusiness(challenge string, method oidc.CodeChallengeMethod) *domain.OIDCCodeChallenge {
	if challenge == "" {
		return nil
	}
	challengeMethod := domain.CodeChallengeMethodPlain
	if method == oidc.CodeChallengeMethodS256 {
		challengeMethod = domain.CodeChallengeMethodS256
	}
	return &domain.OIDCCodeChallenge{
		Challenge: challenge,
		Method:    challengeMethod,
	}
}

func CodeChallengeToOIDC(challenge *domain.OIDCCodeChallenge) *oidc.CodeChallenge {
	if challenge == nil {
		return nil
	}
	challengeMethod := oidc.CodeChallengeMethodPlain
	if challenge.Method == domain.CodeChallengeMethodS256 {
		challengeMethod = oidc.CodeChallengeMethodS256
	}
	return &oidc.CodeChallenge{
		Challenge: challenge.Challenge,
		Method:    challengeMethod,
	}
}

func AMRFromMFAType(mfaType domain.MFAType) string {
	switch mfaType {
	case domain.MFATypeOTP:
		return amrOTP
	case domain.MFATypeU2F,
		domain.MFATypeU2FUserVerification:
		return amrUserPresence
	default:
		return ""
	}
}
