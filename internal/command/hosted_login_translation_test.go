package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/service"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
	"golang.org/x/text/language"
)

func TestSetTranslationEvents(t *testing.T) {
	t.Parallel()

	testCtx := authz.SetCtxData(context.Background(), authz.CtxData{UserID: "test-user"})
	testCtx = service.WithService(testCtx, "test-service")

	tt := []struct {
		testName string

		inputAggregate    eventstore.Aggregate
		inputLanguage     language.Base
		inputTranslations map[string]any

		expectedCommands   []eventstore.Command
		expectedWriteModel *HostedLoginTranslationWriteModel
		expectedError      error
	}{
		{
			testName:          "when aggregate type is instance should return matching write model and instance.hosted_login_translation_set event",
			inputAggregate:    eventstore.Aggregate{ID: "123", Type: instance.AggregateType},
			inputLanguage:     language.MustParseBase("en"),
			inputTranslations: map[string]any{"test": "translation"},
			expectedCommands: []eventstore.Command{
				instance.NewHostedLoginTranslationSetEvent(testCtx, &eventstore.Aggregate{ID: "123", Type: instance.AggregateType}, map[string]any{"test": "translation"}, "en"),
			},
			expectedWriteModel: &HostedLoginTranslationWriteModel{
				WriteModel: eventstore.WriteModel{AggregateID: "123", ResourceOwner: "123"},
			},
		},
		{
			testName:          "when aggregate type is org should return matching write model and org.hosted_login_translation_set event",
			inputAggregate:    eventstore.Aggregate{ID: "123", Type: org.AggregateType},
			inputLanguage:     language.MustParseBase("en"),
			inputTranslations: map[string]any{"test": "translation"},
			expectedCommands: []eventstore.Command{
				org.NewHostedLoginTranslationSetEvent(testCtx, &eventstore.Aggregate{ID: "123", Type: org.AggregateType}, map[string]any{"test": "translation"}, "en"),
			},
			expectedWriteModel: &HostedLoginTranslationWriteModel{
				WriteModel: eventstore.WriteModel{AggregateID: "123", ResourceOwner: "123"},
			},
		},
		{
			testName:          "when aggregate type is neither org nor instance should return invalid argument error",
			inputAggregate:    eventstore.Aggregate{ID: "123"},
			inputLanguage:     language.MustParseBase("en"),
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
