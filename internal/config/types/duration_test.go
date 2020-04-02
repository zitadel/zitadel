package types

import (
	"testing"
	"time"
)

func TestDuration_UnmarshalText(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    time.Duration
	}{
		{
			"ok",
			args{
				data: []byte("10s"),
			},
			false,
			time.Duration(10 * time.Second),
		},
		{
			"error",
			args{
				data: []byte("10"),
			},
			true,
			time.Duration(0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Duration{}
			if err := d.UnmarshalText(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}
			if d.Duration != tt.want {
				t.Errorf("UnmarshalText() got = %v, want %v", d.Duration, tt.want)
			}
		})
	}
}
