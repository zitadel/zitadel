package limits_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/limits"
	"github.com/zitadel/zitadel/internal/api/limits/mock_limits"
)

func TestLoader_Load(t *testing.T) {
	type fields struct {
		querier limits.Querier
	}
	type args struct {
		// If detachContext is false, the context returned from the last call to Load is passed with this call to Load.
		// If detachContext is true, a new context.Background() is passed with this call to Load.
		detachContext bool
	}
	minuteLimits := limits.Limits{AuditLogRetention: gu.Ptr(time.Minute), Block: gu.Ptr(true)}
	hourLimits := limits.Limits{AuditLogRetention: gu.Ptr(time.Hour), Block: gu.Ptr(false)}
	tests := []struct {
		name string
		// The load function is called for each entry in the returned args slice.
		input func(controller *gomock.Controller) (fields, []args)
		want  []limits.Limits
	}{{
		name: "limits on the context are reused",
		input: func(controller *gomock.Controller) (fields, []args) {
			querier := mock_limits.NewMockQuerier(controller)
			querier.EXPECT().Limits(gomock.Any(), gomock.Any()).Return(&minuteLimits, nil)
			return fields{querier: querier}, []args{{detachContext: true}, {detachContext: false}}
		},
		want: []limits.Limits{minuteLimits, minuteLimits},
	}, {
		name: "limits are queried for each unrelated context",
		input: func(controller *gomock.Controller) (fields, []args) {
			querier := mock_limits.NewMockQuerier(controller)
			querier.EXPECT().Limits(gomock.Any(), gomock.Any()).Return(&minuteLimits, nil)
			querier.EXPECT().Limits(gomock.Any(), gomock.Any()).Return(&hourLimits, nil)
			return fields{querier: querier}, []args{{detachContext: true}, {detachContext: true}}
		},
		want: []limits.Limits{minuteLimits, hourLimits},
	}, {
		name: "limits are queried once per context even if the querier returns an error",
		input: func(controller *gomock.Controller) (fields, []args) {
			querier := mock_limits.NewMockQuerier(controller)
			querier.EXPECT().Limits(gomock.Any(), gomock.Any()).Return(nil, errors.New("error from querier"))
			return fields{querier: querier}, []args{{detachContext: true}, {detachContext: false}}
		},
		want: []limits.Limits{{}, {}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, calls := tt.input(gomock.NewController(t))
			ll := limits.NewLoader(f.querier)
			var ctx context.Context
			for i, a := range calls {
				if a.detachContext {
					ctx = context.Background()
				}
				var l limits.Limits
				ctx, l = ll.Load(ctx, "instanceID")
				if !reflect.DeepEqual(l, tt.want[i]) {
					t.Errorf("Load() got = %v, want %v", l, tt.want[i])
				}
			}
		})
	}
}
