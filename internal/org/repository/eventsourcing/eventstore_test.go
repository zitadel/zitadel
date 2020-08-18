package eventsourcing

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/api/authz"
	http_util "github.com/caos/zitadel/internal/api/http"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_mock "github.com/caos/zitadel/internal/eventstore/mock"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
)

type testOrgEventstore struct {
	OrgEventstore
	mockEventstore *es_mock.MockEventstore
}

func newTestEventstore(t *testing.T) *testOrgEventstore {
	ctrl := gomock.NewController(t)
	mock := es_mock.NewMockEventstore(ctrl)
	verificationAlgorithm := crypto.NewMockEncryptionAlgorithm(ctrl)
	verificationGenerator := crypto.NewMockGenerator(ctrl)
	return &testOrgEventstore{
		OrgEventstore: OrgEventstore{
			Eventstore:            mock,
			verificationAlgorithm: verificationAlgorithm,
			verificationGenerator: verificationGenerator,
		},
		mockEventstore: mock}
}

func (es *testOrgEventstore) expectFilterEvents(events []*es_models.Event, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().FilterEvents(gomock.Any(), gomock.Any()).Return(events, err)
	return es
}

func (es *testOrgEventstore) expectPushEvents(startSequence uint64, err error) *testOrgEventstore {
	es.mockEventstore.EXPECT().PushAggregates(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, aggregates ...*es_models.Aggregate) error {
			for _, aggregate := range aggregates {
				for _, event := range aggregate.Events {
					event.Sequence = startSequence
					startSequence++
				}
			}
			return err
		})
	return es
}

func (es *testOrgEventstore) expectAggregateCreator() *testOrgEventstore {
	es.mockEventstore.EXPECT().AggregateCreator().Return(es_models.NewAggregateCreator("test"))
	return es
}

func (es *testOrgEventstore) expectGenerateVerification(r rune) *testOrgEventstore {
	generator, _ := es.verificationGenerator.(*crypto.MockGenerator)
	generator.EXPECT().Length().Return(uint(2))
	generator.EXPECT().Runes().Return([]rune("aa"))
	generator.EXPECT().Alg().Return(es.verificationAlgorithm)
	return es
}

func (es *testOrgEventstore) expectEncrypt() *testOrgEventstore {
	algorithm, _ := es.verificationAlgorithm.(*crypto.MockEncryptionAlgorithm)
	algorithm.EXPECT().Encrypt(gomock.Any()).DoAndReturn(
		func(value []byte) (*crypto.CryptoValue, error) {
			return &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "enc",
				KeyID:      "id",
				Crypted:    value,
			}, nil
		})
	algorithm.EXPECT().Algorithm().Return("enc")
	algorithm.EXPECT().EncryptionKeyID().Return("id")
	return es
}

func (es *testOrgEventstore) expectDecrypt() *testOrgEventstore {
	algorithm, _ := es.verificationAlgorithm.(*crypto.MockEncryptionAlgorithm)
	algorithm.EXPECT().Algorithm().AnyTimes().Return("enc")
	algorithm.EXPECT().DecryptionKeyIDs().Return([]string{"id"})
	algorithm.EXPECT().DecryptString(gomock.Any(), gomock.Any()).DoAndReturn(
		func(value []byte, id string) (string, error) {
			return string(value), nil
		})
	return es
}

func (es *testOrgEventstore) expectVerification(success bool) *testOrgEventstore {
	es.verificationValidator = func(_, _, _ string, _ http_util.CheckType) error {
		if success {
			return nil
		}
		return errors.ThrowInvalidArgument(nil, "id", "invalid token")
	}
	return es
}

