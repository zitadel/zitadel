package handlers

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	es_repo_mock "github.com/zitadel/zitadel/internal/eventstore/repository/mock"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/notification/backchannel"
	"github.com/zitadel/zitadel/internal/notification/channels"
	channel_mock "github.com/zitadel/zitadel/internal/notification/channels/mock"
	"github.com/zitadel/zitadel/internal/notification/handlers/mock"
	"github.com/zitadel/zitadel/internal/notification/messages"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/queue"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
)

func Test_backChannelLogoutWorker_reduceNotificationRequested(t *testing.T) {
	testNow := time.Now()
	sessionLogoutAgg := &sessionlogout.NewAggregate(sessionID, instanceID).Aggregate
	type fields struct {
		es          func(*testing.T) *eventstore.Eventstore
		queue       func(*gomock.Controller) Queue
		commands    func(*gomock.Controller) Commands
		queries     func(*gomock.Controller) Queries
		channel     func(*gomock.Controller) channels.NotificationChannel
		idGenerator id.Generator
	}
	type args struct {
		job *river.Job[*backchannel.LogoutRequest]
	}
	type want struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "too old",
			fields: fields{
				es: expectEventstore(),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					return c
				},
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: time.Now().Add(-1 * time.Hour),
					},
					Args: &backchannel.LogoutRequest{
						Aggregate: sessionLogoutAgg,
					},
				},
			},
			want: want{
				err: new(river.JobCancelError),
			},
		},
		{
			name: "no oidc sessions with back channel logout",
			fields: fields{
				es: expectEventstore(
					expectFilter(),
				),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					return c
				},
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: testNow,
					},
					Args: &backchannel.LogoutRequest{
						Aggregate: sessionLogoutAgg,
						SessionID: sessionID,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "create jobs for oidc sessions with back channel logout",
			fields: fields{
				es: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							sessionlogout.NewBackChannelLogoutRegisteredEvent(
								context.Background(),
								sessionLogoutAgg,
								"oidc-session-id1",
								"user-id",
								"client-id1",
								"back-channel-logout-uri1",
							),
						),
						eventFromEventPusher(
							sessionlogout.NewBackChannelLogoutRegisteredEvent(
								context.Background(),
								sessionLogoutAgg,
								"oidc-session-id2",
								"user-id",
								"client-id2",
								"back-channel-logout-uri2",
							),
						),
					),
				),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					q.EXPECT().Insert(gomock.Any(),
						&backchannel.LogoutRequest{
							Aggregate:            sessionLogoutAgg,
							SessionID:            sessionID,
							TriggeredAtOrigin:    "",
							TriggeringEventType:  "",
							TokenID:              "id1",
							UserID:               "user-id",
							OIDCSessionID:        "oidc-session-id1",
							ClientID:             "client-id1",
							BackChannelLogoutURI: "back-channel-logout-uri1",
						},
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithQueueName(backchannel.QueueName))),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithMaxAttempts(1))),
					).Return(nil)
					q.EXPECT().Insert(gomock.Any(), &backchannel.LogoutRequest{
						Aggregate:            sessionLogoutAgg,
						SessionID:            sessionID,
						TriggeredAtOrigin:    "",
						TriggeringEventType:  "",
						TokenID:              "id2",
						UserID:               "user-id",
						OIDCSessionID:        "oidc-session-id2",
						ClientID:             "client-id2",
						BackChannelLogoutURI: "back-channel-logout-uri2",
					},
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithQueueName(backchannel.QueueName))),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithMaxAttempts(1))),
					).Return(nil)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					return c
				},
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1", "id2"),
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: testNow,
					},
					Args: &backchannel.LogoutRequest{
						Aggregate: sessionLogoutAgg,
						SessionID: sessionID,
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "send logout request",
			fields: fields{
				es: expectEventstore(),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					c.EXPECT().BackChannelLogoutSent(gomock.Any(), sessionID, "oidc-session-id", instanceID).Return(nil)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					q.EXPECT().GetActiveSigningWebKey(gomock.Any()).Return(
						&jose.JSONWebKey{
							Key:       privateKey,
							Algorithm: string(signingAlgorithm),
							Use:       "sig",
						}, nil)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					c.EXPECT().HandleMessage(gomock.Any()).DoAndReturn(
						func(message channels.Message) error {
							logoutMessage, ok := message.(*messages.Form)
							if !ok {
								ctrl.T.Errorf("unexpected message type: %T", message)
							}
							logoutToken, ok := logoutMessage.Serializable.(*LogoutTokenMessage)
							if !ok {
								ctrl.T.Errorf("unexpected serializable type: %T", logoutMessage.Serializable)
							}
							jws, err := jose.ParseSigned(logoutToken.LogoutToken, []jose.SignatureAlgorithm{jose.RS256})
							require.NoError(t, err)
							payload, err := jws.Verify(privateKey.Public())
							require.NoError(t, err)
							var claims oidc.LogoutTokenClaims
							err = json.Unmarshal(payload, &claims)
							require.NoError(t, err)
							assert.Equal(t, "", claims.Issuer)
							assert.Equal(t, "user-id", claims.Subject)
							assert.Equal(t, "client-id", claims.Audience[0])
							assert.WithinRange(t, claims.IssuedAt.AsTime(), testNow.Add(-2*time.Second), time.Now())
							assert.WithinRange(t, claims.Expiration.AsTime(), testNow.Add(-time.Second).Add(time.Hour), time.Now().Add(time.Hour))
							assert.Equal(t, "id1", claims.JWTID)
							assert.NotNil(t, claims.Events["http://schemas.openid.net/event/backchannel-logout"])
							assert.Equal(t, sessionID, claims.SessionID)

							return nil
						})
					return c
				},
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t),
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: testNow,
					},
					Args: &backchannel.LogoutRequest{
						Aggregate:            sessionLogoutAgg,
						SessionID:            sessionID,
						TriggeredAtOrigin:    "",
						TriggeringEventType:  "",
						TokenID:              "id1",
						UserID:               "user-id",
						OIDCSessionID:        "oidc-session-id",
						ClientID:             "client-id",
						BackChannelLogoutURI: "back-channel-logout-uri",
					},
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "job creation failed",
			fields: fields{
				es: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							sessionlogout.NewBackChannelLogoutRegisteredEvent(
								context.Background(),
								sessionLogoutAgg,
								"oidc-session-id1",
								"user-id",
								"client-id1",
								"back-channel-logout-uri1",
							),
						),
					),
				),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					q.EXPECT().Insert(gomock.Any(),
						gomock.Any(),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithQueueName(backchannel.QueueName))),
						gomock.AssignableToTypeOf(reflect.TypeOf(queue.WithMaxAttempts(1))),
					).Return(assert.AnError)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					return c
				},
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "id1"),
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: testNow,
					},
					Args: &backchannel.LogoutRequest{
						Aggregate: sessionLogoutAgg,
						SessionID: sessionID,
					},
				},
			},
			want: want{
				err: assert.AnError,
			},
		},

		{
			name: "send logout request failed",
			fields: fields{
				es: expectEventstore(),
				queue: func(ctrl *gomock.Controller) Queue {
					q := mock.NewMockQueue(ctrl)
					return q
				},
				commands: func(ctrl *gomock.Controller) Commands {
					c := mock.NewMockCommands(ctrl)
					return c
				},
				queries: func(ctrl *gomock.Controller) Queries {
					q := mock.NewMockQueries(ctrl)
					q.EXPECT().GetActiveSigningWebKey(gomock.Any()).Return(
						&jose.JSONWebKey{
							Key:       privateKey,
							Algorithm: string(signingAlgorithm),
							Use:       "sig",
						}, nil)
					return q
				},
				channel: func(ctrl *gomock.Controller) channels.NotificationChannel {
					c := channel_mock.NewMockNotificationChannel(ctrl)
					c.EXPECT().HandleMessage(gomock.Any()).Return(assert.AnError)
					return c
				},
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t),
			},
			args: args{
				job: &river.Job[*backchannel.LogoutRequest]{
					JobRow: &rivertype.JobRow{
						CreatedAt: testNow,
					},
					Args: &backchannel.LogoutRequest{
						Aggregate:            sessionLogoutAgg,
						SessionID:            sessionID,
						TriggeredAtOrigin:    "",
						TriggeringEventType:  "",
						TokenID:              "id1",
						UserID:               "user-id",
						OIDCSessionID:        "oidc-session-id",
						ClientID:             "client-id",
						BackChannelLogoutURI: "back-channel-logout-uri",
					},
				},
			},
			want: want{
				err: assert.AnError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			err := newBackChannelLogoutWorker(
				tt.fields.queries(ctrl),
				tt.fields.commands(ctrl),
				tt.fields.es(t),
				tt.fields.queue(ctrl),
				tt.fields.channel(ctrl),
				tt.fields.idGenerator,
				func() time.Time { return testNow },
			).Work(
				authz.WithInstanceID(context.Background(), instanceID),
				tt.args.job,
			)
			assert.ErrorIs(t, err, tt.want.err)
		})
	}
}

