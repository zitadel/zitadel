//go:build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
)

// Like Polis' webhook suite, verify that SCIM membership changes reach an application.
// Delivery ordering and retry guarantees belong to the Actions worker, not this adapter.
func TestGroups_membershipWebhooks(t *testing.T) {
	instance := integration.NewInstance(CTX)
	ctx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	received := make(chan exec_repo.ContextInfoEvent, 16)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var event exec_repo.ContextInfoEvent
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		select {
		case received <- event:
		default:
			http.Error(w, "too many events", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	target := instance.CreateTarget(ctx, t, integration.TargetName(), server.URL, target_domain.TargetTypeWebhook, true, action.PayloadType_PAYLOAD_TYPE_JSON)
	wait, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(ct *assert.CollectT) {
		_, err := instance.Client.ActionV2.GetTarget(ctx, &action.GetTargetRequest{Id: target.GetId()})
		require.NoError(ct, err)
	}, wait, tick)
	for _, eventType := range []string{"group.user.added", "group.user.removed"} {
		condition := &action.Condition{ConditionType: &action.Condition_Event{Event: &action.EventExecution{Condition: &action.EventExecution_Event{Event: eventType}}}}
		instance.SetExecution(ctx, t, condition, []string{target.GetId()})
		require.EventuallyWithT(t, func(ct *assert.CollectT) {
			response, err := instance.Client.ActionV2.ListExecutions(ctx, &action.ListExecutionsRequest{Filters: []*action.ExecutionSearchFilter{{Filter: &action.ExecutionSearchFilter_InConditionsFilter{InConditionsFilter: &action.InConditionsFilter{Conditions: []*action.Condition{condition}}}}}})
			require.NoError(ct, err)
			require.Len(ct, response.Executions, 1)
			require.Len(ct, response.Executions[0].Targets, 1)
		}, wait, tick)
	}
	orgID := instance.DefaultOrg.GetId()
	userID := instance.CreateHumanUserVerified(ctx, orgID, integration.Email(), integration.Phone()).GetUserId()
	group, err := instance.Client.SCIM.Groups.Create(ctx, orgID, []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:schemas:core:2.0:Group"],"displayName":%q}`, integration.GroupName())))
	require.NoError(t, err)
	for _, op := range []string{"add", "remove"} {
		payload := []byte(fmt.Sprintf(`{"schemas":["urn:ietf:params:scim:api:messages:2.0:PatchOp"],"Operations":[{"op":%q,"path":"members","value":[{"value":%q}]}]}`, op, userID))
		require.NoError(t, instance.Client.SCIM.Groups.Update(ctx, orgID, group.ID, payload))
		// Retrying a completed request must not emit another membership change.
		require.NoError(t, instance.Client.SCIM.Groups.Update(ctx, orgID, group.ID, payload))
		select {
		case event := <-received:
			suffix := "added"
			if op == "remove" {
				suffix = "removed"
			}
			assert.Equal(t, "group.user."+suffix, event.EventType)
			assert.Equal(t, group.ID, event.AggregateID)
			assert.Equal(t, orgID, event.ResourceOwner)
			assert.Equal(t, instance.ID(), event.InstanceID)
			var data struct {
				UserID string `json:"userId"`
			}
			require.NoError(t, json.Unmarshal(event.EventPayload, &data))
			assert.Equal(t, userID, data.UserID)
		case <-time.After(wait):
			t.Fatal("membership webhook was not delivered")
		}
	}
}
