package command

import (
	"context"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/features"
	"github.com/caos/zitadel/internal/repository/instance"
)

func TestSetDefaultFeatures(t *testing.T) {
	type args struct {
		a                        *instance.Aggregate
		tierName                 string
		tierDescription          string
		state                    domain.FeaturesState
		stateDescription         string
		retention                time.Duration
		loginPolicyFactors       bool
		loginPolicyIDP           bool
		loginPolicyPasswordless  bool
		loginPolicyRegistration  bool
		loginPolicyUsernameLogin bool
		loginPolicyPasswordReset bool
		passwordComplexityPolicy bool
		labelPolicyPrivateLabel  bool
		labelPolicyWatermark     bool
		customDomain             bool
		privacyPolicy            bool
		metadataUser             bool
		customTextMessage        bool
		customTextLogin          bool
		lockoutPolicy            bool
		actionsAllowed           domain.ActionsAllowed
		maxActions               int
	}
	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid state",
			args: args{
				a:                        instance.NewAggregate("INSTANCE"),
				tierName:                 "",
				tierDescription:          "",
				state:                    0,
				stateDescription:         "",
				retention:                0,
				loginPolicyFactors:       false,
				loginPolicyIDP:           false,
				loginPolicyPasswordless:  false,
				loginPolicyRegistration:  false,
				loginPolicyUsernameLogin: false,
				loginPolicyPasswordReset: false,
				passwordComplexityPolicy: false,
				labelPolicyPrivateLabel:  false,
				labelPolicyWatermark:     false,
				customDomain:             false,
				privacyPolicy:            false,
				metadataUser:             false,
				customTextMessage:        false,
				customTextLogin:          false,
				lockoutPolicy:            false,
				actionsAllowed:           0,
				maxActions:               0,
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "INSTA-d3r1s", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:                        instance.NewAggregate("INSTANCE"),
				tierName:                 "",
				tierDescription:          "",
				state:                    domain.FeaturesStateActive,
				stateDescription:         "",
				retention:                0,
				loginPolicyFactors:       false,
				loginPolicyIDP:           false,
				loginPolicyPasswordless:  false,
				loginPolicyRegistration:  false,
				loginPolicyUsernameLogin: false,
				loginPolicyPasswordReset: false,
				passwordComplexityPolicy: false,
				labelPolicyPrivateLabel:  false,
				labelPolicyWatermark:     false,
				customDomain:             false,
				privacyPolicy:            false,
				metadataUser:             false,
				customTextMessage:        false,
				customTextLogin:          false,
				lockoutPolicy:            false,
				actionsAllowed:           0,
				maxActions:               0,
			},
			want: Want{
				Commands: []eventstore.Command{
					func() *instance.FeaturesSetEvent {
						event, _ := instance.NewFeaturesSetEvent(context.Background(), &instance.NewAggregate("INSTANCE").Aggregate,
							[]features.FeaturesChanges{
								features.ChangeState(domain.FeaturesStateActive),
							},
						)
						return event
					}(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t, SetDefaultFeatures(
				tt.args.a,
				tt.args.tierName,
				tt.args.tierDescription,
				tt.args.state,
				tt.args.stateDescription,
				tt.args.retention,
				tt.args.loginPolicyFactors,
				tt.args.loginPolicyIDP,
				tt.args.loginPolicyPasswordless,
				tt.args.loginPolicyRegistration,
				tt.args.loginPolicyUsernameLogin,
				tt.args.loginPolicyPasswordReset,
				tt.args.passwordComplexityPolicy,
				tt.args.labelPolicyPrivateLabel,
				tt.args.labelPolicyWatermark,
				tt.args.customDomain,
				tt.args.privacyPolicy,
				tt.args.metadataUser,
				tt.args.customTextMessage,
				tt.args.customTextLogin,
				tt.args.lockoutPolicy,
				tt.args.actionsAllowed,
				tt.args.maxActions,
			), NewMultiFilter().
				Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, nil
				}).
				Filter(),
				tt.want)
		})
	}
}
