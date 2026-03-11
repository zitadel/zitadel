package ratelimit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

const rateLimitStorageKeyVersion = "v1"

// CanonicalRateLimitKey builds the storage key used by all rate limiter backends.
// The operator-controlled template output is treated as a suffix selector, while
// ZITADEL always prefixes the key with instance, rule, and window information to
// prevent cross-tenant and cross-rule collisions.
func CanonicalRateLimitKey(ruleID, instanceID string, window time.Duration, rendered string) string {
	ruleID = normalizeRateLimitKeyPart(ruleID, "rule")
	instanceID = normalizeRateLimitKeyPart(instanceID, "instance")

	sum := sha256.Sum256([]byte(strings.TrimSpace(rendered)))
	return fmt.Sprintf(
		"rl:%s:%s:%s:%d:%s",
		rateLimitStorageKeyVersion,
		instanceID,
		ruleID,
		int(window.Seconds()),
		hex.EncodeToString(sum[:12]),
	)
}

func normalizeRateLimitKeyPart(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
