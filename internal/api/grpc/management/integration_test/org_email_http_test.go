//go:build integration

package management_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

func TestServer_AddOrgEmailProviderHTTP(t *testing.T) {
	tests := []struct {
		name    string
		req     *mgmt_pb.AddOrgEmailProviderHTTPRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.AddOrgEmailProviderHTTPRequest{
				Endpoint:    "http://relay.example.com/email",
				Description: "test http provider",
			},
		},
		{
			name: "missing endpoint",
			req: &mgmt_pb.AddOrgEmailProviderHTTPRequest{
				Description: "missing endpoint",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrgEmailProviderHTTP(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.GetId())
			assert.NotEmpty(t, got.GetSigningKey())
			integration.AssertDetails(t, &mgmt_pb.AddOrgEmailProviderHTTPResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			}, got)
		})
	}
}

func TestServer_UpdateOrgEmailProviderHTTP(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderHTTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderHTTPRequest{
		Endpoint:    "http://relay.example.com/email",
		Description: "http update test",
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	tests := []struct {
		name    string
		req     *mgmt_pb.UpdateOrgEmailProviderHTTPRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.UpdateOrgEmailProviderHTTPRequest{
				Id:          providerID,
				Endpoint:    "http://relay2.example.com/email",
				Description: "updated http provider",
			},
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.UpdateOrgEmailProviderHTTPRequest{
				Id:       "nonexistent-id",
				Endpoint: "http://relay.example.com/email",
			},
			wantErr: true,
		},
		{
			name: "missing endpoint",
			req: &mgmt_pb.UpdateOrgEmailProviderHTTPRequest{
				Id:          providerID,
				Description: "missing endpoint",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrgEmailProviderHTTP(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.UpdateOrgEmailProviderHTTPResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			}, got)
		})
	}
}
