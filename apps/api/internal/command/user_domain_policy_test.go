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

func Test_customDomainPolicy(t *testing.T) {
	type args struct {
		filter preparation.FilterToQueryReducer
		orgID  string
	}
	tests := []struct {
		name    string
		args    args
		want    *OrgDomainPolicyWriteModel
		wantErr bool
	}{
		{
			name: "err from filter",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, zerrors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
				},
				orgID: "id",
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
				orgID: "id",
			},
			want: &OrgDomainPolicyWriteModel{
				PolicyDomainWriteModel: PolicyDomainWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:   "id",
						ResourceOwner: "id",
					},
					State: domain.PolicyStateUnspecified,
				},
			},
			wantErr: false,
		},
		{
			name: "policy found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						org.NewDomainPolicyAddedEvent(
							context.Background(),
							&org.NewAggregate("id").Aggregate,
							true,
							true,
							true,
						),
					}, nil
				},
				orgID: "id",
			},
			want: &OrgDomainPolicyWriteModel{
				PolicyDomainWriteModel: PolicyDomainWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:   "id",
						ResourceOwner: "id",
						Events:        []eventstore.Event{},
					},
					UserLoginMustBeDomain:                  true,
					ValidateOrgDomains:                     true,
					SMTPSenderAddressMatchesInstanceDomain: true,
					State:                                  domain.PolicyStateActive,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := orgDomainPolicy(context.Background(), tt.args.filter, tt.args.orgID)
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
		want    *InstanceDomainPolicyWriteModel
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
			want: &InstanceDomainPolicyWriteModel{
				PolicyDomainWriteModel: PolicyDomainWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					State: domain.PolicyStateUnspecified,
				},
			},
			wantErr: false,
		},
		{
			name: "policy found",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						instance.NewDomainPolicyAddedEvent(
							context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							true,
							true,
							true,
						),
					}, nil
				},
			},
			want: &InstanceDomainPolicyWriteModel{
				PolicyDomainWriteModel: PolicyDomainWriteModel{
					WriteModel: eventstore.WriteModel{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
						Events:        []eventstore.Event{},
						InstanceID:    "INSTANCE",
					},
					UserLoginMustBeDomain:                  true,
					ValidateOrgDomains:                     true,
					SMTPSenderAddressMatchesInstanceDomain: true,
					State:                                  domain.PolicyStateActive,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := instanceDomainPolicy(authz.WithInstanceID(context.Background(), "INSTANCE"), tt.args.filter)
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
		orgID  string
	}
	tests := []struct {
		name    string
		args    args
		want    *PolicyDomainWriteModel
		wantErr bool
	}{
		{
			name: "err from filter custom",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return nil, zerrors.ThrowInternal(nil, "USER-IgYlN", "Errors.Internal")
				},
				orgID: "id",
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
							&org.NewAggregate("id").Aggregate,
							true,
							true,
							true,
						),
					}, nil
				},
				orgID: "id",
			},
			want: &PolicyDomainWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "id",
					ResourceOwner: "id",
					Events:        []eventstore.Event{},
				},
				UserLoginMustBeDomain:                  true,
				ValidateOrgDomains:                     true,
				SMTPSenderAddressMatchesInstanceDomain: true,
				State:                                  domain.PolicyStateActive,
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
				orgID: "id",
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
							instance.NewDomainPolicyAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
								true,
								true,
							),
						}, nil
					}).
					Filter(),
				orgID: "id",
			},
			want: &PolicyDomainWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "INSTANCE",
					ResourceOwner: "INSTANCE",
					Events:        []eventstore.Event{},
					InstanceID:    "INSTANCE",
				},
				UserLoginMustBeDomain:                  true,
				ValidateOrgDomains:                     true,
				SMTPSenderAddressMatchesInstanceDomain: true,
				State:                                  domain.PolicyStateActive,
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
			got, err := domainPolicyWriteModel(authz.WithInstanceID(context.Background(), "INSTANCE"), tt.args.filter, tt.args.orgID)
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
