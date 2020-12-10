package email

import "github.com/caos/zitadel/internal/eventstore/v2"

type WriteModel struct {
	eventstore.WriteModel

	Email string
}
