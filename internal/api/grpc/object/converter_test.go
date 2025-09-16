package object

import (
	"testing"

	"github.com/stretchr/testify/assert"

	object_pb "github.com/zitadel/zitadel/pkg/grpc/object"
)

func Test_ListQueryToModel(t *testing.T) {
	type args struct {
		req *object_pb.ListQuery
	}
	type res struct {
		offset, limit uint64
		asc           bool
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
			name: "all fields empty",
			args: args{
				req: &object_pb.ListQuery{
					Offset: 0,
					Limit:  0,
					Asc:    false,
				},
			},
			res: res{
				offset: 0,
				limit:  0,
				asc:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, limit, asc := ListQueryToModel(tt.args.req)
			assert.Equal(t, tt.res.offset, offset)
			assert.Equal(t, tt.res.limit, limit)
			assert.Equal(t, tt.res.asc, asc)
		})
	}
}
