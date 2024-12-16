package handlers

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

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
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func Test_userNotifierLegacy_reduceInitCodeAdded(t *testing.T) {
	expectMailSubject := "Initialize User"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanInitCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanInitialCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanInitCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanInitialCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:   code,
						Expiry: time.Hour,
					},
				}, w
		},
	}, {
		name: "button url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s/ui/login/user/init?authRequestID=%s&code=%s&loginname=%s&orgID=%s&passwordset=%t&userID=%s", eventOrigin, "", testCode, preferredLoginName, orgID, false, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanInitCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanInitialCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "button url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/user/init?authRequestID=%s&code=%s&loginname=%s&orgID=%s&passwordset=%t&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, "", testCode, preferredLoginName, orgID, false, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanInitCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanInitialCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:   code,
						Expiry: time.Hour,
					},
				}, w
		},
	}, {
		name: "button url without event trigger url with authRequestID",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/user/init?authRequestID=%s&code=%s&loginname=%s&orgID=%s&passwordset=%t&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, authRequestID, testCode, preferredLoginName, orgID, false, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanInitCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanInitialCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:          code,
						Expiry:        time.Hour,
						AuthRequestID: authRequestID,
					},
				}, w
		},
	}}
	// TODO: Why don't we have an url template on user.HumanInitialCodeAddedEvent?
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			queries := mock.NewMockQueries(ctrl)
			commands := mock.NewMockCommands(ctrl)
			f, a, w := tt.test(ctrl, queries, commands)
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reduceInitCodeAdded(a.event)
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

