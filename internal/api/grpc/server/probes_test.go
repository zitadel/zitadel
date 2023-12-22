package server

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestValidator_Healthz(t *testing.T) {
	type fields struct {
		validations map[string]ValidationFunction
	}
	type res struct {
		want   *emptypb.Empty
		hasErr bool
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			"ok",
			fields{},
			res{
				&emptypb.Empty{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				validations: tt.fields.validations,
			}
			got, err := v.Healthz(context.Background(), &emptypb.Empty{})
			if (err != nil) != tt.res.hasErr {
				t.Errorf("Healthz() error = %v, wantErr %v", err, tt.res.hasErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("Healthz() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}

func TestValidator_Ready(t *testing.T) {
	type fields struct {
		validations map[string]ValidationFunction
	}
	type res struct {
		want   *emptypb.Empty
		hasErr bool
	}
	tests := []struct {
		name   string
		fields fields
		res    res
	}{
		{
			"unready error",
			fields{validations: map[string]ValidationFunction{
				"error": func(_ context.Context) error {
					return zerrors.ThrowInternal(nil, "id", "message")
				},
			}},
			res{
				nil,
				true,
			},
		},
		{
			"ready ok",
			fields{validations: map[string]ValidationFunction{
				"ok": func(_ context.Context) error {
					return nil
				},
			}},
			res{
				&emptypb.Empty{},
				false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				validations: tt.fields.validations,
			}
			got, err := v.Ready(context.Background(), &emptypb.Empty{})
			if (err != nil) != tt.res.hasErr {
				t.Errorf("Ready() error = %v, wantErr %v", err, tt.res.hasErr)
				return
			}
			if !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("Ready() got = %v, want %v", got, tt.res.want)
			}
		})
	}
}

func Test_validate(t *testing.T) {
	type args struct {
		validations map[string]ValidationFunction
	}
	type res struct {
		want map[string]any
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"no error empty",
			args{
				validations: map[string]ValidationFunction{
					"ok": func(_ context.Context) error {
						return nil
					},
				},
			},
			res{
				map[string]any{},
			},
		},
		{
			"error in list",
			args{
				validations: map[string]ValidationFunction{
					"ok": func(_ context.Context) error {
						return nil
					},
					"error": func(_ context.Context) error {
						return zerrors.ThrowInternal(nil, "id", "message")
					},
				},
			},
			res{
				map[string]any{
					"error": zerrors.ThrowInternal(nil, "id", "message"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validate(context.Background(), tt.args.validations); !reflect.DeepEqual(got, tt.res.want) {
				t.Errorf("validate() = %v, want %v", got, tt.res.want)
			}
		})
	}
}
