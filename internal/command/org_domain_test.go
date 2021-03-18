package command

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestCommandSide_AddOrgDomain(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		domain *domain.OrgDomain
	}
	type res struct {
		want *domain.OrgDomain
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: &domain.OrgDomain{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain already exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "domain add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainAddedEvent(context.Background(),
								&org.NewAggregate("", "").Aggregate,
								"domain.ch",
							)),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
			res: res{
				want: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddOrgDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
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

func TestCommandSide_GenerateOrgDomainValidation(t *testing.T) {
	type fields struct {
		eventstore      *eventstore.Eventstore
		secretGenerator crypto.Generator
	}
	type args struct {
		ctx    context.Context
		domain *domain.OrgDomain
	}
	type res struct {
		wantToken string
		wantURL   string
		err       func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid validation type, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "domain already verified, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "add dns validation, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerificationAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeDNS,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							)),
						},
					),
				),
				secretGenerator: GetMockSecretGenerator(t),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				wantToken: "a",
				wantURL:   "_zitadel-challenge.domain.ch",
			},
		},
		{
			name: "add http validation, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerificationAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeHTTP,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							)),
						},
					),
				),
				secretGenerator: GetMockSecretGenerator(t),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeHTTP,
				},
			},
			res: res{
				wantToken: "a",
				wantURL:   "https://domain.ch/.well-known/zitadel-challenge/a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore,
				domainVerificationGenerator: tt.fields.secretGenerator,
			}
			token, url, err := r.GenerateOrgDomainValidation(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.wantToken, token)
				assert.Equal(t, tt.res.wantURL, url)
			}
		})
	}
}

func TestCommandSide_ValidateOrgDomain(t *testing.T) {
	type fields struct {
		eventstore           *eventstore.Eventstore
		idGenerator          id.Generator
		secretGenerator      crypto.Generator
		alg                  crypto.EncryptionAlgorithm
		domainValidationFunc func(domain, token, verifier string, checkType http.CheckType) error
	}
	type args struct {
		ctx            context.Context
		domain         *domain.OrgDomain
		claimedUserIDs []string
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "domain already verified, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "no code existing, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid domain verification, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerificationAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeDNS,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerificationFailedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
							)),
						},
					),
				),
				alg:                  crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				domainValidationFunc: invalidDomainVerification,
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain verification, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerificationAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeDNS,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
							)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("domain.ch")),
					),
				),
				alg:                  crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				domainValidationFunc: validDomainVerification,
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "domain verification, claimed users not found, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerificationAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeDNS,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							),
						),
					),
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
							)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("domain.ch")),
					),
				),
				alg:                  crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				domainValidationFunc: validDomainVerification,
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
				claimedUserIDs: []string{"user1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "domain verification, claimed users, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerificationAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
								domain.OrgDomainValidationTypeDNS,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("a"),
								},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							user.NewHumanAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org2").Aggregate,
								"username@domain.ch",
								"firstname",
								"lastname",
								"nickname",
								"displayname",
								language.German,
								domain.GenderUnspecified,
								"email",
								true,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							org.NewOrgIAMPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org2", "org2").Aggregate,
								false))),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainVerifiedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
							)),
							eventFromEventPusher(user.NewDomainClaimedEvent(context.Background(),
								&user.NewAggregate("user1", "org2").Aggregate,
								"tempid@temporary.zitadel.ch",
								"username@domain.ch",
								false,
							)),
						},
						uniqueConstraintsFromEventConstraint(org.NewAddOrgDomainUniqueConstraint("domain.ch")),
						uniqueConstraintsFromEventConstraint(user.NewRemoveUsernameUniqueConstraint("username@domain.ch", "org2", false)),
						uniqueConstraintsFromEventConstraint(user.NewAddUsernameUniqueConstraint("tempid@temporary.zitadel.ch", "org2", false)),
					),
				),
				alg:                  crypto.CreateMockEncryptionAlg(gomock.NewController(t)),
				domainValidationFunc: validDomainVerification,
				idGenerator:          id_mock.NewIDGeneratorExpectIDs(t, "tempid"),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
				claimedUserIDs: []string{"user1"},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:                  tt.fields.eventstore,
				domainVerificationGenerator: tt.fields.secretGenerator,
				domainVerificationAlg:       tt.fields.alg,
				domainVerificationValidator: tt.fields.domainValidationFunc,
				iamDomain:                   "zitadel.ch",
				idGenerator:                 tt.fields.idGenerator,
			}
			got, err := r.ValidateOrgDomain(tt.args.ctx, tt.args.domain, tt.args.claimedUserIDs...)
			if tt.res.err == nil {
				assert.NoError(t, err)
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

func TestCommandSide_SetPrimaryDomain(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		domain *domain.OrgDomain
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "domain not verified, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "set primary, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainPrimarySetEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch",
							)),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetPrimaryOrgDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
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

func TestCommandSide_RemoveOrgDomain(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		domain *domain.OrgDomain
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain:         "domain.ch",
					ValidationType: domain.OrgDomainValidationTypeDNS,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove verified domain, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainPrimarySetEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "remove domain, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch", false,
							)),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "remove verified domain, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewOrgAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"name",
							),
						),
						eventFromEventPusher(
							org.NewDomainAddedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
						eventFromEventPusher(
							org.NewDomainVerifiedEvent(context.Background(),
								&user.NewAggregate("user1", "org1").Aggregate,
								"domain.ch",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(org.NewDomainRemovedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"domain.ch", true,
							)),
						},
						uniqueConstraintsFromEventConstraint(org.NewRemoveOrgDomainUniqueConstraint("domain.ch")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				domain: &domain.OrgDomain{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "org1",
					},
					Domain: "domain.ch",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveOrgDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
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

func invalidDomainVerification(domain, token, verifier string, checkType http.CheckType) error {
	return caos_errs.ThrowInvalidArgument(nil, "HTTP-GH422", "Errors.Internal")
}

func validDomainVerification(domain, token, verifier string, checkType http.CheckType) error {
	return nil
}
