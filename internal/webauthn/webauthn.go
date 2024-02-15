package webauthn

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Config struct {
	DisplayName    string
	ExternalSecure bool
}

type webUser struct {
	*domain.Human
	accountName string
	credentials []webauthn.Credential
}

func (u *webUser) WebAuthnID() []byte {
	return []byte(u.AggregateID)
}

func (u *webUser) WebAuthnName() string {
	if u.accountName != "" {
		return u.accountName
	}
	return u.GetUsername()
}

func (u *webUser) WebAuthnDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.GetUsername()
}

func (u *webUser) WebAuthnIcon() string {
	return ""
}

func (u *webUser) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func (w *Config) BeginRegistration(ctx context.Context, user *domain.Human, accountName string, authType domain.AuthenticatorAttachment, userVerification domain.UserVerificationRequirement, rpID string, webAuthNs ...*domain.WebAuthNToken) (*domain.WebAuthNToken, error) {
	webAuthNServer, err := w.serverFromContext(ctx, rpID, "")
	if err != nil {
		return nil, err
	}
	creds := WebAuthNsToCredentials(webAuthNs, rpID)
	existing := make([]protocol.CredentialDescriptor, len(creds))
	for i, cred := range creds {
		existing[i] = protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		}
	}
	credentialOptions, sessionData, err := webAuthNServer.BeginRegistration(
		&webUser{
			Human:       user,
			accountName: accountName,
			credentials: creds,
		},
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			UserVerification:        UserVerificationFromDomain(userVerification),
			AuthenticatorAttachment: AuthenticatorAttachmentFromDomain(authType),
		}),
		webauthn.WithConveyancePreference(protocol.PreferNoAttestation),
		webauthn.WithExclusions(existing),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "WEBAU-bM8sd", "Errors.User.WebAuthN.BeginRegisterFailed")
	}
	cred, err := json.Marshal(credentialOptions)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "WEBAU-D7cus", "Errors.User.WebAuthN.MarshalError")
	}
	return &domain.WebAuthNToken{
		Challenge:              sessionData.Challenge,
		CredentialCreationData: cred,
		AllowedCredentialIDs:   sessionData.AllowedCredentialIDs,
		UserVerification:       UserVerificationToDomain(sessionData.UserVerification),
		RPID:                   webAuthNServer.Config.RPID,
	}, nil
}

func (w *Config) FinishRegistration(ctx context.Context, user *domain.Human, webAuthN *domain.WebAuthNToken, tokenName string, credData []byte, isLoginUI bool) (*domain.WebAuthNToken, error) {
	if webAuthN == nil {
		return nil, zerrors.ThrowInternal(nil, "WEBAU-5M9so", "Errors.User.WebAuthN.NotFound")
	}
	credentialData, err := protocol.ParseCredentialCreationResponseBody(bytes.NewReader(credData))
	if err != nil {
		logging.WithFields("error", tryExtractProtocolErrMsg(err)).Debug("webauthn credential could not be parsed")
		return nil, zerrors.ThrowInternal(err, "WEBAU-sEr8c", "Errors.User.WebAuthN.ErrorOnParseCredential")
	}
	sessionData := WebAuthNToSessionData(webAuthN)
	webAuthNServer, err := w.serverFromContext(ctx, webAuthN.RPID, credentialData.Response.CollectedClientData.Origin)
	if err != nil {
		return nil, err
	}
	credential, err := webAuthNServer.CreateCredential(
		&webUser{
			Human: user,
		},
		sessionData,
		credentialData)
	if err != nil {
		logging.WithFields("error", tryExtractProtocolErrMsg(err)).Debug("webauthn credential could not be created")
		return nil, zerrors.ThrowInternal(err, "WEBAU-3Vb9s", "Errors.User.WebAuthN.CreateCredentialFailed")
	}

	webAuthN.KeyID = credential.ID
	webAuthN.PublicKey = credential.PublicKey
	webAuthN.AttestationType = credential.AttestationType
	webAuthN.AAGUID = credential.Authenticator.AAGUID
	webAuthN.SignCount = credential.Authenticator.SignCount
	webAuthN.WebAuthNTokenName = tokenName
	webAuthN.RPID = webAuthNServer.Config.RPID
	return webAuthN, nil
}

