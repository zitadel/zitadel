package domain_test

import (
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel/backend/v3/domain"
	domainmock "github.com/zitadel/zitadel/backend/v3/domain/mock"
	noopdb "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/noop"
	"github.com/zitadel/zitadel/internal/api/authz"
	old_domain "github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id"
	id_gen_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/session"
	"github.com/zitadel/zitadel/internal/zerrors"
	session_grpc "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestCreateSessionCommand_Validate(t *testing.T) {
	ctx := authz.NewMockContext("instance-ctx", "", "")
	oldIDConfig := id.GeneratorConfig
	t.Cleanup(func() {
		id.GeneratorConfig = oldIDConfig
	})

	id.Configure(&id.Config{Identification: id.Identification{PrivateIp: id.PrivateIp{Enabled: true}}})
	t.Parallel()

	tt := []struct {
		testName        string
		inputCtx        context.Context
		inputInstanceID string
		inputLifetime   *durationpb.Duration

		expectedInstanceID string
		expectedError      error
	}{
		{
			testName:           "when input instance id is not set should set from context",
			inputCtx:           ctx,
			expectedInstanceID: "instance-ctx",
		},
		{
			testName:           "when input instance id is set should set from input",
			inputInstanceID:    "instance-1",
			expectedInstanceID: "instance-1",
		},
		{
			testName:           "when input lifetime is set as negative value should return invalid argument error",
			inputInstanceID:    "instance-1",
			inputLifetime:      durationpb.New(-1 * time.Second),
			expectedInstanceID: "instance-1",
			expectedError:      zerrors.ThrowInvalidArgument(nil, "DOM-XA5OMq", "Errors.Session.PositiveLifetime"),
		},
		{
			testName:           "when input lifetime is set as positive value should return nil",
			inputInstanceID:    "instance-1",
			inputLifetime:      durationpb.New(1 * time.Second),
			expectedInstanceID: "instance-1",
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			sessionCreator := domain.NewCreateSessionCommand(tc.inputInstanceID, nil, nil, tc.inputLifetime, nil)

			// Test
			err := sessionCreator.Validate(tc.inputCtx, nil)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedInstanceID, sessionCreator.InstanceID)
		})
	}
}

