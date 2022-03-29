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

func TestAddOIDCApp(t *testing.T) {
	type args struct {
		a        *project.Aggregate
		appID    string
		name     string
		clientID string
		filter   preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid appID",
			args: args{
				a:        agg,
				appID:    "",
				name:     "name",
				clientID: "clientID",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid name",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "",
				clientID: "clientID",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-Fef31", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid clientID",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "name",
				clientID: "",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-ghTsJ", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "project not exists",
			args: args{
				a:        agg,
				appID:    "id",
				name:     "name",
				clientID: "clientID",
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowNotFound(nil, "PROJE-5LQ0U", "Errors.Project.NotFound"),
			},
		},
		{
			name: "correct",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "name",
				clientID: "clientID",
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"appID",
						"name",
					),
					project.NewOIDCConfigAddedEvent(ctx, &agg.Aggregate,
						domain.OIDCVersionV1,
						"appID",
						"clientID",
						nil,
						nil,
						nil,
						nil,
						domain.OIDCApplicationTypeWeb,
						domain.OIDCAuthMethodTypeBasic,
						nil,
						false,
						domain.OIDCTokenTypeBearer,
						false,
						false,
						false,
						0,
						nil,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t,
				AddOIDCApp(*tt.args.a,
					domain.OIDCVersionV1,
					tt.args.appID,
					tt.args.name,
					tt.args.clientID,
					nil,
					nil,
					nil,
					nil,
					domain.OIDCApplicationTypeWeb,
					domain.OIDCAuthMethodTypeBasic,
					nil,
					false,
					domain.OIDCTokenTypeBearer,
					false,
					false,
					false,
					0,
					nil,
				), tt.args.filter, tt.want)
		})
	}
}

func TestAddAPIConfig(t *testing.T) {
	type args struct {
		a        *project.Aggregate
		appID    string
		name     string
		clientID string
		filter   preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want Want
	}{
		{
			name: "invalid appID",
			args: args{
				a:        agg,
				appID:    "",
				name:     "name",
				clientID: "clientID",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid name",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "",
				clientID: "clientID",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-F7g21", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid clientID",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "name",
				clientID: "",
			},
			want: Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-XXED5", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "project not exists",
			args: args{
				a:        agg,
				appID:    "id",
				name:     "name",
				clientID: "clientID",
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: Want{
				CreateErr: errors.ThrowNotFound(nil, "PROJE-Sf2gb", "Errors.Project.NotFound"),
			},
		},
		{
			name: "correct",
			args: args{
				a:        agg,
				appID:    "appID",
				name:     "name",
				clientID: "clientID",
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"project",
								false,
								false,
								false,
								domain.PrivateLabelingSettingUnspecified,
							),
						}, nil
					}).
					Filter(),
			},
			want: Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(
						ctx,
						&agg.Aggregate,
						"appID",
						"name",
					),
					project.NewAPIConfigAddedEvent(ctx, &agg.Aggregate,
						"appID",
						"clientID",
						nil,
						domain.APIAuthMethodTypeBasic,
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AssertValidation(t,
				AddAPIApp(*tt.args.a,
					tt.args.appID,
					tt.args.name,
					tt.args.clientID,
					nil,
					domain.APIAuthMethodTypeBasic,
				), tt.args.filter, tt.want)
		})
	}
}
