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
		opts    []option
	}
	tests := []struct {
		name    string
		creator *AggregateCreator
		args    args
		want    *Aggregate
		wantErr bool
	}{
		{
			name:    "no ctxdata and no options",
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
			creator: &AggregateCreator{serviceName: "admin"},
			wantErr: true,
			want:    nil,
			args: args{
				ctx:     context.Background(),
				typ:     "user",
				version: "v1.0.0",
				opts: []option{
					OverwriteEditorUser("hodor"),
					OverwriteResourceOwner("org"),
				},
			},
		},
		{
			name:    "no type error",
			creator: &AggregateCreator{serviceName: "admin"},
			wantErr: true,
			want:    nil,
			args: args{
				ctx:     context.Background(),
				id:      "hodor",
				version: "v1.0.0",
				opts: []option{
					OverwriteEditorUser("hodor"),
					OverwriteResourceOwner("org"),
				},
			},
		},
		{
			name:    "invalid version error",
			creator: &AggregateCreator{serviceName: "admin"},
			wantErr: true,
			want:    nil,
			args: args{
				ctx: context.Background(),
				id:  "hodor",
				typ: "user",
				opts: []option{
					OverwriteEditorUser("hodor"),
					OverwriteResourceOwner("org"),
				},
			},
		},
		{
			name:    "create ok",
			creator: &AggregateCreator{serviceName: "admin"},
			wantErr: false,
			want: &Aggregate{
				ID:            "hodor",
				Events:        make([]*Event, 0, 2),
				typ:           "user",
				version:       "v1.0.0",
				editorService: "admin",
				editorUser:    "hodor",
				resourceOwner: "org",
			},
			args: args{
				ctx:     context.Background(),
				id:      "hodor",
				typ:     "user",
				version: "v1.0.0",
				opts: []option{
					OverwriteEditorUser("hodor"),
					OverwriteResourceOwner("org"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.creator.NewAggregate(tt.args.ctx, tt.args.id, tt.args.typ, tt.args.version, 0, tt.args.opts...)
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
