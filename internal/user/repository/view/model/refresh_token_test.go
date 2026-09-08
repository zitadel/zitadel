package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func mockRefreshTokenRemovedData(tokenID string) []byte {
	data, _ := json.Marshal(&user.HumanRefreshTokenRemovedEvent{TokenID: tokenID})
	return data
}

func TestAppendEventIfMyRefreshTokenAppliesRevocation(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	tests := []struct {
		name  string
		event eventstore.Event
	}{
		{
			name:  "refresh token removed expires token",
			event: &es_models.Event{AggregateID: "user1", Seq: 2, CreationDate: now, Typ: user.HumanRefreshTokenRemovedType, ResourceOwner: "org1", Data: mockRefreshTokenRemovedData("token1")},
		},
		{
			name:  "user deactivated expires token",
			event: &es_models.Event{AggregateID: "user1", Seq: 2, CreationDate: now, Typ: user.UserDeactivatedType, ResourceOwner: "org1"},
		},
		{
			name:  "user locked expires token",
			event: &es_models.Event{AggregateID: "user1", Seq: 2, CreationDate: now, Typ: user.UserLockedType, ResourceOwner: "org1"},
		},
		{
			name:  "user removed expires token",
			event: &es_models.Event{AggregateID: "user1", Seq: 2, CreationDate: now, Typ: user.UserRemovedType, ResourceOwner: "org1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := &RefreshTokenView{ID: "token1", UserID: "user1", ResourceOwner: "org1", Expiration: future}
			if err := view.AppendEventIfMyRefreshToken(tt.event); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !view.Expiration.Equal(now) {
				t.Errorf("expected token to be expired at %v, got %v", now, view.Expiration)
			}
		})
	}
}

func TestAppendEventIfMyRefreshTokenIgnoresOtherToken(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)
	view := &RefreshTokenView{ID: "token1", UserID: "user1", ResourceOwner: "org1", Expiration: future}
	event := &es_models.Event{AggregateID: "user1", Seq: 2, CreationDate: now, Typ: user.HumanRefreshTokenRemovedType, ResourceOwner: "org1", Data: mockRefreshTokenRemovedData("token2")}
	if err := view.AppendEventIfMyRefreshToken(event); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !view.Expiration.Equal(future) {
		t.Errorf("expected token expiration to be unchanged at %v, got %v", future, view.Expiration)
	}
}