func TestCreateSessionCommand_Execute(t *testing.T) {
	oldIDConfig := id.GeneratorConfig
	t.Cleanup(func() {
		id.GeneratorConfig = oldIDConfig
	})

	id.Configure(&id.Config{Identification: id.Identification{PrivateIp: id.PrivateIp{Enabled: true}}})
	t.Parallel()

	idGenErr := errors.New("mock id gen")
	createErr := errors.New("mock create error")

	tt := []struct {
		testName        string
		idGenMock       func(ctrl *gomock.Controller) id.Generator
		sessionRepoMock func(ctr *gomock.Controller) domain.SessionRepository

		inputUserAgent *session_grpc.UserAgent
		inputMetas     map[string][]byte
		inputLifetime  *durationpb.Duration

		expectedError     error
		expectedSessionID *string
	}{
		{
			testName: "when id generation fails should return internal error",
			idGenMock: func(ctrl *gomock.Controller) id.Generator {
				mock := id_gen_mock.NewMockGenerator(ctrl)

				mock.EXPECT().Next().Times(1).Return("", idGenErr)
				return mock
			},
			expectedError: zerrors.ThrowInternal(idGenErr, "DOM-ngXOIK", "failed generating session ID"),
		},
		{
			testName: "when session creation fails should return internal error",
			idGenMock: func(ctrl *gomock.Controller) id.Generator {
				mock := id_gen_mock.NewMockGenerator(ctrl)

				mock.EXPECT().Next().Times(1).Return("session-1", nil)
				return mock
			},
			sessionRepoMock: func(ctr *gomock.Controller) domain.SessionRepository {
				mock := domainmock.NewSessionRepo(ctr)

				session := &domain.Session{
					InstanceID: "instance-1",
					ID:         "session-1",
					UserAgent:  nil,
				}
				mock.EXPECT().Create(gomock.Any(), gomock.Any(), session).Times(1).Return(createErr)

				return mock
			},
			expectedError: zerrors.ThrowInternal(createErr, "DOM-HYKAgF", "failed creating session"),
		},
		{
			testName: "when session creation succeeds should set session ID and return no error",
			idGenMock: func(ctrl *gomock.Controller) id.Generator {
				mock := id_gen_mock.NewMockGenerator(ctrl)

				mock.EXPECT().Next().Times(1).Return("session-1", nil)
				return mock
			},
			sessionRepoMock: func(ctr *gomock.Controller) domain.SessionRepository {
				mock := domainmock.NewSessionRepo(ctr)

				session := &domain.Session{
					InstanceID: "instance-1",
					ID:         "session-1",
					UserAgent: &domain.SessionUserAgent{
						InstanceID:    "instance-1",
						FingerprintID: gu.Ptr("123"),
						Description:   gu.Ptr("some description"),
						IP:            net.ParseIP("127.0.0.1"),
						Header: http.Header{
							"h1": []string{"v1.1", "v1.2"},
							"h2": []string{"v2.1", "v2.2"},
						},
					},
					Metadata: []domain.SessionMetadata{
						{
							Metadata:  domain.Metadata{InstanceID: "instance-1", Key: "meta-1", Value: []byte("value1")},
							SessionID: "session-1",
						},
						{
							Metadata:  domain.Metadata{InstanceID: "instance-1", Key: "meta-2", Value: []byte("value2")},
							SessionID: "session-1",
						},
					},
					Lifetime: 1 * time.Second,
				}
				mock.EXPECT().Create(gomock.Any(), gomock.Any(), session).Times(1).Return(nil)

				return mock
			},
			inputUserAgent: &session_grpc.UserAgent{
				FingerprintId: gu.Ptr("123"),
				Ip:            gu.Ptr("127.0.0.1"),
				Description:   gu.Ptr("some description"),
				Header: map[string]*session_grpc.UserAgent_HeaderValues{
					"h1": {Values: []string{"v1.1", "v1.2"}},
					"h2": {Values: []string{"v2.1", "v2.2"}},
				},
			},
			inputMetas: map[string][]byte{
				"meta-1": []byte("value1"),
				"meta-2": []byte("value2"),
			},
			inputLifetime:     durationpb.New(1 * time.Second),
			expectedSessionID: gu.Ptr("session-1"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			ctrl := gomock.NewController(t)

			var idGenMock id.Generator
			if tc.idGenMock != nil {
				idGenMock = tc.idGenMock(ctrl)
			}
			cmd := domain.NewCreateSessionCommand("instance-1", tc.inputUserAgent, tc.inputMetas, tc.inputLifetime, idGenMock)

			opts := &domain.InvokeOpts{
				Invoker: domain.NewTransactionInvoker(nil),
			}
			domain.WithQueryExecutor(new(noopdb.Pool))(opts)

			if tc.sessionRepoMock != nil {
				domain.WithSessionRepo(tc.sessionRepoMock(ctrl))(opts)
			}

			// Test
			err := cmd.Execute(t.Context(), opts)

			// Verify
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedSessionID, cmd.SessionID)
		})
	}
}

