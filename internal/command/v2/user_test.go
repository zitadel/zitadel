package command

import (
	"context"
	"testing"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

func TestExistsUser(t *testing.T) {
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
			name: "human registered",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						user.NewHumanRegisteredEvent(
							context.Background(),
							&user.NewAggregate("id", "ro").Aggregate,
							"userName",
							"firstName",
							"lastName",
							"nickName",
							"displayName",
							language.German,
							domain.GenderFemale,
							"support@zitadel.ch",
							true,
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
			name: "human added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						user.NewHumanAddedEvent(
							context.Background(),
							&user.NewAggregate("id", "ro").Aggregate,
							"userName",
							"firstName",
							"lastName",
							"nickName",
							"displayName",
							language.German,
							domain.GenderFemale,
							"support@zitadel.ch",
							true,
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
			name: "machine added",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						user.NewMachineAddedEvent(
							context.Background(),
							&user.NewAggregate("id", "ro").Aggregate,
							"userName",
							"name",
							"description",
							true,
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
			name: "user removed",
			args: args{
				filter: func(_ context.Context, _ *eventstore.SearchQueryBuilder) ([]eventstore.Event, error) {
					return []eventstore.Event{
						user.NewMachineAddedEvent(
							context.Background(),
							&user.NewAggregate("removed", "ro").Aggregate,
							"userName",
							"name",
							"description",
							true,
						),
						user.NewUserRemovedEvent(
							context.Background(),
							&user.NewAggregate("removed", "ro").Aggregate,
							"userName",
							nil,
							true,
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
					return nil, errors.ThrowInternal(nil, "USER-Drebn", "Errors.Internal")
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
			gotExists, err := ExistsUser(context.Background(), tt.args.filter, tt.args.id, tt.args.resourceOwner)
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