func Test_userNotifierLegacy_reduceEmailCodeAdded(t *testing.T) {
	expectMailSubject := "Verify email"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
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
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
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
	}, {
		name: "button url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s/ui/login/mail/verification?authRequestID=%s&code=%s&orgID=%s&userID=%s", eventOrigin, "", testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
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
		name: "button url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/mail/verification?authRequestID=%s&code=%s&orgID=%s&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, "", testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
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
	}, {
		name: "button url without event trigger url with authRequestID",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/mail/verification?authRequestID=%s&code=%s&orgID=%s&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, authRequestID, testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanEmailCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:          code,
						Expiry:        time.Hour,
						URLTemplate:   "",
						CodeReturned:  false,
						AuthRequestID: authRequestID,
					},
				}, w
		},
	}, {
		name: "button url with url template and event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			urlTemplate := "https://my.custom.url/org/{{.OrgID}}/user/{{.UserID}}/verify/{{.Code}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("https://my.custom.url/org/%s/user/%s/verify/%s", orgID, userID, testCode)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanEmailVerificationCodeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanEmailCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						URLTemplate:       urlTemplate,
						CodeReturned:      false,
						TriggeredAtOrigin: eventOrigin,
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
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reduceEmailCodeAdded(a.event)
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

func Test_userNotifierLegacy_reducePasswordCodeAdded(t *testing.T) {
	expectMailSubject := "Reset password"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
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
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
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
	}, {
		name: "button url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s/ui/login/password/init?authRequestID=%s&code=%s&orgID=%s&userID=%s", eventOrigin, "", testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
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
		name: "button url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/password/init?authRequestID=%s&code=%s&orgID=%s&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, "", testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
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
	}, {
		name: "button url without event trigger url with authRequestID",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/password/init?authRequestID=%s&code=%s&orgID=%s&userID=%s", externalProtocol, instancePrimaryDomain, externalPort, authRequestID, testCode, orgID, userID)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:          code,
						Expiry:        time.Hour,
						URLTemplate:   "",
						CodeReturned:  false,
						AuthRequestID: authRequestID,
					},
				}, w
		},
	}, {
		name: "button url with url template and event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			urlTemplate := "https://my.custom.url/org/{{.OrgID}}/user/{{.UserID}}/verify/{{.Code}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("https://my.custom.url/org/%s/user/%s/verify/%s", orgID, userID, testCode)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						URLTemplate:       urlTemplate,
						CodeReturned:      false,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "external code",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			expectContent := "We received a password reset request. Please use the button below to reset your password. (Code ) If you didn't ask for this mail, please ignore it."
			w.messageSMS = &messages.SMS{
				SenderPhoneNumber:    "senderNumber",
				RecipientPhoneNumber: lastPhone,
				Content:              expectContent,
			}
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordCodeSent(gomock.Any(), orgID, userID, &senders.CodeGeneratorInfo{ID: smsProviderID, VerificationID: verificationID}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanPasswordCodeAddedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              nil,
						Expiry:            0,
						URLTemplate:       "",
						CodeReturned:      false,
						NotificationType:  domain.NotificationTypeSms,
						GeneratorID:       smsProviderID,
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
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reducePasswordCodeAdded(a.event)
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

func Test_userNotifierLegacy_reduceDomainClaimed(t *testing.T) {
	expectMailSubject := "Domain has been claimed"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().UserDomainClaimedSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.DomainClaimedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().UserDomainClaimedSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.DomainClaimedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
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
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reduceDomainClaimed(a.event)
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

func Test_userNotifierLegacy_reducePasswordlessCodeRequested(t *testing.T) {
	expectMailSubject := "Add Passwordless Login"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanPasswordlessInitCodeSent(gomock.Any(), userID, orgID, codeID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordlessInitCodeRequestedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
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
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanPasswordlessInitCodeSent(gomock.Any(), userID, orgID, codeID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordlessInitCodeRequestedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						ID:           codeID,
						Code:         code,
						Expiry:       time.Hour,
						URLTemplate:  "",
						CodeReturned: false,
					},
				}, w
		},
	}, {
		name: "button url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectContent := fmt.Sprintf("%s/ui/login/login/passwordless/init?userID=%s&orgID=%s&codeID=%s&code=%s", eventOrigin, userID, orgID, codeID, testCode)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanPasswordlessInitCodeSent(gomock.Any(), userID, orgID, codeID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanPasswordlessInitCodeRequestedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
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
	}, {
		name: "button url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectContent := fmt.Sprintf("%s://%s:%d/ui/login/login/passwordless/init?userID=%s&orgID=%s&codeID=%s&code=%s", externalProtocol, instancePrimaryDomain, externalPort, userID, orgID, codeID, testCode)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanPasswordlessInitCodeSent(gomock.Any(), userID, orgID, codeID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &user.HumanPasswordlessInitCodeRequestedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						ID:           codeID,
						Code:         code,
						Expiry:       time.Hour,
						URLTemplate:  "",
						CodeReturned: false,
					},
				}, w
		},
	}, {
		name: "button url with url template and event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			urlTemplate := "https://my.custom.url/org/{{.OrgID}}/user/{{.UserID}}/verify/{{.Code}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("https://my.custom.url/org/%s/user/%s/verify/%s", orgID, userID, testCode)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().HumanPasswordlessInitCodeSent(gomock.Any(), userID, orgID, codeID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &user.HumanPasswordlessInitCodeRequestedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						ID:                codeID,
						Code:              code,
						Expiry:            time.Hour,
						URLTemplate:       urlTemplate,
						CodeReturned:      false,
						TriggeredAtOrigin: eventOrigin,
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
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reducePasswordlessCodeRequested(a.event)
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

func Test_userNotifierLegacy_reducePasswordChanged(t *testing.T) {
	expectMailSubject := "Password of user has changed"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			queries.EXPECT().NotificationPolicyByOrg(gomock.Any(), gomock.Any(), orgID, gomock.Any()).Return(&query.NotificationPolicy{
				PasswordChange: true,
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordChangeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.HumanPasswordChangedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{lastEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			queries.EXPECT().NotificationPolicyByOrg(gomock.Any(), gomock.Any(), orgID, gomock.Any()).Return(&query.NotificationPolicy{
				PasswordChange: true,
			}, nil)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			commands.EXPECT().PasswordChangeSent(gomock.Any(), orgID, userID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &user.HumanPasswordChangedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
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
			stmt, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reducePasswordChanged(a.event)
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

func Test_userNotifierLegacy_reduceOTPEmailChallenged(t *testing.T) {
	expectMailSubject := "Verify One-Time Password"
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s%s/%s/%s", eventOrigin, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{verifiedEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			commands.EXPECT().OTPEmailSent(gomock.Any(), userID, orgID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &session.OTPEmailChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						URLTmpl:           "",
						ReturnCode:        false,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.LogoURL}}"
			expectContent := fmt.Sprintf("%s://%s:%d%s/%s/%s", externalProtocol, instancePrimaryDomain, externalPort, assetsPath, policyID, logoURL)
			w.message = &messages.Email{
				Recipients: []string{verifiedEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, "testcode")
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			commands.EXPECT().OTPEmailSent(gomock.Any(), userID, orgID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &session.OTPEmailChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:       code,
						Expiry:     time.Hour,
						URLTmpl:    "",
						ReturnCode: false,
					},
				}, w
		},
	}, {
		name: "button url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s/otp/verify?loginName=%s&code=%s", eventOrigin, preferredLoginName, testCode)
			w.message = &messages.Email{
				Recipients: []string{verifiedEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			commands.EXPECT().OTPEmailSent(gomock.Any(), userID, orgID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &session.OTPEmailChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						URLTmpl:           "",
						ReturnCode:        false,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "button url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			testCode := "testcode"
			expectContent := fmt.Sprintf("%s://%s:%d/otp/verify?loginName=%s&code=%s", externalProtocol, instancePrimaryDomain, externalPort, preferredLoginName, testCode)
			w.message = &messages.Email{
				Recipients: []string{verifiedEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			commands.EXPECT().OTPEmailSent(gomock.Any(), userID, orgID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
				}, args{
					event: &session.OTPEmailChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:       code,
						Expiry:     time.Hour,
						ReturnCode: false,
					},
				}, w
		},
	}, {
		name: "button url with url template and event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			givenTemplate := "{{.URL}}"
			urlTemplate := "https://my.custom.url/user/{{.LoginName}}/verify"
			testCode := "testcode"
			expectContent := fmt.Sprintf("https://my.custom.url/user/%s/verify", preferredLoginName)
			w.message = &messages.Email{
				Recipients: []string{verifiedEmail},
				Subject:    expectMailSubject,
				Content:    expectContent,
			}
			codeAlg, code := cryptoValue(t, ctrl, testCode)
			expectTemplateWithNotifyUserQueries(queries, givenTemplate)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			commands.EXPECT().OTPEmailSent(gomock.Any(), userID, orgID).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
					userDataCrypto: codeAlg,
					SMSTokenCrypto: nil,
				}, args{
					event: &session.OTPEmailChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              code,
						Expiry:            time.Hour,
						ReturnCode:        false,
						URLTmpl:           urlTemplate,
						TriggeredAtOrigin: eventOrigin,
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
			_, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reduceSessionOTPEmailChallenged(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_userNotifierLegacy_reduceOTPSMSChallenged(t *testing.T) {
	tests := []struct {
		name string
		test func(*gomock.Controller, *mock.MockQueries, *mock.MockCommands) (fields, args, wantLegacy)
	}{{
		name: "asset url with event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			testCode := ""
			expiry := 0 * time.Hour
			expectContent := fmt.Sprintf(`%[1]s is your one-time password for %[2]s. Use it within the next %[3]s.
@%[2]s #%[1]s`, testCode, eventOriginDomain, expiry)
			w.messageSMS = &messages.SMS{
				SenderPhoneNumber:    "senderNumber",
				RecipientPhoneNumber: verifiedPhone,
				Content:              expectContent,
			}
			expectTemplateWithNotifyUserQueriesSMS(queries)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			commands.EXPECT().OTPSMSSent(gomock.Any(), userID, orgID, &senders.CodeGeneratorInfo{ID: smsProviderID, VerificationID: verificationID}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &session.OTPSMSChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:              nil,
						Expiry:            expiry,
						CodeReturned:      false,
						GeneratorID:       smsProviderID,
						TriggeredAtOrigin: eventOrigin,
					},
				}, w
		},
	}, {
		name: "asset url without event trigger url",
		test: func(ctrl *gomock.Controller, queries *mock.MockQueries, commands *mock.MockCommands) (f fields, a args, w wantLegacy) {
			testCode := ""
			expiry := 0 * time.Hour
			expectContent := fmt.Sprintf(`%[1]s is your one-time password for %[2]s. Use it within the next %[3]s.
@%[2]s #%[1]s`, testCode, instancePrimaryDomain, expiry)
			w.messageSMS = &messages.SMS{
				SenderPhoneNumber:    "senderNumber",
				RecipientPhoneNumber: verifiedPhone,
				Content:              expectContent,
			}
			expectTemplateWithNotifyUserQueriesSMS(queries)
			queries.EXPECT().SessionByID(gomock.Any(), gomock.Any(), userID, gomock.Any()).Return(&query.Session{}, nil)
			queries.EXPECT().SearchInstanceDomains(gomock.Any(), gomock.Any()).Return(&query.InstanceDomains{
				Domains: []*query.InstanceDomain{{
					Domain:    instancePrimaryDomain,
					IsPrimary: true,
				}},
			}, nil)
			commands.EXPECT().OTPSMSSent(gomock.Any(), userID, orgID, &senders.CodeGeneratorInfo{ID: smsProviderID, VerificationID: verificationID}).Return(nil)
			return fields{
					queries:  queries,
					commands: commands,
					es: eventstore.NewEventstore(&eventstore.Config{
						Querier: es_repo_mock.NewRepo(t).ExpectFilterEvents().MockQuerier,
					}),
				}, args{
					event: &session.OTPSMSChallengedEvent{
						BaseEvent: *eventstore.BaseEventFromRepo(&repository.Event{
							AggregateID:   userID,
							ResourceOwner: sql.NullString{String: orgID},
							CreationDate:  time.Now().UTC(),
						}),
						Code:         nil,
						Expiry:       expiry,
						CodeReturned: false,
						GeneratorID:  smsProviderID,
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
			_, err := newUserNotifierLegacy(t, ctrl, queries, f, a, w).reduceSessionOTPSMSChallenged(a.event)
			if w.err != nil {
				w.err(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type wantLegacy struct {
	message    *messages.Email
	messageSMS *messages.SMS
	err        assert.ErrorAssertionFunc
}

func newUserNotifierLegacy(t *testing.T, ctrl *gomock.Controller, queries *mock.MockQueries, f fields, a args, w wantLegacy) *userNotifierLegacy {
	queries.EXPECT().NotificationProviderByIDAndType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().Return(&query.DebugNotificationProvider{}, nil)
	smtpAlg, _ := cryptoValue(t, ctrl, "smtppw")
	channel := channel_mock.NewMockNotificationChannel(ctrl)
	if w.err == nil {
		if w.message != nil {
			w.message.TriggeringEvent = a.event
			channel.EXPECT().HandleMessage(w.message).Return(nil)
		}
		if w.messageSMS != nil {
			w.messageSMS.TriggeringEvent = a.event
			channel.EXPECT().HandleMessage(w.messageSMS).DoAndReturn(func(message *messages.SMS) error {
				message.VerificationID = gu.Ptr(verificationID)
				return nil
			})
		}
	}
	return &userNotifierLegacy{
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
	}
}
