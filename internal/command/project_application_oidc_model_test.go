package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestOIDCApplicationWriteModel_NewChangedEvent_AppLinks(t *testing.T) {
	t.Parallel()

	agg := &project.NewAggregate("project-id", "org-id").Aggregate
	base := func() *OIDCApplicationWriteModel {
		return &OIDCApplicationWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   "project-id",
				ResourceOwner: "org-id",
			},
			AppID:                         "app-id",
			IOSTeamID:                     "OLDTEAM",
			IOSBundleID:                   "com.old.app",
			AndroidPackageName:            "com.old.app",
			AndroidSHA256CertFingerprints: []string{"AA:AA"},
		}
	}

	t.Run("omit keeps existing", func(t *testing.T) {
		t.Parallel()
		wm := base()
		event, hasChanged, err := wm.NewChangedEvent(
			context.Background(), agg, "app-id",
			nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil,
		)
		require.NoError(t, err)
		assert.False(t, hasChanged)
		assert.Nil(t, event)
	})

	t.Run("set changes fields", func(t *testing.T) {
		t.Parallel()
		wm := base()
		event, hasChanged, err := wm.NewChangedEvent(
			context.Background(), agg, "app-id",
			nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			gu.Ptr("NEWTEAM"),
			gu.Ptr("com.new.app"),
			gu.Ptr("com.new.app"),
			[]string{"BB:BB"},
		)
		require.NoError(t, err)
		require.True(t, hasChanged)
		require.NotNil(t, event)
		assert.Equal(t, gu.Ptr("NEWTEAM"), event.IOSTeamID)
		assert.Equal(t, gu.Ptr("com.new.app"), event.IOSBundleID)
		assert.Equal(t, gu.Ptr("com.new.app"), event.AndroidPackageName)
		assert.Equal(t, &[]string{"BB:BB"}, event.AndroidSHA256CertFingerprints)
	})

	t.Run("clear with empty values", func(t *testing.T) {
		t.Parallel()
		wm := base()
		event, hasChanged, err := wm.NewChangedEvent(
			context.Background(), agg, "app-id",
			nil, nil, nil, nil, nil, nil, nil, nil,
			nil, nil, nil, nil, nil, nil, nil, nil, nil, nil,
			gu.Ptr(""),
			gu.Ptr(""),
			gu.Ptr(""),
			[]string{},
		)
		require.NoError(t, err)
		require.True(t, hasChanged)
		require.NotNil(t, event)
		assert.Equal(t, gu.Ptr(""), event.IOSTeamID)
		assert.Equal(t, gu.Ptr(""), event.IOSBundleID)
		assert.Equal(t, gu.Ptr(""), event.AndroidPackageName)
		assert.Equal(t, &[]string{}, event.AndroidSHA256CertFingerprints)
	})
}