func TestCreateSessionCommand_Events(t *testing.T) {
	oldIDConfig := id.GeneratorConfig
	t.Cleanup(func() {
		id.GeneratorConfig = oldIDConfig
	})

	id.Configure(&id.Config{Identification: id.Identification{PrivateIp: id.PrivateIp{Enabled: true}}})
	t.Parallel()

	sessionAgg := &session.NewAggregate("session-1", "instance-1").Aggregate
	tt := []struct {
		testName string

		inputUserAgent *session_grpc.UserAgent
		inputMetas     map[string][]byte
		inputLifetime  *durationpb.Duration

		expectedEvents []eventstore.Command
	}{
		{
			testName: "when all params are nil should return session added event with nil user agent",

			expectedEvents: []eventstore.Command{
				session.NewAddedEvent(t.Context(), sessionAgg, nil),
			},
		},
		{
			testName: "when user agent is set should return session added event with user agent",
			inputUserAgent: &session_grpc.UserAgent{
				FingerprintId: gu.Ptr("fingerprint-id"),
				Ip:            gu.Ptr("127.0.0.1"),
				Description:   gu.Ptr("description"),
			},
			expectedEvents: []eventstore.Command{
				session.NewAddedEvent(t.Context(), sessionAgg, &old_domain.UserAgent{
					FingerprintID: gu.Ptr("fingerprint-id"),
					IP:            net.ParseIP("127.0.0.1"),
					Description:   gu.Ptr("description"),
				}),
			},
		},
		{
			testName: "when metas set should return session added and metadata set events",
			inputUserAgent: &session_grpc.UserAgent{
				FingerprintId: gu.Ptr("fingerprint-id"),
				Ip:            gu.Ptr("127.0.0.1"),
				Description:   gu.Ptr("description"),
			},
			inputMetas: map[string][]byte{"meta1": []byte("value1")},
			expectedEvents: []eventstore.Command{
				session.NewAddedEvent(t.Context(), sessionAgg, &old_domain.UserAgent{
					FingerprintID: gu.Ptr("fingerprint-id"),
					IP:            net.ParseIP("127.0.0.1"),
					Description:   gu.Ptr("description"),
				}),
				session.NewMetadataSetEvent(t.Context(), sessionAgg, map[string][]byte{"meta1": []byte("value1")}),
			},
		},
		{
			testName: "when lifetime is set should return session added, metadata set and lifetime set events",
			inputUserAgent: &session_grpc.UserAgent{
				FingerprintId: gu.Ptr("fingerprint-id"),
				Ip:            gu.Ptr("127.0.0.1"),
				Description:   gu.Ptr("description"),
			},
			inputMetas:    map[string][]byte{"meta1": []byte("value1")},
			inputLifetime: durationpb.New(1 * time.Second),
			expectedEvents: []eventstore.Command{
				session.NewAddedEvent(t.Context(), sessionAgg, &old_domain.UserAgent{
					FingerprintID: gu.Ptr("fingerprint-id"),
					IP:            net.ParseIP("127.0.0.1"),
					Description:   gu.Ptr("description"),
				}),
				session.NewMetadataSetEvent(t.Context(), sessionAgg, map[string][]byte{"meta1": []byte("value1")}),
				session.NewLifetimeSetEvent(t.Context(), sessionAgg, 1*time.Second),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			// Given
			cmd := domain.NewCreateSessionCommand("instance-1", tc.inputUserAgent, tc.inputMetas, tc.inputLifetime, nil)
			cmd.SessionID = gu.Ptr("session-1")

			// Test
			events, err := cmd.Events(t.Context(), nil)

			// Verify
			assert.NoError(t, err)
			require.Len(t, events, len(tc.expectedEvents))
			for i, expectedType := range tc.expectedEvents {
				assert.IsType(t, expectedType, events[i])
				switch expectedAssertedType := expectedType.(type) {
				case *session.AddedEvent:
					actualAssertedType, ok := events[i].(*session.AddedEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.UserAgent, actualAssertedType.UserAgent)
				case *session.MetadataSetEvent:
					actualAssertedType, ok := events[i].(*session.MetadataSetEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.Metadata, actualAssertedType.Metadata)
				case *session.LifetimeSetEvent:
					actualAssertedType, ok := events[i].(*session.LifetimeSetEvent)
					require.True(t, ok)
					assert.Equal(t, expectedAssertedType.Lifetime, actualAssertedType.Lifetime)
				}
			}
		})
	}
}
