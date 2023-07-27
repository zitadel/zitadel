package command

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestSessionCommands_getHumanWriteModel(t *testing.T) {
	userAggr := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore        *eventstore.Eventstore
		sessionWriteModel *SessionWriteModel
	}
	type res struct {
		want *HumanWriteModel
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			name: "missing UID",
			fields: fields{
				eventstore:        &eventstore.Eventstore{},
				sessionWriteModel: &SessionWriteModel{},
			},
			res: res{
				want: nil,
				err:  caos_errs.ThrowPreconditionFailed(nil, "COMMAND-eeR2e", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilterError(io.ErrClosedPipe),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			res: res{
				want: nil,
				err:  io.ErrClosedPipe,
			},
		},
		{
			name: "removed user",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
						eventFromEventPusher(
							user.NewUserRemovedEvent(context.Background(),
								userAggr,
								"", nil, true,
							),
						),
					),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			res: res{
				want: nil,
				err:  caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.ie4Ai.NotFound"),
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								userAggr,
								"", "", "", "", "", language.Georgian,
								domain.GenderDiverse, "", true,
							),
						),
					),
				),
				sessionWriteModel: &SessionWriteModel{
					UserID: "user1",
				},
			},
			res: res{
				want: &HumanWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:   "user1",
						ResourceOwner: "org1",
						Events:        []eventstore.Event{},
					},
					PreferredLanguage: language.Georgian,
					Gender:            domain.GenderDiverse,
					UserState:         domain.UserStateActive,
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		s := &SessionCommands{
			eventstore:        tt.fields.eventstore,
			sessionWriteModel: tt.fields.sessionWriteModel,
		}
		got, err := s.gethumanWriteModel(context.Background())
		require.ErrorIs(t, err, tt.res.err)
		assert.Equal(t, tt.res.want, got)
	}
}

