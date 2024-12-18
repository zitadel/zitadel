package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/ui/login"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
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
				commands.EXPECT().NotificationCanceled(gomock.Any(), gomock.Any(), notificationID, instanceID, nil).Return(nil)
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
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().Add(-1 * time.Hour),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
					}, w
			},
		},
		{
			name: "send ok (email)",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().NotificationSent(gomock.Any(), gomock.Any(), notificationID, instanceID).Return(nil)
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
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
					}, w
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
				}
				codeAlg, code := cryptoValue(t, ctrl, testCode)
				expectTemplateWithNotifyUserQueriesSMS(queries)
				commands.EXPECT().NotificationSent(gomock.Any(), gomock.Any(), notificationID, instanceID).Return(nil)
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
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   sessionID,
								AggregateResourceOwner:        instanceID,
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
					Recipients: []string{verifiedEmail},
					Subject:    "Domain has been claimed",
					Content:    expectContent,
				}
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().NotificationSent(gomock.Any(), gomock.Any(), notificationID, instanceID).Return(nil)
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
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				w.sendError = sendError
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().NotificationRetryRequested(gomock.Any(), gomock.Any(), notificationID, instanceID,
					&command.NotificationRetryRequest{
						NotificationRequest: command.NotificationRequest{
							UserID:                        userID,
							UserResourceOwner:             orgID,
							AggregateID:                   "",
							AggregateResourceOwner:        "",
							TriggerOrigin:                 eventOrigin,
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
						BackOff: 1 * time.Second,
						NotifyUser: &query.NotifyUser{
							ID:                 userID,
							ResourceOwner:      orgID,
							LastEmail:          lastEmail,
							VerifiedEmail:      verifiedEmail,
							PreferredLoginName: preferredLoginName,
							LastPhone:          lastPhone,
							VerifiedPhone:      verifiedPhone,
						},
					},
					sendError,
				).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
						maxAttempts:    2,
					},
					argsWorker{
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
					}, w
			},
		},
		{
			name: "send failed (max attempts), cancel",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				w.sendError = sendError
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateWithNotifyUserQueries(queries, givenTemplate)
				commands.EXPECT().NotificationCanceled(gomock.Any(), gomock.Any(), notificationID, instanceID, sendError).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
						maxAttempts:    1,
					},
					argsWorker{
						event: &notification.RequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Seq:           1,
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
					}, w
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			commands := mock.NewMockCommands(ctrl)
			f, a, w := tt.test(ctrl, queries, commands)
			err := newNotificationWorker(t, ctrl, queries, f, a, w).reduceNotificationRequested(
				authz.WithInstanceID(context.Background(), instanceID),
				authz.WithInstanceID(context.Background(), instanceID),
				&sql.Tx{},
				a.event.(*notification.RequestedEvent))
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceNotificationRetry(t *testing.T) {
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
				commands.EXPECT().NotificationCanceled(gomock.Any(), gomock.Any(), notificationID, instanceID, nil).Return(nil)
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
						event: &notification.RetryRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().Add(-1 * time.Hour),
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
							BackOff: 1 * time.Second,
							NotifyUser: &query.NotifyUser{
								ID:                 userID,
								ResourceOwner:      orgID,
								LastEmail:          lastEmail,
								VerifiedEmail:      verifiedEmail,
								PreferredLoginName: preferredLoginName,
								LastPhone:          lastPhone,
								VerifiedPhone:      verifiedPhone,
							},
						},
					}, w
			},
		},
		{
			name: "backoff not done",
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
						event: &notification.RetryRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now(),
								Typ:           notification.RequestedType,
								Seq:           2,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
							BackOff: 10 * time.Second,
							NotifyUser: &query.NotifyUser{
								ID:                 userID,
								ResourceOwner:      orgID,
								LastEmail:          lastEmail,
								VerifiedEmail:      verifiedEmail,
								PreferredLoginName: preferredLoginName,
								LastPhone:          lastPhone,
								VerifiedPhone:      verifiedPhone,
							},
						},
					}, w
			},
		},
		{
			name: "send ok",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateQueries(queries, givenTemplate)
				commands.EXPECT().NotificationSent(gomock.Any(), gomock.Any(), notificationID, instanceID).Return(nil)
				commands.EXPECT().InviteCodeSent(gomock.Any(), orgID, userID).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						maxAttempts:    3,
					},
					argsWorker{
						event: &notification.RetryRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().Add(-2 * time.Second),
								Typ:           notification.RequestedType,
								Seq:           2,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
							BackOff: 1 * time.Second,
							NotifyUser: &query.NotifyUser{
								ID:                 userID,
								ResourceOwner:      orgID,
								LastEmail:          lastEmail,
								VerifiedEmail:      verifiedEmail,
								PreferredLoginName: preferredLoginName,
								LastPhone:          lastPhone,
								VerifiedPhone:      verifiedPhone,
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
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				w.sendError = sendError
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateQueries(queries, givenTemplate)
				commands.EXPECT().NotificationRetryRequested(gomock.Any(), gomock.Any(), notificationID, instanceID,
					&command.NotificationRetryRequest{
						NotificationRequest: command.NotificationRequest{
							UserID:                        userID,
							UserResourceOwner:             orgID,
							AggregateID:                   "",
							AggregateResourceOwner:        "",
							TriggerOrigin:                 eventOrigin,
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
						BackOff: 1 * time.Second,
						NotifyUser: &query.NotifyUser{
							ID:                 userID,
							ResourceOwner:      orgID,
							LastEmail:          lastEmail,
							VerifiedEmail:      verifiedEmail,
							PreferredLoginName: preferredLoginName,
							LastPhone:          lastPhone,
							VerifiedPhone:      verifiedPhone,
						},
					},
					sendError,
				).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
						maxAttempts:    3,
					},
					argsWorker{
						event: &notification.RetryRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().Add(-2 * time.Second),
								Typ:           notification.RequestedType,
								Seq:           2,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
							BackOff: 1 * time.Second,
							NotifyUser: &query.NotifyUser{
								ID:                 userID,
								ResourceOwner:      orgID,
								LastEmail:          lastEmail,
								VerifiedEmail:      verifiedEmail,
								PreferredLoginName: preferredLoginName,
								LastPhone:          lastPhone,
								VerifiedPhone:      verifiedPhone,
							},
						},
					}, w
			},
		},
		{
			name: "send failed (max attempts), cancel",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fieldsWorker, a argsWorker, w wantWorker) {
				givenTemplate := "{{.LogoURL}}"
				expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
				w.message = &messages.Email{
					Recipients: []string{lastEmail},
					Subject:    "Invitation to APP",
					Content:    expectContent,
				}
				w.sendError = sendError
				codeAlg, code := cryptoValue(t, ctrl, "testcode")
				expectTemplateQueries(queries, givenTemplate)
				commands.EXPECT().NotificationCanceled(gomock.Any(), gomock.Any(), notificationID, instanceID, sendError).Return(nil)
				return fieldsWorker{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
						userDataCrypto: codeAlg,
						now:            testNow,
						backOff:        testBackOff,
						maxAttempts:    2,
					},
					argsWorker{
						event: &notification.RetryRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   notificationID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().Add(-2 * time.Second),
								Seq:           2,
								Typ:           notification.RequestedType,
							}),
							Request: notification.Request{
								UserID:                        userID,
								UserResourceOwner:             orgID,
								AggregateID:                   "",
								AggregateResourceOwner:        "",
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
							BackOff: 1 * time.Second,
							NotifyUser: &query.NotifyUser{
								ID:                 userID,
								ResourceOwner:      orgID,
								LastEmail:          lastEmail,
								VerifiedEmail:      verifiedEmail,
								PreferredLoginName: preferredLoginName,
								LastPhone:          lastPhone,
								VerifiedPhone:      verifiedPhone,
							},
						},
					}, w
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			commands := mock.NewMockCommands(ctrl)
			f, a, w := tt.test(ctrl, queries, commands)
			err := newNotificationWorker(t, ctrl, queries, f, a, w).reduceNotificationRetry(
				authz.WithInstanceID(context.Background(), instanceID),
				authz.WithInstanceID(context.Background(), instanceID),
				&sql.Tx{},
				a.event.(*notification.RetryRequestedEvent),
			)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newNotificationWorker(t *testing.T, ctrl *gomock.Controller, queries *mock.MockQueries, f fieldsWorker, a argsWorker, w wantWorker) *NotificationWorker {
	queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
	smtpAlg, _ := cryptoValue(t, ctrl, "smtppw")
	channel := channel_mock.NewMockNotificationChannel(ctrl)
	if w.err == nil {
		if w.message != nil {
			w.message.TriggeringEvent = a.event
			channel.EXPECT().HandleMessage(w.message).Return(w.sendError)
		}
		if w.messageSMS != nil {
			w.messageSMS.TriggeringEvent = a.event
			channel.EXPECT().HandleMessage(w.messageSMS).DoAndReturn(func(message *messages.SMS) error {
				message.VerificationID = gu.Ptr(verificationID)
				return w.sendError
			})
		}
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
			BulkLimit:           10,
			RequeueEvery:        2 * time.Second,
			TransactionDuration: 5 * time.Second,
			MaxAttempts:         f.maxAttempts,
			MaxTtl:              5 * time.Minute,
			MinRetryDelay:       1 * time.Second,
			MaxRetryDelay:       10 * time.Second,
			RetryDelayFactor:    2,
		},
		now:     f.now,
		backOff: f.backOff,
	}
}

func TestNotificationWorker_exponentialBackOff(t *testing.T) {
	type fields struct {
		config WorkerConfig
	}
	type args struct {
		current time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantMin time.Duration
		wantMax time.Duration
	}{
		{
			name: "less than min, min - 1.5*min",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 0,
			},
			wantMin: 1000 * time.Millisecond,
			wantMax: 1500 * time.Millisecond,
		},
		{
			name: "current, 1.5*current - max",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 4 * time.Second,
			},
			wantMin: 4000 * time.Millisecond,
			wantMax: 5000 * time.Millisecond,
		},
		{
			name: "max, max",
			fields: fields{
				config: WorkerConfig{
					MinRetryDelay:    1 * time.Second,
					MaxRetryDelay:    5 * time.Second,
					RetryDelayFactor: 1.5,
				},
			},
			args: args{
				current: 5 * time.Second,
			},
			wantMin: 5000 * time.Millisecond,
			wantMax: 5000 * time.Millisecond,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &NotificationWorker{
				config: tt.fields.config,
			}
			b := w.exponentialBackOff(tt.args.current)
			assert.GreaterOrEqual(t, b, tt.wantMin)
			assert.LessOrEqual(t, b, tt.wantMax)
		})
	}
}
