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
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestNotificationQueries_GetActiveEmailConfig(t *testing.T) {
	t.Parallel()
	const instId = "instance-1"
	const cryptAlg = "alg1"
	const keyId = "key1"
	const pwd = "password"
	const clientSecret = "client-secret"

	encryptedPwd := []byte("password-encrypted")
	encryptedClientSecret := []byte("client-secret-encrypted")

	var (
		errDecryptionFailed = errors.New("decryption failed")
		errDatabaseError    = errors.New("database error")
	)

	tt := []struct {
		name          string
		smtpConfig    *query.SMTPConfig
		smtpConfigErr error
		setupCrypto   func(m *crypto.MockEncryptionAlgorithm)
		expected      *email.Config
		expectedErr   error
	}{
		{
			name: "smtp config with plain auth (user and password)",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp1",
				Description: "smtp config with plain auth (user and password)",
				SMTPConfig: &query.SMTP{
					TLS:            true,
					SenderAddress:  "sender@example.com",
					SenderName:     "sender",
					ReplyToAddress: "reply@example.com",
					Host:           "mail.com",
					User:           "mail-user",
					PlainAuth: &query.PlainAuth{
						Password: &crypto.CryptoValue{
							CryptoType: 0,
							Algorithm:  cryptAlg,
							KeyID:      keyId,
							Crypted:    encryptedPwd,
						},
					},
				},
			},
			setupCrypto: func(m *crypto.MockEncryptionAlgorithm) {
				m.EXPECT().Algorithm().Return(cryptAlg)
				m.EXPECT().DecryptionKeyIDs().Return([]string{keyId})
				m.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return(pwd, nil)
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp1",
					Description: "smtp config with plain auth (user and password)",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
						PlainAuth: &smtp.PlainAuthConfig{
							User:     "mail-user",
							Password: pwd,
						},
					},
					Tls:            true,
					From:           "sender@example.com",
					FromName:       "sender",
					ReplyToAddress: "reply@example.com",
				},
			},
		},
		{
			name: "smtp config with user but no plain auth config (fallback to plain auth without password)",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp2",
				Description: "smtp config with user but no plain auth config (fallback to plain auth without password)",
				SMTPConfig: &query.SMTP{
					TLS:            true,
					SenderAddress:  "sender@example.com",
					SenderName:     "sender",
					ReplyToAddress: "reply@example.com",
					Host:           "mail.com",
					User:           "mail-user",
				},
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp2",
					Description: "smtp config with user but no plain auth config (fallback to plain auth without password)",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
						PlainAuth: &smtp.PlainAuthConfig{
							User: "mail-user",
						},
					},
					Tls:            true,
					From:           "sender@example.com",
					FromName:       "sender",
					ReplyToAddress: "reply@example.com",
				},
			},
		},
		{
			name: "smtp config without authentication",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp3",
				Description: "smtp config without authentication",
				SMTPConfig: &query.SMTP{
					TLS:            true,
					SenderAddress:  "sender@example.com",
					SenderName:     "sender",
					ReplyToAddress: "reply@example.com",
					Host:           "mail.com",
				},
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp3",
					Description: "smtp config without authentication",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
					},
					Tls:            true,
					From:           "sender@example.com",
					FromName:       "sender",
					ReplyToAddress: "reply@example.com",
				},
			},
		},
		{
			name: "smtp config with xoauth2 auth and client credentials",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp4",
				Description: "smtp config with xoauth2 auth and client credentials",
				SMTPConfig: &query.SMTP{
					SenderAddress: "sender@example.com",
					SenderName:    "sender",
					Host:          "mail.com",
					User:          "mail-user",
					XOAuth2Auth: &query.XOAuth2Auth{
						TokenEndpoint: "https://token.example.com/token",
						Scopes:        []string{"https://mail.example.com/.default"},
						ClientCredentials: &query.XOAuthClientCredentials{
							ClientId: "client-id",
							ClientSecret: &crypto.CryptoValue{
								Algorithm: cryptAlg,
								KeyID:     keyId,
								Crypted:   encryptedClientSecret,
							},
						},
					},
				},
			},
			setupCrypto: func(m *crypto.MockEncryptionAlgorithm) {
				m.EXPECT().Algorithm().Return(cryptAlg)
				m.EXPECT().DecryptionKeyIDs().Return([]string{keyId})
				m.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return(clientSecret, nil)
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp4",
					Description: "smtp config with xoauth2 auth and client credentials",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
						XOAuth2Auth: &smtp.XOAuth2AuthConfig{
							User:          "mail-user",
							TokenEndpoint: "https://token.example.com/token",
							Scopes:        []string{"https://mail.example.com/.default"},
							ClientCredentialsAuth: &smtp.OAuth2ClientCredentials{
								ClientId:     "client-id",
								ClientSecret: clientSecret,
							},
						},
					},
					From:     "sender@example.com",
					FromName: "sender",
				},
			},
		},
		{
			name: "smtp config with xoauth2 auth without client credentials",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp5",
				Description: "smtp config with xoauth2 auth without client credentials",
				SMTPConfig: &query.SMTP{
					SenderAddress: "sender@example.com",
					SenderName:    "sender",
					Host:          "mail.com",
					User:          "mail-user",
					XOAuth2Auth: &query.XOAuth2Auth{
						TokenEndpoint: "https://token.example.com/token",
						Scopes:        []string{"scope1"},
					},
				},
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp5",
					Description: "smtp config with xoauth2 auth without client credentials",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
						XOAuth2Auth: &smtp.XOAuth2AuthConfig{
							User:          "mail-user",
							TokenEndpoint: "https://token.example.com/token",
							Scopes:        []string{"scope1"},
						},
					},
					From:     "sender@example.com",
					FromName: "sender",
				},
			},
		},
		{
			name: "plain auth password decryption error",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp6",
				Description: "plain auth password decryption error",
				SMTPConfig: &query.SMTP{
					Host: "mail.com",
					User: "mail-user",
					PlainAuth: &query.PlainAuth{
						Password: &crypto.CryptoValue{
							Algorithm: cryptAlg,
							KeyID:     keyId,
							Crypted:   encryptedPwd,
						},
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
			name: "xoauth2 client credentials decryption error",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp7",
				Description: "xoauth2 client credentials decryption error",
				SMTPConfig: &query.SMTP{
					Host: "mail.com",
					User: "mail-user",
					XOAuth2Auth: &query.XOAuth2Auth{
						TokenEndpoint: "https://token.example.com/token",
						ClientCredentials: &query.XOAuthClientCredentials{
							ClientId: "client-id",
							ClientSecret: &crypto.CryptoValue{
								CryptoType: 0,
								Algorithm:  cryptAlg,
								KeyID:      keyId,
								Crypted:    encryptedClientSecret,
							},
						},
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
			smtpConfig: &query.SMTPConfig{
				ID:          "http1",
				Description: "http webhook config",
				HTTPConfig: &query.HTTP{
					Endpoint:   "https://example.com/notify",
					SigningKey: "my-signing-key",
				},
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "http1",
					Description: "http webhook config",
				},
				WebhookConfig: &webhook.Config{
					CallURL:    "https://example.com/notify",
					Method:     http.MethodPost,
					SigningKey: "my-signing-key",
				},
			},
		},
		{
			name: "neither smtp nor http config set",
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp8",
				Description: "neither smtp nor http config set",
			},
			expectedErr: zerrors.ThrowNotFound(nil, "", ""),
		},
		{
			name:          "SMTPConfigActive query error",
			smtpConfigErr: errDatabaseError,
			expectedErr:   errDatabaseError,
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
			queryMock.EXPECT().SMTPConfigActive(gomock.Any(), instId).Return(tc.smtpConfig, tc.smtpConfigErr)

			notificationQueries := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext domain", uint16(1234), false, "filepath", nil, cryptAlgMock, nil, nil)
			cfg, err := notificationQueries.GetActiveEmailConfig(ctx)

			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr != nil {
				return
			}

			assert.EqualValues(t, tc.expected, cfg)
		})
	}
}
