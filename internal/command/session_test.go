package command

import (
	"context"
	"errors"
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
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_CreateSession(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		createToken func() (*crypto.CryptoValue, string, error)
	}
	type args struct {
		ctx      context.Context
		checks   []SessionCheck
		metadata map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  func(err error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
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
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "generator failed"))
				},
			},
		},
		{
			"eventstore failed",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
				eventstore: eventstoreExpect(t,
					expectFilterError(caos_errs.ThrowInternal(nil, "id", "filter failed")),
				),
			},
			args{
				ctx: context.Background(),
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "filter failed"))
				},
			},
		},
		{
			"empty session",
			fields{
				idGenerator: mock.NewIDGeneratorExpectIDs(t, "sessionID"),
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPush(
						eventPusherToEvents(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate),
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								}),
						),
					),
				),
				createToken: func() (*crypto.CryptoValue, string, error) {
					return &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("token"),
						},
						"token",
						nil
				},
			},
			args{
				ctx: authz.NewMockContext("", "org1", ""),
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
				eventstore:   tt.fields.eventstore,
				idGenerator:  tt.fields.idGenerator,
				tokenCreator: tt.fields.createToken,
			}
			got, err := c.CreateSession(tt.args.ctx, tt.args.checks, tt.args.metadata)
			if tt.res.err == nil {
				require.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommands_UpdateSession(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		sessionAlg crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx          context.Context
		sessionID    string
		sessionToken string
		checks       []SessionCheck
		metadata     map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  func(err error) bool
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
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "filter failed"))
				},
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								})),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"))
				},
			},
		},
		{
			"no change",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								})),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
				eventstore: tt.fields.eventstore,
				sessionAlg: tt.fields.sessionAlg,
			}
			got, err := c.UpdateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken, tt.args.checks, tt.args.metadata)
			if tt.res.err == nil {
				require.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommands_updateSession(t *testing.T) {
	testNow := time.Now()
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		checks   *SessionChecks
		metadata map[string][]byte
	}
	type res struct {
		want *SessionChanged
		err  func(err error) bool
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
				checks: &SessionChecks{
					sessionWriteModel: &SessionWriteModel{State: domain.SessionStateTerminated},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowPreconditionFailed(nil, "COMAND-SAjeh", "Errors.Session.Terminated"))
				},
			},
		},
		{
			"check failed",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionChecks{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					checks: []SessionCheck{
						func(ctx context.Context, cmd *SessionChecks) error {
							return caos_errs.ThrowInternal(nil, "id", "check failed")
						},
					},
				},
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "check failed"))
				},
			},
		},
		{
			"no change",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				checks: &SessionChecks{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					checks:            []SessionCheck{},
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
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								}),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				checks: &SessionChecks{
					sessionWriteModel: NewSessionWriteModel("sessionID", "org1"),
					checks: []SessionCheck{
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
									&crypto.CryptoValue{
										CryptoType: crypto.TypeHash,
										Algorithm:  "hash",
										KeyID:      "",
										Crypted:    []byte("password"),
									}, false, ""),
							),
						),
					),
					createToken: func() (*crypto.CryptoValue, string, error) {
						return &crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("token"),
							},
							"token",
							nil
					},
					userPasswordAlg: crypto.CreateMockHashAlg(gomock.NewController(t)),
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
			if tt.res.err == nil {
				require.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommands_TerminateSession(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		sessionAlg crypto.EncryptionAlgorithm
	}
	type args struct {
		ctx          context.Context
		sessionID    string
		sessionToken string
	}
	type res struct {
		want *SessionChanged
		err  func(err error) bool
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
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "filter failed"))
				},
			},
		},
		{
			"invalid session token",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								})),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "invalid",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowPermissionDenied(nil, "COMMAND-sGr42", "Errors.Session.Token.Invalid"))
				},
			},
		},
		{
			"not active",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								})),
						eventFromEventPusher(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
		{
			"push failed",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								}),
						),
					),
					expectPushFailed(
						caos_errs.ThrowInternal(nil, "id", "pushed failed"),
						eventPusherToEvents(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
			},
			args{
				ctx:          context.Background(),
				sessionID:    "sessionID",
				sessionToken: "token",
			},
			res{
				err: func(err error) bool {
					return errors.Is(err, caos_errs.ThrowInternal(nil, "id", "pushed failed"))
				},
			},
		},
		{
			"terminate",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
						eventFromEventPusher(
							session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("token"),
								}),
						),
					),
					expectPush(
						eventPusherToEvents(
							session.NewTerminateEvent(context.Background(), &session.NewAggregate("sessionID", "org1").Aggregate)),
					),
				),
				sessionAlg: crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
				sessionAlg: tt.fields.sessionAlg,
			}
			got, err := c.TerminateSession(tt.args.ctx, tt.args.sessionID, tt.args.sessionToken)
			if tt.res.err == nil {
				require.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
