package migration

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
)

func TestStepStates_Reduce_failedKeepsLastRun(t *testing.T) {
	agg := eventstore.NewAggregate(t.Context(), SystemAggregateID, SystemAggregate, "v1")
	done := &SetupStep{
		BaseEvent: eventstore.BaseEvent{EventType: repeatableDoneType, Agg: agg},
		Name:      "79_backfill_unique_constraint_owners",
		LastRun:   map[string]any{"version": "v1.0.0", "finalized": false},
	}
	failed := &SetupStep{
		BaseEvent: eventstore.BaseEvent{EventType: failedType, Agg: agg},
		Name:      "79_backfill_unique_constraint_owners",
		LastRun:   map[string]any{"version": "v2.0.0", "finalized": true},
	}

	states := &StepStates{}
	states.AppendEvents(done, failed)
	require.NoError(t, states.Reduce())

	step := states.byName("79_backfill_unique_constraint_owners")
	require.NotNil(t, step)
	assert.Equal(t, StepFailed, step.state)
	lastRun, ok := step.LastRun.(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "v1.0.0", lastRun["version"])
	assert.Equal(t, false, lastRun["finalized"])
}

func TestStepStates_lastByState(t *testing.T) {
	now := time.Now()
	past := now.Add(-10 * time.Millisecond)
	tests := []struct {
		name   string
		fields *StepStates
		arg    StepState
		want   *Step
	}{
		{
			name:   "no events reduced invalid state",
			fields: &StepStates{},
			arg:    -1,
		},
		{
			name:   "no events reduced by valid state",
			fields: &StepStates{},
			arg:    StepDone,
		},
		{
			name: "no state found",
			fields: &StepStates{
				Steps: []*Step{
					{
						SetupStep: &SetupStep{
							Name: "done",
						},
						state: StepDone,
					},
					{
						SetupStep: &SetupStep{
							Name: "failed",
						},
						state: StepFailed,
					},
				},
			},
			arg: StepStarted,
		},
		{
			name: "found",
			fields: &StepStates{
				Steps: []*Step{
					{
						SetupStep: &SetupStep{
							BaseEvent: eventstore.BaseEvent{
								Creation: past,
							},
						},
						state: StepStarted,
					},
					{
						SetupStep: &SetupStep{
							BaseEvent: eventstore.BaseEvent{
								Creation: now,
							},
						},
						state: StepStarted,
					},
				},
			},
			arg: StepStarted,
			want: &Step{
				state: StepStarted,
				SetupStep: &SetupStep{
					BaseEvent: eventstore.BaseEvent{
						Creation: now,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StepStates{
				ReadModel: tt.fields.ReadModel,
				Steps:     tt.fields.Steps,
			}
			if gotStep := s.lastByState(tt.arg); !reflect.DeepEqual(gotStep, tt.want) {
				t.Errorf("StepStates.lastByState() = %v, want %v", *gotStep, *tt.want)
			}
		})
	}
}
