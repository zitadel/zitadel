package eventstore

import (
	"encoding/json"
	"reflect"
	"testing"
)

var _ Event = (*event)(nil)

type event struct {
	Asdf string
}

type storageEvent struct{}

// Embedded implements StorageEvent.
func (s storageEvent) Embedded() EmbeddedEvent {
	return EmbeddedEvent{}
}

// Unmarshal implements StorageEvent.
func (s storageEvent) Unmarshal(ptr any) error {
	return json.Unmarshal([]byte(`{"asdf": "asdf"}`), ptr)
}

func TestUnmarshal(t *testing.T) {
	type args struct {
		from StorageEvent
	}
	tests := []struct {
		name    string
		args    args
		want    *event
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				from: &storageEvent{},
			},
			want: &event{
				Asdf: "asdf",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal[event](tt.args.from)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unmarshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
