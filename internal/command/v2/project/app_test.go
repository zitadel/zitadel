package project

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

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
