//go:build integration

package system_test

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_Limits_AuditLogRetention(t *testing.T) {
	_, instanceID, iamOwnerCtx := Tester.UseIsolatedInstance(t, CTX, SystemCTX)
	userID, projectID, appID, projectGrantID := seedObjects(iamOwnerCtx, t)
	beforeTime := time.Now()
	farPast := timestamppb.New(beforeTime.Add(-10 * time.Hour).UTC())
	zeroCounts := &eventCounts{}
	seededCount := requireEventually(t, iamOwnerCtx, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "seeded events are > 0", assert.Greater, zeroCounts)
	}, "wait for seeded event assertions to pass")
	produceEvents(iamOwnerCtx, t, userID, appID, projectID, projectGrantID)
	addedCount := requireEventually(t, iamOwnerCtx, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "added events are > seeded events", assert.Greater, seededCount)
	}, "wait for added event assertions to pass")
	_, err := Tester.Client.System.SetLimits(SystemCTX, &system.SetLimitsRequest{
		InstanceId:        instanceID,
		AuditLogRetention: durationpb.New(time.Now().Sub(beforeTime)),
	})
	require.NoError(t, err)
	var limitedCounts *eventCounts
	requireEventually(t, iamOwnerCtx, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "limited events < added events", assert.Less, addedCount)
		counts.assertAll(t, c, "limited events > 0", assert.Greater, zeroCounts)
		limitedCounts = counts
	}, "wait for limited event assertions to pass")
	listedEvents, err := Tester.Client.Admin.ListEvents(iamOwnerCtx, &admin.ListEventsRequest{CreationDateFilter: &admin.ListEventsRequest_From{
		From: farPast,
	}})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(listedEvents.GetEvents()), limitedCounts.all, "ListEvents with from query older than retention doesn't return more events")
	listedEvents, err = Tester.Client.Admin.ListEvents(iamOwnerCtx, &admin.ListEventsRequest{CreationDateFilter: &admin.ListEventsRequest_Range{Range: &admin.ListEventsRequestCreationDateRange{
		Since: farPast,
	}}})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(listedEvents.GetEvents()), limitedCounts.all, "ListEvents with since query older than retention doesn't return more events")
	_, err = Tester.Client.System.ResetLimits(SystemCTX, &system.ResetLimitsRequest{
		InstanceId: instanceID,
	})
	require.NoError(t, err)
	requireEventually(t, iamOwnerCtx, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "with reset limit, added events are > seeded events", assert.Greater, seededCount)
	}, "wait for reset event assertions to pass")
}

func requireEventually(
	t *testing.T,
	ctx context.Context,
	userID, projectID, appID, projectGrantID string,
	assertCounts func(assert.TestingT, *eventCounts),
	msg string,
) (counts *eventCounts) {
	countTimeout := 30 * time.Second
	assertTimeout := countTimeout + time.Second
	countCtx, cancel := context.WithTimeout(ctx, countTimeout)
	defer cancel()
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		counts = countEvents(countCtx, t, userID, projectID, appID, projectGrantID)
		assertCounts(c, counts)
	}, assertTimeout, time.Second, msg)
	return counts
}

var runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(resourceType string, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return "test" + resourceType + "-" + string(b)
}

func seedObjects(ctx context.Context, t *testing.T) (string, string, string, string) {
	t.Helper()
	project, err := Tester.Client.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: randomString("project", 5),
	})
	require.NoError(t, err)
	app, err := Tester.Client.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		Name:      randomString("app", 5),
		ProjectId: project.GetId(),
	})
	org, err := Tester.Client.Mgmt.AddOrg(ctx, &management.AddOrgRequest{
		Name: randomString("org", 5),
	})
	require.NoError(t, err)
	role := randomString("role", 5)
	require.NoError(t, err)
	_, err = Tester.Client.Mgmt.AddProjectRole(ctx, &management.AddProjectRoleRequest{
		ProjectId:   project.GetId(),
		RoleKey:     role,
		DisplayName: role,
	})
	require.NoError(t, err)
	projectGrant, err := Tester.Client.Mgmt.AddProjectGrant(ctx, &management.AddProjectGrantRequest{
		ProjectId:    project.GetId(),
		GrantedOrgId: org.GetId(),
		RoleKeys:     []string{role},
	})
	require.NoError(t, err)
	user, err := Tester.Client.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	userID := user.GetUser().GetId()
	requireUserEvent(ctx, t, userID)
	return userID, project.GetId(), app.GetAppId(), projectGrant.GetGrantId()
}

