package oidc

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/zitadel/logging"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	http_utils "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
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
	list := make([]string, 0)
	if a.PasswordVerified {
		list = append(list, Password, PWD)
	}
	if len(a.MFAsVerified) > 0 {
		list = append(list, MFA)
		for _, mfa := range a.MFAsVerified {
			if amrMFA := AMRFromMFAType(mfa); amrMFA != "" {
				list = append(list, amrMFA)
			}
		}
	}
	return list
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

func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return ResponseModeToOIDC(a.oidc().ResponseMode)
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

func (a *AuthRequest) oidc() *domain.AuthRequestOIDC {
	return a.Request.(*domain.AuthRequestOIDC)
}

func AuthRequestFromBusiness(authReq *domain.AuthRequest) (_ *AuthRequest, err error) {
	if _, ok := authReq.Request.(*domain.AuthRequestOIDC); !ok {
		return nil, zerrors.ThrowInvalidArgument(nil, "OIDC-Haz7A", "auth request is not of type oidc")
	}
	return &AuthRequest{authReq}, nil
}

func CreateAuthRequestToBusiness(ctx context.Context, authReq *oidc.AuthRequest, userAgentID, userID string, audience []string) *domain.AuthRequest {
	return &domain.AuthRequest{
		CreationDate:        time.Now(),
		AgentID:             userAgentID,
		BrowserInfo:         ParseBrowserInfoFromContext(ctx),
		ApplicationID:       authReq.ClientID,
		CallbackURI:         authReq.RedirectURI,
		TransferState:       authReq.State,
		Prompt:              PromptToBusiness(authReq.Prompt),
		PossibleLOAs:        ACRValuesToBusiness(authReq.ACRValues),
		UiLocales:           UILocalesToBusiness(authReq.UILocales),
		LoginHint:           authReq.LoginHint,
		SelectedIDPConfigID: GetSelectedIDPIDFromScopes(authReq.Scopes),
		MaxAuthAge:          MaxAgeToBusiness(authReq.MaxAge),
		UserID:              userID,
		InstanceID:          authz.GetInstance(ctx).InstanceID(),
		Audience:            audience,
		Request: &domain.AuthRequestOIDC{
			Scopes:        authReq.Scopes,
			ResponseType:  ResponseTypeToBusiness(authReq.ResponseType),
			ResponseMode:  ResponseModeToBusiness(authReq.ResponseMode),
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

func PromptToBusiness(oidcPrompt []string) []domain.Prompt {
	prompts := make([]domain.Prompt, 0, len(oidcPrompt))
	for _, oidcPrompt := range oidcPrompt {
		switch oidcPrompt {
		case oidc.PromptNone:
			prompts = append(prompts, domain.PromptNone)
		case oidc.PromptLogin:
			prompts = append(prompts, domain.PromptLogin)
		case oidc.PromptConsent:
			prompts = append(prompts, domain.PromptConsent)
		case oidc.PromptSelectAccount:
			prompts = append(prompts, domain.PromptSelectAccount)
		case "create": //this prompt is not final yet, so not implemented in oidc lib
			prompts = append(prompts, domain.PromptCreate)
		}
	}
	return prompts
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

func GetSelectedIDPIDFromScopes(scopes oidc.SpaceDelimitedArray) string {
	for _, scope := range scopes {
		if strings.HasPrefix(scope, domain.SelectIDPScope) {
			return strings.TrimPrefix(scope, domain.SelectIDPScope)
		}
	}
	return ""
}

func MaxAgeToBusiness(maxAge *uint) *time.Duration {
	if maxAge == nil {
		return nil
	}
	dur := time.Duration(*maxAge) * time.Second
	return &dur
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

// ResponseModeToBusiness returns the OIDCResponseMode enum value from the domain package.
// An empty or invalid value defaults to unspecified.
func ResponseModeToBusiness(responseMode oidc.ResponseMode) domain.OIDCResponseMode {
	if responseMode == "" {
		return domain.OIDCResponseModeUnspecified
	}
	out, err := domain.OIDCResponseModeString(string(responseMode))
	logging.OnError(err).Debugln("invalid oidc response_mode, using default")
	return out
}

// ResponseModeToOIDC return the oidc string representation of the enum value from the domain package.
// When responseMode is `0 - unspecified`, an empty string is returned.
// This allows the oidc package to pick the appropriate response mode based on the response type.
func ResponseModeToOIDC(responseMode domain.OIDCResponseMode) oidc.ResponseMode {
	if responseMode == domain.OIDCResponseModeUnspecified || !responseMode.IsAOIDCResponseMode() {
		return ""
	}
	return oidc.ResponseMode(responseMode.String())
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
	case domain.MFATypeTOTP,
		domain.MFATypeOTPSMS,
		domain.MFATypeOTPEmail:
		return OTP
	case domain.MFATypeU2F,
		domain.MFATypeU2FUserVerification:
		return UserPresence
	default:
		return ""
	}
}

func RefreshTokenRequestFromBusiness(tokenView *model.RefreshTokenView) op.RefreshTokenRequest {
	return &RefreshTokenRequest{tokenView}
}

type RefreshTokenRequest struct {
	*model.RefreshTokenView
}

func (r *RefreshTokenRequest) GetAMR() []string {
	return r.AuthMethodsReferences
}

func (r *RefreshTokenRequest) GetAudience() []string {
	return r.Audience
}

func (r *RefreshTokenRequest) GetAuthTime() time.Time {
	return r.AuthTime
}

func (r *RefreshTokenRequest) GetClientID() string {
	return r.ClientID
}

func (r *RefreshTokenRequest) GetScopes() []string {
	return r.Scopes
}

func (r *RefreshTokenRequest) GetSubject() string {
	return r.UserID
}

func (r *RefreshTokenRequest) SetCurrentScopes(scopes []string) {
	r.Scopes = scopes
}

func (r *RefreshTokenRequest) GetActor() *oidc.ActorClaims {
	return actorDomainToClaims(r.Actor)
}
