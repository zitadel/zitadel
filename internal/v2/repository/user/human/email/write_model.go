package email

import "github.com/caos/zitadel/internal/eventstore/v2"

type HumanEmailWriteModel struct {
	eventstore.WriteModel

	Email string
}
