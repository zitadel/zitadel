package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type Aggregate struct {
	eventstore.Aggregate

	SetUpStarted Step
	SetUpDone    Step
}
