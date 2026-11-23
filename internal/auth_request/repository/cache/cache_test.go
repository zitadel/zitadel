package cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/domain"
)

func TestAuthRequestCache_returnsCopies(t *testing.T) {
	ctx := context.Background()
	c := Start(nil, 10)
	stored := &domain.AuthRequest{ID: "id", Code: "code", UserID: "userA"}
	c.CacheAuthRequest(ctx, stored)
	stored.UserID = "mutated after caching"

	byID, ok := c.getCachedByID(ctx, "id")
	require.True(t, ok)
	byID.UserID = "userB" // e.g. a concurrent login-name switch

	byCode, ok := c.getCachedByCode(ctx, "code")
	require.True(t, ok)
	assert.Equal(t, "userA", byCode.UserID)
	again, _ := c.getCachedByID(ctx, "id")
	assert.Equal(t, "userA", again.UserID)
}
