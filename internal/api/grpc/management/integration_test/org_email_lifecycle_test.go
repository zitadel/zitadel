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
	"github.com/zitadel/zitadel/pkg/grpc/settings"
)

func TestServer_ActivateOrgEmailProvider(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "activate-test@example.com",
		SenderName:    "Activate Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "activate test",
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
		req     *mgmt_pb.ActivateOrgEmailProviderRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.ActivateOrgEmailProviderRequest{
				Id: providerID,
			},
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.ActivateOrgEmailProviderRequest{
				Id: "nonexistent-id",
			},
			wantErr: true,
		},
		{
			name: "empty id",
			req: &mgmt_pb.ActivateOrgEmailProviderRequest{
				Id: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.ActivateOrgEmailProvider(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.ActivateOrgEmailProviderResponse{
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
				assert.Equal(ct, settings.EmailProviderState_EMAIL_PROVIDER_ACTIVE, getResp.GetConfig().GetState())
			}, retryDuration, tick)
		})
	}
}

func TestServer_DeactivateOrgEmailProvider(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "deactivate-test@example.com",
		SenderName:    "Deactivate Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "deactivate test",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "test-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()

	_, err = Client.ActivateOrgEmailProvider(OrgCTX, &mgmt_pb.ActivateOrgEmailProviderRequest{
		Id: providerID,
	})
	require.NoError(t, err)

	tests := []struct {
		name    string
		req     *mgmt_pb.DeactivateOrgEmailProviderRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.DeactivateOrgEmailProviderRequest{
				Id: providerID,
			},
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.DeactivateOrgEmailProviderRequest{
				Id: "nonexistent-id",
			},
			wantErr: true,
		},
		{
			name: "empty id",
			req: &mgmt_pb.DeactivateOrgEmailProviderRequest{
				Id: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.DeactivateOrgEmailProvider(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.DeactivateOrgEmailProviderResponse{
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
				assert.Equal(ct, settings.EmailProviderState_EMAIL_PROVIDER_INACTIVE, getResp.GetConfig().GetState())
			}, retryDuration, tick)
		})
	}
}

func TestServer_RemoveOrgEmailProvider(t *testing.T) {
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "remove-test@example.com",
		SenderName:    "Remove Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "remove test",
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
		req     *mgmt_pb.RemoveOrgEmailProviderRequest
		wantErr bool
	}{
		{
			name: "success",
			req: &mgmt_pb.RemoveOrgEmailProviderRequest{
				Id: providerID,
			},
		},
		{
			name: "already removed",
			req: &mgmt_pb.RemoveOrgEmailProviderRequest{
				Id: providerID,
			},
			wantErr: true,
		},
		{
			name: "nonexistent id",
			req: &mgmt_pb.RemoveOrgEmailProviderRequest{
				Id: "nonexistent-id",
			},
			wantErr: true,
		},
		{
			name: "empty id",
			req: &mgmt_pb.RemoveOrgEmailProviderRequest{
				Id: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RemoveOrgEmailProvider(OrgCTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, &mgmt_pb.RemoveOrgEmailProviderResponse{
				Details: &object.ObjectDetails{
					ResourceOwner: Instance.DefaultOrg.Id,
					ChangeDate:    timestamppb.Now(),
				},
			}, got)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
			assert.EventuallyWithT(t, func(ct *assert.CollectT) {
				_, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
					Id: providerID,
				})
				require.Error(ct, err)
			}, retryDuration, tick)
		})
	}
}