func produceEvents(ctx context.Context, t *testing.T, machineID, appID, projectID, grantID string) {
	t.Helper()
	_, err := Tester.Client.Mgmt.UpdateOrg(ctx, &management.UpdateOrgRequest{
		Name: randomString("org", 5),
	})
	require.NoError(t, err)
	_, err = Tester.Client.Mgmt.UpdateProject(ctx, &management.UpdateProjectRequest{
		Id:   projectID,
		Name: randomString("project", 5),
	})
	require.NoError(t, err)
	_, err = Tester.Client.Mgmt.UpdateApp(ctx, &management.UpdateAppRequest{
		AppId:     appID,
		ProjectId: projectID,
		Name:      randomString("app", 5),
	})
	require.NoError(t, err)
	requireUserEvent(ctx, t, machineID)
	_, err = Tester.Client.Mgmt.UpdateProjectGrant(ctx, &management.UpdateProjectGrantRequest{
		ProjectId: projectID,
		GrantId:   grantID,
	})
	require.NoError(t, err)
}

func requireUserEvent(ctx context.Context, t *testing.T, machineID string) {
	_, err := Tester.Client.Mgmt.UpdateMachine(ctx, &management.UpdateMachineRequest{
		UserId: machineID,
		Name:   randomString("machine", 5),
	})
	require.NoError(t, err)
}

type eventCounts struct {
	all, myUser, aUser, grant, project, app, org int
}

func (e *eventCounts) assertAll(t *testing.T, c assert.TestingT, name string, compare assert.ComparisonAssertionFunc, than *eventCounts) {
	t.Run(name, func(t *testing.T) {
		compare(c, e.all, than.all, "ListEvents")
		compare(c, e.myUser, than.myUser, "ListMyUserChanges")
		compare(c, e.aUser, than.aUser, "ListUserChanges")
		compare(c, e.grant, than.grant, "ListProjectGrantChanges")
		compare(c, e.project, than.project, "ListProjectChanges")
		compare(c, e.app, than.app, "ListAppChanges")
		compare(c, e.org, than.org, "ListOrgChanges")
	})
}

func countEvents(ctx context.Context, t *testing.T, userID, projectID, appID, grantID string) *eventCounts {
	t.Helper()
	counts := new(eventCounts)
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Admin.ListEvents(ctx, &admin.ListEventsRequest{})
		require.NoError(t, err)
		counts.all = len(result.GetEvents())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Auth.ListMyUserChanges(ctx, &auth.ListMyUserChangesRequest{})
		require.NoError(t, err)
		counts.myUser = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Mgmt.ListUserChanges(ctx, &management.ListUserChangesRequest{UserId: userID})
		require.NoError(t, err)
		counts.aUser = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Mgmt.ListAppChanges(ctx, &management.ListAppChangesRequest{ProjectId: projectID, AppId: appID})
		require.NoError(t, err)
		counts.app = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Mgmt.ListOrgChanges(ctx, &management.ListOrgChangesRequest{})
		require.NoError(t, err)
		counts.org = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Mgmt.ListProjectChanges(ctx, &management.ListProjectChangesRequest{ProjectId: projectID})
		require.NoError(t, err)
		counts.project = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := Tester.Client.Mgmt.ListProjectGrantChanges(ctx, &management.ListProjectGrantChangesRequest{ProjectId: projectID, GrantId: grantID})
		require.NoError(t, err)
		counts.grant = len(result.GetResult())
	}()
	wg.Wait()
	return counts
}