func TestOrgEventstore_OrgByID(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "no aggregate id in input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no events found success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, nil)},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name:   "filter fail",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message"))},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6},
			}, nil)},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "hodor-org", ChangeDate: time.Now(), CreationDate: time.Now()}},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.OrgByID(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_DeactivateOrg(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already inactive error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.DeactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ReactivateOrg(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx   context.Context
		orgID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgInactiveEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "org already active error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil)},
			args: args{
				ctx:   authz.NewMockContext("user", "org"),
				orgID: "hodor",
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.ReactivateOrg(tt.args.ctx, tt.args.orgID)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_GenerateOrgDomainValidation(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		token string
		url   string
		isErr func(error) bool
	}
	type args struct {
		ctx    context.Context
		domain *org_model.OrgDomain
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no domain",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: nil,
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "validation type invalid error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil).
				expectGenerateVerification(65).
				expectEncrypt(),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org"},
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "org doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsNotFound,
			},
		},
		{
			name: "domain doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "domain already verified error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerifiedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent()}, nil).
				expectGenerateVerification(65).
				expectEncrypt().
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				token: "",
				url:   "",
				isErr: errors.IsInternal,
			},
		},
		{
			name: "push correct http",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent()}, nil).
				expectGenerateVerification(65).
				expectEncrypt().
				expectAggregateCreator().
				expectPushEvents(6, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				token: "aa",
				url:   "https://hodor.org/.well-known/zitadel-challenge/aa",
				isErr: nil,
			},
		},
		{
			name: "push correct dns",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent()}, nil).
				expectGenerateVerification(65).
				expectEncrypt().
				expectAggregateCreator().
				expectPushEvents(6, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeDNS},
			},
			res: res{
				token: "aa",
				url:   "_zitadel-challenge.hodor.org",
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, url, err := tt.fields.Eventstore.GenerateOrgDomainValidation(tt.args.ctx, tt.args.domain)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			assert.Equal(t, tt.res.token, token)
			assert.Equal(t, tt.res.url, url)
		})
	}
}

func TestOrgEventstore_ValidateOrgDomain(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr func(error) bool
	}
	type args struct {
		ctx    context.Context
		domain *org_model.OrgDomain
		users  func(ctx context.Context, domain string) ([]*es_models.Aggregate, error)
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no domain",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: nil,
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "org doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsNotFound,
			},
		},
		{
			name: "domain doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "domain already verified error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerifiedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "domain validation not created error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "verification fails",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerificationAddedEvent("token")}, nil).
				expectDecrypt().
				expectVerification(false).
				expectAggregateCreator().
				expectPushEvents(0, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "verification and push fails",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerificationAddedEvent("token")}, nil).
				expectDecrypt().
				expectVerification(false).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "(user) aggregate fails",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerificationAddedEvent("token")}, nil).
				expectDecrypt().
				expectVerification(true).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
				users: func(ctx context.Context, domain string) ([]*es_models.Aggregate, error) {
					return nil, errors.ThrowInternal(nil, "id", "internal error")
				},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerificationAddedEvent("token")}, nil).
				expectDecrypt().
				expectVerification(true).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerificationAddedEvent("token")}, nil).
				expectDecrypt().
				expectVerification(true).
				expectAggregateCreator().
				expectPushEvents(6, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.Eventstore.ValidateOrgDomain(tt.args.ctx, tt.args.domain, tt.args.users)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
		})
	}
}

func TestOrgEventstore_SetPrimaryOrgDomain(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr func(error) bool
	}
	type args struct {
		ctx    context.Context
		domain *org_model.OrgDomain
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no domain",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: nil,
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "org doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsNotFound,
			},
		},
		{
			name: "domain doesn't exist error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "domain not verified error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent()}, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerifiedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test")),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{orgCreatedEvent(), orgDomainAddedEvent(), orgDomainVerifiedEvent()}, nil).
				expectAggregateCreator().
				expectPushEvents(6, nil),
			},
			args: args{
				ctx:    authz.NewMockContext("org", "user"),
				domain: &org_model.OrgDomain{ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org"}, Domain: "hodor.org", ValidationType: org_model.OrgDomainValidationTypeHTTP},
			},
			res: res{
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.Eventstore.SetPrimaryOrgDomain(tt.args.ctx, tt.args.domain)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
		})
	}
}

