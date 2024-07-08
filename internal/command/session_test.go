package command

import (
	"context"
	"io"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSessionCommands_getHumanWriteModel(t *testing.T) {
	userAggr := &user.NewAggregate("user1", "org1").Aggregate

	type fields struct {
		eventstore        func(*testing.T) *eventstore.Eventstore
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
				eventstore:        expectEventstore(),
				sessionWriteModel: &SessionWriteModel{},
			},
			res: res{
				want: nil,
				err:  zerrors.ThrowPreconditionFailed(nil, "COMMAND-eeR2e", "Errors.User.UserIDMissing"),
			},
		},
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
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
				eventstore: expectEventstore(
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
				err:  zerrors.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.User.NotFound"),
			},
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
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
			eventstore:        tt.fields.eventstore(t),
			sessionWriteModel: tt.fields.sessionWriteModel,
		}
		got, err := s.gethumanWriteModel(context.Background())
		require.ErrorIs(t, err, tt.res.err)
		assert.Equal(t, tt.res.want, got)
	}
}

func TestCommands_CreateSession(t *testing.T) {
	type fields struct {
		idGenerator  id_generator.Generator
		tokenCreator func(sessionID string) (string, string, error)
	}
	type args struct {
		ctx       context.Context
		checks    []SessionCommand
		metadata  map[string][]byte
		userAgent *domain.UserAgent
		lifetime  time.Duration
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
				idGenerator: mock.NewIDGeneratorExpectError(t, zerrors.ThrowInternal(nil, "id", "generator failed")),
			},
			args{
				ctx: context.Background(),
			},
			[]expect{},
			res{
				err: zerrors.ThrowInternal(nil, "id", "generator failed"),
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
				expectFilterError(zerrors.ThrowInternal(nil, "id", "filter failed")),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"negative lifetime",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
				tokenCreator: func(sessionID string) (string, string, error) {
					return "tokenID",
						"token",
						nil
				},
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				lifetime: -10 * time.Minute,
			},
			[]expect{
				expectFilter(),
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-asEG4", "Errors.Session.PositiveLifetime"),
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
				ctx: authz.NewMockContext("instance1", "", ""),
				userAgent: &domain.UserAgent{
					FingerprintID: gu.Ptr("fp1"),
					IP:            net.ParseIP("1.2.3.4"),
					Description:   gu.Ptr("firefox"),
					Header:        http.Header{"foo": []string{"bar"}},
				},
				lifetime: 10 * time.Minute,
			},
			[]expect{
				expectFilter(),
				expectPush(
					session.NewAddedEvent(context.Background(),
						&session.NewAggregate("sessionID", "instance1").Aggregate,
						&domain.UserAgent{
							FingerprintID: gu.Ptr("fp1"),
							IP:            net.ParseIP("1.2.3.4"),
							Description:   gu.Ptr("firefox"),
							Header:        http.Header{"foo": []string{"bar"}},
						},
					),
					session.NewLifetimeSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate, 10*time.Minute),
					session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
						"tokenID",
					),
				),
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{ResourceOwner: "instance1"},
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
				eventstore:          expectEventstore(tt.expect...)(t),
				sessionTokenCreator: tt.fields.tokenCreator,
			}
			id_generator.SetGenerator(tt.fields.idGenerator)
			got, err := c.CreateSession(tt.args.ctx, tt.args.checks, tt.args.metadata, tt.args.userAgent, tt.args.lifetime)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCommands_UpdateSession(t *testing.T) {
	type fields struct {
		eventstore    func(*testing.T) *eventstore.Eventstore
		tokenVerifier func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
	}
	type args struct {
		ctx       context.Context
		sessionID string
		checks    []SessionCommand
		metadata  map[string][]byte
		lifetime  time.Duration
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
				eventstore: expectEventstore(
					expectFilterError(zerrors.ThrowInternal(nil, "id", "filter failed")),
				),
			},
			args{
				ctx: context.Background(),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"no change",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID")),
					),
				),
				tokenVerifier: func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error) {
					return nil
				},
			},
			args{
				ctx:       context.Background(),
				sessionID: "sessionID",
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "instance1",
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
				eventstore:           tt.fields.eventstore(t),
				sessionTokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := c.UpdateSession(tt.args.ctx, tt.args.sessionID, tt.args.checks, tt.args.metadata, tt.args.lifetime)
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
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		checks   *SessionCommands
		metadata map[string][]byte
		lifetime time.Duration
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
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: &SessionWriteModel{State: domain.SessionStateTerminated},
				},
			},
			res{
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Hewfq", "Errors.Session.Terminated"),
			},
		},
		{
			"check failed",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
							return nil, zerrors.ThrowInternal(nil, "id", "check failed")
						},
					},
				},
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "check failed"),
			},
		},
		{
			"no change",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands:   []SessionCommand{},
				},
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{},
					ID:            "sessionID",
					NewToken:      "",
				},
			},
		},
		{
			"negative lifetime",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands:   []SessionCommand{},
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					now: func() time.Time {
						return testNow
					},
				},
				lifetime: -10 * time.Minute,
			},
			res{
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-asEG4", "Errors.Session.PositiveLifetime"),
			},
		},
		{
			"lifetime set",
			fields{
				eventstore: expectEventstore(
					expectPush(
						session.NewLifetimeSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							10*time.Minute,
						),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"tokenID",
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands:   []SessionCommand{},
					createToken: func(sessionID string) (string, string, error) {
						return "tokenID",
							"token",
							nil
					},
					now: func() time.Time {
						return testNow
					},
				},
				lifetime: 10 * time.Minute,
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "instance1",
					},
					ID:       "sessionID",
					NewToken: "token",
				},
			},
		},
		{
			"set user, invalid password",
			fields{
				eventstore: expectEventstore(
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
					expectFilter(), // recheck
					expectFilter(
						org.NewLockoutPolicyAddedEvent(context.Background(), &org.NewAggregate("org1").Aggregate, 0, 0, false),
					),
					expectPush(
						user.NewHumanPasswordCheckFailedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckPassword("invalid password"),
					},
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
				err: zerrors.ThrowInvalidArgument(nil, "COMMAND-3M0fs", "Errors.User.Password.Invalid"),
			},
		},
		{
			"set user, password, metadata and token",
			fields{
				eventstore: expectEventstore(
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
					expectFilter(), // recheck
					expectPush(
						session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"userID", "org1", testNow, &language.Afrikaans,
						),
						user.NewHumanPasswordCheckSucceededEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate, nil),
						session.NewPasswordCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							testNow,
						),
						session.NewMetadataSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							map[string][]byte{"key": []byte("value")},
						),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"tokenID",
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckPassword("password"),
					},
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
						ResourceOwner: "instance1",
					},
					ID:       "sessionID",
					NewToken: "token",
				},
			},
		},
		{
			"set user, intent not successful",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
						),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckIntent("intent", "aW50ZW50"),
					},
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
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-Df4bw", "Errors.Intent.NotSucceeded"),
			},
		},
		{
			"set user, intent not for user",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
						),
						eventFromEventPusher(
							idpintent.NewSucceededEvent(context.Background(),
								&idpintent.NewAggregate("id", "instance1").Aggregate,
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
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckIntent("intent", "aW50ZW50"),
					},
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
				err: zerrors.ThrowPreconditionFailed(nil, "COMMAND-O8xk3w", "Errors.Intent.OtherUser"),
			},
		},
		{
			"set user, intent incorrect token",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckIntent("intent2", "aW50ZW50"),
					},
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
				err: zerrors.ThrowPermissionDenied(nil, "CRYPTO-CRYPTO", "Errors.Intent.InvalidToken"),
			},
		},
		{
			"set user, intent, metadata and token",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
						),
						eventFromEventPusher(
							idpintent.NewSucceededEvent(context.Background(),
								&idpintent.NewAggregate("id", "instance1").Aggregate,
								nil,
								"idpUserID",
								"idpUsername",
								"userID",
								nil,
								"",
							),
						),
					),
					expectPush(
						session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"userID", "org1", testNow, &language.Afrikaans),
						session.NewIntentCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							testNow),
						session.NewMetadataSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							map[string][]byte{"key": []byte("value")}),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"tokenID"),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckIntent("intent", "aW50ZW50"),
					},
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
						ResourceOwner: "instance1",
					},
					ID:       "sessionID",
					NewToken: "token",
				},
			},
		},
		{
			"set user, intent (user not linked yet)",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"username", "", "", "", "", language.English, domain.GenderUnspecified, "", false),
						),
						eventFromEventPusher(
							idpintent.NewStartedEvent(context.Background(),
								&idpintent.NewAggregate("id", "instance1").Aggregate,
								nil,
								nil,
								"idpID",
							),
						),
						eventFromEventPusher(
							idpintent.NewSucceededEvent(context.Background(),
								&idpintent.NewAggregate("id", "instance1").Aggregate,
								nil,
								"idpUserID",
								"idpUsername",
								"",
								nil,
								"",
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewUserIDPLinkAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
								"idpID",
								"idpUsername",
								"idpUserID",
							),
						),
					),
					expectPush(
						session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"userID", "org1", testNow, &language.Afrikaans),
						session.NewIntentCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							testNow),
						session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"tokenID"),
					),
				),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1", &language.Afrikaans),
						CheckIntent("intent", "aW50ZW50"),
					},
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
			},
			res{
				want: &SessionChanged{
					ObjectDetails: &domain.ObjectDetails{
						ResourceOwner: "instance1",
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
				eventstore: tt.fields.eventstore(t),
			}
			tt.args.checks.eventstore = c.eventstore
			got, err := c.updateSession(tt.args.ctx, tt.args.checks, tt.args.metadata, tt.args.lifetime)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCheckTOTP(t *testing.T) {
	ctx := authz.NewMockContext("instance1", "org1", "user1")

	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	key, err := domain.NewTOTPKey("example.com", "user1")
	require.NoError(t, err)
	secret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	require.NoError(t, err)

	sessAgg := &session.NewAggregate("session1", "instance1").Aggregate
	userAgg := &user.NewAggregate("user1", "org1").Aggregate
	orgAgg := &org.NewAggregate("org1").Aggregate

	code, err := totp.GenerateCode(key.Secret(), testNow)
	require.NoError(t, err)

	type fields struct {
		sessionWriteModel *SessionWriteModel
		eventstore        func(*testing.T) *eventstore.Eventstore
	}

	tests := []struct {
		name              string
		code              string
		fields            fields
		wantEventCommands []eventstore.Command
		wantErrorCommands []eventstore.Command
		wantErr           error
	}{
		{
			name: "missing userID",
			code: code,
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					aggregate: sessAgg,
				},
				eventstore: expectEventstore(),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "COMMAND-8N9ds", "Errors.User.UserIDMissing"),
		},
		{
			name: "filter error",
			code: code,
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "otp not ready error",
			code: code,
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
					),
				),
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-3Mif9s", "Errors.User.MFA.OTP.NotReady"),
		},
		{
			name: "otp verify error",
			code: "foobar",
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(ctx, userAgg, "agent1"),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(org.NewLockoutPolicyAddedEvent(ctx, orgAgg, 0, 0, false)),
					),
				),
			},
			wantErrorCommands: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(ctx, userAgg, nil),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "otp verify error, locked",
			code: "foobar",
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(ctx, userAgg, "agent1"),
						),
					),
					expectFilter(), // recheck
					expectFilter(
						eventFromEventPusher(org.NewLockoutPolicyAddedEvent(ctx, orgAgg, 1, 1, false)),
					),
				),
			},
			wantErrorCommands: []eventstore.Command{
				user.NewHumanOTPCheckFailedEvent(ctx, userAgg, nil),
				user.NewUserLockedEvent(ctx, userAgg),
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "ok",
			code: code,
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(ctx, userAgg, "agent1"),
						),
					),
					expectFilter(), // recheck
				),
			},
			wantEventCommands: []eventstore.Command{
				user.NewHumanOTPCheckSucceededEvent(ctx, userAgg, nil),
				session.NewTOTPCheckedEvent(ctx, sessAgg, testNow),
			},
		},
		{
			name: "ok, but locked in the meantime",
			code: code,
			fields: fields{
				sessionWriteModel: &SessionWriteModel{
					UserID:        "user1",
					UserCheckedAt: testNow,
					aggregate:     sessAgg,
				},
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							user.NewHumanOTPAddedEvent(ctx, userAgg, secret),
						),
						eventFromEventPusher(
							user.NewHumanOTPVerifiedEvent(ctx, userAgg, "agent1"),
						),
					),
					expectFilter(
						user.NewUserLockedEvent(ctx, userAgg),
					),
				),
			},
			wantErr: zerrors.ThrowPreconditionFailed(nil, "COMMAND-SF3fg", "Errors.User.Locked"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &SessionCommands{
				sessionWriteModel: tt.fields.sessionWriteModel,
				eventstore:        tt.fields.eventstore(t),
				totpAlg:           cryptoAlg,
				now:               func() time.Time { return testNow },
			}
			gotCmds, err := CheckTOTP(tt.code)(ctx, cmd)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantErrorCommands, gotCmds)
			assert.Equal(t, tt.wantEventCommands, cmd.eventCommands)
		})
	}
}

