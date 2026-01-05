package sessionv2

import (
	"time"

	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func sessionToPb(s *domain.Session) *session.Session {
	return &session.Session{
		Id:             s.ID,
		CreationDate:   timestamppb.New(s.CreatedAt),
		ChangeDate:     timestamppb.New(s.UpdatedAt),
		Factors:        factorsToPb(s.Factors),
		Metadata:       metadataToPb(s.Metadata),
		UserAgent:      userAgentToPb(s.UserAgent),
		ExpirationDate: expirationToPb(s.Expiration),
	}
}

func expirationToPb(expiration time.Time) *timestamppb.Timestamp {
	if expiration.IsZero() {
		return nil
	}
	return timestamppb.New(expiration)
}

func factorsToPb(factors domain.SessionFactors) *session.Factors {
	user := userFactorToPb(factors.GetUserFactor())
	if user == nil {
		return nil
	}
	return &session.Factors{
		User:     user,
		Password: passwordFactorToPb(factors.GetPasswordFactor()),
		WebAuthN: webAuthNFactorToPb(factors.GetPasskeyFactor()),
		Intent:   intentFactorToPb(factors.GetIDPIntentFactor()),
		Totp:     totpFactorToPb(factors.GetTOTPFactor()),
		OtpSms:   otpFactorToPb(factors.GetOTPSMSFactor().LastVerifiedAt),
		OtpEmail: otpFactorToPb(factors.GetOTPEmailFactor().LastVerifiedAt),
		//RecoveryCode: recoveryCodeFactorToPb(factors.GetRecoveryCodeFactor()), // (@grvijayan) todo
	}
}

func userFactorToPb(factor *domain.SessionFactorUser) *session.UserFactor {
	if factor.UserID == "" || factor.LastVerifiedAt.IsZero() {
		return nil
	}
	return &session.UserFactor{
		VerifiedAt: timestamppb.New(factor.LastVerifiedAt),
		Id:         factor.UserID,
		//LoginName:      factor.LoginName,
		//DisplayName:    factor.DisplayName,
		//OrganizationId: factor.ResourceOwner,
	}
}

func passwordFactorToPb(factor *domain.SessionFactorPassword) *session.PasswordFactor {
	if factor.LastVerifiedAt.IsZero() {
		return nil
	}
	return &session.PasswordFactor{
		VerifiedAt: timestamppb.New(factor.LastVerifiedAt),
	}
}

func webAuthNFactorToPb(factor *domain.SessionFactorPasskey) *session.WebAuthNFactor {
	if factor.LastVerifiedAt.IsZero() {
		return nil
	}
	return &session.WebAuthNFactor{
		VerifiedAt:   timestamppb.New(factor.LastVerifiedAt),
		UserVerified: factor.UserVerified,
	}
}

func intentFactorToPb(factor *domain.SessionFactorIdentityProviderIntent) *session.IntentFactor {
	if factor.LastVerifiedAt.IsZero() {
		return nil
	}
	return &session.IntentFactor{
		VerifiedAt: timestamppb.New(factor.LastVerifiedAt),
	}
}

func totpFactorToPb(factor *domain.SessionFactorTOTP) *session.TOTPFactor {
	if factor.LastVerifiedAt.IsZero() {
		return nil
	}
	return &session.TOTPFactor{
		VerifiedAt: timestamppb.New(factor.LastVerifiedAt),
	}
}

func otpFactorToPb(lastVerifiedAt time.Time) *session.OTPFactor {
	if lastVerifiedAt.IsZero() {
		return nil
	}
	return &session.OTPFactor{
		VerifiedAt: timestamppb.New(lastVerifiedAt),
	}
}

func metadataToPb(metadata []domain.SessionMetadata) map[string][]byte {
	if len(metadata) == 0 {
		return nil
	}
	result := make(map[string][]byte, len(metadata))
	for _, md := range metadata {
		result[md.Key] = md.Value
	}
	return result
}

func userAgentToPb(ua *domain.SessionUserAgent) *session.UserAgent {
	if ua.FingerprintID == nil &&
		len(ua.IP) == 0 &&
		ua.Description == nil &&
		ua.Header == nil {
		return nil
	}
	out := &session.UserAgent{
		FingerprintId: ua.FingerprintID,
		Description:   ua.Description,
	}
	if ua.IP != nil {
		out.Ip = gu.Ptr(ua.IP.String())
	}
	if ua.Header == nil {
		return out
	}
	out.Header = make(map[string]*session.UserAgent_HeaderValues, len(ua.Header))
	for k, v := range ua.Header {
		out.Header[k] = &session.UserAgent_HeaderValues{
			Values: v,
		}
	}
	return out
}
