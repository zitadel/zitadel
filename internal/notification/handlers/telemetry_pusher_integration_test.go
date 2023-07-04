//go:build integration

package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TelemetryPusher(t *testing.T) {
	bodies := make(chan []byte, 0)
	t.Log("testing against instance with primary domain", PrimaryDomain)
	mockServer := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
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
	awaitMilestone(t, bodies, "InstanceCreated")
	project, err := MgmtClient.AddProject(IAMOwnerCtx, &management.AddProjectRequest{Name: "integration"})
	if err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, "ProjectCreated")
	if _, err = MgmtClient.AddOIDCApp(IAMOwnerCtx, &management.AddOIDCAppRequest{
		ProjectId: project.GetId(),
		Name:      "integration",
	}); err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, "ApplicationCreated")
	if _, err = SystemClient.RemoveInstance(SystemUserCTX, &system.RemoveInstanceRequest{InstanceId: InstanceID}); err != nil {
		t.Fatal(err)
	}
	awaitMilestone(t, bodies, "InstanceDeleted")
}

func awaitMilestone(t *testing.T, bodies chan []byte, expectMilestoneType string) {
	for {
		select {
		case body := <-bodies:
			plain := new(bytes.Buffer)
			if err := json.Indent(plain, body, "", "  "); err != nil {
				t.Fatal(err)
			}
			t.Log("received milestone", plain.String())
			milestone := struct {
				Type          string
				PrimaryDomain string
			}{}
			if err := json.Unmarshal(body, &milestone); err != nil {
				t.Error(err)
			}
			if milestone.Type == expectMilestoneType && milestone.PrimaryDomain == PrimaryDomain {
				return
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("timed out waiting for milestone")
		}
	}
}