func TestCommands_TerminateSession(t *testing.T) {
	type fields struct {
		eventstore      func(t *testing.T) *eventstore.Eventstore
		tokenVerifier   func(ctx context.Context, sessionToken, sessionID, tokenID string) (err error)
		checkPermission domain.PermissionCheck
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
				eventstore: expectEventstore(
					expectFilterError(zerrors.ThrowInternal(nil, "id", "filter failed")),
				),
			},
			args{
				ctx: context.Background(),
			},
			res{
				err: zerrors.ThrowInternal(nil, "id", "filter failed"),
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID")),
					),
				),
				tokenVerifier: newMockTokenVerifierInvalid(),
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				err: zerrors.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"missing permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID")),
					),
				),
				checkPermission: newMockPermissionCheckNotAllowed(),
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "",
			},
			res{
				err: zerrors.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
			},
		},
		{
			"not active",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID")),
						eventFromEventPusher(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate)),
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
					ResourceOwner: "instance1",
				},
			},
		},
		{
			"push failed",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID"),
						),
					),
					expectPushFailed(
						zerrors.ThrowInternal(nil, "id", "pushed failed"),
						session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate),
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
				err: zerrors.ThrowInternal(nil, "id", "pushed failed"),
			},
		},
		{
			"terminate with token",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID"),
						),
					),
					expectPush(
						session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate),
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
					ResourceOwner: "instance1",
				},
			},
		},
		{
			"terminate own session",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"user1", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID"),
						),
					),
					expectPush(
						session.NewTerminateEvent(authz.NewMockContext("instance1", "org1", "user1"), &session.NewAggregate("sessionID", "instance1").Aggregate),
					),
				),
			},
			args{
				ctx:          authz.NewMockContext("instance1", "org1", "user1"),
				sessionID:    "sessionID",
				sessionToken: "",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			"terminate with permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"userID", "org1", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
								"tokenID"),
						),
					),
					expectPush(
						session.NewTerminateEvent(authz.NewMockContext("instance1", "org1", "admin1"), &session.NewAggregate("sessionID", "instance1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:          authz.NewMockContext("instance1", "org1", "admin1"),
				sessionID:    "sessionID",
				sessionToken: "",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
		{
			"terminate session owned by org with permission",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(),
								&session.NewAggregate("sessionID", "instance1").Aggregate,
								&domain.UserAgent{
									FingerprintID: gu.Ptr("fp1"),
									IP:            net.ParseIP("1.2.3.4"),
									Description:   gu.Ptr("firefox"),
									Header:        http.Header{"foo": []string{"bar"}},
								},
							),
						),
						eventFromEventPusher(
							session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "org2").Aggregate,
								"userID", "", testNow, &language.Afrikaans),
						),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org2").Aggregate,
								"tokenID"),
						),
					),
					expectFilter(
						user.NewHumanAddedEvent(context.Background(), &user.NewAggregate("userID", "org1").Aggregate,
							"username", "firstname", "lastname", "nickname", "displayname", language.English, domain.GenderUnspecified, "email", false),
					),
					expectPush(
						session.NewTerminateEvent(authz.NewMockContext("instance1", "org1", "admin1"), &session.NewAggregate("sessionID", "instance1").Aggregate),
					),
				),
				checkPermission: newMockPermissionCheckAllowed(),
			},
			args{
				ctx:          authz.NewMockContext("instance1", "org1", "admin1"),
				sessionID:    "sessionID",
				sessionToken: "",
			},
			res{
				want: &domain.ObjectDetails{
					ResourceOwner: "instance1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:           tt.fields.eventstore(t),
				sessionTokenVerifier: tt.fields.tokenVerifier,
				checkPermission:      tt.fields.checkPermission,
			}
			got, err := c.TerminateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}
