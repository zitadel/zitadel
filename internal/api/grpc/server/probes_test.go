package server

import (
	"context"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/caos/zitadel/internal/errors"
)

func TestValidator_Healthz(t *testing.T) {
	type fields struct {
		validations map[string]ValidationFunction
	}
	type args struct {
		in0 context.Context
		e   *empty.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			"ok",
			fields{},
			args{
				e: &empty.Empty{},
			},
			&empty.Empty{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				validations: tt.fields.validations,
			}
			got, err := v.Healthz(tt.args.in0, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Healthz() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Healthz() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidator_Ready(t *testing.T) {
	type fields struct {
		validations map[string]ValidationFunction
	}
	type args struct {
		ctx context.Context
		e   *empty.Empty
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{
			"unready error",
			fields{validations: map[string]ValidationFunction{
				"error": func(_ context.Context) error {
					return errors.ThrowInternal(nil, "id", "message")
				},
			}},
			args{
				ctx: context.Background(),
				e:   &empty.Empty{},
			},
			nil,
			true,
		},
		{
			"ready ok",
			fields{validations: map[string]ValidationFunction{
				"ok": func(_ context.Context) error {
					return nil
				},
			}},
			args{
				ctx: context.Background(),
				e:   &empty.Empty{},
			},
			&empty.Empty{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				validations: tt.fields.validations,
			}
			got, err := v.Ready(tt.args.ctx, tt.args.e)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ready() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ready() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		ctx         context.Context
		validations map[string]ValidationFunction
	}
	tests := []struct {
		name string
		args args
		want map[string]error
	}{
		{
			"no error empty",
			args{
				ctx: context.Background(),
				validations: map[string]ValidationFunction{
					"ok": func(_ context.Context) error {
						return nil
					},
				},
			},
			map[string]error{},
		},
		{
			"error in list",
			args{
				ctx: context.Background(),
				validations: map[string]ValidationFunction{
					"ok": func(_ context.Context) error {
						return nil
					},
					"error": func(_ context.Context) error {
						return errors.ThrowInternal(nil, "id", "message")
					},
				},
			},
			map[string]error{
				"error": errors.ThrowInternal(nil, "id", "message"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validate(tt.args.ctx, tt.args.validations); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
