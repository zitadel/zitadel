package handlers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	channel_mock "github.com/zitadel/zitadel/internal/notification/channels/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/smtp"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/notification"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	notificationID = "notificationID"
)

func Test_userNotifier_reduceNotificationRequested(t *testing.T) {
	testNow := time.Now
	testBackOff := func(current time.Duration) time.Duration {
		return time.Second
	}
	sendError := errors.New("send error")
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fieldsWorker, argsWorker, wantWorker)
	}{
		{
			name: "too old",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now().Add(-1 * time.Hour),
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            notificationID,
									ResourceOwner: instanceID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     user.HumanInviteCodeAddedType,
								MessageType:                   domain.InviteUserMessageType,
								NotificationType:              domain.NotificationTypeEmail,
								URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
								CodeExpiry:                    1 * time.Hour,
								Code:                          code,
								UnverifiedNotificationChannel: true,
								IsOTP:                         false,
								RequiresPreviousDomain:        false,
								Args: &domain.NotificationArguments{
									ApplicationName: "APP",
								},
							},
						},
					},
					wantWorker{
						err: func(tt assert.TestingT, err error, i ...interface{}) bool {
							return errors.Is(err, new(river.JobCancelError))
						},
					}
			},
		},
		{
			name: "send ok (email)",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients:          []string{lastEmail},
					Subject:             "Invitation to APP",
					Content:             expectContent,
					TriggeringEventType: user.HumanInviteCodeAddedType,
				}
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().InviteCodeSent(gomock.Any(), orgID, userID).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            userID,
									ResourceOwner: orgID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     user.HumanInviteCodeAddedType,
								MessageType:                   domain.InviteUserMessageType,
								NotificationType:              domain.NotificationTypeEmail,
								URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
								CodeExpiry:                    1 * time.Hour,
								Code:                          code,
								UnverifiedNotificationChannel: true,
								IsOTP:                         false,
								RequiresPreviousDomain:        false,
								Args: &domain.NotificationArguments{
									ApplicationName: "APP",
								},
							},
						},
					},
					w
			},
		},
		{
			name: "send ok (sms with external provider)",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				expiry := 0 * time.Hour
				testCode := ""
				expectContent := fmt.Sprintf(`%[1]s is your one-time password for %[2]s. Use it within the next %[3]s.
@%[2]s #%[1]s`, testCode, eventOriginDomain, expiry)
				w.messageSMS = &messages.SMS{
					SenderPhoneNumber:    "senderNumber",
					RecipientPhoneNumber: verifiedPhone,
					Content:              expectContent,
					TriggeringEventType:  session.OTPSMSChallengedType,
					InstanceID:           instanceID,
					JobID:                "1",
					UserID:               userID,
				}
				codeAlg, code := cryptoValue(t, ctrl, testCode)
				expectTemplateWithNotifyUserQueriesSMS(queries)
				commands.EXPECT().OTPSMSSent(gomock.Any(), sessionID, instanceID, &senders.CodeGeneratorInfo{
					ID:             smsProviderID,
					VerificationID: verificationID,
				}).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
								ID:        1,
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            sessionID,
									ResourceOwner: instanceID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     session.OTPSMSChallengedType,
								MessageType:                   domain.VerifySMSOTPMessageType,
								NotificationType:              domain.NotificationTypeSms,
								URLTemplate:                   "",
								CodeExpiry:                    expiry,
								Code:                          code,
								UnverifiedNotificationChannel: false,
								IsOTP:                         true,
								RequiresPreviousDomain:        false,
								Args: &domain.NotificationArguments{
									Origin: eventOrigin,
									Domain: eventOriginDomain,
									Expiry: expiry,
								},
							},
						},
					}, w
			},
		},
		{
			name: "previous domain",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients:          []string{verifiedEmail},
					Subject:             "Domain has been claimed",
					Content:             expectContent,
					TriggeringEventType: user.UserDomainClaimedType,
				}
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().UserDomainClaimedSent(gomock.Any(), orgID, userID).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: nil,
						now:            testNow,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            userID,
									ResourceOwner: orgID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     user.UserDomainClaimedType,
								MessageType:                   domain.DomainClaimedMessageType,
								NotificationType:              domain.NotificationTypeEmail,
								URLTemplate:                   login.LoginLink(eventOrigin, orgID),
								CodeExpiry:                    0,
								Code:                          nil,
								UnverifiedNotificationChannel: false,
								IsOTP:                         false,
								RequiresPreviousDomain:        true,
								Args: &domain.NotificationArguments{
									TempUsername: "tempUsername",
								},
							},
						},
					}, w
			},
		},
		{
			name: "send failed, retry",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients:          []string{lastEmail},
					Subject:             "Invitation to APP",
					Content:             expectContent,
					TriggeringEventType: user.HumanInviteCodeAddedType,
				}
				w.sendError = sendError
				w.err = func(tt assert.TestingT, err error, i ...interface{}) bool {
					return errors.Is(err, sendError)
				}
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								ID:        1,
								CreatedAt: time.Now(),
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            notificationID,
									ResourceOwner: instanceID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     user.HumanInviteCodeAddedType,
								MessageType:                   domain.InviteUserMessageType,
								NotificationType:              domain.NotificationTypeEmail,
								URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
								CodeExpiry:                    1 * time.Hour,
								Code:                          code,
								UnverifiedNotificationChannel: true,
								IsOTP:                         false,
								RequiresPreviousDomain:        false,
								Args: &domain.NotificationArguments{
									ApplicationName: "APP",
								},
							},
						},
					},
					w
			},
		},
		{
			name: "send failed (max attempts), cancel",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients:          []string{lastEmail},
					Subject:             "Invitation to APP",
					Content:             expectContent,
					TriggeringEventType: user.HumanInviteCodeAddedType,
				}
				w.sendError = sendError
				w.err = func(tt assert.TestingT, err error, i ...interface{}) bool {
					return err != nil
				}

				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
					},
					argsWorker{
						job: &river.Job[*notification.Request]{
							JobRow: &rivertype.JobRow{
								CreatedAt: time.Now(),
							},
							Args: &notification.Request{
								Aggregate: &eventstore.Aggregate{
									InstanceID:    instanceID,
									ID:            userID,
									ResourceOwner: orgID,
								},
								UserID:                        userID,
								UserResourceOwner:             orgID,
								TriggeredAtOrigin:             eventOrigin,
								EventType:                     user.HumanInviteCodeAddedType,
								MessageType:                   domain.InviteUserMessageType,
								NotificationType:              domain.NotificationTypeEmail,
								URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
								CodeExpiry:                    1 * time.Hour,
								Code:                          code,
								UnverifiedNotificationChannel: true,
								IsOTP:                         false,
								RequiresPreviousDomain:        false,
								Args: &domain.NotificationArguments{
									ApplicationName: "APP",
								},
							},
						},
					},
					w
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			commands := mock.NewMockCommands(ctrl)
			f, a, w := tt.test(ctrl, queries, commands)
			err := newNotificationWorker(t, ctrl, queries, f, w).Work(
				authz.WithInstanceID(context.Background(), instanceID),
				a.job,
			)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newNotificationWorker(t *testing.T, ctrl *gomock.Controller, queries *mock.MockQueries, f fieldsWorker, w wantWorker) *NotificationWorker {
	queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
	smtpAlg, _ := cryptoValue(t, ctrl, "smtppw")
	channel := channel_mock.NewMockNotificationChannel(ctrl)
	if w.message != nil {
		channel.EXPECT().HandleMessage(w.message).Return(w.sendError)
	}
	if w.messageSMS != nil {
		channel.EXPECT().HandleMessage(w.messageSMS).DoAndReturn(func(message *messages.SMS) error {
			message.VerificationID = gu.Ptr(verificationID)
			return w.sendError
		})
	}
	return &NotificationWorker{
		commands: f.commands,
		queries: NewNotificationQueries(
			f.queries,
			f.es,
			externalDomain,
			externalPort,
			externalSecure,
			"",
			f.userDataCrypto,
			smtpAlg,
			f.SMSTokenCrypto,
		),
		channels: &notificationChannels{
			Chain: *senders.ChainChannels(channel),
			EmailConfig: &email.Config{
				ProviderConfig: &email.Provider{
					ID:          "emailProviderID",
					Description: "description",
				},
				SMTPConfig: &smtp.Config{
					SMTP: smtp.SMTP{
						Host:     "host",
						User:     "user",
						Password: "password",
					},
					Tls:            true,
					From:           "from",
					FromName:       "fromName",
					ReplyToAddress: "replyToAddress",
				},
				WebhookConfig: nil,
			},
			SMSConfig: &sms.Config{
				ProviderConfig: &sms.Provider{
					ID:          "smsProviderID",
					Description: "description",
				},
				TwilioConfig: &twilio.Config{
					SID:              "sid",
					Token:            "token",
					SenderNumber:     "senderNumber",
					VerifyServiceSID: "verifyServiceSID",
				},
			},
		},
		config: WorkerConfig{
			Workers:             1,
			TransactionDuration: 5 * time.Second,
			MaxTtl:              5 * time.Minute,
		},
		now: f.now,
	}
}
