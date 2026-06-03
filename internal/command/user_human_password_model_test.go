package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func TestHumanPasswordWriteModel_PreviousHashesAccumulate(t *testing.T) {
	agg := &user.NewAggregate("user1", "org1").Aggregate

	t.Run("three sequential events accumulate most-recent-first", func(t *testing.T) {
		wm := NewHumanPasswordWriteModel("user1", "org1")
		wm.Events = []eventstore.Event{
			user.NewHumanPasswordChangedEvent(context.Background(), agg, "h0", false, ""),
			user.NewHumanPasswordChangedEvent(context.Background(), agg, "h1", false, ""),
			user.NewHumanPasswordChangedEvent(context.Background(), agg, "h2", false, ""),
		}
		assert.NoError(t, wm.Reduce())
		assert.Equal(t, "h2", wm.EncodedHash)
		assert.Equal(t, []string{"h1", "h0"}, wm.PreviousHashes)
	})

	t.Run("legacy Secret-only event is skipped and does not add an empty entry", func(t *testing.T) {
		wm := NewHumanPasswordWriteModel("user1", "org1")

		legacyEvt := &user.HumanPasswordChangedEvent{
			BaseEvent: *eventstore.NewBaseEventForPush(context.Background(), agg, user.HumanPasswordChangedType),
			Secret: &crypto.CryptoValue{
				CryptoType: crypto.TypeEncryption,
				Algorithm:  "enc",
				KeyID:      "id",
				Crypted:    []byte("secret"),
			},
			EncodedHash:    "",
			ChangeRequired: false,
		}

		wm.Events = []eventstore.Event{
			user.NewHumanPasswordChangedEvent(context.Background(), agg, "h0", false, ""),
			legacyEvt,
			user.NewHumanPasswordChangedEvent(context.Background(), agg, "h2", false, ""),
		}
		assert.NoError(t, wm.Reduce())
		assert.Equal(t, "h2", wm.EncodedHash)
		for _, ph := range wm.PreviousHashes {
			assert.NotEmpty(t, ph, "PreviousHashes must not contain empty strings from legacy events")
		}
	})
}
