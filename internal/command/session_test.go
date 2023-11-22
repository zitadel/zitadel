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
				err:  caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Df4b3", "Errors.User.NotFound"),
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
				err: caos_errs.ThrowInvalidArgument(nil, "COMMAND-asEG4", "Errors.Session.PositiveLifetime"),
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
				eventstore:          eventstoreExpect(t, tt.expect...),
				idGenerator:         tt.fields.idGenerator,
				sessionTokenCreator: tt.fields.tokenCreator,
			}
			got, err := c.CreateSession(tt.args.ctx, tt.args.checks, tt.args.metadata, tt.args.userAgent, tt.args.lifetime)
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
		lifetime     time.Duration
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
				err: caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
			},
		},
		{
			"no change",
			fields{
				eventstore: eventstoreExpect(t,
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
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
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
				eventstore:           tt.fields.eventstore,
				sessionTokenVerifier: tt.fields.tokenVerifier,
			}
			got, err := c.UpdateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken, tt.args.checks, tt.args.metadata, tt.args.lifetime)
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionCommands{
					sessionWriteModel: &SessionWriteModel{State: domain.SessionStateTerminated},
				},
			},
			res{
				err: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Hewfq", "Errors.Session.Terminated"),
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
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
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
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands:   []SessionCommand{},
				},
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
		{
			"negative lifetime",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands:   []SessionCommand{},
					eventstore:        eventstoreExpect(t),
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
				err: caos_errs.ThrowInvalidArgument(nil, "COMMAND-asEG4", "Errors.Session.PositiveLifetime"),
			},
		},
		{
			"lifetime set",
			fields{
				eventstore: eventstoreExpect(t,
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
					eventstore:        eventstoreExpect(t),
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
			"set user, password, metadata and token",
			fields{
				eventstore: eventstoreExpect(t,
					expectPush(
						session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"userID", "org1", testNow,
						),
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
						CheckUser("userID", "org1"),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1"),
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
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1"),
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
				ctx: authz.NewMockContext("instance1", "", ""),
				checks: &SessionCommands{
					sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
					sessionCommands: []SessionCommand{
						CheckUser("userID", "org1"),
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
						session.NewUserCheckedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
							"userID", "org1", testNow),
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
						CheckUser("userID", "org1"),
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
				eventstore: tt.fields.eventstore,
			}
			got, err := c.updateSession(tt.args.ctx, tt.args.checks, tt.args.metadata, tt.args.lifetime)
			require.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.want, got)
		})
	}
}

func TestCheckTOTP(t *testing.T) {
	ctx := authz.NewMockContext("instance1", "org1", "user1")

	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	key, secret, err := domain.NewTOTPKey("example.com", "user1", cryptoAlg)
	require.NoError(t, err)

	sessAgg := &session.NewAggregate("session1", "instance1").Aggregate
	userAgg := &user.NewAggregate("user1", "org1").Aggregate

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
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Neil7", "Errors.User.UserIDMissing"),
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
			wantErr: caos_errs.ThrowPreconditionFailed(nil, "COMMAND-eej1U", "Errors.User.MFA.OTP.NotReady"),
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
				),
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
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
				),
			},
			wantEventCommands: []eventstore.Command{
				session.NewTOTPCheckedEvent(ctx, sessAgg, testNow),
			},
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
			err := CheckTOTP(tt.code)(ctx, cmd)
			require.ErrorIs(t, err, tt.wantErr)
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
				err: caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"),
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
				err: caos_errs.ThrowPermissionDenied(nil, "AUTHZ-HKJD33", "Errors.PermissionDenied"),
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
						caos_errs.ThrowInternal(nil, "id", "pushed failed"),
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
				err: caos_errs.ThrowInternal(nil, "id", "pushed failed"),
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
								"user1", "org1", testNow),
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
								"userID", "org1", testNow),
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