func newBackChannelLogoutWorker(queries Queries, commands Commands, es *eventstore.Eventstore, queue Queue, channel channels.NotificationChannel, idGenerator id.Generator, testNow func() time.Time) *BackChannelLogoutWorker {
	return &BackChannelLogoutWorker{
		commands: commands,
		queries: NewNotificationQueries(
			queries,
			es,
			externalDomain,
			externalPort,
			externalSecure,
			"",
			nil,
			nil,
			nil,
		),
		eventstore: es,
		queue:      queue,
		channels: &notificationChannels{
			Chain: *senders.ChainChannels(channel),
		},
		config: &BackChannelLogoutWorkerConfig{
			Workers:             1,
			TransactionDuration: 5 * time.Second,
			MaxTtl:              5 * time.Minute,
			MaxAttempts:         1,
			TokenLifetime:       time.Hour,
		},
		now:         testNow,
		idGenerator: idGenerator,
	}
}

type expect func(mockRepository *es_repo_mock.MockRepository)

func expectEventstore(expects ...expect) func(*testing.T) *eventstore.Eventstore {
	return func(t *testing.T) *eventstore.Eventstore {
		m := es_repo_mock.NewRepo(t)
		for _, e := range expects {
			e(m)
		}
		es := eventstore.NewEventstore(
			&eventstore.Config{
				Querier: m.MockQuerier,
				Pusher:  m.MockPusher,
			},
		)
		return es
	}
}

func expectFilter(events ...eventstore.Event) expect {
	return func(m *es_repo_mock.MockRepository) {
		m.ExpectFilterEvents(events...)
	}
}

func eventFromEventPusher(event eventstore.Command) *repository.Event {
	data, _ := eventstore.EventData(event)
	return &repository.Event{
		InstanceID:    event.Aggregate().InstanceID,
		ID:            "",
		Seq:           0,
		CreationDate:  time.Time{},
		Typ:           event.Type(),
		Data:          data,
		EditorUser:    event.Creator(),
		Version:       event.Aggregate().Version,
		AggregateID:   event.Aggregate().ID,
		AggregateType: event.Aggregate().Type,
		ResourceOwner: sql.NullString{String: event.Aggregate().ResourceOwner, Valid: event.Aggregate().ResourceOwner != ""},
		Constraints:   event.UniqueConstraints(),
	}
}

var (
	privateKey = func() *rsa.PrivateKey {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return privateKey
	}()
	signingAlgorithm = jose.RS256
)
