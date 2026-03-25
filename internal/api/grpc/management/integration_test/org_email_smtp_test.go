//go:build integration

package management_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
)

func TestServer_AddOrgEmailProviderSMTP(t *testing.T) {
	tests := []struct {
		name    string
		req     *mgmt_pb.AddOrgEmailProviderSMTPRequest
		wantErr bool
	}{
		{
			name: "success with plain auth",
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderAddress: "noreply@example.com",
				SenderName:    "Test Sender",
				Tls:           true,
				Host:          "smtp.example.com:587",
				User:          "smtpuser",
				Description:   "test smtp provider",
				Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
					Plain: &mgmt_pb.OrgSMTPPlainAuth{
						Password: "test-password",
					},
				},
			},
		},
		{
			name: "success with no auth",
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderAddress: "noreply@example.com",
				SenderName:    "Test Sender No Auth",
				Tls:           false,
				Host:          "smtp.example.com:25",
				Description:   "test smtp no auth",
				Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_None{
					None: &mgmt_pb.OrgSMTPNoAuth{},
				},
			},
		},
		{
			name: "missing sender address",
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderName:  "Test Sender",
				Host:        "smtp.example.com:587",
				Description: "missing sender",
			},
			wantErr: true,
		},
		{
			name: "missing host",
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderAddress: "noreply@example.com",
				SenderName:    "Test Sender",
				Description:   "missing host",
			},
			wantErr: true,
		},
		{
			name: "missing sender name",
			req: &mgmt_pb.AddOrgEmailProviderSMTPRequest{
				SenderAddress: "noreply@example.com",
				Host:          "smtp.example.com:587",
				Description:   "missing sender name",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrgEmailProviderSMTP(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, got.GetId())
			integration.AssertDetails(t, &mgmt_pb.AddOrgEmailProviderSMTPResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			}, got)
		})
	}
}

func TestServer_GetOrgEmailProviderById(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "get-by-id@example.com",
		SenderName:    "Get By ID",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "get by id test",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "test-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	tests := []struct {
		name    string
		req     *mgmt_pb.GetOrgEmailProviderByIdRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.GetOrgEmailProviderByIdRequest{
				Id: providerID,
			},
		},
		{
			name: "not found",
			req: &mgmt_pb.GetOrgEmailProviderByIdRequest{
				Id: "nonexistent-id",
			},
			wantErr: true,
		},
		{
			name: "empty id",
			req: &mgmt_pb.GetOrgEmailProviderByIdRequest{
				Id: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				got, err := Client.GetOrgEmailProviderById(OrgCTX, tt.req)
				if tt.wantErr {
					require.Error(ct, err)
					return
				}
				require.NoError(ct, err)
				cfg := got.GetConfig()
				require.NotNil(ct, cfg)
				assert.Equal(ct, providerID, cfg.GetId())
				assert.Equal(ct, "get by id test", cfg.GetDescription())
				smtp := cfg.GetSmtp()
				require.NotNil(ct, smtp)
				assert.Equal(ct, "get-by-id@example.com", smtp.GetSenderAddress())
				assert.Equal(ct, "Get By ID", smtp.GetSenderName())
				assert.Equal(ct, true, smtp.GetTls())
				assert.Equal(ct, "smtp.example.com:587", smtp.GetHost())
				assert.Equal(ct, "smtpuser", smtp.GetUser())
			}, retryDuration, tick)
		})
	}
}

func TestServer_ListOrgEmailProviders(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "list-test@example.com",
		SenderName:    "List Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "list test provider",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "test-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		got, err := Client.ListOrgEmailProviders(OrgCTX, &mgmt_pb.ListOrgEmailProvidersRequest{})
		require.NoError(ct, err)
		require.NotNil(ct, got.GetDetails())
		results := got.GetResult()
		if !assert.GreaterOrEqual(ct, len(results), 1) {
			return
		}
		found := false
		for _, p := range results {
			if p.GetId() == providerID {
				found = true
				assert.Equal(ct, "list test provider", p.GetDescription())
			}
		}
		assert.True(ct, found, "added provider not found in list")
	}, retryDuration, tick)
}

func TestServer_UpdateOrgEmailProviderSMTP(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "update-test@example.com",
		SenderName:    "Update Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "before update",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "test-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	tests := []struct {
		name    string
		req     *mgmt_pb.UpdateOrgEmailProviderSMTPRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:            providerID,
				SenderAddress: "updated@example.com",
				SenderName:    "Updated Sender",
				Tls:           false,
				Host:          "smtp2.example.com:465",
				User:          "newuser",
				Description:   "after update",
				Auth: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Plain{
					Plain: &mgmt_pb.OrgSMTPPlainAuth{
						Password: "new-password",
					},
				},
			},
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:            "nonexistent-id",
				SenderAddress: "updated@example.com",
				SenderName:    "Updated Sender",
				Host:          "smtp2.example.com:465",
			},
			wantErr: true,
		},
		{
			name: "missing sender address",
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
				Id:         providerID,
				SenderName: "Updated Sender",
				Host:       "smtp2.example.com:465",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrgEmailProviderSMTP(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.UpdateOrgEmailProviderSMTPResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			}, got)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				getResp, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
					Id: providerID,
				})
				require.NoError(ct, err)
				cfg := getResp.GetConfig()
				require.NotNil(ct, cfg)
				smtp := cfg.GetSmtp()
				require.NotNil(ct, smtp)
				assert.Equal(ct, "updated@example.com", smtp.GetSenderAddress())
				assert.Equal(ct, "Updated Sender", smtp.GetSenderName())
				assert.Equal(ct, false, smtp.GetTls())
				assert.Equal(ct, "smtp2.example.com:465", smtp.GetHost())
				assert.Equal(ct, "newuser", smtp.GetUser())
				assert.Equal(ct, "after update", cfg.GetDescription())
			}, retryDuration, tick)
		})
	}
}

func TestServer_UpdateOrgEmailProviderSMTPPassword(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "password-test@example.com",
		SenderName:    "Password Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "password update test",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "original-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	tests := []struct {
		name    string
		req     *mgmt_pb.UpdateOrgEmailProviderSMTPPasswordRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPPasswordRequest{
				Id:       providerID,
				Password: "new-password",
			},
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.UpdateOrgEmailProviderSMTPPasswordRequest{
				Id:       "nonexistent-id",
				Password: "new-password",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateOrgEmailProviderSMTPPassword(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.UpdateOrgEmailProviderSMTPPasswordResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			}, got)
		})
	}
}
