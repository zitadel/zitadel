package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
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
			expected: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "smtp1",
					Description: "smtp config with authentication (user and password)",
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
		t.Run(tc.smtpConfig.Description, func(t *testing.T) {
			ctx := authz.NewMockContext(instId, "org-1", "user-1")

			queryMock := mock.NewMockQueries(ctrl)
			queryMock.EXPECT().SMTPConfigActive(gomock.Any(), instId).Return(tc.smtpConfig, nil)

			notificationQueries := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext domain", uint16(1234), false, "filepath", nil, cryptAlgMock, nil, true)
			cfg, err := notificationQueries.GetActiveEmailConfig(ctx, "")
			assert.NoError(t, err)
			assert.EqualValues(t, tc.expected, cfg)
		})
	}
}

func TestNotificationQueries_GetActiveEmailConfig_OrgFallback(t *testing.T) {
	const instId = "instance-1"
	const orgId = "org-1"

	instanceSMTPConfig := &query.SMTPConfig{
		ID:          "instance-smtp",
		Description: "instance SMTP",
		SMTPConfig: &query.SMTP{
			Host:          "instance-mail.com",
			SenderAddress: "noreply@instance.com",
			SenderName:    "Instance",
		},
	}
	orgSMTPConfig := &query.SMTPConfig{
		ID:          "org-smtp",
		Description: "org SMTP",
		SMTPConfig: &query.SMTP{
			Host:          "org-mail.com",
			SenderAddress: "noreply@org.com",
			SenderName:    "Org",
		},
	}

	t.Run("fallback enabled + org not found → uses instance SMTP", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := authz.NewMockContext(instId, orgId, "user-1")
		queryMock := mock.NewMockQueries(ctrl)
		queryMock.EXPECT().OrgSMTPConfigActive(gomock.Any(), orgId).Return(nil, zerrors.ThrowNotFound(nil, "QUERY-test", "not found"))
		queryMock.EXPECT().SMTPConfigActive(gomock.Any(), instId).Return(instanceSMTPConfig, nil)

		nq := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext", uint16(443), true, "", nil, nil, nil, true)
		cfg, err := nq.GetActiveEmailConfig(ctx, orgId)
		require.NoError(t, err)
		assert.Equal(t, "instance-smtp", cfg.ProviderConfig.ID)
	})

	t.Run("fallback disabled + org not found → returns error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := authz.NewMockContext(instId, orgId, "user-1")
		queryMock := mock.NewMockQueries(ctrl)
		queryMock.EXPECT().OrgSMTPConfigActive(gomock.Any(), orgId).Return(nil, zerrors.ThrowNotFound(nil, "QUERY-test", "not found"))

		nq := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext", uint16(443), true, "", nil, nil, nil, false)
		_, err := nq.GetActiveEmailConfig(ctx, orgId)
		require.Error(t, err)
		assert.True(t, zerrors.IsNotFound(err))
	})

	t.Run("org SMTP exists → uses org SMTP regardless of fallback", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := authz.NewMockContext(instId, orgId, "user-1")
		queryMock := mock.NewMockQueries(ctrl)
		queryMock.EXPECT().OrgSMTPConfigActive(gomock.Any(), orgId).Return(orgSMTPConfig, nil)

		nq := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext", uint16(443), true, "", nil, nil, nil, true)
		cfg, err := nq.GetActiveEmailConfig(ctx, orgId)
		require.NoError(t, err)
		assert.Equal(t, "org-smtp", cfg.ProviderConfig.ID)
	})

	t.Run("real error propagates regardless of fallback", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := authz.NewMockContext(instId, orgId, "user-1")
		queryMock := mock.NewMockQueries(ctrl)
		queryMock.EXPECT().OrgSMTPConfigActive(gomock.Any(), orgId).Return(nil, zerrors.ThrowInternal(nil, "QUERY-test", "db failure"))

		nq := NewNotificationQueries(queryMock, &eventstore.Eventstore{}, "ext", uint16(443), true, "", nil, nil, nil, true)
		_, err := nq.GetActiveEmailConfig(ctx, orgId)
		require.Error(t, err)
		assert.False(t, zerrors.IsNotFound(err))
	})
}