func TestServer_OrgEmailProviderSMTP_FullLifecycle(t *testing.T) {
	// 1. Add SMTP provider
	addResp, err := Client.AddOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "lifecycle@example.com",
		SenderName:    "Lifecycle Test",
		Tls:           true,
		Host:          "smtp.example.com:587",
		User:          "smtpuser",
		Description:   "lifecycle test",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "test-password",
			},
		},
	})
	require.NoError(t, err)
	providerID := addResp.GetId()
	assert.NotEmpty(t, providerID)

	// 2. Get by ID and verify config
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		getResp, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
			Id: providerID,
		})
		require.NoError(ct, err)
		cfg := getResp.GetConfig()
		require.NotNil(ct, cfg)
		assert.Equal(ct, providerID, cfg.GetId())
		assert.Equal(ct, "lifecycle test", cfg.GetDescription())
	}, retryDuration, tick)

	// 3. List and verify it appears
	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		listResp, err := Client.ListOrgEmailProviders(OrgCTX, &mgmt_pb.ListOrgEmailProvidersRequest{})
		require.NoError(ct, err)
		found := false
		for _, p := range listResp.GetResult() {
			if p.GetId() == providerID {
				found = true
			}
		}
		assert.True(ct, found, "provider not found in list")
	}, retryDuration, tick)

	// 4. Activate
	_, err = Client.ActivateOrgEmailProvider(OrgCTX, &mgmt_pb.ActivateOrgEmailProviderRequest{
		Id: providerID,
	})
	require.NoError(t, err)

	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		getResp, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
			Id: providerID,
		})
		require.NoError(ct, err)
		assert.Equal(ct, settings.EmailProviderState_EMAIL_PROVIDER_ACTIVE, getResp.GetConfig().GetState())
	}, retryDuration, tick)

	// 5. Update SMTP config
	_, err = Client.UpdateOrgEmailProviderSMTP(OrgCTX, &mgmt_pb.UpdateOrgEmailProviderSMTPRequest{
		Id:            providerID,
		SenderAddress: "updated-lifecycle@example.com",
		SenderName:    "Updated Lifecycle",
		Tls:           false,
		Host:          "smtp2.example.com:465",
		User:          "newuser",
		Description:   "updated lifecycle",
		Auth: &mgmt_pb.UpdateOrgEmailProviderSMTPRequest_Plain{
			Plain: &mgmt_pb.OrgSMTPPlainAuth{
				Password: "updated-password",
			},
		},
	})
	require.NoError(t, err)

	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		getResp, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
			Id: providerID,
		})
		require.NoError(ct, err)
		smtp := getResp.GetConfig().GetSmtp()
		require.NotNil(ct, smtp)
		assert.Equal(ct, "updated-lifecycle@example.com", smtp.GetSenderAddress())
		assert.Equal(ct, "Updated Lifecycle", smtp.GetSenderName())
	}, retryDuration, tick)

	// 6. Update password
	_, err = Client.UpdateOrgEmailProviderSMTPPassword(OrgCTX, &mgmt_pb.UpdateOrgEmailProviderSMTPPasswordRequest{
		Id:       providerID,
		Password: "final-password",
	})
	require.NoError(t, err)

	// 7. Deactivate
	_, err = Client.DeactivateOrgEmailProvider(OrgCTX, &mgmt_pb.DeactivateOrgEmailProviderRequest{
		Id: providerID,
	})
	require.NoError(t, err)

	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		getResp, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
			Id: providerID,
		})
		require.NoError(ct, err)
		assert.Equal(ct, settings.EmailProviderState_EMAIL_PROVIDER_INACTIVE, getResp.GetConfig().GetState())
	}, retryDuration, tick)

	// 8. Remove
	_, err = Client.RemoveOrgEmailProvider(OrgCTX, &mgmt_pb.RemoveOrgEmailProviderRequest{
		Id: providerID,
	})
	require.NoError(t, err)

	// 9. Verify removed
	retryDuration, tick = integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
	assert.EventuallyWithT(t, func(ct *assert.CollectT) {
		_, err := Client.GetOrgEmailProviderById(OrgCTX, &mgmt_pb.GetOrgEmailProviderByIdRequest{
			Id: providerID,
		})
		require.Error(ct, err)
	}, retryDuration, tick)
}

func TestServer_OrgEmailProvider_PermissionDenied(t *testing.T) {
	_, err := Client.AddOrgEmailProviderSMTP(CTX, &mgmt_pb.AddOrgEmailProviderSMTPRequest{
		SenderAddress: "perm-test@example.com",
		SenderName:    "Perm Test",
		Host:          "smtp.example.com:587",
		Auth: &mgmt_pb.AddOrgEmailProviderSMTPRequest_None{
			None: &mgmt_pb.OrgSMTPNoAuth{},
		},
	})
	require.Error(t, err)

	_, err = Client.ListOrgEmailProviders(CTX, &mgmt_pb.ListOrgEmailProvidersRequest{})
	require.Error(t, err)

	_, err = Client.AddOrgEmailProviderHTTP(CTX, &mgmt_pb.AddOrgEmailProviderHTTPRequest{
		Endpoint: "http://relay.example.com/email",
	})
	require.Error(t, err)

	_, err = Client.ActivateOrgEmailProvider(CTX, &mgmt_pb.ActivateOrgEmailProviderRequest{
		Id: "some-id",
	})
	require.Error(t, err)

	_, err = Client.DeactivateOrgEmailProvider(CTX, &mgmt_pb.DeactivateOrgEmailProviderRequest{
		Id: "some-id",
	})
	require.Error(t, err)

	_, err = Client.RemoveOrgEmailProvider(CTX, &mgmt_pb.RemoveOrgEmailProviderRequest{
		Id: "some-id",
	})
	require.Error(t, err)
}
