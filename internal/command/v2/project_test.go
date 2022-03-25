package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func TestAddProject(t *testing.T) {
	type args struct {
		a                      *project.Aggregate
		name                   string
		owner                  string
		privateLabelingSetting domain.PrivateLabelingSetting
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid name",
			args: args{
				a:                      agg,
				name:                   "",
				owner:                  "owner",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid private labeling setting",
			args: args{
				a:                      agg,
				name:                   "name",
				owner:                  "owner",
				privateLabelingSetting: -1,
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid owner",
			args: args{
				a:                      agg,
				name:                   "name",
				owner:                  "",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "correct",
			args: args{
				a:                      agg,
				name:                   "ZITADEL",
				owner:                  "CAOS AG",
				privateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					project.NewProjectAddedEvent(ctx, &agg.Aggregate,
						"ZITADEL",
						false,
						false,
						false,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					),
					project.NewProjectMemberAddedEvent(ctx, &agg.Aggregate,
						"CAOS AG",
						domain.RoleProjectOwner),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, AddProject(tt.args.a, tt.args.name, tt.args.owner, false, false, false, tt.args.privateLabelingSetting), nil, tt.want)
		})
	}
}

func TestExistsProject(t *testing.T) {
	type args struct {
		filter        preparation.FilterToQueryReducer
		id            string
		resourceOwner string
	}
	tests := []struct {
		name       string
		args       args
		wantExists bool
		wantErr    bool
	}{
		{
			name: "no events",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{}, nil
				},
				id:            "id",
				resourceOwner: "ro",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "project added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"name",
							false,
							false,
							false,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						),
					}, nil
				},
				id:            "id",
				resourceOwner: "ro",
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: "project removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						project.NewProjectAddedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"name",
							false,
							false,
							false,
							domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
						),
						project.NewProjectRemovedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"name",
						),
					}, nil
				},
				id:            "id",
				resourceOwner: "ro",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "error durring filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, errors.ThrowInternal(nil, "PROJE-Op26p", "Errors.Internal")
				},
				id:            "id",
				resourceOwner: "ro",
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := ExistsProject(context.Background(), tt.args.filter, tt.args.id, tt.args.resourceOwner)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExistsUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotExists != tt.wantExists {
				t.Errorf("ExistsUser() = %v, want %v", gotExists, tt.wantExists)
			}
		})
	}
}
