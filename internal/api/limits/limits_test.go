package limits_test

import (
	"context"
	"github.com/zitadel/zitadel/internal/api/limits"
	"reflect"
	"testing"
	"time"
)

func TestLoader_Load(t *testing.T) {
	type fields struct {
		querier limits.QuerierFunc
	}
	type args struct {
		ctx        context.Context
		instanceID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   context.Context
		want1  limits.Limits
	}{{
		name: "",
		fields: fields{
			querier: func(ctx context.Context, resourceOwner string) (limits limits.Limits, err error) {
				return mockLimits{}, nil
			},
		},
		args:  args{},
		want:  nil,
		want1: nil,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := limits.NewLoader(tt.fields.querier)
			got, got1 := l.Load(tt.args.ctx, tt.args.instanceID)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Load() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

type mockLimits struct {
	auditLogRetention time.Duration
	block             bool
}

func (m mockLimits) GetAuditLogRetention() *time.Duration {
	return &m.auditLogRetention
}

func (m mockLimits) DoBlock() *bool {
	return &m.block
}
