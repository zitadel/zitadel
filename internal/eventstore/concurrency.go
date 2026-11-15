package eventstore

import "github.com/zitadel/zitadel/internal/zerrors"

// ErrConcurrentModification requires reloading state before retrying the command.
var ErrConcurrentModification = zerrors.CreateZitadelError(zerrors.KindAborted, nil, "EVENTSTORE-concurrent", "Resource changed concurrently; retry the operation", 0)
