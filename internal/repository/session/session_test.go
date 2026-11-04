package session

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

func TestTerminateEventMapper(t *testing.T) {
	event := &repository.Event{
		Seq:           15,
		Typ:           TerminateType,
		AggregateType: AggregateType,
		AggregateID:   "agg-id",
		InstanceID:    "instance-id",
		Version:       "v1",
		Data: []byte(`{
			"triggerOrigin": "https://login.example.com"
		}`),
	}

	got, err := TerminateEventMapper(event)
	require.NoError(t, err)

	terminate, ok := got.(*TerminateEvent)
	require.True(t, ok, "expected *TerminateEvent")
	assert.Equal(t, "https://login.example.com", terminate.TriggerOrigin())
	assert.Equal(t, "https://login.example.com", terminate.TriggeredAtOrigin)
}
