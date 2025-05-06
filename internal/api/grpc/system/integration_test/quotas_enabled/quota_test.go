//go:build integration

package quotas_enabled_test

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
	"github.com/zitadel/zitadel/internal/integration/sink"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	quota_pb "github.com/zitadel/zitadel/pkg/grpc/quota"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

var callURL = sink.CallURL(sink.ChannelQuota)

func TestServer_QuotaNotification_Limit(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	amount := 10
	percent := 50
	percentAmount := amount * percent / 100

	setQuota(t, instance.Instance.Id, amount, true, []*quota_pb.Notification{
		{Percent: uint32(percent), Repeat: true, CallUrl: callURL},
		{Percent: 100, Repeat: true, CallUrl: callURL},
	})

	sub := sink.Subscribe(CTX, sink.ChannelQuota)
	defer sub.Close()

	for i := 0; i < percentAmount; i++ {
		_, err := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, sub, quota.RequestsAllAuthenticated, percent)

	for i := 0; i < (amount - percentAmount); i++ {
		_, err := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, sub, quota.RequestsAllAuthenticated, 100)

	_, limitErr := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
	require.Error(t, limitErr)
}

func setQuota(t *testing.T, instanceID string, amount int, limit bool, notifications []*quota_pb.Notification) {
	_, err := integration.SystemClient().SetQuota(CTX, &system.SetQuotaRequest{
		InstanceId:    instanceID,
		Unit:          quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
		From:          timestamppb.Now(),
		ResetInterval: durationpb.New(time.Minute * 5),
		Amount:        uint64(amount),
		Limit:         limit,
		Notifications: notifications,
	})
	require.NoError(t, err)

	// wait for some time as there is an eventual consistency until the quota is applied and used in the interceptor
	time.Sleep(time.Second * 5)
}

func TestServer_QuotaNotification_NoLimit(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	amount := 10
	percent := 50
	percentAmount := amount * percent / 100

	setQuota(t, instance.Instance.Id, amount, false, []*quota_pb.Notification{
		{Percent: uint32(percent), Repeat: false, CallUrl: callURL},
		{Percent: 100, Repeat: true, CallUrl: callURL},
	})

	sub := sink.Subscribe(CTX, sink.ChannelQuota)
	defer sub.Close()

	for i := 0; i < percentAmount; i++ {
		_, err := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, sub, quota.RequestsAllAuthenticated, percent)

	for i := 0; i < (amount - percentAmount); i++ {
		_, err := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, sub, quota.RequestsAllAuthenticated, 100)

	for i := 0; i < amount; i++ {
		_, err := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
		require.NoErrorf(t, err, "error in %d call of %d", i, percentAmount)
	}
	awaitNotification(t, sub, quota.RequestsAllAuthenticated, 200)

	_, limitErr := instance.Client.Admin.GetDefaultOrg(iamCTX, &admin.GetDefaultOrgRequest{})
	require.NoError(t, limitErr)
}

func awaitNotification(t *testing.T, sub *sink.Subscription, unit quota.Unit, percent int) {
	for {
		select {
		case req, ok := <-sub.Recv():
			require.True(t, ok, "channel closed")

			plain := new(bytes.Buffer)
			if err := json.Indent(plain, req.Body, "", "  "); err != nil {
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
			if err := json.Unmarshal(req.Body, &event); err != nil {
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
	instance := integration.NewInstance(CTX)

	got, err := integration.SystemClient().SetQuota(CTX, &system.SetQuotaRequest{
		InstanceId:    instance.Instance.Id,
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
	require.Equal(t, got.Details.ResourceOwner, instance.Instance.Id)

	gotAlreadyExisting, errAlreadyExisting := integration.SystemClient().SetQuota(CTX, &system.SetQuotaRequest{
		InstanceId:    instance.Instance.Id,
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
	require.Equal(t, gotAlreadyExisting.Details.ResourceOwner, instance.Instance.Id)

	gotRemove, errRemove := integration.SystemClient().RemoveQuota(CTX, &system.RemoveQuotaRequest{
		InstanceId: instance.Instance.Id,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.NoError(t, errRemove)
	require.Equal(t, gotRemove.Details.ResourceOwner, instance.Instance.Id)

	gotRemoveAlready, errRemoveAlready := integration.SystemClient().RemoveQuota(CTX, &system.RemoveQuotaRequest{
		InstanceId: instance.Instance.Id,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.Error(t, errRemoveAlready)
	require.Nil(t, gotRemoveAlready)
}
