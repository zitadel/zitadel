//go:build integration

package session_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client session.SessionServiceClient
	User   *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.SessionV2

		CTX, _ = Tester.WithSystemAuthorization(ctx, integration.OrgOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		Tester.RegisterUserPasskey(CTX, User.GetUserId())
		return m.Run()
	}())
}

func verifyCurrentSession(t testing.TB, id, token string, sequence uint64, window time.Duration, metadata map[string][]byte, factors ...wantFactor) (s *session.Session) {
	require.NotEmpty(t, id)
	require.NotEmpty(t, token)

retry:
	for {
		resp, err := Client.GetSession(CTX, &session.GetSessionRequest{
			SessionId:    id,
			SessionToken: &token,
		})
		if err == nil {
			s = resp.GetSession()
			break retry
		}
		if status.Convert(err).Code() == codes.NotFound {
			select {
			case <-CTX.Done():
				t.Fatal(CTX.Err(), err)
			case <-time.After(time.Second):
				t.Log("retrying GetSession")
				continue
			}
		}
		require.NoError(t, err)
	}

	assert.Equal(t, id, s.GetId())
	assert.WithinRange(t, s.GetCreationDate().AsTime(), time.Now().Add(-window), time.Now().Add(window))
	assert.WithinRange(t, s.GetChangeDate().AsTime(), time.Now().Add(-window), time.Now().Add(window))
	assert.Equal(t, sequence, s.GetSequence())
	assert.Equal(t, metadata, s.GetMetadata())
	verifyFactors(t, s.GetFactors(), window, factors)
	return s
}

type wantFactor int

const (
	wantUserFactor wantFactor = iota
	wantPasswordFactor
	wantPasskeyFactor
)

func verifyFactors(t testing.TB, factors *session.Factors, window time.Duration, want []wantFactor) {
	for _, w := range want {
		switch w {
		case wantUserFactor:
			uf := factors.GetUser()
			assert.NotNil(t, uf)
			assert.WithinRange(t, uf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
			assert.Equal(t, User.GetUserId(), uf.GetId())
		case wantPasswordFactor:
			pf := factors.GetPassword()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		case wantPasskeyFactor:
			pf := factors.GetPasskey()
			assert.NotNil(t, pf)
			assert.WithinRange(t, pf.GetVerifiedAt().AsTime(), time.Now().Add(-window), time.Now().Add(window))
		}
	}
}

func TestServer_CreateSession(t *testing.T) {
	tests := []struct {
		name        string
		req         *session.CreateSessionRequest
		want        *session.CreateSessionResponse
		wantErr     bool
		wantFactors []wantFactor
	}{
		{
			name: "empty session",
			req: &session.CreateSessionRequest{
				Metadata: map[string][]byte{"foo": []byte("bar")},
			},
			want: &session.CreateSessionResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "with user",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					User: &session.CheckUser{
						Search: &session.CheckUser_UserId{
							UserId: User.GetUserId(),
						},
					},
				},
				Metadata: map[string][]byte{"foo": []byte("bar")},
			},
			want: &session.CreateSessionResponse{
				Details: &object.Details{
					ResourceOwner: Tester.Organisation.ID,
				},
			},
			wantFactors: []wantFactor{wantUserFactor},
		},
		{
			name: "password without user error",
			req: &session.CreateSessionRequest{
				Checks: &session.Checks{
					Password: &session.CheckPassword{
						Password: "Difficult",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "passkey without user error",
			req: &session.CreateSessionRequest{
				Challenges: []session.ChallengeKind{
					session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateSession(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)

			verifyCurrentSession(t, got.GetSessionId(), got.GetSessionToken(), got.GetDetails().GetSequence(), time.Minute, tt.req.GetMetadata(), tt.wantFactors...)
		})
	}
}

func TestServer_CreateSession_passkey(t *testing.T) {
	// create new session with user and request the passkey challenge
	createResp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: User.GetUserId(),
				},
			},
		},
		Challenges: []session.ChallengeKind{
			session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil)

	assertionData, err := Tester.WebAuthN.CreateAssertionResponse(createResp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
	require.NoError(t, err)

	// update the session with passkey assertion data
	updateResp, err := Client.SetSession(CTX, &session.SetSessionRequest{
		SessionId:    createResp.GetSessionId(),
		SessionToken: createResp.GetSessionToken(),
		Checks: &session.Checks{
			Passkey: &session.CheckPasskey{
				CredentialAssertionData: assertionData,
			},
		},
	})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), updateResp.GetSessionToken(), updateResp.GetDetails().GetSequence(), time.Minute, nil, wantUserFactor, wantPasskeyFactor)
}

func TestServer_SetSession_flow(t *testing.T) {
	var wantFactors []wantFactor

	// create new, empty session
	createResp, err := Client.CreateSession(CTX, &session.CreateSessionRequest{})
	require.NoError(t, err)
	verifyCurrentSession(t, createResp.GetSessionId(), createResp.GetSessionToken(), createResp.GetDetails().GetSequence(), time.Minute, nil, wantFactors...)
	sessionToken := createResp.GetSessionToken()

	t.Run("check user", func(t *testing.T) {
		wantFactors = append(wantFactors, wantUserFactor)
		resp, err := Client.SetSession(CTX, &session.SetSessionRequest{
			SessionId:    createResp.GetSessionId(),
			SessionToken: sessionToken,
			Checks: &session.Checks{
				User: &session.CheckUser{
					Search: &session.CheckUser_UserId{
						UserId: User.GetUserId(),
					},
				},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, wantFactors...)
		sessionToken = resp.GetSessionToken()
	})

	t.Run("check passkey", func(t *testing.T) {
		resp, err := Client.SetSession(CTX, &session.SetSessionRequest{
			SessionId:    createResp.GetSessionId(),
			SessionToken: sessionToken,
			Challenges: []session.ChallengeKind{
				session.ChallengeKind_CHALLENGE_KIND_PASSKEY,
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, wantFactors...)
		sessionToken = resp.GetSessionToken()

		wantFactors = append(wantFactors, wantPasskeyFactor)
		assertionData, err := Tester.WebAuthN.CreateAssertionResponse(resp.GetChallenges().GetPasskey().GetPublicKeyCredentialRequestOptions())
		require.NoError(t, err)

		resp, err = Client.SetSession(CTX, &session.SetSessionRequest{
			SessionId:    createResp.GetSessionId(),
			SessionToken: sessionToken,
			Checks: &session.Checks{
				Passkey: &session.CheckPasskey{
					CredentialAssertionData: assertionData,
				},
			},
		})
		require.NoError(t, err)
		verifyCurrentSession(t, createResp.GetSessionId(), resp.GetSessionToken(), resp.GetDetails().GetSequence(), time.Minute, nil, wantFactors...)
	})
}
