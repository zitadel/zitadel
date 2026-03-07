package command

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/risk"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// newTestRiskService creates a real risk.Service with deterministic rules only (no LLM).
func newTestRiskService(t *testing.T) *risk.Service {
	t.Helper()
	svc, err := risk.New(risk.Config{
		Enabled:               true,
		FailOpen:              false,
		FailureBurstThreshold: 5,
		HistoryWindow:         time.Hour,
		ContextChangeWindow:   15 * time.Minute,
		MaxSignalsPerUser:     50,
		MaxSignalsPerSession:  20,
	}, nil, nil)
	require.NoError(t, err)
	return svc
}

// makeSessionChecks builds a SessionCommands for a create_session with the given UA parameters.
// Callers must also set sessionWriteModel.UserID and sessionWriteModel.UserResourceOwner.
func makeSessionChecks(c *Commands, sessionID, fp, ip, ua string, now time.Time) *SessionCommands {
	userAgent := &domain.UserAgent{
		FingerprintID: gu.Ptr(fp),
		IP:            net.ParseIP(ip),
		Description:   gu.Ptr(ua),
	}
	checks := &SessionCommands{
		sessionWriteModel: NewSessionWriteModel(sessionID, "instance1"),
		sessionCommands:   []SessionCommand{},
		eventstore:        c.eventstore,
		createToken: func(string) (string, string, error) {
			return "tokenID-" + sessionID, "token-" + sessionID, nil
		},
		now:              func() time.Time { return now },
		operation:        sessionOperationCreate,
		currentUserAgent: userAgent,
	}
	checks.Start(context.Background(), userAgent)
	return checks
}

type fakeRiskEvaluator struct {
	enabled     bool
	failOpen    bool
	decision    risk.Decision
	evaluateErr error
	recorded    []risk.Signal
	findings    [][]risk.Finding
}

func (f *fakeRiskEvaluator) Enabled() bool { return f.enabled }

func (f *fakeRiskEvaluator) FailOpen() bool { return f.failOpen }

func (f *fakeRiskEvaluator) Evaluate(context.Context, risk.Signal) (risk.Decision, error) {
	if f.evaluateErr != nil {
		return risk.Decision{}, f.evaluateErr
	}
	return f.decision, nil
}

func (f *fakeRiskEvaluator) Record(_ context.Context, signal risk.Signal, findings []risk.Finding) error {
	f.recorded = append(f.recorded, signal)
	f.findings = append(f.findings, append([]risk.Finding(nil), findings...))
	return nil
}

func TestCommands_updateSession_blockedByRisk(t *testing.T) {
	t.Parallel()

	evaluator := &fakeRiskEvaluator{
		enabled: true,
		decision: risk.Decision{
			Allow:    false,
			Findings: []risk.Finding{{Name: "context_drift", Block: true}},
		},
	}
	c := &Commands{eventstore: expectEventstore()(t), riskEvaluator: evaluator}

	checks := &SessionCommands{
		sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
		sessionCommands:   []SessionCommand{},
		now: func() time.Time {
			return time.Now().UTC()
		},
		operation: sessionOperationCreate,
		currentUserAgent: &domain.UserAgent{
			FingerprintID: gu.Ptr("fp1"),
			IP:            net.ParseIP("2.2.2.2"),
			Description:   gu.Ptr("safari"),
		},
	}
	checks.sessionWriteModel.UserID = "user1"

	got, err := c.updateSession(authz.NewMockContext("instance1", "", ""), checks, nil, 0)
	require.ErrorIs(t, err, zerrors.ThrowPermissionDenied(nil, "COMMAND-RISK0", "Errors.PermissionDenied"))
	assert.Nil(t, got)
	require.Len(t, evaluator.recorded, 1)
	assert.Equal(t, risk.OutcomeBlocked, evaluator.recorded[0].Outcome)
}

