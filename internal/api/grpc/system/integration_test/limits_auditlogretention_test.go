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

	"github.com/zitadel/zitadel/v2/internal/integration"
	"github.com/zitadel/zitadel/v2/pkg/grpc/admin"
	"github.com/zitadel/zitadel/v2/pkg/grpc/auth"
	"github.com/zitadel/zitadel/v2/pkg/grpc/management"
	"github.com/zitadel/zitadel/v2/pkg/grpc/system"
)

func TestServer_Limits_AuditLogRetention(t *testing.T) {
	t.Parallel()

	isoInstance := integration.NewInstance(CTX)
	iamOwnerCtx := isoInstance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	userID, projectID, appID, projectGrantID := seedObjects(iamOwnerCtx, t, isoInstance.Client)
	beforeTime := time.Now()
	farPast := timestamppb.New(beforeTime.Add(-10 * time.Hour).UTC())
	zeroCounts := &eventCounts{}
	seededCount := requireEventually(t, iamOwnerCtx, isoInstance.Client, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "seeded events are > 0", assert.Greater, zeroCounts)
	}, "wait for seeded event assertions to pass")
	produceEvents(iamOwnerCtx, t, isoInstance.Client, userID, appID, projectID, projectGrantID)
	addedCount := requireEventually(t, iamOwnerCtx, isoInstance.Client, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "added events are > seeded events", assert.Greater, seededCount)
	}, "wait for added event assertions to pass")
	_, err := integration.SystemClient().SetLimits(CTX, &system.SetLimitsRequest{
		InstanceId:        isoInstance.ID(),
		AuditLogRetention: durationpb.New(time.Now().Sub(beforeTime)),
	})
	require.NoError(t, err)
	var limitedCounts *eventCounts
	requireEventually(t, iamOwnerCtx, isoInstance.Client, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "limited events < added events", assert.Less, addedCount)
		counts.assertAll(t, c, "limited events > 0", assert.Greater, zeroCounts)
		limitedCounts = counts
	}, "wait for limited event assertions to pass")
	listedEvents, err := isoInstance.Client.Admin.ListEvents(iamOwnerCtx, &admin.ListEventsRequest{CreationDateFilter: &admin.ListEventsRequest_From{
		From: farPast,
	}})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(listedEvents.GetEvents()), limitedCounts.all, "ListEvents with from query older than retention doesn't return more events")
	listedEvents, err = isoInstance.Client.Admin.ListEvents(iamOwnerCtx, &admin.ListEventsRequest{CreationDateFilter: &admin.ListEventsRequest_Range{Range: &admin.ListEventsRequestCreationDateRange{
		Since: farPast,
	}}})
	require.NoError(t, err)
	assert.LessOrEqual(t, len(listedEvents.GetEvents()), limitedCounts.all, "ListEvents with since query older than retention doesn't return more events")
	_, err = integration.SystemClient().ResetLimits(CTX, &system.ResetLimitsRequest{
		InstanceId: isoInstance.ID(),
	})
	require.NoError(t, err)
	requireEventually(t, iamOwnerCtx, isoInstance.Client, userID, projectID, appID, projectGrantID, func(c assert.TestingT, counts *eventCounts) {
		counts.assertAll(t, c, "with reset limit, added events are > seeded events", assert.Greater, seededCount)
	}, "wait for reset event assertions to pass")
}

