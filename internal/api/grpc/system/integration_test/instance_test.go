//go:build integration

package system_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/instance"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	system_pb "github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_ListInstances(t *testing.T) {
	isoInstance := integration.NewInstance(CTX)

	tests := []struct {
		name    string
		req     *system_pb.ListInstancesRequest
		want    []*instance.Instance
		wantErr bool
	}{
		{
			name: "empty query error",
			req: &system_pb.ListInstancesRequest{
				Queries: []*instance.Query{{}},
			},
			wantErr: true,
		},
		{
			name: "non-existing id",
			req: &system_pb.ListInstancesRequest{
				Queries: []*instance.Query{{
					Query: &instance.Query_IdQuery{
						IdQuery: &instance.IdsQuery{
							Ids: []string{"foo"},
						},
					},
				}},
			},
			want: []*instance.Instance{},
		},
		{
			name: "get 1 by id",
			req: &system_pb.ListInstancesRequest{
				Query: &object.ListQuery{
					Limit: 1,
				},
				Queries: []*instance.Query{{
					Query: &instance.Query_IdQuery{
						IdQuery: &instance.IdsQuery{
							Ids: []string{isoInstance.ID()},
						},
					},
				}},
			},
			want: []*instance.Instance{{
				Id: isoInstance.ID(),
			}},
		},
		{
			name: "non-existing domain",
			req: &system_pb.ListInstancesRequest{
				Queries: []*instance.Query{{
					Query: &instance.Query_DomainQuery{
						DomainQuery: &instance.DomainsQuery{
							Domains: []string{"foo"},
						},
					},
				}},
			},
			want: []*instance.Instance{},
		},
		{
			name: "get 1 by domain",
			req: &system_pb.ListInstancesRequest{
				Query: &object.ListQuery{
					Limit: 1,
				},
				Queries: []*instance.Query{{
					Query: &instance.Query_DomainQuery{
						DomainQuery: &instance.DomainsQuery{
							Domains: []string{isoInstance.Domain},
						},
					},
				}},
			},
			want: []*instance.Instance{{
				Id: isoInstance.ID(),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := integration.SystemClient().ListInstances(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			got := resp.GetResult()
			require.Len(t, got, len(tt.want))
			for i := 0; i < len(tt.want); i++ {
				assert.Equalf(t, tt.want[i].GetId(), got[i].GetId(), "instance[%d] id", i)
			}
		})
	}
}
