package execution

import (
	"time"

	"github.com/zitadel/zitadel/internal/execution/target"
)

type TargetType uint

const (
	TargetTypeWebhook TargetType = iota
	TargetTypeCall
	TargetTypeAsync
)

type Target interface {
	GetExecutionID() string
	GetTargetID() string
	IsInterruptOnError() bool
	GetEndpoint() string
	GetTargetType() target.TargetType
	GetTimeout() time.Duration
	GetSigningKey() string
}