func requireEventually(
	t *testing.T,
	ctx context.Context,
	cc *integration.Client,
	userID, projectID, appID, projectGrantID string,
	assertCounts func(assert.TestingT, *eventCounts),
	msg string,
) (counts *eventCounts) {
	countTimeout := 30 * time.Second
	assertTimeout := countTimeout + time.Second
	countCtx, cancel := context.WithTimeout(ctx, countTimeout)
	defer cancel()
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		counts = countEvents(countCtx, c, cc, userID, projectID, appID, projectGrantID)
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

func seedObjects(ctx context.Context, t *testing.T, cc *integration.Client) (string, string, string, string) {
	t.Helper()
	project, err := cc.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: randomString("project", 5),
	})
	require.NoError(t, err)
	app, err := cc.Mgmt.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		Name:      randomString("app", 5),
		ProjectId: project.GetId(),
	})
	org, err := cc.Mgmt.AddOrg(ctx, &management.AddOrgRequest{
		Name: randomString("org", 5),
	})
	require.NoError(t, err)
	role := randomString("role", 5)
	require.NoError(t, err)
	_, err = cc.Mgmt.AddProjectRole(ctx, &management.AddProjectRoleRequest{
		ProjectId:   project.GetId(),
		RoleKey:     role,
		DisplayName: role,
	})
	require.NoError(t, err)
	projectGrant, err := cc.Mgmt.AddProjectGrant(ctx, &management.AddProjectGrantRequest{
		ProjectId:    project.GetId(),
		GrantedOrgId: org.GetId(),
		RoleKeys:     []string{role},
	})
	require.NoError(t, err)
	user, err := cc.Auth.GetMyUser(ctx, &auth.GetMyUserRequest{})
	require.NoError(t, err)
	userID := user.GetUser().GetId()
	requireUserEvent(ctx, t, cc, userID)
	return userID, project.GetId(), app.GetAppId(), projectGrant.GetGrantId()
}

func produceEvents(ctx context.Context, t *testing.T, cc *integration.Client, machineID, appID, projectID, grantID string) {
	t.Helper()
	_, err := cc.Mgmt.UpdateOrg(ctx, &management.UpdateOrgRequest{
		Name: randomString("org", 5),
	})
	require.NoError(t, err)
	_, err = cc.Mgmt.UpdateProject(ctx, &management.UpdateProjectRequest{
		Id:   projectID,
		Name: randomString("project", 5),
	})
	require.NoError(t, err)
	_, err = cc.Mgmt.UpdateApp(ctx, &management.UpdateAppRequest{
		AppId:     appID,
		ProjectId: projectID,
		Name:      randomString("app", 5),
	})
	require.NoError(t, err)
	requireUserEvent(ctx, t, cc, machineID)
	_, err = cc.Mgmt.UpdateProjectGrant(ctx, &management.UpdateProjectGrantRequest{
		ProjectId: projectID,
		GrantId:   grantID,
	})
	require.NoError(t, err)
}

func requireUserEvent(ctx context.Context, t *testing.T, cc *integration.Client, machineID string) {
	_, err := cc.Mgmt.UpdateMachine(ctx, &management.UpdateMachineRequest{
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

func countEvents(ctx context.Context, t assert.TestingT, cc *integration.Client, userID, projectID, appID, grantID string) *eventCounts {
	counts := new(eventCounts)
	var wg sync.WaitGroup
	wg.Add(7)
	go func() {
		defer wg.Done()
		result, err := cc.Admin.ListEvents(ctx, &admin.ListEventsRequest{})
		assert.NoError(t, err)
		counts.all = len(result.GetEvents())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Auth.ListMyUserChanges(ctx, &auth.ListMyUserChangesRequest{})
		assert.NoError(t, err)
		counts.myUser = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Mgmt.ListUserChanges(ctx, &management.ListUserChangesRequest{UserId: userID})
		assert.NoError(t, err)
		counts.aUser = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Mgmt.ListAppChanges(ctx, &management.ListAppChangesRequest{ProjectId: projectID, AppId: appID})
		assert.NoError(t, err)
		counts.app = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Mgmt.ListOrgChanges(ctx, &management.ListOrgChangesRequest{})
		assert.NoError(t, err)
		counts.org = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Mgmt.ListProjectChanges(ctx, &management.ListProjectChangesRequest{ProjectId: projectID})
		assert.NoError(t, err)
		counts.project = len(result.GetResult())
	}()
	go func() {
		defer wg.Done()
		result, err := cc.Mgmt.ListProjectGrantChanges(ctx, &management.ListProjectGrantChangesRequest{ProjectId: projectID, GrantId: grantID})
		assert.NoError(t, err)
		counts.grant = len(result.GetResult())
	}()
	wg.Wait()
	return counts
}
