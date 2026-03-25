package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_ChangeOrgSMTPConfigHTTP(t *testing.T) {
	type fields struct {
		eventstore                  func(t *testing.T) *eventstore.Eventstore
		newEncryptedCodeWithDefault encryptedCodeWithDefaultFunc
		defaultSecretGenerators     *SecretGenerators
	}
	type args struct {
		http *ChangeOrgSMTPConfigHTTP
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "resourceowner empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				http: &ChangeOrgSMTPConfigHTTP{},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-k7QCG1", "Errors.ResourceOwnerMissing"))
				},
			},
		},
		{
			name: "id empty, precondition error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				http: &ChangeOrgSMTPConfigHTTP{
					ResourceOwner: "ORG1",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-2MHkV1", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "smtp not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				http: &ChangeOrgSMTPConfigHTTP{
					ResourceOwner: "ORG1",
					ID:            "ID",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-xIrdl1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
		{
			name: "no changes, returns details",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgSMTPConfigHTTPAddedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
								"test",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
				),
			},
			args: args{
				http: &ChangeOrgSMTPConfigHTTP{
					ResourceOwner: "ORG1",
					ID:            "ID",
					Description:   "test",
					Endpoint:      "endpoint",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
		{
			name: "smtp config http change, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							org.NewOrgSMTPConfigHTTPAddedEvent(
								context.Background(),
								&org.NewAggregate("ORG1").Aggregate,
								"ID",
								"",
								"endpoint",
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						newOrgSMTPConfigHTTPChangedEvent(
							context.Background(),
							"ORG1",
							"ID",
							"test",
							"endpoint2",
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("87654321"),
							},
						),
					),
				),
				newEncryptedCodeWithDefault: mockEncryptedCodeWithDefault("87654321", time.Hour),
				defaultSecretGenerators:     &SecretGenerators{},
			},
			args: args{
				http: &ChangeOrgSMTPConfigHTTP{
					ResourceOwner:        "ORG1",
					ID:                   "ID",
					Description:          "test",
					Endpoint:             "endpoint2",
					ExpirationSigningKey: true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "ORG1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore(t),
				newEncryptedCodeWithDefault: tt.fields.newEncryptedCodeWithDefault,
				defaultSecretGenerators:     tt.fields.defaultSecretGenerators,
			}
			err := r.ChangeOrgSMTPConfigHTTP(context.Background(), tt.args.http)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, tt.args.http.Details)
			}
		})
	}
}
