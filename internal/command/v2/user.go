package command

import (
	"context"
	"time"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

func ExistsUser(ctx context.Context, filter preparation.FilterToQueryReducer, id, resourceOwner string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(resourceOwner).
		OrderAsc().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(id).
		EventTypes(
			user.HumanRegisteredType,
			user.UserV1RegisteredType,
			user.HumanAddedType,
			user.UserV1AddedType,
			user.MachineAddedEventType,
			user.UserRemovedType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch event.(type) {
		case *user.HumanRegisteredEvent, *user.HumanAddedEvent, *user.MachineAddedEvent:
			exists = true
		case *user.UserRemovedEvent:
			exists = false
		}
	}

	return exists, nil
}

func newUserCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (value *crypto.CryptoValue, expiry time.Duration, err error) {
	return newCryptoCode(ctx, filter, domain.SecretGeneratorTypeInitCode, alg)
}
