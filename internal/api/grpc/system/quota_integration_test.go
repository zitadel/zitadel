//go:build integration

package system_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	quota_pb "github.com/zitadel/zitadel/pkg/grpc/quota"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_QuotaNotification(t *testing.T) {
	bodies := make(chan []byte, 0)
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		if r.Header.Get("single-value") != "single-value" {
			t.Error("single-value header not set")
		}
		if reflect.DeepEqual(r.Header.Get("multi-value"), "multi-value-1,multi-value-2") {
			t.Error("single-value header not set")
		}
		bodies <- body
		w.WriteHeader(http.StatusOK)
	}))
	host := "localhost:8082"
	listener, err := net.Listen("tcp", host)
	if err != nil {
		t.Fatal(err)
	}
	mockServer.Listener = listener
	mockServer.Start()
	t.Cleanup(mockServer.Close)

	_, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(CTX, SystemCTX)
	amount := 10
	percent := 50
	percentAmount := amount / percent * 100

	_, err = Client.AddQuota(SystemCTX, &system.AddQuotaRequest{
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
				CallUrl: "http://" + host,
			},
			{
				Percent: 100,
				Repeat:  true,
				CallUrl: "http://" + host,
			},
		},
	})
	require.NoError(t, err)

	for i := 0; i < percentAmount; i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoError(t, err)
	}
	awaitNotification(t, bodies, quota.RequestsAllAuthenticated, percent)

	for i := 0; i < (amount - percentAmount); i++ {
		_, err := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
		require.NoError(t, err)
	}
	awaitNotification(t, bodies, quota.RequestsAllAuthenticated, 100)

	_, limitErr := Tester.Client.Admin.GetDefaultOrg(iamOwnerCtx, &admin.GetDefaultOrgRequest{})
	require.Error(t, limitErr)
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
			var event *quota.NotificationDueEvent
			if err := json.Unmarshal(body, event); err != nil {
				t.Error(err)
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

	got, err := Client.AddQuota(SystemCTX, &system.AddQuotaRequest{
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
				CallUrl: "url",
			},
		},
	})
	require.NoError(t, err)
	integration.AssertObjectDetails(t, &system.AddQuotaResponse{
		Details: &object.ObjectDetails{
			ChangeDate:    timestamppb.Now(),
			ResourceOwner: instanceID,
		},
	}, got)

	gotAlreadyExisting, errAlreadyExisting := Client.AddQuota(SystemCTX, &system.AddQuotaRequest{
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
				CallUrl: "url",
			},
		},
	})
	require.Error(t, errAlreadyExisting)
	require.Nil(t, gotAlreadyExisting)

	gotRemove, errRemove := Client.RemoveQuota(SystemCTX, &system.RemoveQuotaRequest{
		InstanceId: instanceID,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.NoError(t, errRemove)
	integration.AssertObjectDetails(t, &system.RemoveQuotaResponse{
		Details: &object.ObjectDetails{
			ChangeDate:    timestamppb.Now(),
			ResourceOwner: instanceID,
		},
	}, gotRemove)

	gotRemoveAlready, errRemoveAlready := Client.RemoveQuota(SystemCTX, &system.RemoveQuotaRequest{
		InstanceId: instanceID,
		Unit:       quota_pb.Unit_UNIT_REQUESTS_ALL_AUTHENTICATED,
	})
	require.Error(t, errRemoveAlready)
	require.Nil(t, gotRemoveAlready)
}
