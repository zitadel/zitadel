package activity

import (
	"context"

	"github.com/zitadel/logging"
)

type storageInfo struct {
	ResourceOwner string
	UserID        string
	Trigger       TriggerMethod
}

type storageInfoKey struct{}

// CreateStorageInfoContext creates a new channel for StorageInfo and returns
// a new context which can be used with [SetStorageInfo].
func CreateStorageInfoContext(ctx context.Context) context.Context {
	c := make(chan storageInfo, 1)
	return context.WithValue(ctx, storageInfoKey{}, c)
}

// SetStorageInfo sends the info back on a channel in the context.
// It may only be called once during a request cycle.
// Subsequent calls may drop the passed info if the channel is already filled an entry.
func SetStorageInfo(ctx context.Context, resourceOwner, userID string, trigger TriggerMethod) {
	if c, ok := ctx.Value(storageInfoKey{}).(chan storageInfo); ok {
		info := storageInfo{
			ResourceOwner: resourceOwner,
			UserID:        userID,
			Trigger:       trigger,
		}
		select {
		case c <- info:
		default:
			logging.New().WithField("info", info).Debug("storage info channel blocked, dropped info")
		}
	}
	logging.Debug("no storage info channel in context")
}

// getStorageInfo receives info from a channel in the context.
// If there is no channel in the context or the channel does not
// contain an entry, an empty [StorageInfo] is returned.
func getStorageInfo(ctx context.Context) (info storageInfo) {
	if c, ok := ctx.Value(storageInfoKey{}).(chan storageInfo); ok {
		select {
		case info = <-c:
		default:
			logging.New().WithField("info", info).Debug("storage info channel blocked, dropped info")
		}
		return info
	}
	logging.Debug("no storage info channel in context")
	return info
}
