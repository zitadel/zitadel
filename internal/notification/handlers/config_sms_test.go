package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotificationQueries_GetActiveSMSConfig(t *testing.T) {
	t.Parallel()

	const instId = "instance-1"
	const cryptAlg = "alg1"
	const keyId = "key1"
	const token = "twilio-token"

	encryptedToken := []byte("twilio-token-encrypted")

	var (
		errDecryptionFailed = errors.New("decryption failed")
		errDatabaseError    = errors.New("database error")
	)

	tt := []struct {
		name         string
		smsConfig    *query.SMSConfig
		smsConfigErr error
		setupCrypto  func(m *crypto.MockEncryptionAlgorithm)
		expected     *sms.Config
		expectedErr  error
	}{
		{
			name: "twilio config with token",
			smsConfig: &query.SMSConfig{
				ID:          "twilio1",
				Description: "twilio config with token",
				TwilioConfig: &query.Twilio{
					SID:              "sid-1",
					SenderNumber:     "+10005550100",
					VerifyServiceSID: "verify-sid-1",
					Token: &crypto.CryptoValue{
						Algorithm: cryptAlg,
						KeyID:     keyId,
						Crypted:   encryptedToken,
					},
				},
			},
			setupCrypto: func(m *crypto.MockEncryptionAlgorithm) {
				m.EXPECT().Algorithm().Return(cryptAlg)
				m.EXPECT().DecryptionKeyIDs().Return([]string{keyId})
				m.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return(token, nil)
			},
			expected: &sms.Config{
				ProviderConfig: &sms.Provider{
					ID:          "twilio1",
					Description: "twilio config with token",
				},
				TwilioConfig: &twilio.Config{
					SID:              "sid-1",
					Token:            token,
					SenderNumber:     "+10005550100",
					VerifyServiceSID: "verify-sid-1",
				},
			},
		},
		{
			name: "twilio config with nil token",
			smsConfig: &query.SMSConfig{
				ID:          "twilio2",
				Description: "twilio config with nil token",
				TwilioConfig: &query.Twilio{
					SID:          "sid-2",
					SenderNumber: "+10005550100",
				},
			},
			expectedErr: zerrors.ThrowNotFound(nil, "", ""),
		},
		{
			name: "twilio token decryption error",
			smsConfig: &query.SMSConfig{
				ID:          "twilio3",
				Description: "twilio token decryption error",
				TwilioConfig: &query.Twilio{
					SID:          "sid-3",
					SenderNumber: "+10005550100",
					Token: &crypto.CryptoValue{
						Algorithm: cryptAlg,
						KeyID:     keyId,
						Crypted:   encryptedToken,
					},
				},
			},
			setupCrypto: func(m *crypto.MockEncryptionAlgorithm) {
				m.EXPECT().Algorithm().Return(cryptAlg)
				m.EXPECT().DecryptionKeyIDs().Return([]string{keyId})
				m.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return("", errDecryptionFailed)
			},
			expectedErr: errDecryptionFailed,
		},
		{
			name: "http webhook config",
			smsConfig: &query.SMSConfig{
				ID:          "http1",
				Description: "http webhook config",
				HTTPConfig: &query.HTTP{
					Endpoint:   "https://example.com/sms",
					SigningKey: "my-signing-key",
				},
			},
			expected: &sms.Config{
				ProviderConfig: &sms.Provider{
					ID:          "http1",
					Description: "http webhook config",
				},
				WebhookConfig: &webhook.Config{
					CallURL:    "https://example.com/sms",
					Method:     http.MethodPost,
					Headers:    nil,
					SigningKey: "my-signing-key",
				},
			},
		},
		{
			name: "neither twilio nor http config set",
			smsConfig: &query.SMSConfig{
				ID:          "sms1",
				Description: "neither twilio nor http config set",
			},
			expectedErr: zerrors.ThrowNotFound(nil, "", ""),
		},
		{
			name:         "SMSProviderConfigActive query error",
			smsConfigErr: errDatabaseError,
			expectedErr:  errDatabaseError,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			ctx := authz.NewMockContext(instId, "org-1", "user-1")

			cryptAlgMock := crypto.NewMockEncryptionAlgorithm(ctrl)
			if tc.setupCrypto != nil {
				tc.setupCrypto(cryptAlgMock)
			}

			queryMock := mock.NewMockQueries(ctrl)
			queryMock.EXPECT().SMSProviderConfigActive(gomock.Any(), instId).Return(tc.smsConfig, tc.smsConfigErr)

			notificationQueries := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext domain", uint16(1234), false, "filepath", nil, nil, cryptAlgMock, nil)
			cfg, err := notificationQueries.GetActiveSMSConfig(ctx)

			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr != nil {
				return
			}

			assert.EqualValues(t, tc.expected, cfg)
		})
	}
}
