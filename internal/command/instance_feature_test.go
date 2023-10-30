package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/feature"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommands_SetBooleanInstanceFeature(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx   context.Context
		f     domain.Feature
		value bool
	}
	type res struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"unknown feature",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				f:     domain.FeatureUnspecified,
				value: true,
			},
			res{
				err: errors.ThrowPreconditionFailed(nil, "FEAT-AS4k1", "Errors.Feature.InvalidValue"),
			},
		},
		{
			"wrong type",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID("instanceID",
							// as there's currently no other [feature.SetEventType] than [feature.Boolean],
							// we need to use a completely other event type to demonstrate the behaviour
							instance.NewInstanceAddedEvent(context.Background(), &instance.NewAggregate("instanceID").Aggregate,
								"instance",
							),
						),
					),
				),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				f:     domain.FeatureLoginDefaultOrg,
				value: true,
			},
			res{
				err: errors.ThrowPreconditionFailed(nil, "FEAT-SDfjk", "Errors.Feature.TypeNotSupported"),
			},
		},
		{
			"first set",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
					expectPush(
						feature.NewSetEvent[feature.Boolean](context.Background(), &feature.NewAggregate("featureID", "instanceID").Aggregate,
							feature.EventTypeFromFeature(domain.FeatureLoginDefaultOrg),
							feature.Boolean{Boolean: true},
						),
					),
				),
				idGenerator: mock.ExpectID(t, "featureID"),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				f:     domain.FeatureLoginDefaultOrg,
				value: true,
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
		{
			"update flag",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID("instanceID",
							feature.NewSetEvent[feature.Boolean](context.Background(), &feature.NewAggregate("featureID", "instanceID").Aggregate,
								feature.EventTypeFromFeature(domain.FeatureLoginDefaultOrg),
								feature.Boolean{Boolean: true},
							),
						),
					),
					expectPush(
						feature.NewSetEvent[feature.Boolean](context.Background(), &feature.NewAggregate("featureID", "instanceID").Aggregate,
							feature.EventTypeFromFeature(domain.FeatureLoginDefaultOrg),
							feature.Boolean{Boolean: false},
						),
					),
				),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				f:     domain.FeatureLoginDefaultOrg,
				value: false,
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
		{
			"no change",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusherWithInstanceID("instanceID",
							feature.NewSetEvent[feature.Boolean](context.Background(), &feature.NewAggregate("featureID", "instanceID").Aggregate,
								feature.EventTypeFromFeature(domain.FeatureLoginDefaultOrg),
								feature.Boolean{Boolean: true},
							),
						),
					),
				),
			},
			args{
				ctx:   authz.WithInstanceID(context.Background(), "instanceID"),
				f:     domain.FeatureLoginDefaultOrg,
				value: true,
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instanceID",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
			}
			got, err := c.SetBooleanInstanceFeature(tt.args.ctx, tt.args.f, tt.args.value)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.details, got)
		})
	}
}
