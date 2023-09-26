package handlers

import (
	"context"
	"testing"
	"time"

	statik_fs "github.com/rakyll/statik/fs"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/query"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func Test_userNotifier_reduceEmailCodeAdded(t *testing.T) {
	type fields struct {
		queries             *mock.MockQueries
		commands            *mock.MockCommands
		es                  *eventstore.Eventstore
		defaultAssetsPrefix func(context.Context) string
		externalDomain      string
		externalPort        uint16
		externalSecure      bool
		fileSystemPath      string
		userDataCrypto      crypto.EncryptionAlgorithm
		SMTPPasswordCrypto  crypto.EncryptionAlgorithm
		SMSTokenCrypto      crypto.EncryptionAlgorithm
		otpEmailTmpl        string
	}
	type args struct {
		event eventstore.Event
	}
	queries := mock.NewMockQueries(gomock.NewController(t))
	tests := []struct {
		name    string
		given   func() (fields, args)
		want    *handler.Statement
		wantErr assert.ErrorAssertionFunc
	}{{
		name: "test suite",
		given: func() (fields, args) {
			codeAlg, code := cryptoValue(t, "testcode")
			smtpAlg, smtppw := cryptoValue(t, "smtppw")
			queries.EXPECT().SMTPConfigByAggregateID(gomock.Any(), gomock.Any()).Return(&query.SMTPConfig{Password: smtppw}, nil)
			queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
			return fields{
					queries:  queries,
					commands: mock.NewMockCommands(gomock.NewController(t)),
					es: eventstore.NewEventstore(eventstore.TestConfig(
						es_repo_mock.NewRepo(t).ExpectFilterEvents(),
					)),
					defaultAssetsPrefix: func(ctx context.Context) string {
						return "irrelevant assets prefix"
					},
					externalDomain:     "",
					externalPort:       0,
					externalSecure:     false,
					fileSystemPath:     "",
					userDataCrypto:     codeAlg,
					SMTPPasswordCrypto: smtpAlg,
					SMSTokenCrypto:     nil,
					otpEmailTmpl:       "",
				}, args{
					event: &user.HumanEmailCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							CreationDate: time.Now().UTC(),
						}),
						Code:               code,
						Expiry:             time.Hour,
						URLTemplate:        "",
						CodeReturned:       false,
						TriggeredAtBaseURL: "https://tiggered.here",
					},
				}
		},
		want:    nil,
		wantErr: assert.NoError,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs, err := statik_fs.NewWithNamespace("notification")
			assert.NoError(t, err)
			f, a := tt.given()
			expectTemplateQueries(f.queries)
			u := &userNotifier{
				commands: f.commands,
				queries: NewNotificationQueries(
					f.queries,
					f.es,
					f.externalDomain,
					f.externalPort,
					f.externalSecure,
					f.fileSystemPath,
					f.userDataCrypto,
					f.SMTPPasswordCrypto,
					f.SMSTokenCrypto,
					fs,
				),
				assetsPrefix: f.defaultAssetsPrefix,
				otpEmailTmpl: f.otpEmailTmpl,
			}
			_, err = u.reduceEmailCodeAdded(a.event)
			tt.wantErr(t, err)
		})
	}
}

func expectTemplateQueries(queries *mock.MockQueries) {
	queries.EXPECT().ActiveLabelPolicyByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.LabelPolicy{}, nil)
	queries.EXPECT().MailTemplateByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.MailTemplate{}, nil)
	queries.EXPECT().GetNotifyUserByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.NotifyUser{}, nil)
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
