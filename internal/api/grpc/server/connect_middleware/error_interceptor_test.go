package connect_middleware

import (
	"context"
	"errors"
	"testing"

	"connectrpc.com/connect"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
)

func Test_toConnectError(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     connect.AnyRequest
		handler func(t *testing.T) connect.UnaryFunc
	}
	tests := []struct {
		name     string
		args     args
		want     any
		wantCode connect.Code
	}{
		{
			name: "no error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: emptyMockHandler(&connect.Response[struct{}]{}, authz.CtxData{}),
			},
			want: &connect.Response[struct{}]{},
		},
		{
			name: "error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: errorMockHandler(),
			},
			want:     nil,
			wantCode: connect.CodeFailedPrecondition,
		},
		{
			name: "panic with string",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: panicMockHandler("test panic"),
			},
			want:     nil,
			wantCode: connect.CodeInternal,
		},
		{
			name: "panic with error",
			args: args{
				ctx:     context.Background(),
				req:     &mockReq[struct{}]{},
				handler: panicMockHandler(errors.New("oops")),
			},
			want:     nil,
			wantCode: connect.CodeInternal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toConnectError(tt.args.ctx, tt.args.req, tt.args.handler(t))
			if tt.wantCode != 0 {
				var connectErr *connect.Error
				require.ErrorAs(t, err, &connectErr)
				assert.Equal(t, tt.wantCode, connectErr.Code())
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
