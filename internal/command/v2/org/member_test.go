package org

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

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
			gotExists, err := IsMember(context.Background(), tt.args.filter, tt.args.orgID, tt.args.userID)
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
