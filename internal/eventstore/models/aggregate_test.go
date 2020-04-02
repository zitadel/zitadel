package models

import (
	"testing"
)

func TestAggregate_AppendEvent(t *testing.T) {
	type fields struct {
		aggregate *Aggregate
	}
	type args struct {
		typ     EventType
		payload interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Aggregate
		wantErr bool
	}{
		{
			name:    "no event type error",
			fields:  fields{aggregate: &Aggregate{}},
			args:    args{},
			want:    &Aggregate{},
			wantErr: true,
		},
		{
			name:    "invalid payload error",
			fields:  fields{aggregate: &Aggregate{}},
			args:    args{typ: "user", payload: 134},
			want:    &Aggregate{},
			wantErr: true,
		},
		{
			name:    "event added",
			fields:  fields{aggregate: &Aggregate{Events: []*Event{}}},
			args:    args{typ: "user.deactivated"},
			want:    &Aggregate{Events: []*Event{&Event{Type: "user.deactivated"}}},
			wantErr: false,
		},
		{
			name:    "event added",
			fields:  fields{aggregate: &Aggregate{Events: []*Event{&Event{}}}},
			args:    args{typ: "user.deactivated"},
			want:    &Aggregate{Events: []*Event{&Event{}, &Event{Type: "user.deactivated"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.aggregate.AppendEvent(tt.args.typ, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregate.AppendEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.fields.aggregate.Events) != len(got.Events) {
				t.Errorf("events len should be %d but was %d", len(tt.fields.aggregate.Events), len(got.Events))
			}
		})
	}
}

func TestAggregate_Validate(t *testing.T) {
	type fields struct {
		aggregate *Aggregate
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "aggregate nil error",
			wantErr: true,
		},
		{
			name:    "aggregate empty error",
			wantErr: true,
			fields:  fields{aggregate: &Aggregate{}},
		},
		{
			name:    "no id error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				typ:            "user",
				version:        "v1.0.0",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
				Events: []*Event{&Event{
					AggregateType:    "user",
					AggregateVersion: "v1.0.0",
					EditorOrg:        "org",
					EditorService:    "management",
					EditorUser:       "hodor",
					ResourceOwner:    "org",
					Type:             "born",
				}},
			}},
		},
		{
			name:    "no type error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				id:             "aggID",
				version:        "v1.0.0",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
				Events: []*Event{&Event{
					AggregateID:      "hodor",
					AggregateVersion: "v1.0.0",
					EditorOrg:        "org",
					EditorService:    "management",
					EditorUser:       "hodor",
					ResourceOwner:    "org",
					Type:             "born",
				}},
			}},
		},
		{
			name:    "no events error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				id:             "aggID",
				typ:            "user",
				version:        "v1.0.0",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
			}},
		},
		{
			name:    "invalid event error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				id:             "aggID",
				typ:            "user",
				version:        "v1.0.0",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
				Events:         []*Event{&Event{}},
			}},
		},
		{
			name:    "invalid version error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				id:             "aggID",
				typ:            "user",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
				Events: []*Event{&Event{
					AggregateID:   "hodor",
					AggregateType: "user",
					EditorOrg:     "org",
					EditorService: "management",
					EditorUser:    "hodor",
					ResourceOwner: "org",
					Type:          "born",
				}},
			}},
		},
		{
			name:    "validation ok",
			wantErr: false,
			fields: fields{aggregate: &Aggregate{
				id:             "aggID",
				typ:            "user",
				version:        "v1.0.0",
				editorOrg:      "org",
				editorService:  "svc",
				editorUser:     "hodor",
				resourceOwner:  "org",
				latestSequence: 5,
				Events: []*Event{&Event{
					AggregateID:      "hodor",
					AggregateType:    "user",
					AggregateVersion: "v1.0.0",
					EditorOrg:        "org",
					EditorService:    "management",
					EditorUser:       "hodor",
					ResourceOwner:    "org",
					Type:             "born",
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.aggregate.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Aggregate.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
