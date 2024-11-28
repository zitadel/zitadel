package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/notification/channels/email"
	"github.com/zitadel/zitadel/internal/notification/channels/set"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	orgID                   = "org1"
	policyID                = "policy1"
	userID                  = "user1"
	codeID                  = "event1"
	logoURL                 = "logo.png"
	instanceID              = "instanceID"
	sessionID               = "sessionID"
	eventOrigin             = "https://triggered.here"
	eventOriginDomain       = "triggered.here"
	assetsPath              = "/assets/v1"
	preferredLoginName      = "loginName1"
	lastEmail               = "last@email.com"
	verifiedEmail           = "verified@email.com"
	lastPhone               = "+41797654321"
	verifiedPhone           = "+41791234567"
	instancePrimaryDomain   = "primary.domain"
	externalDomain          = "external.domain"
	externalPort            = 3000
	externalSecure          = false
	externalProtocol        = "http"
	defaultOTPEmailTemplate = "/otp/verify?loginName={{.LoginName}}&code={{.Code}}"
	authRequestID           = "authRequestID"
	smsProviderID           = "smsProviderID"
	emailProviderID         = "emailProviderID"
	verificationID          = "verificationID"
)

func Test_userNotifier_reduceInitCodeAdded(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     eventOrigin,
					URLTemplate: fmt.Sprintf("%s/ui/login/user/init?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&passwordset={{.PasswordSet}}&authRequestID=%s",
						eventOrigin, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInitialCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InitCodeMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInitialCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInitialCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate: fmt.Sprintf("%s://%s:%d/ui/login/user/init?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&passwordset={{.PasswordSet}}&authRequestID=%s",
						externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInitialCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InitCodeMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInitialCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInitialCodeAddedType,
							}),
							Code:          code,
							Expiry:        time.Hour,
							AuthRequestID: authRequestID,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceInitCodeAdded(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceEmailCodeAdded(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     eventOrigin,
					URLTemplate: fmt.Sprintf("%s/ui/login/mail/verification?userID=%s&code={{.Code}}&orgID=%s&authRequestID=%s",
						eventOrigin, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanEmailCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.VerifyEmailMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanEmailCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanEmailCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)

				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate: fmt.Sprintf("%s://%s:%d/ui/login/mail/verification?userID=%s&code={{.Code}}&orgID=%s&authRequestID=%s",
						externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanEmailCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.VerifyEmailMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanEmailCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanEmailCodeAddedType,
							}),
							Code:          code,
							Expiry:        time.Hour,
							URLTemplate:   "",
							CodeReturned:  false,
							AuthRequestID: authRequestID,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testcode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &user.HumanEmailCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanEmailCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      true,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceEmailCodeAdded(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reducePasswordCodeAdded(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     eventOrigin,
					URLTemplate: fmt.Sprintf("%s/ui/login/password/init?userID=%s&code={{.Code}}&orgID=%s&authRequestID=%s",
						eventOrigin, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanPasswordCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordResetMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "asset url without event trigger url",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate: fmt.Sprintf("%s://%s:%d/ui/login/password/init?userID=%s&code={{.Code}}&orgID=%s&authRequestID=%s",
						externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanPasswordCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordResetMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordCodeAddedType,
							}),
							Code:          code,
							Expiry:        time.Hour,
							URLTemplate:   "",
							CodeReturned:  false,
							AuthRequestID: authRequestID,
						},
					}, w
			},
		},
		{
			name: "external code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     eventOrigin,
					URLTemplate: fmt.Sprintf("%s/ui/login/password/init?userID=%s&code={{.Code}}&orgID=%s&authRequestID=%s",
						eventOrigin, userID, orgID, authRequestID),
					Code:                          nil,
					CodeExpiry:                    0,
					EventType:                     user.HumanPasswordCodeAddedType,
					NotificationType:              domain.NotificationTypeSms,
					MessageType:                   domain.PasswordResetMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          &domain.NotificationArguments{AuthRequestID: authRequestID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordCodeAddedType,
							}),
							Code:              nil,
							Expiry:            0,
							URLTemplate:       "",
							CodeReturned:      false,
							NotificationType:  domain.NotificationTypeSms,
							GeneratorID:       smsProviderID,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testcode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordCodeAddedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordCodeAddedType,
							}),
							Code:              code,
							Expiry:            1 * time.Hour,
							URLTemplate:       "",
							CodeReturned:      true,
							NotificationType:  domain.NotificationTypeSms,
							GeneratorID:       smsProviderID,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reducePasswordCodeAdded(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceDomainClaimed(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{{
		name: "with event trigger",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
			commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
				UserID:            userID,
				UserResourceOwner: orgID,
				TriggerOrigin:     eventOrigin,
				URLTemplate: fmt.Sprintf("%s/ui/login/login?orgID=%s",
					eventOrigin, orgID),
				Code:                          nil,
				CodeExpiry:                    0,
				EventType:                     user.UserDomainClaimedType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.DomainClaimedMessageType,
				UnverifiedNotificationChannel: true,
				Args:                          &domain.NotificationArguments{TempUsername: "newUsername"},
				AggregateID:                   "",
				AggregateResourceOwner:        "",
				IsOTP:                         false,
				RequiresPreviousDomain:        true,
			}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.DomainClaimedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							InstanceID:    instanceID,
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
							Typ:           user.UserDomainClaimedType,
						}),
						TriggeredAtOrigin: eventOrigin,
						UserName:          "newUsername",
					},
				}, w
		},
	}, {
		name: "without event trigger",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
				UserID:            userID,
				UserResourceOwner: orgID,
				TriggerOrigin:     fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
				URLTemplate: fmt.Sprintf("%s://%s:%d/ui/login/login?orgID=%s",
					externalProtocol, instancePrimaryDomain, externalPort, orgID),
				Code:                          nil,
				CodeExpiry:                    0,
				EventType:                     user.UserDomainClaimedType,
				NotificationType:              domain.NotificationTypeEmail,
				MessageType:                   domain.DomainClaimedMessageType,
				UnverifiedNotificationChannel: true,
				Args:                          &domain.NotificationArguments{TempUsername: "newUsername"},
				AggregateID:                   "",
				AggregateResourceOwner:        "",
				IsOTP:                         false,
				RequiresPreviousDomain:        true,
			}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.DomainClaimedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							InstanceID:    instanceID,
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
							Typ:           user.UserDomainClaimedType,
						}),
						UserName: "newUsername",
					},
				}, w
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			commands := mock.NewMockCommands(ctrl)
			f, a, w := tt.test(ctrl, queries, commands)
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceDomainClaimed(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reducePasswordlessCodeRequested(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   fmt.Sprintf("%s/ui/login/login/passwordless/init?userID=%s&orgID=%s&codeID=%s&code={{.Code}}", eventOrigin, userID, orgID, codeID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanPasswordlessInitCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordlessRegistrationMessageType,
					UnverifiedNotificationChannel: false,
					Args:                          &domain.NotificationArguments{CodeID: codeID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordlessInitCodeRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordlessInitCodeAddedType,
							}),
							ID:                codeID,
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testCode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate:                   fmt.Sprintf("%s://%s:%d/ui/login/login/passwordless/init?userID=%s&orgID=%s&codeID=%s&code={{.Code}}", externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, codeID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanPasswordlessInitCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordlessRegistrationMessageType,
					UnverifiedNotificationChannel: false,
					Args:                          &domain.NotificationArguments{CodeID: codeID},
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordlessInitCodeRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordlessInitCodeAddedType,
							}),
							ID:           codeID,
							Code:         code,
							Expiry:       time.Hour,
							URLTemplate:  "",
							CodeReturned: false,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testcode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordlessInitCodeRequestedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordlessInitCodeAddedType,
							}),
							ID:                codeID,
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      true,
							TriggeredAtOrigin: eventOrigin,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reducePasswordlessCodeRequested(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reducePasswordChanged(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				queries.EXPECT().NotificationPolicyByOrg(gomock.Any(), gomock.Any(), orgID, gomock.Any()).Return(&query.NotificationPolicy{
					PasswordChange: true,
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   fmt.Sprintf("%s/ui/console?login_hint={{.PreferredLoginName}}", eventOrigin),
					Code:                          nil,
					CodeExpiry:                    0,
					EventType:                     user.HumanPasswordChangedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordChangeMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          nil,
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordChangedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordChangedType,
							}),
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				queries.EXPECT().NotificationPolicyByOrg(gomock.Any(), gomock.Any(), orgID, gomock.Any()).Return(&query.NotificationPolicy{
					PasswordChange: true,
				}, nil)
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:            userID,
					UserResourceOwner: orgID,
					TriggerOrigin:     fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate: fmt.Sprintf("%s://%s:%d/ui/console?login_hint={{.PreferredLoginName}}",
						externalProtocol, instancePrimaryDomain, externalPort),
					Code:                          nil,
					CodeExpiry:                    0,
					EventType:                     user.HumanPasswordChangedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.PasswordChangeMessageType,
					UnverifiedNotificationChannel: true,
					Args:                          nil,
					AggregateID:                   "",
					AggregateResourceOwner:        "",
					IsOTP:                         false,
					RequiresPreviousDomain:        false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordChangedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordChangedType,
							}),
						},
					}, w
			},
		}, {
			name: "no notification",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				queries.EXPECT().NotificationPolicyByOrg(gomock.Any(), gomock.Any(), orgID, gomock.Any()).Return(&query.NotificationPolicy{
					PasswordChange: false,
				}, nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanPasswordChangedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanPasswordChangedType,
							}),
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reducePasswordChanged(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceOTPEmailChallenged(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "url with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testCode")
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   fmt.Sprintf("%s/otp/verify?loginName={{.LoginName}}&code={{.Code}}", eventOrigin),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     session.OTPEmailChallengedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.VerifyEmailOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    eventOriginDomain,
						Expiry:    1 * time.Hour,
						Origin:    eventOrigin,
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPEmailChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPEmailChallengedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTmpl:           "",
							ReturnCode:        false,
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testCode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate:                   fmt.Sprintf("%s://%s:%d/otp/verify?loginName={{.LoginName}}&code={{.Code}}", externalProtocol, instancePrimaryDomain, externalPort),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     session.OTPEmailChallengedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.VerifyEmailOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    instancePrimaryDomain,
						Expiry:    1 * time.Hour,
						Origin:    fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPEmailChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPEmailChallengedType,
							}),
							Code:       code,
							Expiry:     time.Hour,
							ReturnCode: false,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testCode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &session.OTPEmailChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPEmailChallengedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTmpl:           "",
							ReturnCode:        true,
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "url template",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testCode")
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   "/verify-otp?sessionID={{.SessionID}}",
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     session.OTPEmailChallengedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.VerifyEmailOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    eventOriginDomain,
						Expiry:    1 * time.Hour,
						Origin:    eventOrigin,
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPEmailChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPEmailChallengedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTmpl:           "/verify-otp?sessionID={{.SessionID}}",
							ReturnCode:        false,
							TriggeredAtOrigin: eventOrigin,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceSessionOTPEmailChallenged(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceOTPSMSChallenged(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				testCode := "testcode"
				_, code := cryptoValue(t, ctrl, testCode)
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   "",
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     session.OTPSMSChallengedType,
					NotificationType:              domain.NotificationTypeSms,
					MessageType:                   domain.VerifySMSOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    eventOriginDomain,
						Expiry:    1 * time.Hour,
						Origin:    eventOrigin,
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPSMSChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPSMSChallengedType,
							}),
							Code:              code,
							Expiry:            1 * time.Hour,
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				testCode := "testcode"
				_, code := cryptoValue(t, ctrl, testCode)
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate:                   "",
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     session.OTPSMSChallengedType,
					NotificationType:              domain.NotificationTypeSms,
					MessageType:                   domain.VerifySMSOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    instancePrimaryDomain,
						Expiry:    1 * time.Hour,
						Origin:    fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPSMSChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPSMSChallengedType,
							}),
							Code:         code,
							Expiry:       1 * time.Hour,
							CodeReturned: false,
						},
					}, w
			},
		},
		{
			name: "external code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), sessionID, gomock.Any()).Return(&query.Session{
					ID:            sessionID,
					ResourceOwner: instanceID,
					UserFactor: query.SessionUserFactor{
						UserID:        userID,
						ResourceOwner: orgID,
					},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   "",
					Code:                          nil,
					CodeExpiry:                    0,
					EventType:                     session.OTPSMSChallengedType,
					NotificationType:              domain.NotificationTypeSms,
					MessageType:                   domain.VerifySMSOTPMessageType,
					UnverifiedNotificationChannel: false,
					Args: &domain.NotificationArguments{
						Domain:    eventOriginDomain,
						Expiry:    0 * time.Hour,
						Origin:    eventOrigin,
						SessionID: sessionID,
					},
					AggregateID:            sessionID,
					AggregateResourceOwner: instanceID,
					IsOTP:                  true,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &session.OTPSMSChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPSMSChallengedType,
							}),
							Code:              nil,
							Expiry:            0,
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testCode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &session.OTPSMSChallengedEvent{
							BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   sessionID,
								ResourceOwner: sql.NullString{String: instanceID},
								CreationDate:  time.Now().UTC(),
								Typ:           session.OTPSMSChallengedType,
							}),
							Code:              code,
							Expiry:            1 * time.Hour,
							CodeReturned:      true,
							TriggeredAtOrigin: eventOrigin,
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceSessionOTPSMSChallenged(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifier_reduceInviteCodeAdded(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, want)
	}{
		{
			name: "with event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInviteCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InviteUserMessageType,
					UnverifiedNotificationChannel: true,
					Args: &domain.NotificationArguments{
						ApplicationName: "ZITADEL",
						AuthRequestID:   authRequestID,
					},
					AggregateID:            "",
					AggregateResourceOwner: "",
					IsOTP:                  false,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInviteCodeAddedEvent{
							BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInviteCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "without event trigger",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testCode")
				queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
					Domains: []*query.InstanceDomain{{
						Domain:    instancePrimaryDomain,
						IsPrimary: true,
					}},
				}, nil)
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 fmt.Sprintf("%s://%s:%d", externalProtocol, instancePrimaryDomain, externalPort),
					URLTemplate:                   fmt.Sprintf("%s://%s:%d/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInviteCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InviteUserMessageType,
					UnverifiedNotificationChannel: true,
					Args: &domain.NotificationArguments{
						ApplicationName: "ZITADEL",
						AuthRequestID:   authRequestID,
					},
					AggregateID:            "",
					AggregateResourceOwner: "",
					IsOTP:                  false,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInviteCodeAddedEvent{
							BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInviteCodeAddedType,
							}),
							Code:          code,
							Expiry:        time.Hour,
							URLTemplate:   "",
							CodeReturned:  false,
							AuthRequestID: authRequestID,
						},
					}, w
			},
		},
		{
			name: "return code",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				w.noOperation = true
				_, code := cryptoValue(t, ctrl, "testcode")
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).MockQuerier,
						}),
					}, args{
						event: &user.HumanInviteCodeAddedEvent{
							BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInviteCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      true,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "url template",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   "/passwordless-init?userID={{.UserID}}",
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInviteCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InviteUserMessageType,
					UnverifiedNotificationChannel: true,
					Args: &domain.NotificationArguments{
						ApplicationName: "ZITADEL",
						AuthRequestID:   authRequestID,
					},
					AggregateID:            "",
					AggregateResourceOwner: "",
					IsOTP:                  false,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInviteCodeAddedEvent{
							BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInviteCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "/passwordless-init?userID={{.UserID}}",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
						},
					}, w
			},
		},
		{
			name: "application name",
			test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w want) {
				_, code := cryptoValue(t, ctrl, "testcode")
				commands.EXPECT().RequestNotification(gomock.Any(), orgID, &command.NotificationRequest{
					UserID:                        userID,
					UserResourceOwner:             orgID,
					TriggerOrigin:                 eventOrigin,
					URLTemplate:                   fmt.Sprintf("%s/ui/login/user/invite?userID=%s&loginname={{.LoginName}}&code={{.Code}}&orgID=%s&authRequestID=%s", eventOrigin, userID, orgID, authRequestID),
					Code:                          code,
					CodeExpiry:                    time.Hour,
					EventType:                     user.HumanInviteCodeAddedType,
					NotificationType:              domain.NotificationTypeEmail,
					MessageType:                   domain.InviteUserMessageType,
					UnverifiedNotificationChannel: true,
					Args: &domain.NotificationArguments{
						ApplicationName: "APP",
						AuthRequestID:   authRequestID,
					},
					AggregateID:            "",
					AggregateResourceOwner: "",
					IsOTP:                  false,
					RequiresPreviousDomain: false,
				}).Return(nil)
				return fields{
						queries:  queries,
						commands: commands,
						es: eventstore.NewEventstore(&eventstore.Config{
							Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
						}),
					}, args{
						event: &user.HumanInviteCodeAddedEvent{
							BaseEvent: eventstore.BaseEventFromRepo(&repository.Event{
								InstanceID:    instanceID,
								AggregateID:   userID,
								ResourceOwner: sql.NullString{String: orgID},
								CreationDate:  time.Now().UTC(),
								Typ:           user.HumanInviteCodeAddedType,
							}),
							Code:              code,
							Expiry:            time.Hour,
							URLTemplate:       "",
							CodeReturned:      false,
							TriggeredAtOrigin: eventOrigin,
							AuthRequestID:     authRequestID,
							ApplicationName:   "APP",
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
			stmt, err := newUserNotifier(t, ctrl, queries, f, a, w).reduceInviteCodeAdded(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
			if w.noOperation {
				assert.Nil(t, stmt.Execute)
				return
			}
			err = stmt.Execute(nil, "")
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type fields struct {
	queries        *mock.MockQueries
	commands       *mock.MockCommands
	es             *eventstore.Eventstore
	userDataCrypto crypto.EncryptionAlgorithm
	SMSTokenCrypto crypto.EncryptionAlgorithm
}
type fieldsWorker struct {
	queries        *mock.MockQueries
	commands       *mock.MockCommands
	es             *eventstore.Eventstore
	userDataCrypto crypto.EncryptionAlgorithm
	SMSTokenCrypto crypto.EncryptionAlgorithm
	now            nowFunc
	backOff        func(current time.Duration) time.Duration
	maxAttempts    uint8
}
type args struct {
	event eventstore.Event
}
type argsWorker struct {
	event eventstore.Event
}
type want struct {
	noOperation bool
	err         assert.ErrorAssertionFunc
}
type wantWorker struct {
	message    *messages.Email
	messageSMS *messages.SMS
	sendError  error
	err        assert.ErrorAssertionFunc
}

func newUserNotifier(t *testing.T, ctrl *gomock.Controller, queries *mock.MockQueries, f fields, a args, w want) *userNotifier {
	queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
	smtpAlg, _ := cryptoValue(t, ctrl, "smtppw")
	return &userNotifier{
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
		otpEmailTmpl: defaultOTPEmailTemplate,
	}
}

var _ types.ChannelChains = (*notificationChannels)(nil)

type notificationChannels struct {
	senders.Chain
	EmailConfig *email.Config
	SMSConfig   *sms.Config
}

func (c *notificationChannels) Email(context.Context) (*senders.Chain, *email.Config, error) {
	return &c.Chain, c.EmailConfig, nil
}

func (c *notificationChannels) SMS(context.Context) (*senders.Chain, *sms.Config, error) {
	return &c.Chain, c.SMSConfig, nil
}

func (c *notificationChannels) Webhook(context.Context, webhook.Config) (*senders.Chain, error) {
	return &c.Chain, nil
}

func (c *notificationChannels) SecurityTokenEvent(context.Context, set.Config) (*senders.Chain, error) {
	return &c.Chain, nil
}

func expectTemplateQueries(queries *mock.MockQueries, template string) {
	queries.EXPECT().GetInstanceRestrictions(gomock.Any()).Return(query.Restrictions{
		AllowedLanguages: []language.Tag{language.English},
	}, nil)
	queries.EXPECT().ActiveLabelPolicyByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.LabelPolicy{
		ID: policyID,
		Light: query.Theme{
			LogoURL: logoURL,
		},
	}, nil)
	queries.EXPECT().MailTemplateByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.MailTemplate{Template: []byte(template)}, nil)
	queries.EXPECT().GetDefaultLanguage(gomock.Any()).Return(language.English)
	queries.EXPECT().CustomTextListByTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(&query.CustomTexts{}, nil)
}

func expectTemplateWithNotifyUserQueries(queries *mock.MockQueries, template string) {
	queries.EXPECT().GetNotifyUserByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.NotifyUser{
		ID:                 userID,
		ResourceOwner:      orgID,
		LastEmail:          lastEmail,
		VerifiedEmail:      verifiedEmail,
		PreferredLoginName: preferredLoginName,
		LastPhone:          lastPhone,
		VerifiedPhone:      verifiedPhone,
	}, nil)
	expectTemplateQueries(queries, template)
}

func expectTemplateWithNotifyUserQueriesSMS(queries *mock.MockQueries) {
	queries.EXPECT().GetNotifyUserByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.NotifyUser{
		ID:                 userID,
		ResourceOwner:      orgID,
		LastEmail:          lastEmail,
		VerifiedEmail:      verifiedEmail,
		PreferredLoginName: preferredLoginName,
		LastPhone:          lastPhone,
		VerifiedPhone:      verifiedPhone,
	}, nil)
	queries.EXPECT().GetInstanceRestrictions(gomock.Any()).Return(query.Restrictions{
		AllowedLanguages: []language.Tag{language.English},
	}, nil)
	queries.EXPECT().ActiveLabelPolicyByOrg(gomock.Any(), gomock.Any(), gomock.Any()).Return(&query.LabelPolicy{
		ID: policyID,
		Light: query.Theme{
			LogoURL: logoURL,
		},
	}, nil)
	queries.EXPECT().GetDefaultLanguage(gomock.Any()).Return(language.English)
	queries.EXPECT().CustomTextListByTemplate(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(2).Return(&query.CustomTexts{}, nil)
}

func cryptoValue(t *testing.T, ctrl *gomock.Controller, value string) (*crypto.MockEncryptionAlgorithm, *crypto.CryptoValue) {
	encAlg := crypto.NewMockEncryptionAlgorithm(ctrl)
	encAlg.EXPECT().Algorithm().AnyTimes().Return("enc")
	encAlg.EXPECT().EncryptionKeyID().AnyTimes().Return("id")
	encAlg.EXPECT().DecryptionKeyIDs().AnyTimes().Return([]string{"id"})
	encAlg.EXPECT().DecryptString(gomock.Any(), gomock.Any()).AnyTimes().Return(value, nil)
	encAlg.EXPECT().Encrypt(gomock.Any()).AnyTimes().Return(make([]byte, 0), nil)
	code, err := crypto.Encrypt(nil, encAlg)
	assert.NoError(t, err)
	return encAlg, code
}
