package command

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/v2"
)

func TestSetTranslationEvents(t *testing.T) {
	t.Parallel()

	testCtx := authz.SetCtxData(context.Background(), authz.CtxData{UserID: "test-user"})
	testCtx = service.WithService(testCtx, "test-service")

	tt := []struct {
		testName string

		inputAggregate    eventstore.Aggregate
		inputLanguage     language.Tag
		inputTranslations map[string]any

		expectedCommands   []eventstore.Command
		expectedWriteModel *HostedLoginTranslationWriteModel
		expectedError      error
	}{
		{
			testName:          "when aggregate type is instance should return matching write model and instance.hosted_login_translation_set event",
			inputAggregate:    eventstore.Aggregate{ID: "123", Type: instance.AggregateType},
			inputLanguage:     language.MustParse("en-US"),
			inputTranslations: map[string]any{"test": "translation"},
			expectedCommands: []eventstore.Command{
				instance.NewHostedLoginTranslationSetEvent(testCtx, &eventstore.Aggregate{ID: "123", Type: instance.AggregateType}, map[string]any{"test": "translation"}, language.MustParse("en-US")),
			},
			expectedWriteModel: &HostedLoginTranslationWriteModel{
				WriteModel: eventstore.WriteModel{AggregateID: "123", ResourceOwner: "123"},
			},
		},
		{
			testName:          "when aggregate type is org should return matching write model and org.hosted_login_translation_set event",
			inputAggregate:    eventstore.Aggregate{ID: "123", Type: org.AggregateType},
			inputLanguage:     language.MustParse("en-GB"),
			inputTranslations: map[string]any{"test": "translation"},
			expectedCommands: []eventstore.Command{
				org.NewHostedLoginTranslationSetEvent(testCtx, &eventstore.Aggregate{ID: "123", Type: org.AggregateType}, map[string]any{"test": "translation"}, language.MustParse("en-GB")),
			},
			expectedWriteModel: &HostedLoginTranslationWriteModel{
				WriteModel: eventstore.WriteModel{AggregateID: "123", ResourceOwner: "123"},
			},
		},
		{
			testName:          "when aggregate type is neither org nor instance should return invalid argument error",
			inputAggregate:    eventstore.Aggregate{ID: "123"},
			inputLanguage:     language.MustParse("en-US"),
			inputTranslations: map[string]any{"test": "translation"},
			expectedError:     zerrors.ThrowInvalidArgument(nil, "COMMA-0aw7In", "Errors.Arguments.LevelType.Invalid"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := Commands{}

			// When
			events, writeModel, err := c.setTranslationEvents(testCtx, tc.inputAggregate, tc.inputLanguage, tc.inputTranslations)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedWriteModel, writeModel)

			require.Len(t, events, len(tc.expectedCommands))
			assert.ElementsMatch(t, tc.expectedCommands, events)
		})
	}
}

func TestSetHostedLoginTranslation(t *testing.T) {
	t.Parallel()

	testCtx := authz.SetCtxData(context.Background(), authz.CtxData{UserID: "test-user"})
	testCtx = service.WithService(testCtx, "test-service")
	testCtx = authz.WithInstanceID(testCtx, "instance-id")

	testTranslation := map[string]any{"test": "translation", "translation": "2"}
	protoTranslation, err := structpb.NewStruct(testTranslation)
	require.NoError(t, err)

	hashTestTranslation := md5.Sum(fmt.Append(nil, testTranslation))
	require.NotEmpty(t, hashTestTranslation)

	tt := []struct {
		testName string

		mockPush func(*testing.T) *eventstore.Eventstore

		inputReq *settings.SetHostedLoginTranslationRequest

		expectedError  error
		expectedResult *settings.SetHostedLoginTranslationResponse
	}{
		{
			testName: "when locale is malformed should return invalid argument error",
			mockPush: func(t *testing.T) *eventstore.Eventstore { return &eventstore.Eventstore{} },
			inputReq: &settings.SetHostedLoginTranslationRequest{
				Level:  &settings.SetHostedLoginTranslationRequest_Instance{},
				Locale: "123",
			},

			expectedError: zerrors.ThrowInvalidArgument(nil, "COMMA-xmjATA", "Errors.Arguments.Locale.Invalid"),
		},
		{
			testName: "when locale is unknown should return invalid argument error",
			mockPush: func(t *testing.T) *eventstore.Eventstore { return &eventstore.Eventstore{} },
			inputReq: &settings.SetHostedLoginTranslationRequest{
				Level:  &settings.SetHostedLoginTranslationRequest_Instance{},
				Locale: "root",
			},

			expectedError: zerrors.ThrowInvalidArgument(nil, "COMMA-xmjATA", "Errors.Arguments.Locale.Invalid"),
		},
		{
			testName: "when event pushing fails should return internal error",

			mockPush: expectEventstore(expectPushFailed(
				errors.New("mock push failed"),
				instance.NewHostedLoginTranslationSetEvent(
					testCtx, &eventstore.Aggregate{
						ID:            "instance-id",
						Type:          instance.AggregateType,
						ResourceOwner: "instance-id",
						InstanceID:    "instance-id",
						Version:       instance.AggregateVersion,
					},
					testTranslation,
					language.MustParse("it-CH"),
				),
			)),

			inputReq: &settings.SetHostedLoginTranslationRequest{
				Level:        &settings.SetHostedLoginTranslationRequest_Instance{},
				Locale:       "it-CH",
				Translations: protoTranslation,
			},

			expectedError: zerrors.ThrowInternal(errors.New("mock push failed"), "COMMA-i8nqFl", "Errors.Internal"),
		},
		{
			testName: "when request is valid should return expected response",

			mockPush: expectEventstore(expectPush(
				org.NewHostedLoginTranslationSetEvent(
					testCtx, &eventstore.Aggregate{
						ID:            "org-id",
						Type:          org.AggregateType,
						ResourceOwner: "org-id",
						InstanceID:    "",
						Version:       org.AggregateVersion,
					},
					testTranslation,
					language.MustParse("it-CH"),
				),
			)),

			inputReq: &settings.SetHostedLoginTranslationRequest{
				Level:        &settings.SetHostedLoginTranslationRequest_OrganizationId{OrganizationId: "org-id"},
				Locale:       "it-CH",
				Translations: protoTranslation,
			},

			expectedResult: &settings.SetHostedLoginTranslationResponse{
				Etag: hex.EncodeToString(hashTestTranslation[:]),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			// Given
			c := Commands{
				eventstore: tc.mockPush(t),
			}

			// When
			res, err := c.SetHostedLoginTranslation(testCtx, tc.inputReq)

			// Verify
			require.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}
