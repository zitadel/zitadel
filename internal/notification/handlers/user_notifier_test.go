package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/notification/messages"

	"github.com/golang/mock/gomock"
	statik_fs "github.com/rakyll/statik/fs"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	channel_mock "github.com/zitadel/zitadel/internal/notification/channels/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	logoURL               = "logo.png"
	policyID              = "policy1"
	userID                = "user1"
	orgID                 = "org1"
	eventOrigin           = "https://triggered.here"
	assetsPath            = "/assets/v1"
	instancePrimaryDomain = "primary.domain"
	externalDomain        = "external.domain"
	externalPort          = 3000
	externalSecure        = false
	lastEmail             = "last@email.com"
)

func Test_userNotifier_reduceEmailCodeAdded(t *testing.T) {
	type fields struct {
		queries            *mock.MockQueries
		commands           *mock.MockCommands
		es                 *eventstore.Eventstore
		userDataCrypto     crypto.EncryptionAlgorithm
		SMTPPasswordCrypto crypto.EncryptionAlgorithm
		SMSTokenCrypto     crypto.EncryptionAlgorithm
	}
	type args struct {
		event eventstore.Event
	}
	type want struct {
		message messages.Email
		err     assert.ErrorAssertionFunc
	}
	queries := mock.NewMockQueries(gomock.NewController(t))
	tests := []struct {
		name string
		test func() (fields, args, want)
	}{{
		name: "asset url with event trigger url",
		test: func() (f fields, a args, w want) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = messages.Email{
				Recipients: []string{lastEmail},
				Subject:    "Verify email",
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, "testcode")
			smtpAlg, _ := cryptoValue(t, "smtppw")
			queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
			expectTemplateQueries(queries, givenTemplate)
			commands := mock.NewMockCommands(gomock.NewController(t))
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(),
					)),
					userDataCrypto:     codeAlg,
					SMTPPasswordCrypto: smtpAlg,
					SMSTokenCrypto:     nil,
				}, args{
					event: &user.HumanEmailCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						URLTemplate:       "",
						CodeReturned:      false,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func() (f fields, a args, w want) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("http://%s:%d%s/%s/%s", instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = messages.Email{
				Recipients: []string{lastEmail},
				Subject:    "Verify email",
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, "testcode")
			smtpAlg, _ := cryptoValue(t, "smtppw")
			queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateQueries(queries, givenTemplate)
			commands := mock.NewMockCommands(gomock.NewController(t))
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(),
					)),
					userDataCrypto:     codeAlg,
					SMTPPasswordCrypto: smtpAlg,
				}, args{
					event: &user.HumanEmailCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:         code,
						Expiry:       time.Hour,
						URLTemplate:  "",
						CodeReturned: false,
					},
				}, w
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := statik_fs.NewWithNamespace("notification")
			assert.NoError(t, err)
			f, a, w := tt.test()
			channel := channel_mock.NewMockNotificationChannel(gomock.NewController(t))
			if w.err == nil {
				w.message.TriggeringEvent = a.event
				channel.EXPECT().HandleMessage(&w.message).Return(nil)
			}
			u := &userNotifier{
				commands: f.commands,
				queries: NewNotificationQueries(
					f.queries,
					f.es,
					externalDomain,
					externalPort,
					externalSecure,
					"",
					f.userDataCrypto,
					f.SMTPPasswordCrypto,
					f.SMSTokenCrypto,
					fs,
				),
				channels: &channels{Chain: *senders.ChainChannels(channel)},
			}
			_, err = u.reduceEmailCodeAdded(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

var _ types.ChannelChains = (*channels)(nil)

type channels struct {
	senders.Chain
}

func (c *channels) Email(context.Context) (*senders.Chain, *smtp.Config, error) {
	return &c.Chain, nil, nil
}

func (c *channels) SMS(context.Context) (*senders.Chain, *twilio.Config, error) {
	return &c.Chain, nil, nil
}

func (c *channels) Webhook(context.Context, webhook.Config) (*senders.Chain, error) {
	return &c.Chain, nil
}

func expectTemplateQueries(queries *mock.MockQueries, template string) {
	queries.EXPECT().ActiveLabelPolicyByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.LabelPolicy{
		ID: policyID,
		Light: query.Theme{
			LogoURL: logoURL,
		},
	}, nil)
	queries.EXPECT().MailTemplateByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.MailTemplate{Template: []byte(template)}, nil)
	queries.EXPECT().GetNotifyUserByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.NotifyUser{LastEmail: lastEmail}, nil)
	queries.EXPECT().GetDefaultLanguage(gomock.Any()).Return(language.English)
	queries.EXPECT().CustomTextListByTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(&query.CustomTexts{}, nil)
}

func cryptoValue(t *testing.T, value string) (*crypto.MockEncryptionAlgorithm, *crypto.CryptoValue) {
	encAlg := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
	encAlg.EXPECT().Algorithm().AnyTimes().Return("enc")
	encAlg.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	encAlg.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	encAlg.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return(value, nil)
	encAlg.EXPECT().Encrypt(gomock.Any()).AnyTimes().Return(make([]byte, 0), nil)
	code, err := crypto.Encrypt(nil, encAlg)
	assert.NoError(t, err)
	return encAlg, code
}
