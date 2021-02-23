package models

import (
	"reflect"
	"testing"
)

func Test_eventData(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "from bytes",
			args:    args{[]byte(`{"hodor":"asdf"}`)},
			want:    []byte(`{"hodor":"asdf"}`),
			wantErr: false,
		},
		{
			name: "from pointer",
			args: args{&struct {
				Hodor string `json:"hodor"`
			}{Hodor: "asdf"}},
			want:    []byte(`{"hodor":"asdf"}`),
			wantErr: false,
		},
		{
			name: "from struct",
			args: args{struct {
				Hodor string `json:"hodor"`
			}{Hodor: "asdf"}},
			want:    []byte(`{"hodor":"asdf"}`),
			wantErr: false,
		},
		{
			name: "from map",
			args: args{
				map[string]interface{}{"hodor": "asdf"},
			},
			want:    []byte(`{"hodor":"asdf"}`),
			wantErr: false,
		},
		{
			name:    "from nil",
			args:    args{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "invalid data",
			args:    args{876},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := eventData(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("eventData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eventData() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestEvent_Validate(t *testing.T) {
	type fields struct {
		event *Event
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "event nil",
			wantErr: true,
		},
		{
			name:    "event empty",
			fields:  fields{event: &Event{}},
			wantErr: true,
		},
		{
			name: "no aggregate id",
			fields: fields{event: &Event{
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				EditorUser:       "hodor",
				ResourceOwner:    "org",
				Type:             "born",
			}},
			wantErr: true,
		},
		{
			name: "no aggregate type",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				EditorUser:       "hodor",
				ResourceOwner:    "org",
				Type:             "born",
			}},
			wantErr: true,
		},
		{
			name: "no aggregate version",
			fields: fields{event: &Event{
				AggregateID:   "hodor",
				AggregateType: "user",
				EditorService: "management",
				EditorUser:    "hodor",
				ResourceOwner: "org",
				Type:          "born",
			}},
			wantErr: true,
		},
		{
			name: "no editor service",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorUser:       "hodor",
				ResourceOwner:    "org",
				Type:             "born",
			}},
			wantErr: true,
		},
		{
			name: "no editor user",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				ResourceOwner:    "org",
				Type:             "born",
			}},
			wantErr: true,
		},
		{
			name: "no resource owner",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				EditorUser:       "hodor",
				Type:             "born",
			}},
			wantErr: true,
		},
		{
			name: "no type",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				EditorUser:       "hodor",
				ResourceOwner:    "org",
			}},
			wantErr: true,
		},
		{
			name: "all fields set",
			fields: fields{event: &Event{
				AggregateID:      "hodor",
				AggregateType:    "user",
				AggregateVersion: "v1.0.0",
				EditorService:    "management",
				EditorUser:       "hodor",
				ResourceOwner:    "org",
				Type:             "born",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.event.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Event.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
