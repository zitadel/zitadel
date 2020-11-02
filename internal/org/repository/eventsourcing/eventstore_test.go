package eventsourcing

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	iam_model "github.com/caos/zitadel/internal/iam/model"

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
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.LastSequence != tt.res.changes.LastSequence && org.Name != tt.res.org.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
		})
	}
}

func TestAddIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add idp, ok",
			args: args{
				es:  GetMockChangesOrgWithCrypto(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					Type:        iam_model.IDPConfigTypeOIDC,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID:           "ClientID",
						ClientSecretString: "ClientSecret",
						Issuer:             "Issuer",
						Scopes:             []string{"scope"},
					},
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					Name: "Name",
					Type: iam_model.IDPConfigTypeOIDC,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
						Issuer:   "Issuer",
						Scopes:   []string{"scope"},
					},
				},
			},
		},
		{
			name: "invalid idp config",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID:           "ClientID",
						ClientSecretString: "ClientSecret",
						Issuer:             "Issuer",
						Scopes:             []string{"scope"},
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddIDPConfig(tt.args.ctx, tt.args.idp)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.IDPConfigID == "" {
				t.Errorf("result has no id")
			}
			if result.OIDCConfig.IDPConfigID == "" {
				t.Errorf("result has no id")
			}
			if result.OIDCConfig == nil && result.OIDCConfig.ClientSecret == nil {
				t.Errorf("result has no client secret")
			}
			if result.Name != tt.res.result.Name {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
			if result.OIDCConfig.ClientID != tt.res.result.OIDCConfig.ClientID {
				t.Errorf("got wrong result key: expected: %v, actual: %v ", tt.res.result.OIDCConfig.ClientID, result.OIDCConfig.ClientID)
			}
		})
	}
}

func TestChangeIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change idp, ok",
			args: args{
				es:  GetMockChangesOrgWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "NameChanged",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "NameChanged",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "invalid idp",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeIDPConfig(tt.args.ctx, tt.args.idp)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IdpConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.Name != tt.res.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.result.Name, result.Name)
			}
		})
	}
}

func TestRemoveIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove idp, ok",
			args: args{
				es:  GetMockChangesOrgWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
		},
		{
			name: "no IDPConfigID",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing idp not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveIDPConfig(tt.args.ctx, tt.args.idp)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err")
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
func TestDeactivateIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "deactivate, ok",
			args: args{
				es:  GetMockChangesOrgWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					State:       iam_model.IDPConfigStateInactive,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "no idp id",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.DeactivateIDPConfig(tt.args.ctx, tt.args.idp.AggregateID, tt.args.idp.IDPConfigID)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
		})
	}
}

func TestReactivateIdpConfiguration(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es  *OrgEventstore
		ctx context.Context
		idp *iam_model.IDPConfig
	}
	type res struct {
		result  *iam_model.IDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "reactivate, ok",
			args: args{
				es:  GetMockChangesOrgWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
				},
			},
			res: res{
				result: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					State:       iam_model.IDPConfigStateActive,
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
		},
		{
			name: "no idp id",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1}},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				idp: &iam_model.IDPConfig{ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 1},
					IDPConfigID: "IDPConfigID",
					Name:        "Name",
					OIDCConfig: &iam_model.OIDCIDPConfig{
						ClientID: "ClientID",
					},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ReactivateIDPConfig(tt.args.ctx, tt.args.idp.AggregateID, tt.args.idp.IDPConfigID)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.State != tt.res.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.res.result.State, result.State)
			}
		})
	}
}

func TestChangeOIDCIDPConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		config *iam_model.OIDCIDPConfig
	}
	type res struct {
		result  *iam_model.OIDCIDPConfig
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change oidc config, ok",
			args: args{
				es:  GetMockChangesOrgWithOIDCIdp(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientIDChange",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				result: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientIDChange",
				},
			},
		},
		{
			name: "invalid config",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "idp not existing",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				config: &iam_model.OIDCIDPConfig{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IDPConfigID: "IDPConfigID",
					ClientID:    "ClientID",
					Issuer:      "Issuer",
					Scopes:      []string{"scope"},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeIDPOIDCConfig(tt.args.ctx, tt.args.config)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AggregateID == "" {
				t.Errorf("result has no id")
			}
			if result.IDPConfigID != tt.res.result.IDPConfigID {
				t.Errorf("got wrong result AppID: expected: %v, actual: %v ", tt.res.result.IDPConfigID, result.IDPConfigID)
			}
			if result.ClientID != tt.res.result.ClientID {
				t.Errorf("got wrong result responsetype: expected: %v, actual: %v ", tt.res.result.ClientID, result.ClientID)
			}
		})
	}
}

func TestAddLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.LoginPolicy
	}
	type res struct {
		result  *iam_model.LoginPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add login policy, ok",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister: true,
				},
			},
			res: res{
				result: &iam_model.LoginPolicy{
					ObjectRoot:    es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister: true,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddLoginPolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && result.AllowRegister != tt.res.result.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.res.result.AllowRegister, result.AllowRegister)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangeLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.LoginPolicy
	}
	type res struct {
		result  *iam_model.LoginPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add login policy, ok",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister:         true,
					AllowExternalIdp:      false,
					AllowUsernamePassword: false,
				},
			},
			res: res{
				result: &iam_model.LoginPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister:         true,
					AllowExternalIdp:      false,
					AllowUsernamePassword: false,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeLoginPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.AllowRegister != tt.res.result.AllowRegister {
				t.Errorf("got wrong result AllowRegister: expected: %v, actual: %v ", tt.res.result.AllowRegister, result.AllowRegister)
			}
			if result.AllowUsernamePassword != tt.res.result.AllowUsernamePassword {
				t.Errorf("got wrong result AllowUsernamePassword: expected: %v, actual: %v ", tt.res.result.AllowUsernamePassword, result.AllowUsernamePassword)
			}
			if result.AllowExternalIdp != tt.res.result.AllowExternalIdp {
				t.Errorf("got wrong result AllowExternalIDP: expected: %v, actual: %v ", tt.res.result.AllowExternalIdp, result.AllowExternalIdp)
			}
		})
	}
}

func TestRemoveLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.LoginPolicy
	}
	type res struct {
		result  *iam_model.LoginPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove login policy, ok",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot:            es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					AllowRegister:         true,
					AllowExternalIdp:      false,
					AllowUsernamePassword: false,
				},
			},
			res: res{},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LoginPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveLoginPolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong result should not get err: %v", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddIdpProviderToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *OrgEventstore
		ctx      context.Context
		provider *iam_model.IDPProvider
	}
	type res struct {
		result  *iam_model.IDPProvider
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add idp to login policy, ok",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IdpConfigID2",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				result: &iam_model.IDPProvider{IdpConfigID: "IdpConfigID2"},
			},
		},
		{
			name: "add idp to login policy, already existing",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IDPConfigID",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IdpConfigID2",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddIDPProviderToLoginPolicy(tt.args.ctx, tt.args.provider)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.IdpConfigID != tt.res.result.IdpConfigID {
				t.Errorf("got wrong result IDPConfigID: expected: %v, actual: %v ", tt.res.result.IdpConfigID, result.IdpConfigID)
			}
			if result.Type != tt.res.result.Type {
				t.Errorf("got wrong result Type: expected: %v, actual: %v ", tt.res.result.Type, result.Type)
			}
		})
	}
}

func TestRemoveIdpProviderFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es       *OrgEventstore
		ctx      context.Context
		provider *iam_model.IDPProvider
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove idp to login policy, ok",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IDPConfigID",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{},
		},
		{
			name: "remove idp to login policy, not existing",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IdpConfigID2",
					Type:        iam_model.IDPProviderTypeSystem,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				provider: &iam_model.IDPProvider{
					ObjectRoot:  es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					IdpConfigID: "IdpConfigID2",
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveIDPProviderFromLoginPolicy(tt.args.ctx, tt.args.provider)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddSecondFactorToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *OrgEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.SecondFactorType
	}
	type res struct {
		result  iam_model.SecondFactorType
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add second factor to login policy, ok",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				result: iam_model.SecondFactorTypeOTP,
			},
		},
		{
			name: "add second factor to login policy, already existing",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicyWithMFA(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid mfa",
			args: args{
				es:          GetMockChangesOrgOK(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockChangesOrgNoEvents(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddSecondFactorToLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result != tt.res.result {
				t.Errorf("got wrong result : expected: %v, actual: %v ", tt.res.result, result)
			}
		})
	}
}

func TestRemoveSecondFactorFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *OrgEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.SecondFactorType
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove second factor from login policy, ok",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicyWithMFA(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{},
		},
		{
			name: "remove second factor from  login policy, not existing",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),

				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:          GetMockChangesOrgOK(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.SecondFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockChangesOrgNoEvents(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.SecondFactorTypeOTP,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveSecondFactorFromLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddMultiFactorToLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *OrgEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.MultiFactorType
	}
	type res struct {
		result  iam_model.MultiFactorType
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add mfa to login policy, ok",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				result: iam_model.MultiFactorTypeU2FWithPIN,
			},
		},
		{
			name: "add mfa to login policy, already existing",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicyWithMFA(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "invalid mfa",
			args: args{
				es:          GetMockChangesOrgOK(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockChangesOrgNoEvents(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddMultiFactorToLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result != tt.res.result {
				t.Errorf("got wrong result : expected: %v, actual: %v ", tt.res.result, result)
			}
		})
	}
}

func TestRemoveMultiFactorFromLoginPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es          *OrgEventstore
		ctx         context.Context
		aggregateID string
		mfa         iam_model.MultiFactorType
	}
	type res struct {
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove mfa from login policy, ok",
			args: args{
				es:          GetMockChangesOrgWithLoginPolicyWithMFA(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{},
		},
		{
			name: "remove mfa from  login policy, not existing",
			args: args{
				es:  GetMockChangesOrgWithLoginPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),

				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "invalid provider",
			args: args{
				es:          GetMockChangesOrgOK(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "AggregateID",
				mfa:         iam_model.MultiFactorTypeUnspecified,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:          GetMockChangesOrgNoEvents(ctrl),
				ctx:         authz.NewMockContext("orgID", "userID"),
				aggregateID: "Test",
				mfa:         iam_model.MultiFactorTypeU2FWithPIN,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemoveMultiFactorFromLoginPolicy(tt.args.ctx, tt.args.aggregateID, tt.args.mfa)

			if !tt.res.wantErr && err != nil {
				t.Errorf("should not get err: %v ", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddLabelPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.LabelPolicy
	}
	type res struct {
		result  *iam_model.LabelPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add label policy, ok",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor: "000000",
				},
			},
			res: res{
				result: &iam_model.LabelPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor: "000000",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddLabelPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.PrimaryColor != tt.res.result.PrimaryColor {
				t.Errorf("got wrong result PrimaryColor: expected: %v, actual: %v ", tt.res.result.PrimaryColor, result.PrimaryColor)
			}
		})
	}
}

func TestChangeLabelPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.LabelPolicy
	}
	type res struct {
		result  *iam_model.LabelPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add label policy, ok",
			args: args{
				es:  GetMockChangesOrgWithLabelPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor:   "000000",
					SecondaryColor: "FFFFFF",
				},
			},
			res: res{
				result: &iam_model.LabelPolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					PrimaryColor:   "000000",
					SecondaryColor: "FFFFFF",
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.LabelPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangeLabelPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.PrimaryColor != tt.res.result.PrimaryColor {
				t.Errorf("got wrong result PrimaryColor: expected: %v, actual: %v ", tt.res.result.PrimaryColor, result.PrimaryColor)
			}
			if result.SecondaryColor != tt.res.result.SecondaryColor {
				t.Errorf("got wrong result SecondaryColor: expected: %v, actual: %v ", tt.res.result.SecondaryColor, result.SecondaryColor)
			}
		})
	}
}
func TestAddPasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicy
	}
	type res struct {
		result  *iam_model.PasswordComplexityPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password complexity policy, ok",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:    10,
					HasLowercase: true,
					HasUppercase: true,
					HasSymbol:    true,
					HasNumber:    true,
				},
			},
			res: res{
				result: &iam_model.PasswordComplexityPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:    10,
					HasLowercase: true,
					HasUppercase: true,
					HasSymbol:    true,
					HasNumber:    true,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MinLength != tt.res.result.MinLength {
				t.Errorf("got wrong result MinLength: expected: %v, actual: %v ", tt.res.result.MinLength, result.MinLength)
			}
			if result.HasLowercase != tt.res.result.HasLowercase {
				t.Errorf("got wrong result HasLowercase: expected: %v, actual: %v ", tt.res.result.HasLowercase, result.HasLowercase)
			}
			if result.HasUppercase != tt.res.result.HasUppercase {
				t.Errorf("got wrong result HasUppercase: expected: %v, actual: %v ", tt.res.result.HasUppercase, result.HasUppercase)
			}
			if result.HasNumber != tt.res.result.HasNumber {
				t.Errorf("got wrong result HasNumber: expected: %v, actual: %v ", tt.res.result.HasNumber, result.HasNumber)
			}
			if result.HasSymbol != tt.res.result.HasSymbol {
				t.Errorf("got wrong result HasSymbol: expected: %v, actual: %v ", tt.res.result.HasSymbol, result.HasSymbol)
			}
		})
	}
}

func TestChangePasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicy
	}
	type res struct {
		result  *iam_model.PasswordComplexityPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password complexity policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordComplexityPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:    5,
					HasLowercase: false,
					HasUppercase: false,
					HasSymbol:    false,
					HasNumber:    false,
				},
			},
			res: res{
				result: &iam_model.PasswordComplexityPolicy{
					ObjectRoot:   es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:    5,
					HasLowercase: false,
					HasUppercase: false,
					HasSymbol:    false,
					HasNumber:    false,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MinLength:  10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordComplexityPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MinLength != tt.res.result.MinLength {
				t.Errorf("got wrong result MinLength: expected: %v, actual: %v ", tt.res.result.MinLength, result.MinLength)
			}
			if result.HasLowercase != tt.res.result.HasLowercase {
				t.Errorf("got wrong result HasLowercase: expected: %v, actual: %v ", tt.res.result.HasLowercase, result.HasLowercase)
			}
			if result.HasUppercase != tt.res.result.HasUppercase {
				t.Errorf("got wrong result HasUppercase: expected: %v, actual: %v ", tt.res.result.HasUppercase, result.HasUppercase)
			}
			if result.HasNumber != tt.res.result.HasNumber {
				t.Errorf("got wrong result HasNumber: expected: %v, actual: %v ", tt.res.result.HasNumber, result.HasNumber)
			}
			if result.HasSymbol != tt.res.result.HasSymbol {
				t.Errorf("got wrong result HasSymbol: expected: %v, actual: %v ", tt.res.result.HasSymbol, result.HasSymbol)
			}
		})
	}
}

func TestRemovePasswordComplexityPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordComplexityPolicy
	}
	type res struct {
		result  *iam_model.PasswordComplexityPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove password compelxity policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordComplexityPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordComplexityPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemovePasswordComplexityPolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong result should not get err: %v", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddPasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordAgePolicy
	}
	type res struct {
		result  *iam_model.PasswordAgePolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password age policy, ok",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
			res: res{
				result: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordAgePolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAgeDays != tt.res.result.MaxAgeDays {
				t.Errorf("got wrong result MaxAgeDays: expected: %v, actual: %v ", tt.res.result.MaxAgeDays, result.MaxAgeDays)
			}
			if result.ExpireWarnDays != tt.res.result.ExpireWarnDays {
				t.Errorf("got wrong result.ExpireWarnDays: expected: %v, actual: %v ", tt.res.result.ExpireWarnDays, result.ExpireWarnDays)
			}
		})
	}
}

func TestChangePasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordAgePolicy
	}
	type res struct {
		result  *iam_model.PasswordAgePolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password age policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordAgePolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     5,
					ExpireWarnDays: 5,
				},
			},
			res: res{
				result: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     5,
					ExpireWarnDays: 5,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot:     es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAgeDays:     10,
					ExpireWarnDays: 10,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordAgePolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAgeDays != tt.res.result.MaxAgeDays {
				t.Errorf("got wrong result MaxAgeDays: expected: %v, actual: %v ", tt.res.result.MaxAgeDays, result.MaxAgeDays)
			}
			if result.ExpireWarnDays != tt.res.result.ExpireWarnDays {
				t.Errorf("got wrong result.ExpireWarnDays: expected: %v, actual: %v ", tt.res.result.ExpireWarnDays, result.ExpireWarnDays)
			}
		})
	}
}

func TestRemovePasswordAgePolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordAgePolicy
	}
	type res struct {
		result  *iam_model.PasswordAgePolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove password age policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordAgePolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordAgePolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemovePasswordAgePolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong result should not get err: %v", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestAddPasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordLockoutPolicy
	}
	type res struct {
		result  *iam_model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "add password lockout policy, ok",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				result: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing org not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.AddPasswordLockoutPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAttempts != tt.res.result.MaxAttempts {
				t.Errorf("got wrong result MaxAttempts: expected: %v, actual: %v ", tt.res.result.MaxAttempts, result.MaxAttempts)
			}
			if result.ShowLockOutFailures != tt.res.result.ShowLockOutFailures {
				t.Errorf("got wrong result.ShowLockOutFailures: expected: %v, actual: %v ", tt.res.result.ShowLockOutFailures, result.ShowLockOutFailures)
			}
		})
	}
}

func TestChangePasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordLockoutPolicy
	}
	type res struct {
		result  *iam_model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "change password lockout policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordLockoutPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         5,
					ShowLockOutFailures: false,
				},
			},
			res: res{
				result: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         5,
					ShowLockOutFailures: false,
				},
			},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot:          es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.ChangePasswordLockoutPolicy(tt.args.ctx, tt.args.policy)
			if (tt.res.wantErr && !tt.res.errFunc(err)) || (err != nil && !tt.res.wantErr) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			if tt.res.wantErr && tt.res.errFunc(err) {
				return
			}
			if result.MaxAttempts != tt.res.result.MaxAttempts {
				t.Errorf("got wrong result MaxAttempts: expected: %v, actual: %v ", tt.res.result.MaxAttempts, result.MaxAttempts)
			}
			if result.ShowLockOutFailures != tt.res.result.ShowLockOutFailures {
				t.Errorf("got wrong result.ShowLockOutFailures: expected: %v, actual: %v ", tt.res.result.ShowLockOutFailures, result.ShowLockOutFailures)
			}
		})
	}
}

func TestRemovePasswordLockoutPolicy(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es     *OrgEventstore
		ctx    context.Context
		policy *iam_model.PasswordLockoutPolicy
	}
	type res struct {
		result  *iam_model.PasswordLockoutPolicy
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "remove password lockout policy, ok",
			args: args{
				es:  GetMockChangesOrgWithPasswordLockoutPolicy(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{},
		},
		{
			name: "invalid policy",
			args: args{
				es:  GetMockChangesOrgOK(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot: es_models.ObjectRoot{Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "existing iam not found",
			args: args{
				es:  GetMockChangesOrgNoEvents(ctrl),
				ctx: authz.NewMockContext("orgID", "userID"),
				policy: &iam_model.PasswordLockoutPolicy{
					ObjectRoot: es_models.ObjectRoot{AggregateID: "AggregateID", Sequence: 0},
				},
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.es.RemovePasswordLockoutPolicy(tt.args.ctx, tt.args.policy)

			if !tt.res.wantErr && err != nil {
				t.Errorf("got wrong result should not get err: %v", err)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
