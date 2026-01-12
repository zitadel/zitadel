package sessionv2

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func Test_sessionToPb(t *testing.T) {
	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Minute)
	verifiedAt := updatedAt.Add(time.Minute)
	expiredAt := verifiedAt.Add(5 * time.Minute)

	s := &domain.Session{
		ID:         "session-1",
		InstanceID: "instance-1",
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
		Factors: domain.SessionFactors{
			&domain.SessionFactorUser{
				UserID:         "user-1",
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorPasskey{
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorPassword{
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorIdentityProviderIntent{
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorTOTP{
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorOTPSMS{
				LastVerifiedAt: verifiedAt,
			},
			&domain.SessionFactorOTPEmail{
				LastVerifiedAt: verifiedAt,
			},
		},
		Metadata: []domain.SessionMetadata{
			{
				Metadata: domain.Metadata{
					InstanceID: "instance-1",
					Key:        "key-1",
					Value:      []byte(`{"some":"value"}`),
				},
				SessionID: "session-1",
			},
			{
				Metadata: domain.Metadata{
					InstanceID: "instance-1",
					Key:        "key-2",
					Value:      []byte(`{"other":"value"}`),
				},
				SessionID: "session-1",
			},
		},
		UserID: "user-1",
		UserAgent: &domain.SessionUserAgent{
			FingerprintID: gu.Ptr("fingerprint-1"),
			Description:   gu.Ptr("ua-1"),
		},
		Expiration: expiredAt,
	}

	wantSession := &session.Session{
		Id:           "session-1",
		CreationDate: timestamppb.New(createdAt),
		ChangeDate:   timestamppb.New(updatedAt),
		Factors: &session.Factors{
			User: &session.UserFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
				Id:         "user-1",
			},
			Password: &session.PasswordFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
			WebAuthN: &session.WebAuthNFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
			Intent: &session.IntentFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
			Totp: &session.TOTPFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
			OtpSms: &session.OTPFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
			OtpEmail: &session.OTPFactor{
				VerifiedAt: timestamppb.New(verifiedAt),
			},
		},
		Metadata: map[string][]byte{
			"key-1": []byte(`{"some":"value"}`),
			"key-2": []byte(`{"other":"value"}`),
		},
		UserAgent: &session.UserAgent{
			FingerprintId: gu.Ptr("fingerprint-1"),
			Description:   gu.Ptr("ua-1"),
		},
		ExpirationDate: timestamppb.New(expiredAt),
	}

	got := sessionToPb(s)
	assert.Equal(t, wantSession, got)
}
