package project

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func TestAddApp(t *testing.T) {
	type args struct {
		a      *project.Aggregate
		id     string
		name   string
		filter preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid id",
			args: args{
				a:    agg,
				id:   "",
				name: "name",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-0wTYg", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid name",
			args: args{
				a:    agg,
				id:   "id",
				name: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-P7gKR", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "project not exists",
			args: args{
				a:    agg,
				id:   "id",
				name: "name",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				CreateErr: errors.ThrowAlreadyExists(nil, "PROJE-5LQ0U", "Errors.Project.NotFound"),
			},
		},
		{
			name: "correct",
			args: args{
				a:    agg,
				id:   "appID",
				name: "Console",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewProjectAddedEvent(
								ctx,
								&agg.Aggregate,
								"ZITADEL",
								false,
								false,
								false,
								domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
							),
						}, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					project.NewApplicationAddedEvent(ctx, &agg.Aggregate,
						"appID",
						"Console",
					),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, AddApp(tt.args.a, tt.args.id, tt.args.name), tt.args.filter, tt.want)
		})
	}
}

func TestAddOIDCConfig(t *testing.T) {
	type args struct {
		a        *project.Aggregate
		appID    string
		clientID string
		filter   preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid appID",
			args: args{
				a:        agg,
				appID:    "",
				clientID: "clientID",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-NnavI", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid clientID",
			args: args{
				a:        agg,
				appID:    "appID",
				clientID: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-ghTsJ", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "app not exists",
			args: args{
				a:        agg,
				appID:    "id",
				clientID: "name",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				CreateErr: errors.ThrowNotFound(nil, "PROJE-sLDbG", "Errors.Project.Application.NotFound"),
			},
		},
		{
			name: "correct",
			args: args{
				a:        agg,
				appID:    "appID",
				clientID: "clientID",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewApplicationAddedEvent(
								ctx,
								&agg.Aggregate,
								"appID",
								"Console",
							),
						}, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
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
			preparation.AssertValidation(t,
				AddOIDCConfig(*tt.args.a,
					domain.OIDCVersionV1,
					tt.args.appID,
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
		clientID string
		filter   preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := project.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "invalid appID",
			args: args{
				a:        agg,
				appID:    "",
				clientID: "clientID",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-XHsKt", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "invalid clientID",
			args: args{
				a:        agg,
				appID:    "appID",
				clientID: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "PROJE-XXED5", "Errors.Invalid.Argument"),
			},
		},
		{
			name: "app not exists",
			args: args{
				a:        agg,
				appID:    "id",
				clientID: "name",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				CreateErr: errors.ThrowNotFound(nil, "PROJE-EpG1p", "Errors.Project.Application.NotFound"),
			},
		},
		{
			name: "correct",
			args: args{
				a:        agg,
				appID:    "appID",
				clientID: "clientID",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							project.NewApplicationAddedEvent(
								ctx,
								&agg.Aggregate,
								"appID",
								"Console",
							),
						}, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
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
			preparation.AssertValidation(t,
				AddAPIConfig(*tt.args.a,
					tt.args.appID,
					tt.args.clientID,
					nil,
					domain.APIAuthMethodTypeBasic,
				), tt.args.filter, tt.want)
		})
	}
}

func TestExistsApp(t *testing.T) {
	type args struct {
		filter        preparation.FilterToQueryReducer
		appID         string
		projectID     string
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
				appID:         "appID",
				projectID:     "projectID",
				resourceOwner: "ro",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "app added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						project.NewApplicationAddedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"appID",
							"name",
						),
					}, nil
				},
				appID:         "appID",
				projectID:     "projectID",
				resourceOwner: "ro",
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: "app removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						project.NewApplicationAddedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"appID",
							"name",
						),
						project.NewApplicationRemovedEvent(
							context.Background(),
							&project.NewAggregate("id", "ro").Aggregate,
							"appID",
							"name",
						),
					}, nil
				},
				appID:         "appID",
				projectID:     "projectID",
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
				appID:         "appID",
				projectID:     "projectID",
				resourceOwner: "ro",
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := ExistsApp(context.Background(), tt.args.filter, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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
