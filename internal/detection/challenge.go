package detection

import (
	"context"

	"github.com/zitadel/zitadel/internal/captcha"
)

// ChallengeVerifier verifies interactive challenge tokens (e.g. captcha).
type ChallengeVerifier interface {
	// VerifyCaptcha verifies a captcha token. Returns true when captcha is
	// not configured or verification succeeds.
	VerifyCaptcha(ctx context.Context, token string, remoteIP string) (bool, error)
	// CaptchaVerifier returns the configured verifier, or nil.
	CaptchaVerifier() captcha.CaptchaVerifier
}
