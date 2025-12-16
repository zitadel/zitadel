package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/query"
	"go.uber.org/mock/gomock"
)

func TestNotificationQueries_GetActiveEmailConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const instId = "instance-1"
	const cryptAlg = "alg1"
	const keyId = "key1"
	const pwd = "password"

	encryptedPwd := []byte("password-encrypted")
	cryptAlgMock := crypto.NewMockEncryptionAlgorithm(ctrl)
	cryptAlgMock.EXPECT().DecryptString(gomock.Any(), gomock.Any()).Return(pwd, nil)
	cryptAlgMock.EXPECT().Algorithm().Return(cryptAlg)
	cryptAlgMock.EXPECT().DecryptionKeyIDs().Return([]string{keyId})

	tt := []struct {
		smtpConfig *query.SMTPConfig
		expected   *email.Config
	}{
		{
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp1",
				Description: "smtp config with authentication (user and password)",
				SMTPConfig: &query.SMTP{
					TLS:            true,
					SenderAddress:  "sender@example.com",
					SenderName:     "sender",
					ReplyToAddress: "reply@example.com",
					Host:           "mail.com",
					User:           "mail-user",
					Password: &crypto.CryptoValue{
						CryptoType: 0,
						Algorithm:  cryptAlg,
						KeyID:      keyId,
						Crypted:    encryptedPwd,
					},
				},
			},
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp1",
					Description: "smtp config with authentication (user and password)",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host:     "mail.com",
						User:     "mail-user",
						Password: pwd,
					},
					Tls:            true,
					From:           "sender@example.com",
					FromName:       "sender",
					ReplyToAddress: "reply@example.com",
				},
			},
		},
		{
			smtpConfig: &query.SMTPConfig{
				ID:          "smtp2",
				Description: "smtp config with authentication (user and no password)",
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
					Description: "smtp config with authentication (user and no password)",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host: "mail.com",
						User: "mail-user",
					},
					Tls:            true,
					From:           "sender@example.com",
					FromName:       "sender",
					ReplyToAddress: "reply@example.com",
				},
			},
		},
		{
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
	}

	for _, tc := range tt {
		ctx := authz.NewMockContext(instId, "org-1", "user-1")

		queryMock := mock.NewMockQueries(ctrl)
		queryMock.EXPECT().SMTPConfigActive(gomock.Any(), instId).Return(tc.smtpConfig, nil)

		notificationQueries := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext domain", uint16(1234), false, "filepath", nil, cryptAlgMock, nil)
		cfg, err := notificationQueries.GetActiveEmailConfig(ctx)
		assert.NoError(t, err)
		assert.EqualValues(t, tc.expected, cfg)
	}
}
