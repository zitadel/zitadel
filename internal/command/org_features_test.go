package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/features"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func TestCommandSide_SetOrgFeatures(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
		features      *domain.Features
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
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					PasswordComplexityPolicy: false,
					LabelPolicy:              false,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no change, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
						),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					PasswordComplexityPolicy: false,
					LabelPolicy:              false,
					CustomDomain:             false,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "set with default policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"primary",
								"secondary",
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					PasswordComplexityPolicy: false,
					LabelPolicy:              false,
					CustomDomain:             false,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "set with custom policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					//NewOrgFeaturesWriteModel
					expectFilter(),
					//begin ensureOrgSettingsToFeatures
					//begin setAllowedLoginPolicy
					//orgLoginPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
						eventFromEventPusher(
							org.NewLoginPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeNotAllowed,
							),
						),
					),
					//getDefaultLoginPolicy
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
								true,
								true,
								true,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					//begin setDefaultAuthFactorsInCustomLoginPolicy
					//orgLoginPolicyAuthFactorsWriteModel
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicySecondFactorAddedEvent(context.Background(), &iam.NewAggregate().Aggregate, domain.SecondFactorTypeU2F),
						),
						eventFromEventPusher(
							iam.NewLoginPolicyMultiFactorAddedEvent(context.Background(), &iam.NewAggregate().Aggregate, domain.MultiFactorTypeU2FWithPIN),
						),
						eventFromEventPusher(
							org.NewLoginPolicySecondFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeOTP),
						),
					),
					//addSecondFactorToLoginPolicy
					expectFilter(),
					//removeSecondFactorFromLoginPolicy
					expectFilter(),
					//addMultiFactorToLoginPolicy
					expectFilter(),
					//end setDefaultAuthFactorsInCustomLoginPolicy
					//orgPasswordComplexityPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							iam.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
						eventFromEventPusher(
							org.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								7,
								false,
								false,
								false,
								false,
							),
						),
					),
					//orgLabelPolicyWriteModelByID
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"primary",
								"secondary",
								false,
							),
						),
						eventFromEventPusher(
							org.NewLabelPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"custom",
								"secondary",
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewLoginPolicySecondFactorRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeOTP),
							),
							eventFromEventPusher(
								org.NewLoginPolicySecondFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.SecondFactorTypeU2F),
							),
							eventFromEventPusher(
								org.NewLoginPolicyMultiFactorAddedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate, domain.MultiFactorTypeU2FWithPIN),
							),
							eventFromEventPusher(
								newLoginPolicyChangedEvent(context.Background(), "org1", true, true, true, true, domain.PasswordlessTypeAllowed),
							),
							eventFromEventPusher(
								org.NewPasswordComplexityPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								org.NewLabelPolicyRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
							eventFromEventPusher(
								newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
				features: &domain.Features{
					TierName:                 "Test",
					State:                    domain.FeaturesStateActive,
					AuditLogRetention:        time.Hour,
					LoginPolicyFactors:       false,
					LoginPolicyIDP:           false,
					LoginPolicyPasswordless:  false,
					LoginPolicyRegistration:  false,
					LoginPolicyUsernameLogin: false,
					PasswordComplexityPolicy: false,
					LabelPolicy:              false,
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
			got, err := r.SetOrgFeatures(tt.args.ctx, tt.args.resourceOwner, tt.args.features)
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

func TestCommandSide_RemoveOrgFeatures(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		resourceOwner string
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
			name: "resourceowner missing, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "no features set, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove with default policies, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							newFeaturesSetEvent(context.Background(), "org1", "Test", domain.FeaturesStateActive, time.Hour)),
					),
					expectFilter(
						eventFromEventPusher(
							newIAMFeaturesSetEvent(context.Background(), "Default", domain.FeaturesStateActive, time.Hour)),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLoginPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								false,
								false,
								false,
								false,
								domain.PasswordlessTypeAllowed,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								8,
								false,
								false,
								false,
								false,
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							iam.NewLabelPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"primary",
								"secondary",
								false,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewFeaturesRemovedEvent(context.Background(), &org.NewAggregate("org1", "org1").Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				resourceOwner: "org1",
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
			got, err := r.RemoveOrgFeatures(tt.args.ctx, tt.args.resourceOwner)
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

func newIAMFeaturesSetEvent(ctx context.Context, tierName string, state domain.FeaturesState, auditLog time.Duration) *iam.FeaturesSetEvent {
	event, _ := iam.NewFeaturesSetEvent(
		ctx,
		&iam.NewAggregate().Aggregate,
		[]features.FeaturesChanges{
			features.ChangeTierName(tierName),
			features.ChangeState(state),
			features.ChangeAuditLogRetention(auditLog),
		},
	)
	return event
}

func newFeaturesSetEvent(ctx context.Context, orgID string, tierName string, state domain.FeaturesState, auditLog time.Duration) *org.FeaturesSetEvent {
	event, _ := org.NewFeaturesSetEvent(
		ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		[]features.FeaturesChanges{
			features.ChangeTierName(tierName),
			features.ChangeState(state),
			features.ChangeAuditLogRetention(auditLog),
		},
	)
	return event
}
