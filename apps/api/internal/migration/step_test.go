package migration

import (
	"reflect"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
)

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