func (w *Config) BeginLogin(ctx context.Context, user *domain.Human, userVerification domain.UserVerificationRequirement, rpID string, webAuthNs ...*domain.WebAuthNToken) (*domain.WebAuthNLogin, error) {
	webAuthNServer, err := w.serverFromContext(ctx, rpID, "")
	if err != nil {
		return nil, err
	}
	assertion, sessionData, err := webAuthNServer.BeginLogin(&webUser{
		Human:       user,
		credentials: WebAuthNsToCredentials(webAuthNs, rpID),
	}, webauthn.WithUserVerification(UserVerificationFromDomain(userVerification)))
	if err != nil {
		logging.WithFields("error", tryExtractProtocolErrMsg(err)).Debug("webauthn login could not be started")
		return nil, zerrors.ThrowInternal(err, "WEBAU-4G8sw", "Errors.User.WebAuthN.BeginLoginFailed")
	}
	cred, err := json.Marshal(assertion)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "WEBAU-2M0s9", "Errors.User.WebAuthN.MarshalError")
	}
	return &domain.WebAuthNLogin{
		Challenge:               sessionData.Challenge,
		CredentialAssertionData: cred,
		AllowedCredentialIDs:    sessionData.AllowedCredentialIDs,
		UserVerification:        userVerification,
		RPID:                    webAuthNServer.Config.RPID,
	}, nil
}

func (w *Config) FinishLogin(ctx context.Context, user *domain.Human, webAuthN *domain.WebAuthNLogin, credData []byte, webAuthNs ...*domain.WebAuthNToken) (*webauthn.Credential, error) {
	assertionData, err := protocol.ParseCredentialRequestResponseBody(bytes.NewReader(credData))
	if err != nil {
		logging.WithFields("error", tryExtractProtocolErrMsg(err)).Debug("webauthn assertion could not be parsed")
		return nil, zerrors.ThrowInternal(err, "WEBAU-ADgv4", "Errors.User.WebAuthN.ValidateLoginFailed")
	}
	webUser := &webUser{
		Human:       user,
		credentials: WebAuthNsToCredentials(webAuthNs, webAuthN.RPID),
	}
	webAuthNServer, err := w.serverFromContext(ctx, webAuthN.RPID, assertionData.Response.CollectedClientData.Origin)
	if err != nil {
		return nil, err
	}
	credential, err := webAuthNServer.ValidateLogin(webUser, WebAuthNLoginToSessionData(webAuthN), assertionData)
	if err != nil {
		logging.WithFields("error", tryExtractProtocolErrMsg(err)).Debug("webauthn assertion failed")
		return nil, zerrors.ThrowInternal(err, "WEBAU-3M9si", "Errors.User.WebAuthN.ValidateLoginFailed")
	}

	if credential.Authenticator.CloneWarning {
		return credential, zerrors.ThrowInternal(nil, "WEBAU-4M90s", "Errors.User.WebAuthN.CloneWarning")
	}
	return credential, nil
}

func (w *Config) serverFromContext(ctx context.Context, id, origin string) (*webauthn.WebAuthn, error) {
	config := w.config(id, origin)
	if id == "" {
		config = w.configFromContext(ctx)
	}
	webAuthn, err := webauthn.New(config)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "WEBAU-UX9ta", "Errors.User.WebAuthN.ServerConfig")
	}
	return webAuthn, nil
}

func (w *Config) configFromContext(ctx context.Context) *webauthn.Config {
	instance := authz.GetInstance(ctx)
	return &webauthn.Config{
		RPDisplayName: w.DisplayName,
		RPID:          instance.RequestedDomain(),
		RPOrigins:     []string{http.BuildOrigin(instance.RequestedHost(), w.ExternalSecure)},
	}
}

func (w *Config) config(id, origin string) *webauthn.Config {
	return &webauthn.Config{
		RPDisplayName: w.DisplayName,
		RPID:          id,
		RPOrigins:     []string{origin},
	}
}

func tryExtractProtocolErrMsg(err error) string {
	var e *protocol.Error
	if errors.As(err, &e) {
		return e.Details + ": " + e.DevInfo
	}
	return e.Error()
}
