//go:build integration

package system_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	quota_pb "github.com/zitadel/zitadel/pkg/grpc/quota"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

var callURL = "http://localhost:" + integration.PortQuotaServer

func TestServer_QuotaNotification_Limit(t *testing.T) {
	_, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(CTX, SystemCTX)
	amount := 10
	percent := 50
	percentAmount := amount * percent / 100

	_, err := Tester.Client.System.SetQuota(SystemCTX, &system.SetQuotaRequest{
		InstanceId:    instanceID,
		Unit:          quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
		From:          timestamppb.Now(),
		ResetInterval: durationpb.New(time.Minute * 5),
		Amount:        uint64(amount),
		Limit:         true,
		Notifications: []*quota_pb.Notification{
			{
				Percent: uint32(percent),
				Repeat:  true,
				CallUrl: callURL,
			},
			{
				Percent: 100,
				Repeat:  true,
				CallUrl: callURL,
			},
		},
	})
	require.NoError(t, err)

	for i := 0; i < percentAmount; i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, Tester.QuotaNotificationChan, quota.RequestsAllAuthenticated, percent)

	for i := 0; i < (amount - percentAmount); i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, Tester.QuotaNotificationChan, quota.RequestsAllAuthenticated, 100)

	_, limitErr := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
	require.Error(t, limitErr)
}

func TestServer_QuotaNotification_NoLimit(t *testing.T) {
	_, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(CTX, SystemCTX)
	amount := 10
	percent := 50
	percentAmount := amount * percent / 100

	_, err := Tester.Client.System.SetQuota(SystemCTX, &system.SetQuotaRequest{
		InstanceId:    instanceID,
		Unit:          quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
		From:          timestamppb.Now(),
		ResetInterval: durationpb.New(time.Minute * 5),
		Amount:        uint64(amount),
		Limit:         false,
		Notifications: []*quota_pb.Notification{
			{
				Percent: uint32(percent),
				Repeat:  false,
				CallUrl: callURL,
			},
			{
				Percent: 100,
				Repeat:  true,
				CallUrl: callURL,
			},
		},
	})
	require.NoError(t, err)

	for i := 0; i < percentAmount; i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, Tester.QuotaNotificationChan, quota.RequestsAllAuthenticated, percent)

	for i := 0; i < (amount - percentAmount); i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, Tester.QuotaNotificationChan, quota.RequestsAllAuthenticated, 100)

	for i := 0; i < amount; i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, Tester.QuotaNotificationChan, quota.RequestsAllAuthenticated, 200)

	_, limitErr := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
	require.NoError(t, limitErr)
}

func awaitNotification(t *testing.T, bodies chan []byte, unit quota.Unit, percent int) {
	for {
		select {
		case body := <-bodies:
			plain := new(bytes.Buffer)
			if err := json.Indent(plain, body, "", "  "); err != nil {
				t.Fatal(err)
			}
			t.Log("received notificationDueEvent", plain.String())
			event := struct {
				Unit        quota.Unit `json:"unit"`
				ID          string     `json:"id"`
				CallURL     string     `json:"callURL"`
				PeriodStart time.Time  `json:"periodStart"`
				Threshold   uint16     `json:"threshold"`
				Usage       uint64     `json:"usage"`
			}{}
			if err := json.Unmarshal(body, &event); err != nil {
				t.Error(err)
			}
			if event.ID == "" {
				continue
			}
			if event.Unit == unit && event.Threshold == uint16(percent) {
				return
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("timed out waiting for unit %s and percent %d", strconv.Itoa(int(unit)), percent)
		}
	}
}

func TestServer_AddAndRemoveQuota(t *testing.T) {
	_, instanceID, _ := Tester.UseIsolatedInstance(CTX, SystemCTX)

	got, err := Tester.Client.System.SetQuota(SystemCTX, &system.SetQuotaRequest{
		InstanceId:    instanceID,
		Unit:          quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
		From:          timestamppb.Now(),
		ResetInterval: durationpb.New(time.Minute),
		Amount:        10,
		Limit:         true,
		Notifications: []*quota_pb.Notification{
			{
				Percent: 20,
				Repeat:  true,
				CallUrl: callURL,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, got.Details.ResourceOwner, instanceID)

	gotAlreadyExisting, errAlreadyExisting := Tester.Client.System.SetQuota(SystemCTX, &system.SetQuotaRequest{
		InstanceId:    instanceID,
		Unit:          quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
		From:          timestamppb.Now(),
		ResetInterval: durationpb.New(time.Minute),
		Amount:        10,
		Limit:         true,
		Notifications: []*quota_pb.Notification{
			{
				Percent: 20,
				Repeat:  true,
				CallUrl: callURL,
			},
		},
	})
	require.Nil(t, errAlreadyExisting)
	require.Equal(t, gotAlreadyExisting.Details.ResourceOwner, instanceID)

	gotRemove, errRemove := Tester.Client.System.RemoveQuota(SystemCTX, &system.RemoveQuotaRequest{
		InstanceId: instanceID,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.NoError(t, errRemove)
	require.Equal(t, gotRemove.Details.ResourceOwner, instanceID)

	gotRemoveAlready, errRemoveAlready := Tester.Client.System.RemoveQuota(SystemCTX, &system.RemoveQuotaRequest{
		InstanceId: instanceID,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.Error(t, errRemoveAlready)
	require.Nil(t, gotRemoveAlready)
}