func TestCommands_updateSession_riskFailOpen(t *testing.T) {
	t.Parallel()

	evaluator := &fakeRiskEvaluator{enabled: true, failOpen: true, evaluateErr: errors.New("boom")}
	testNow := time.Now().UTC()
	c := &Commands{
		eventstore: expectEventstore(
			expectPush(
				session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
					&domain.UserAgent{FingerprintID: gu.Ptr("fp1"), IP: net.ParseIP("1.2.3.4"), Description: gu.Ptr("firefox")},
				),
				session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate, "tokenID"),
			),
		)(t),
		riskEvaluator: evaluator,
	}

	checks := &SessionCommands{
		sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
		sessionCommands:   []SessionCommand{},
		eventstore:        c.eventstore,
		createToken: func(sessionID string) (string, string, error) {
			return "tokenID", "token", nil
		},
		now: func() time.Time {
			return testNow
		},
		operation: sessionOperationCreate,
		currentUserAgent: &domain.UserAgent{
			FingerprintID: gu.Ptr("fp1"),
			IP:            net.ParseIP("1.2.3.4"),
			Description:   gu.Ptr("firefox"),
		},
	}
	checks.Start(context.Background(), checks.currentUserAgent)
	checks.sessionWriteModel.UserID = "user1"
	checks.sessionWriteModel.UserResourceOwner = "org1"

	got, err := c.updateSession(authz.NewMockContext("instance1", "", ""), checks, nil, 0)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "token", got.NewToken)
	require.Len(t, evaluator.recorded, 1)
	assert.Equal(t, risk.OutcomeSuccess, evaluator.recorded[0].Outcome)
}

func TestCommands_updateSession_recordsRiskFindingsOnSuccess(t *testing.T) {
	t.Parallel()

	evaluator := &fakeRiskEvaluator{
		enabled: true,
		decision: risk.Decision{
			Allow:    true,
			Findings: []risk.Finding{{Name: "llm_high_risk", Source: "llm", Message: "model observed a risky pattern"}},
		},
	}
	testNow := time.Now().UTC()
	c := &Commands{
		eventstore: expectEventstore(
			expectPush(
				session.NewAddedEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate,
					&domain.UserAgent{FingerprintID: gu.Ptr("fp1"), IP: net.ParseIP("1.2.3.4"), Description: gu.Ptr("firefox")},
				),
				session.NewTokenSetEvent(context.Background(), &session.NewAggregate("sessionID", "instance1").Aggregate, "tokenID"),
			),
		)(t),
		riskEvaluator: evaluator,
	}

	checks := &SessionCommands{
		sessionWriteModel: NewSessionWriteModel("sessionID", "instance1"),
		sessionCommands:   []SessionCommand{},
		eventstore:        c.eventstore,
		createToken: func(sessionID string) (string, string, error) {
			return "tokenID", "token", nil
		},
		now: func() time.Time {
			return testNow
		},
		operation: sessionOperationCreate,
		currentUserAgent: &domain.UserAgent{
			FingerprintID: gu.Ptr("fp1"),
			IP:            net.ParseIP("1.2.3.4"),
			Description:   gu.Ptr("firefox"),
		},
	}
	checks.Start(context.Background(), checks.currentUserAgent)
	checks.sessionWriteModel.UserID = "user1"
	checks.sessionWriteModel.UserResourceOwner = "org1"

	got, err := c.updateSession(authz.NewMockContext("instance1", "", ""), checks, nil, 0)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Len(t, evaluator.recorded, 1)
	require.Len(t, evaluator.findings, 1)
	require.Len(t, evaluator.findings[0], 1)
	assert.Equal(t, "llm_high_risk", evaluator.findings[0][0].Name)
}

