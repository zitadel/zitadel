package command

import (
	"context"
	"reflect"
	"testing"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_customPasswordComplexityPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *PasswordComplexityPolicyWriteModel
		wantErr bool
	}{
		{
			name: "err from filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, zerrors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
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
						org.NewPasswordComplexityPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id").Aggregate,
							8,
							true,
							true,
							true,
							true,
						),
					}, nil
				},
			},
			want: &PasswordComplexityPolicyWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "id",
					ResourceOwner: "id",
					Events:        []eventstore.Event{},
				},
				MinLength:    8,
				HasLowercase: true,
				HasUppercase: true,
				HasNumber:    true,
				HasSymbol:    true,
				State:        domain.PolicyStateActive,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := customPasswordComplexityPolicy(context.Background(), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("customPasswordComplexityPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("customPasswordComplexityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_defaultPasswordComplexityPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *PasswordComplexityPolicyWriteModel
		wantErr bool
	}{
		{
			name: "err from filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, zerrors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
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
						instance.NewPasswordComplexityPolicyAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							8,
							true,
							true,
							true,
							true,
						),
					}, nil
				},
			},
			want: &PasswordComplexityPolicyWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "INSTANCE",
					ResourceOwner: "INSTANCE",
					Events:        []eventstore.Event{},
					InstanceID:    "INSTANCE",
				},
				MinLength:    8,
				HasLowercase: true,
				HasUppercase: true,
				HasNumber:    true,
				HasSymbol:    true,
				State:        domain.PolicyStateActive,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := defaultPasswordComplexityPolicy(authz.WithInstanceID(context.Background(), "INSTANCE"), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultPasswordComplexityPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultPasswordComplexityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_passwordComplexityPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
	}
	tests := []struct {
		name    string
		args    args
		want    *PasswordComplexityPolicyWriteModel
		wantErr bool
	}{
		{
			name: "err from filter custom",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, zerrors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
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
						org.NewPasswordComplexityPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id").Aggregate,
							8,
							true,
							true,
							true,
							true,
						),
					}, nil
				},
			},
			want: &PasswordComplexityPolicyWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "id",
					ResourceOwner: "id",
					Events:        []eventstore.Event{},
				},
				MinLength:    8,
				HasLowercase: true,
				HasUppercase: true,
				HasNumber:    true,
				HasSymbol:    true,
				State:        domain.PolicyStateActive,
			},
			wantErr: false,
		},
		{
			name: "err from filter default",
			args: args{
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, zerrors.ThrowInternal(nil, "USER-6HnsD", "Errors.Internal")
					}).
					Filter(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "default found",
			args: args{
				filter: NewMultiFilter().
					Append(func(ctx context.Context, queryFactory *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return nil, nil
					}).
					Append(func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
						return []eventstore.Event{
							instance.NewPasswordComplexityPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								8,
								true,
								true,
								true,
								true,
							),
						}, nil
					}).
					Filter(),
			},
			want: &PasswordComplexityPolicyWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "INSTANCE",
					ResourceOwner: "INSTANCE",
					Events:        []eventstore.Event{},
					InstanceID:    "INSTANCE",
				},
				MinLength:    8,
				HasLowercase: true,
				HasUppercase: true,
				HasNumber:    true,
				HasSymbol:    true,
				State:        domain.PolicyStateActive,
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
			got, err := passwordComplexityPolicyWriteModel(authz.WithInstanceID(context.Background(), "INSTANCE"), tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultPasswordComplexityPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("defaultPasswordComplexityPolicy() = %v, want %v", got, tt.want)
			}
		})
	}
}
