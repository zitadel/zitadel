package models

import (
	"context"
	"reflect"
	"testing"
)

func TestAggregateCreator_NewAggregate(t *testing.T) {
	type args struct {
		ctx     context.Context
		id      string
		typ     AggregateType
		version Version
	}
	tests := []struct {
		name    string
		creator *AggregateCreator
		args    args
		want    *Aggregate
		wantErr bool
	}{
		{
			name:    "invalid CtxData error",
			creator: &AggregateCreator{serviceName: "admin"},
			wantErr: true,
			want:    nil,
			args: args{
				ctx:     context.Background(),
				id:      "hodor",
				typ:     "user",
				version: "v1.0.0",
			},
		},
		{
			name:    "no id error",
			creator: &AggregateCreator{serviceName: "admin", ignoreCtxData: true},
			wantErr: true,
			want:    nil,
			args: args{
				ctx:     context.Background(),
				typ:     "user",
				version: "v1.0.0",
			},
		},
		{
			name:    "no type error",
			creator: &AggregateCreator{serviceName: "admin", ignoreCtxData: true},
			wantErr: true,
			want:    nil,
			args: args{
				ctx:     context.Background(),
				id:      "hodor",
				version: "v1.0.0",
			},
		},
		{
			name:    "invalid version error",
			creator: &AggregateCreator{serviceName: "admin", ignoreCtxData: true},
			wantErr: true,
			want:    nil,
			args: args{
				ctx: context.Background(),
				id:  "hodor",
				typ: "user",
			},
		},
		{
			name:    "create ok",
			creator: &AggregateCreator{serviceName: "admin", ignoreCtxData: true},
			wantErr: false,
			want: &Aggregate{
				ID:            "hodor",
				Events:        make([]*Event, 0, 2),
				Type:          "user",
				Version:       "v1.0.0",
				editorService: "admin",
			},
			args: args{
				ctx:     context.Background(),
				id:      "hodor",
				typ:     "user",
				version: "v1.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.creator.NewAggregate(tt.args.ctx, tt.args.id, tt.args.typ, tt.args.version, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("AggregateCreator.NewAggregate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AggregateCreator.NewAggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}