func TestOrgEventstore_OrgMemberByIDs(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no aggregate id in input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4}, UserID: "asdf"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no aggregate id in input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "asdf"}},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name:   "no events found success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, nil)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 4,
				isErr:            nil,
			},
		},
		{
			name:   "filter fail",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{}, errors.ThrowInternal(nil, "EVENT-SAa1O", "message"))},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants"}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "new events found and added success",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: model.OrgMemberChanged},
			}, nil)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "banana"},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "not member of org error",
			fields: fields{Eventstore: newTestEventstore(t).expectFilterEvents([]*es_models.Event{
				{Sequence: 6, Data: []byte("{\"userId\": \"banana\", \"roles\": [\"bananaa\"]}"), Type: model.OrgMemberAdded},
				{Sequence: 7, Data: []byte("{\"userId\": \"apple\"}"), Type: model.OrgMemberRemoved},
			}, nil)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{ObjectRoot: es_models.ObjectRoot{Sequence: 4, AggregateID: "plants", ChangeDate: time.Now(), CreationDate: time.Now()}, UserID: "apple"},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.OrgMemberByIDs(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_AddOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		expectedSequence uint64
		isErr            func(error) bool
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "push failed",
			fields: fields{Eventstore: newTestEventstore(t).
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
				}, nil).
				expectAggregateCreator().
				expectPushEvents(0, errors.ThrowInternal(nil, "EVENT-S8WzW", "test"))},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "push correct",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
				}, nil).
				expectPushEvents(6, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 6,
				isErr:            nil,
			},
		},
		{
			name: "member already exists error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						Type:     model.OrgMemberAdded,
						Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence: 6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
		{
			name: "member deleted success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectPushEvents(10, nil).
				expectFilterEvents([]*es_models.Event{
					{
						Type:     model.OrgMemberAdded,
						Data:     []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence: 6,
					},
					{
						Type:     model.OrgMemberRemoved,
						Data:     []byte(`{"userId": "hodor"}`),
						Sequence: 10,
					},
				}, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 10,
				isErr:            nil,
			},
		},
		{
			name: "org not exists error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents(nil, nil).
				expectPushEvents(0, errors.ThrowAlreadyExists(nil, "EVENT-yLTI6", "weiss nöd wie teste")),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{
						Sequence:    4,
						AggregateID: "hodor-org",
					},
					UserID: "hodor",
					Roles:  []string{"nix"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorAlreadyExists,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Eventstore.AddOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_ChangeOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr            func(error) bool
		expectedSequence uint64
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsPreconditionFailed,
			},
		},
		{
			name: "member not found error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsNotFound,
			},
		},
		{
			name: "member found no changes error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "push error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 0,
				isErr:            errors.IsInternal,
			},
		},
		{
			name: "change success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(7, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master of desaster"},
				},
			},
			res: res{
				expectedSequence: 7,
				isErr:            nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &OrgEventstore{
				Eventstore: tt.fields.Eventstore,
			}
			got, err := es.ChangeOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if got == nil && tt.res.expectedSequence != 0 {
				t.Errorf("org should not be nil but was %v", got)
				t.FailNow()
			}
			if tt.res.expectedSequence != 0 && tt.res.expectedSequence != got.Sequence {
				t.Errorf("org should have sequence %d but had %d", tt.res.expectedSequence, got.Sequence)
			}
		})
	}
}