func TestCommands_CreateSession(t *testing.T) {
	type fields struct {
		idGenerator  id.Generator
		tokenCreator func(sessionID string) (string, string, error)
	}
	type args struct {
		ctx      context.Context
		checks   []SessionCommand
		domain   string
		metadata map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		expect []expect
		res    res
	}{
		{
			"id generator fails",
			fields{
				idGenerator: mock.NewIDGeneratorExpectError(t, caos_errs.ThrowInternal(nil, "id", "generator failed")),
			},
			args{
				ctx: context.Background(),
			},
			[]expect{},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "generator failed"),
			},
		},
		{
			"eventstore failed",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
			},
			args{
				ctx: context.Background(),
			},
			[]expect{
				expectFilterError(caos_errs.ThrowInternal(nil, "id", "filter failed")),
			},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"empty session",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
				tokenCreator: func(sessionID string) (string, string, error) {
					return "tokenID",
						"token",
						nil
				},
			},
			args{
				ctx: authz.NewMockContext("", "org1", ""),
			},
			[]expect{
				expectFilter(),
				expectPush(
					eventPusherToEvents(
						session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, ""),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
							"tokenID",
						),
					),
				),
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{ResourceOwner: "org1"},
					ID:            "sessionID",
					NewToken:      "token",
				},
			},
		},
		{
			"empty session with domain",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
				tokenCreator: func(sessionID string) (string, string, error) {
					return "tokenID",
						"token",
						nil
				},
			},
			args{
				ctx:    authz.NewMockContext("", "org1", ""),
				domain: "domain.tld",
			},
			[]expect{
				expectFilter(),
				expectPush(
					eventPusherToEvents(
						session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld"),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
							"tokenID",
						),
					),
				),
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{ResourceOwner: "org1"},
					ID:            "sessionID",
					NewToken:      "token",
				},
			},
		},
		// the rest is tested in the Test_updateSession
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          eventstoreExpect(t, tt.expect...),
				idGenerator:         tt.fields.idGenerator,
				sessionTokenCreator: tt.fields.tokenCreator,
			}
			got, err := c.CreateSession(tt.args.ctx, tt.args.checks, tt.args.domain, tt.args.metadata)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCommands_UpdateSession(t *testing.T) {
	type fields struct {
		eventstore    *eventstore.Eventstore
		tokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
	}
	type args struct {
		ctx          context.Context
		sessionID    string
		sessionToken string
		checks       []SessionCommand
		metadata     map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"eventstore failed",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilterError(caos_errs.ThrowInternal(nil, "id", "filter failed")),
				),
			},
			args{
				ctx: context.Background(),
			},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				err: caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"no change",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
					ID:       "sessionID",
					NewToken: "",
				},
			},
		},
		// the rest is tested in the Test_updateSession
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore,
				sessionTokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := c.UpdateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken, tt.args.checks, tt.args.metadata)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCommands_updateSession(t *testing.T) {
	decryption := func(err error) crypto.EncryptionAlgorithm {
		mCrypto := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
		mCrypto.EXPECT().EncryptionKeyID().Return("id")
		mCrypto.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
			func(code []byte, keyID string) (string, error) {
				if err != nil {
					return "", err
				}
				return string(code), nil
			})
		return mCrypto
	}

	testNow := time.Now()
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		checks   *SessionCommands
		metadata map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"terminated",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: &SessionWriteModel{State: domain.SessionStateTerminated},
				},
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMAND-SAjeh", "Errors.Session.Terminated"),
			},
		},
		{
			"check failed",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						func(ctx context.Context, cmd *SessionCommands) error {
							return caos_errs.ThrowInternal(nil, "id", "check failed")
						},
					},
				},
			},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "check failed"),
			},
		},
		{
			"no change",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands:   []SessionCommand{},
				},
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
					ID:       "sessionID",
					NewToken: "",
				},
			},
		},
		{
			"set user, password, metadata and token",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						eventPusherToEvents(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"userID", testNow),
							session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								testNow),
							session.NewMetadataSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								map[string][]byte{"key": []byte("value")}),
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID"),
						CheckPassword("password"),
					},
					eventstore: eventstoreExpect(t,
						expectFilter(
							eventFromEventPusher(
								user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
									"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
							),
							eventFromEventPusher(
								user.NewHumanPasswordChangedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
									"$plain$x$password", false, ""),
							),
						),
					),
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					hasher: mockPasswordHasher("x"),
					now: func() time.Time {
						return testNow
					},
				},
				metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
					ID:       "sessionID",
					NewToken: "token",
				},
			},
		},
		{
			"set user, intent not successful",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID"),
						CheckIntent("intent", "aW50ZW50"),
					},
					eventstore: eventstoreExpect(t,
						expectFilter(
							eventFromEventPusher(
								user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
									"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
							),
						),
					),
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					intentAlg: decryption(nil),
					now: func() time.Time {
						return testNow
					},
				},
				metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4bw", "Errors.Intent.NotSucceeded"),
			},
		},
		{
			"set user, intent not for user",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID"),
						CheckIntent("intent", "aW50ZW50"),
					},
					eventstore: eventstoreExpect(t,
						expectFilter(
							eventFromEventPusher(
								user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
									"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
							),
							eventFromEventPusher(
								idpintent.NewSucceededEvent(context.Background(), &idpintent.NewAggregate("intent", "org1").Aggregate,
									nil,
									"idpUserID",
									"idpUserName",
									"userID2",
									nil,
									"",
								),
							),
						),
					),
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					intentAlg: decryption(nil),
					now: func() time.Time {
						return testNow
					},
				},
				metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-O8xk3w", "Errors.Intent.OtherUser"),
			},
		},
		{
			"set user, intent incorrect token",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID"),
						CheckIntent("intent2", "aW50ZW50"),
					},
					eventstore: eventstoreExpect(t),
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					intentAlg: decryption(nil),
					now: func() time.Time {
						return testNow
					},
				},
				metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
			res{
				err: caos_errs.ThrowPermissionDenied(nil, "CRYPTO-CRYPTO", "Errors.Intent.InvalidToken"),
			},
		},
		{
			"set user, intent, metadata and token",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						eventPusherToEvents(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"userID", testNow),
							session.NewIntentCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								testNow),
							session.NewMetadataSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								map[string][]byte{"key": []byte("value")}),
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID"),
						CheckIntent("intent", "aW50ZW50"),
					},
					eventstore: eventstoreExpect(t,
						expectFilter(
							eventFromEventPusher(
								user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
									"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
							),
							eventFromEventPusher(
								idpintent.NewSucceededEvent(context.Background(), &idpintent.NewAggregate("intent", "org1").Aggregate,
									nil,
									"idpUserID",
									"idpUsername",
									"userID",
									nil,
									"",
								),
							),
						),
					),
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					intentAlg: decryption(nil),
					now: func() time.Time {
						return testNow
					},
				},
				metadata: map[string][]byte{
					"key": []byte("value"),
				},
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
					ID:       "sessionID",
					NewToken: "token",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := c.updateSession(tt.args.ctx, tt.args.checks, tt.args.metadata)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCommands_TerminateSession(t *testing.T) {
	type fields struct {
		eventstore    *eventstore.Eventstore
		tokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
	}
	type args struct {
		ctx          context.Context
		sessionID    string
		sessionToken string
	}
	type res struct {
		want *domain.ObjectDetails
		err  error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"eventstore failed",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilterError(caos_errs.ThrowInternal(nil, "id", "filter failed")),
				),
			},
			args{
				ctx: context.Background(),
			},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid")
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				err: caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"not active",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID")),
						eventFromEventPusher(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push failed",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID"),
						),
					),
					expectPushFailed(
						caos_errs.ThrowInternal(nil, "id", "pushed failed"),
						eventPusherToEvents(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				err: caos_errs.ThrowInternal(nil, "id", "pushed failed"),
			},
		},
		{
			"terminate",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate, "domain.tld")),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								"tokenID"),
						),
					),
					expectPush(
						eventPusherToEvents(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore,
				sessionTokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := c.TerminateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}
