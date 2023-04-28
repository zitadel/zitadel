package models

import (
	"testing"
	"time"
)

func TestObjectRoot_AppendEvent(t *testing.T) {
	type fields struct {
		ID           string
		Sequence     uint64
		CreationDate time.Time
		ChangeDate   time.Time
	}
	type args struct {
		event     *Event
		isNewRoot bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"new root",
			fields{},
			args{
				&Event{
					AggregateID:  "aggID",
					Sequence:     34555,
					CreationDate: time.Now(),
				},
				true,
			},
		},
		{
			"existing root",
			fields{
				"agg",
				234,
				time.Now().Add(-24 * time.Hour),
				time.Now().Add(-12 * time.Hour),
			},
			args{
				&Event{
					AggregateID:      "agg",
					Sequence:         34555425,
					CreationDate:     time.Now(),
					PreviousSequence: 22,
				},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &ObjectRoot{
				AggregateID:  tt.fields.ID,
				Sequence:     tt.fields.Sequence,
				CreationDate: tt.fields.CreationDate,
				ChangeDate:   tt.fields.ChangeDate,
			}
			o.AppendEvent(tt.args.event)
			if tt.args.isNewRoot {
				if !o.CreationDate.Equal(tt.args.event.CreationDate) {
					t.Error("creationDate should be equal to event on new root")
				}
			} else {
				if o.CreationDate.Equal(o.ChangeDate) {
					t.Error("creationDate and changedate should differ")
				}
			}
			if o.Sequence != tt.args.event.Sequence {
				t.Errorf("sequence not equal to event: event: %d root: %d", tt.args.event.Sequence, o.Sequence)
			}
			if !o.ChangeDate.Equal(tt.args.event.CreationDate) {
				t.Errorf("changedate should be equal to event creation date:  event: %v root: %v", tt.args.event.CreationDate, o.ChangeDate)
			}
		})
	}
}
