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

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestCommands_AddUserPasskeyCode(t *testing.T) {
	type fields struct {
		newCode     cryptoCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     newCryptoCodeWithExpiry,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					)),
					expectPush([]*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"123", &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("passkey1"),
								}, time.Minute, "", false,
							),
						),
					}),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newCode:     tt.fields.newCode,
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.AddUserPasskeyCode(context.Background(), tt.args.userID, tt.args.resourceOwner, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_AddUserPasskeyCodeURLTemplate(t *testing.T) {
	type fields struct {
		newCode     cryptoCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
		urlTmpl       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.ObjectDetails
		wantErr error
	}{
		{
			name: "template error",
			fields: fields{
				newCode:    newCryptoCodeWithExpiry,
				eventstore: eventstoreExpect(t),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "{{",
			},
			wantErr: caos_errs.ThrowInvalidArgument(nil, "DOMAIN-oGh5e", "Errors.User.InvalidURLTemplate"),
		},
		{
			name: "id generator error",
			fields: fields{
				newCode:     newCryptoCodeWithExpiry,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.ResourceOwner}}&codeID={{.CodeID}}&code={{.Code}}",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					)),
					expectPush([]*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"123", &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("passkey1"),
								},
								time.Minute,
								"https://example.com/passkey/register?userID={{.UserID}}&orgID={{.ResourceOwner}}&codeID={{.CodeID}}&code={{.Code}}",
								false,
							),
						),
					}),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
				urlTmpl:       "https://example.com/passkey/register?userID={{.UserID}}&orgID={{.ResourceOwner}}&codeID={{.CodeID}}&code={{.Code}}",
			},
			want: &domain.ObjectDetails{
				ResourceOwner: "org1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newCode:     tt.fields.newCode,
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.AddUserPasskeyCodeURLTemplate(context.Background(), tt.args.userID, tt.args.resourceOwner, crypto.CreateMockEncryptionAlg(gomock.NewController(t)), tt.args.urlTmpl)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_AddUserPasskeyCodeReturn(t *testing.T) {
	type fields struct {
		newCode     cryptoCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.PasskeyCodeDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     newCryptoCodeWithExpiry,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					)),
					expectPush([]*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"123", &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("passkey1"),
								}, time.Minute, "", true,
							),
						),
					}),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.PasskeyCodeDetails{
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				CodeID: "123",
				Code:   "passkey1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newCode:     tt.fields.newCode,
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.AddUserPasskeyCodeReturn(context.Background(), tt.args.userID, tt.args.resourceOwner, crypto.CreateMockEncryptionAlg(gomock.NewController(t)))
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommands_addUserPasskeyCode(t *testing.T) {
	type fields struct {
		newCode     cryptoCodeFunc
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		userID        string
		resourceOwner string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.PasskeyCodeDetails
		wantErr error
	}{
		{
			name: "id generator error",
			fields: fields{
				newCode:     newCryptoCodeWithExpiry,
				eventstore:  eventstoreExpect(t),
				idGenerator: id_mock.NewIDGeneratorExpectError(t, io.ErrClosedPipe),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "crypto error",
			fields: fields{
				newCode:     newCryptoCodeWithExpiry,
				eventstore:  eventstoreExpect(t, expectFilterError(io.ErrClosedPipe)),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "filter query error",
			fields: fields{
				newCode: newCryptoCodeWithExpiry,
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(testSecretGeneratorAddedEvent(domain.SecretGeneratorTypePasswordlessInitCode))),
					expectFilterError(io.ErrClosedPipe),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "push error",
			fields: fields{
				newCode: mockCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					)),
					expectPushFailed(io.ErrClosedPipe, []*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"123", &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("passkey1"),
								}, time.Minute, "", false,
							),
						),
					}),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			fields: fields{
				newCode: mockCode("passkey1", time.Minute),
				eventstore: eventstoreExpect(t,
					expectFilter(eventFromEventPusher(
						user.NewHumanAddedEvent(context.Background(),
							&user.NewAggregate("user1", "org1").Aggregate,
							"username",
							"firstname",
							"lastname",
							"nickname",
							"displayname",
							language.German,
							domain.GenderUnspecified,
							"email@test.ch",
							true,
						),
					)),
					expectPush([]*repository.Event{
						eventFromEventPusher(
							user.NewHumanPasswordlessInitCodeRequestedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"123", &crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("passkey1"),
								}, time.Minute, "", false,
							),
						),
					}),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "123"),
			},
			args: args{
				userID:        "user1",
				resourceOwner: "org1",
			},
			want: &domain.PasskeyCodeDetails{
				ObjectDetails: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				CodeID: "123",
				Code:   "passkey1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				newCode:     tt.fields.newCode,
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.addUserPasskeyCode(context.Background(), tt.args.userID, tt.args.resourceOwner, crypto.CreateMockEncryptionAlg(gomock.NewController(t)), "", false)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
