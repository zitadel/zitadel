//go:build integration

package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TelemetryPushMilestones(t *testing.T) {
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
	listener, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		t.Fatal(err)
	}
	mockServer.Listener = listener
	mockServer.Start()
	t.Cleanup(mockServer.Close)
	primaryDomain, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(CTX, SystemCTX)
	t.Log("testing against instance with primary domain", primaryDomain)
	awaitMilestone(t, bodies, primaryDomain, "InstanceCreated")
	project, err := Tester.Client.Mgmt.AddProject(iamOwnerCtx, &management.AddProjectRequest{Name: "integration"})
	if err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, primaryDomain, "ProjectCreated")
	if _, err = Tester.Client.Mgmt.AddOIDCApp(iamOwnerCtx, &management.AddOIDCAppRequest{
		ProjectId: project.GetId(),
		Name:      "integration",
	}); err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, primaryDomain, "ApplicationCreated")
	// TODO: trigger and await milestone AuthenticationSucceededOnInstance
	// TODO: trigger and await milestone AuthenticationSucceededOnApplication
	if _, err = Tester.Client.System.RemoveInstance(SystemCTX, &system.RemoveInstanceRequest{InstanceId: instanceID}); err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, primaryDomain, "InstanceDeleted")
}

func awaitMilestone(t *testing.T, bodies chan []byte, primaryDomain, expectMilestoneType string) {
	for {
		select {
		case body := <-bodies:
			plain := new(bytes.Buffer)
			if err := json.Indent(plain, body, "", "  "); err != nil {
				t.Fatal(err)
			}
			t.Log("received milestone", plain.String())
			milestone := struct {
				Type          string `json:"type"`
				PrimaryDomain string `json:"primaryDomain"`
			}{}
			if err := json.Unmarshal(body, &milestone); err != nil {
				t.Error(err)
			}
			if milestone.Type == expectMilestoneType && milestone.PrimaryDomain == primaryDomain {
				return
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("timed out waiting for milestone %s in domain %s", expectMilestoneType, primaryDomain)
		}
	}
}
