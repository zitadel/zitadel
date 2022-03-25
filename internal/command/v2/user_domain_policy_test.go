package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/org"
)

func Test_customDomainPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyDomainWriteModel
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
						org.NewDomainPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id", "ro").Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyDomainWriteModel{
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
			got, err := orgDomainPolicy(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("customDomainPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customDomainPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultDomainPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyDomainWriteModel
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
						instance.NewDomainPolicyAddedEvent(
							context.Background(),
							&instance.NewAggregate().Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyDomainWriteModel{
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
			got, err := instanceDomainPolicy(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultDomainPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultDomainPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_DomainPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *command.PolicyDomainWriteModel
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
						org.NewDomainPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id", "ro").Aggregate,
							true,
						),
					}, nil
				},
			},
			want: &command.PolicyDomainWriteModel{
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
							instance.NewDomainPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate().Aggregate,
								true,
							),
						}, nil
					}).
					Filter(),
			},
			want: &command.PolicyDomainWriteModel{
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
			got, err := domainPolicyWriteModel(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultDomainPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultDomainPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
