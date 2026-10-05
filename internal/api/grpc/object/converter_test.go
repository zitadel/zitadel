package object

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/zerrors"
	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func Test_ListQueryToModel(t *testing.T) {
	defaults := systemdefaults.SystemDefaults{
		DefaultQueryLimit: 100,
		MaxQueryLimit:     1000,
	}
	type args struct {
		req *object_pb.ListQuery
	}
	type res struct {
		offset, limit uint64
		asc           bool
		err           error
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "all fields filled",
			args: args{
				req: &object_pb.ListQuery{
					Offset: 100,
					Limit:  100,
					Asc:    true,
				},
			},
			res: res{
				offset: 100,
				limit:  100,
				asc:    true,
			},
		},
		{
			name: "all fields empty, max limit as default",
			args: args{
				req: &object_pb.ListQuery{
					Offset: 0,
					Limit:  0,
					Asc:    false,
				},
			},
			res: res{
				offset: 0,
				limit:  1000,
				asc:    false,
			},
		},
		{
			name: "nil query, max limit as default",
			args: args{
				req: nil,
			},
			res: res{
				offset: 0,
				limit:  1000,
				asc:    false,
			},
		},
		{
			name: "limit exceeds max",
			args: args{
				req: &object_pb.ListQuery{
					Offset: 10,
					Limit:  1001,
					Asc:    true,
				},
			},
			res: res{
				err: zerrors.ThrowInvalidArgument(nil, "QUERY-4M0fs", "Errors.Query.LimitExceeded"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit, asc, err := ListQueryToModel(defaults, tt.args.req)
			assert.ErrorIs(t, err, tt.res.err)
			assert.Equal(t, tt.res.offset, offset)
			assert.Equal(t, tt.res.limit, limit)
			assert.Equal(t, tt.res.asc, asc)
		})
	}
}
