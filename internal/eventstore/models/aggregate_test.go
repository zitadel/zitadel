package models

import (
	"testing"

	"github.com/caos/zitadel/internal/errors"
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
			name:   "event added",
			fields: fields{aggregate: &Aggregate{Events: []*Event{}}},
			args:   args{typ: "user.deactivated"},
			want: &Aggregate{Events: []*Event{
				{Type: "user.deactivated"},
			}},
			wantErr: false,
		},
		{
			name: "event added",
			fields: fields{aggregate: &Aggregate{Events: []*Event{
				{},
			}}},
			args: args{typ: "user.deactivated"},
			want: &Aggregate{Events: []*Event{
				{},
				{Type: "user.deactivated"},
			}},
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
				typ:              "user",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Events: []*Event{
					{
						AggregateType:    "user",
						AggregateVersion: "v1.0.0",
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
				ID:               "aggID",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Events: []*Event{
					{
						AggregateID:      "hodor",
						AggregateVersion: "v1.0.0",
						EditorService:    "management",
						EditorUser:       "hodor",
						ResourceOwner:    "org",
						Type:             "born",
					}},
			}},
		},
		{
			name:    "invalid version error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				ID:               "aggID",
				typ:              "user",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Events: []*Event{
					{
						AggregateID:   "hodor",
						AggregateType: "user",
						EditorService: "management",
						EditorUser:    "hodor",
						ResourceOwner: "org",
						Type:          "born",
					}},
			}},
		},
		{
			name:    "no query in precondition error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				ID:               "aggID",
				typ:              "user",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Precondition: &precondition{
					Validation: func(...*Event) error { return nil },
				},
				Events: []*Event{
					{
						AggregateID:      "hodor",
						AggregateType:    "user",
						AggregateVersion: "v1.0.0",
						EditorService:    "management",
						EditorUser:       "hodor",
						ResourceOwner:    "org",
						Type:             "born",
					}},
			}},
		},
		{
			name:    "no func in precondition error",
			wantErr: true,
			fields: fields{aggregate: &Aggregate{
				ID:               "aggID",
				typ:              "user",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Precondition: &precondition{
					Query: NewSearchQuery().AggregateIDFilter("hodor"),
				},
				Events: []*Event{
					{
						AggregateID:      "hodor",
						AggregateType:    "user",
						AggregateVersion: "v1.0.0",
						EditorService:    "management",
						EditorUser:       "hodor",
						ResourceOwner:    "org",
						Type:             "born",
					}},
			}},
		},
		{
			name:    "validation without precondition ok",
			wantErr: false,
			fields: fields{aggregate: &Aggregate{
				ID:               "aggID",
				typ:              "user",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Events: []*Event{
					{
						AggregateID:      "hodor",
						AggregateType:    "user",
						AggregateVersion: "v1.0.0",
						EditorService:    "management",
						EditorUser:       "hodor",
						ResourceOwner:    "org",
						Type:             "born",
					}},
			}},
		},
		{
			name:    "validation with precondition ok",
			wantErr: false,
			fields: fields{aggregate: &Aggregate{
				ID:               "aggID",
				typ:              "user",
				version:          "v1.0.0",
				editorService:    "svc",
				editorUser:       "hodor",
				resourceOwner:    "org",
				PreviousSequence: 5,
				Precondition: &precondition{
					Validation: func(...*Event) error { return nil },
					Query:      NewSearchQuery().AggregateIDFilter("hodor"),
				},
				Events: []*Event{
					{
						AggregateID:      "hodor",
						AggregateType:    "user",
						AggregateVersion: "v1.0.0",
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
			err := tt.fields.aggregate.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Aggregate.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && !errors.IsPreconditionFailed(err) {
				t.Errorf("error must extend precondition failed: %v", err)
			}
		})
	}
}

func TestAggregate_SetPrecondition(t *testing.T) {
	type fields struct {
		aggregate *Aggregate
	}
	type args struct {
		query        *SearchQuery
		validateFunc func(...*Event) error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Aggregate
	}{
		{
			name:   "set precondition",
			fields: fields{aggregate: &Aggregate{}},
			args: args{
				query:        &SearchQuery{},
				validateFunc: func(...*Event) error { return nil },
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := tt.fields.aggregate.SetPrecondition(tt.args.query, tt.args.validateFunc)
			if got.Precondition == nil {
				t.Error("precondition must not be nil")
				t.FailNow()
			}
			if got.Precondition.Query == nil {
				t.Error("query of precondition must not be nil")
			}
			if got.Precondition.Validation == nil {
				t.Error("precondition func must not be nil")
			}
		})
	}
}