// TestCommands_updateSession_threeSessionsSameContextAllowed verifies that three
// consecutive create_session calls from the same user with identical IP and
// user-agent are all allowed. The real risk.Service is injected to exercise the
// deterministic rules (contextDrift, failureBurst) end-to-end.
func TestCommands_updateSession_threeSessionsSameContextAllowed(t *testing.T) {
	t.Parallel()

	riskSvc := newTestRiskService(t)
	base := time.Now().UTC()
	ua := &domain.UserAgent{FingerprintID: gu.Ptr("fp1"), IP: net.ParseIP("1.2.3.4"), Description: gu.Ptr("chrome")}

	sessions := []struct {
		id string
	}{
		{"session1"},
		{"session2"},
		{"session3"},
	}

	// Build eventstore expectations: each session pushes AddedEvent + TokenSetEvent.
	var expects []expect
	for _, s := range sessions {
		expects = append(expects, expectPush(
			session.NewAddedEvent(context.Background(), &session.NewAggregate(s.id, "instance1").Aggregate, ua),
			session.NewTokenSetEvent(context.Background(), &session.NewAggregate(s.id, "instance1").Aggregate, "tokenID-"+s.id),
		))
	}

	c := &Commands{
		eventstore:    expectEventstore(expects...)(t),
		riskEvaluator: riskSvc,
	}
	ctx := authz.NewMockContext("instance1", "", "")

	for i, s := range sessions {
		checks := makeSessionChecks(c, s.id, "fp1", "1.2.3.4", "chrome", base.Add(time.Duration(i)*time.Minute))
		checks.sessionWriteModel.UserID = "user1"
		checks.sessionWriteModel.UserResourceOwner = "org1"

		got, err := c.updateSession(ctx, checks, nil, 0)
		require.NoErrorf(t, err, "session %d should be allowed", i+1)
		require.NotNilf(t, got, "session %d should return a result", i+1)
		assert.Equalf(t, "token-"+s.id, got.NewToken, "session %d should have a token", i+1)

		// Record the successful signal so subsequent sessions see the history.
		require.NoError(t, riskSvc.Record(ctx, checks.riskSignal(ctx, "", risk.OutcomeSuccess), nil))
	}
}

// TestCommands_updateSession_contextDriftBlocksThirdSession verifies that the
// third create_session is blocked when both the IP and user-agent changed
// compared to the two preceding successful logins from the same user.
//
// Session 1 (chrome / 1.2.3.4) → allowed, signal recorded.
// Session 2 (chrome / 1.2.3.4) → allowed, signal recorded.
// Session 3 (safari / 9.9.9.9) → contextDrift fires → blocked.
func TestCommands_updateSession_contextDriftBlocksThirdSession(t *testing.T) {
	t.Parallel()

	riskSvc := newTestRiskService(t)
	base := time.Now().UTC()
	uaChrome := &domain.UserAgent{FingerprintID: gu.Ptr("fp1"), IP: net.ParseIP("1.2.3.4"), Description: gu.Ptr("chrome")}

	// Sessions 1 and 2 succeed and record signals; session 3 is blocked.
	passSessions := []struct{ id string }{{"session1"}, {"session2"}}
	var expects []expect
	for _, s := range passSessions {
		expects = append(expects, expectPush(
			session.NewAddedEvent(context.Background(), &session.NewAggregate(s.id, "instance1").Aggregate, uaChrome),
			session.NewTokenSetEvent(context.Background(), &session.NewAggregate(s.id, "instance1").Aggregate, "tokenID-"+s.id),
		))
	}
	// Session 3 is blocked before any Push — no extra expectation needed.
	c := &Commands{
		eventstore:    expectEventstore(expects...)(t),
		riskEvaluator: riskSvc,
	}
	ctx := authz.NewMockContext("instance1", "", "")

	for i, s := range passSessions {
		checks := makeSessionChecks(c, s.id, "fp1", "1.2.3.4", "chrome", base.Add(time.Duration(i)*time.Minute))
		checks.sessionWriteModel.UserID = "user1"
		checks.sessionWriteModel.UserResourceOwner = "org1"

		got, err := c.updateSession(ctx, checks, nil, 0)
		require.NoErrorf(t, err, "session %d should be allowed", i+1)
		require.NotNil(t, got)

		require.NoError(t, riskSvc.Record(ctx, checks.riskSignal(ctx, "", risk.OutcomeSuccess), nil))
	}

	// Session 3: different IP and user-agent → contextDrift must block.
	checks3 := makeSessionChecks(c, "session3", "fp2", "9.9.9.9", "safari", base.Add(2*time.Minute))
	checks3.sessionWriteModel.UserID = "user1"
	checks3.sessionWriteModel.UserResourceOwner = "org1"

	got3, err3 := c.updateSession(ctx, checks3, nil, 0)
	require.ErrorIs(t, err3, zerrors.ThrowPermissionDenied(nil, "COMMAND-RISK0", "Errors.PermissionDenied"),
		"session 3 should be blocked by contextDrift")
	assert.Nil(t, got3)
}
