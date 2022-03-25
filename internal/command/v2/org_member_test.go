package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestAddMember(t *testing.T) {
	type args struct {
		a      *org.Aggregate
		userID string
		roles  []string
		filter preparation.FilterToQueryReducer
	}

	ctx := context.Background()
	agg := org.NewAggregate("test", "test")

	tests := []struct {
		name string
		args args
		want preparation.Want
	}{
		{
			name: "no user id",
			args: args{
				a:      agg,
				userID: "",
			},
			want: preparation.Want{
				ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument"),
			},
		},
		// {
		// 	name: "TODO: invalid roles",
		// 	args: args{
		// 		a:      agg,
		// 		userID: "",
		// 		roles:  []string{""},
		// 	},
		// 	want: preparation.Want{
		// 		ValidationErr: errors.ThrowInvalidArgument(nil, "ORG-4Mlfs", "Errors.Invalid.Argument"),
		// 	},
		// },
		{
			name: "user not exists",
			args: args{
				a:      agg,
				userID: "userID",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				CreateErr: errors.ThrowNotFound(nil, "ORG-GoXOn", "Errors.User.NotFound"),
			},
		},
		{
			name: "already member",
			args: args{
				a:      agg,
				userID: "userID",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							user.NewMachineAddedEvent(
								ctx,
								&user.NewAggregate("id", "ro").Aggregate,
								"userName",
								"name",
								"description",
								true,
							),
						}, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							org.NewMemberAddedEvent(
								ctx,
								&org.NewAggregate("id", "ro").Aggregate,
								"userID",
							),
						}, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				CreateErr: errors.ThrowAlreadyExists(nil, "ORG-poWwe", "Errors.Org.Member.AlreadyExists"),
			},
		},
		{
			name: "correct",
			args: args{
				a:      agg,
				userID: "userID",
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							user.NewMachineAddedEvent(
								ctx,
								&user.NewAggregate("id", "ro").Aggregate,
								"userName",
								"name",
								"description",
								true,
							),
						}, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Filter(),
			},
			want: preparation.Want{
				Commands: []eventstore.Command{
					org.NewMemberAddedEvent(ctx, &agg.Aggregate, "userID"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preparation.AssertValidation(t, AddOrgMember(tt.args.a, tt.args.userID, tt.args.roles...), tt.args.filter, tt.want)
		})
	}
}

func TestIsMember(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
		orgID  string
		userID string
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
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "member added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID", "ro").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			name: "member removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID", "ro").Aggregate,
							"userID",
						),
						org.NewMemberRemovedEvent(
							context.Background(),
							&org.NewAggregate("orgID", "ro").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			name: "member cascade removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewMemberAddedEvent(
							context.Background(),
							&org.NewAggregate("orgID", "ro").Aggregate,
							"userID",
						),
						org.NewMemberCascadeRemovedEvent(
							context.Background(),
							&org.NewAggregate("orgID", "ro").Aggregate,
							"userID",
						),
					}, nil
				},
				orgID:  "orgID",
				userID: "userID",
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
				orgID:  "orgID",
				userID: "userID",
			},
			wantExists: false,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExists, err := IsOrgMember(context.Background(), tt.args.filter, tt.args.orgID, tt.args.userID)
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