func TestOrgEventstore_RemoveOrgMember(t *testing.T) {
	type fields struct {
		Eventstore *testOrgEventstore
	}
	type res struct {
		isErr func(error) bool
	}
	type args struct {
		ctx    context.Context
		member *org_model.OrgMember
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input member",
			fields: fields{Eventstore: newTestEventstore(t)},
			args: args{
				ctx:    authz.NewMockContext("user", "org"),
				member: nil,
			},
			res: res{
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "member not found error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "brudi", "roles": ["master of desaster"]}`),
						Sequence:    6,
					},
				}, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
					Roles:      []string{"master"},
				},
			},
			res: res{
				isErr: nil,
			},
		},
		{
			name: "push error",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(0, errors.ThrowInternal(nil, "PEVENT-3wqa2", "test")),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: errors.IsInternal,
			},
		},
		{
			name: "remove success",
			fields: fields{Eventstore: newTestEventstore(t).
				expectAggregateCreator().
				expectFilterEvents([]*es_models.Event{
					{
						AggregateID: "hodor-org",
						Type:        model.OrgAdded,
						Sequence:    4,
						Data:        []byte("{}"),
					},
					{
						AggregateID: "hodor-org",
						Type:        model.OrgMemberAdded,
						Data:        []byte(`{"userId": "hodor", "roles": ["master"]}`),
						Sequence:    6,
					},
				}, nil).
				expectPushEvents(7, nil),
			},
			args: args{
				ctx: authz.NewMockContext("user", "org"),
				member: &org_model.OrgMember{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "hodor-org", Sequence: 5},
					UserID:     "hodor",
				},
			},
			res: res{
				isErr: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &OrgEventstore{
				Eventstore: tt.fields.Eventstore,
			}
			err := es.RemoveOrgMember(tt.args.ctx, tt.args.member)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got:%T %v", err, err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
		})
	}
}

func orgCreatedEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             []byte(`{"name": "hodor-org"}`),
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         32,
		Type:             model.OrgAdded,
	}
}

func orgDomainAddedEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             []byte(`{"domain":"hodor.org"}`),
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         33,
		Type:             model.OrgDomainAdded,
	}
}

func orgDomainVerificationAddedEvent(token string) *es_models.Event {
	data, _ := json.Marshal(&org_model.OrgDomain{
		Domain:         "hodor.org",
		ValidationType: org_model.OrgDomainValidationTypeDNS,
		ValidationCode: &crypto.CryptoValue{
			CryptoType: crypto.TypeEncryption,
			Algorithm:  "enc",
			KeyID:      "id",
			Crypted:    []byte(token),
		},
	})
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             data,
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         34,
		Type:             model.OrgDomainVerificationAdded,
	}
}

func orgDomainVerifiedEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             []byte(`{"domain":"hodor.org"}`),
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         35,
		Type:             model.OrgDomainVerified,
	}
}

func orgInactiveEvent() *es_models.Event {
	return &es_models.Event{
		AggregateID:      "hodor-org",
		AggregateType:    model.OrgAggregate,
		AggregateVersion: "v1",
		CreationDate:     time.Now().Add(-1 * time.Minute),
		Data:             nil,
		EditorService:    "testsvc",
		EditorUser:       "testuser",
		ID:               "sdlfö4t23kj",
		ResourceOwner:    "hodor-org",
		Sequence:         52,
		Type:             model.OrgDeactivated,
	}
}

func TestChangesOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es           *OrgEventstore
		id           string
		lastSequence uint64
		limit        uint64
	}
	type res struct {
		changes *org_model.OrgChanges
		org     *model.Org
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:           GetMockChangesOrgOK(ctrl),
				id:           "1",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				changes: &org_model.OrgChanges{Changes: []*org_model.OrgChange{{EventType: "", Sequence: 1, ModifierId: ""}}, LastSequence: 1},
				org:     &model.Org{Name: "MusterOrg"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:           GetMockChangesOrgNoEvents(ctrl),
				id:           "2",
				lastSequence: 0,
				limit:        0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.OrgChanges(nil, tt.args.id, tt.args.lastSequence, tt.args.limit, false)

			org := &model.Org{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, org)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && org.Name != tt.res.org.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
