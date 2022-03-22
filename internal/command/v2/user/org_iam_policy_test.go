package user

import (
	"context"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

func Test_customOrgIAMPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyOrgIAMWriteModel
		wantErr bool
	}{
		{
			name: "err from filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, errors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no events",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{}, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "policy found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewOrgIAMPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id", "ro").Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyOrgIAMWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "id",
					ResourceOwner: "ro",
					Events:        []eventstore.Event{},
				},
				UserLoginMustBeDomain: true,
				State:                 domain.PolicyStateActive,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customOrgIAMPolicy(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("customOrgIAMPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customOrgIAMPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultOrgIAMPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyOrgIAMWriteModel
		wantErr bool
	}{
		{
			name: "err from filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, errors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "no events",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{}, nil
				},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "policy found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						iam.NewOrgIAMPolicyAddedEvent(
							context.Background(),
							&iam.NewAggregate().Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyOrgIAMWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "IAM",
					ResourceOwner: "IAM",
					Events:        []eventstore.Event{},
				},
				UserLoginMustBeDomain: true,
				State:                 domain.PolicyStateActive,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultOrgIAMPolicy(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultOrgIAMPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultOrgIAMPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_OrgIAMPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyOrgIAMWriteModel
		wantErr bool
	}{
		{
			name: "err from filter custom",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, errors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "custom found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewOrgIAMPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id", "ro").Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyOrgIAMWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "id",
					ResourceOwner: "ro",
					Events:        []eventstore.Event{},
				},
				UserLoginMustBeDomain: true,
				State:                 domain.PolicyStateActive,
			},
			wantErr: false,
		},
		{
			name: "err from filter default",
			args: args{
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, errors.ThrowInternal(nil, "USER-6HnsD", "Errors.Internal")
					}).
					Filter(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "default found",
			args: args{
				filter: preparation.NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Append(func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							iam.NewOrgIAMPolicyAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
							),
						}, nil
					}).
					Filter(),
			},
			want: &command.PolicyOrgIAMWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "IAM",
					ResourceOwner: "IAM",
					Events:        []eventstore.Event{},
				},
				UserLoginMustBeDomain: true,
				State:                 domain.PolicyStateActive,
			},
			wantErr: false,
		},
		{
			name: "no policy found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, nil
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orgIAMPolicyWriteModel(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultOrgIAMPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultOrgIAMPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
