package senders

import (
	"context"
	"github.com/zitadel/zitadel/internal/notification/handlers"

	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/notification/channels/log"
)

func connectToDebugChannels[T channels.Message](ctx context.Context, queries *handlers.NotificationQueries) []channels.NotificationChannel[T] {
	var (
		providers []channels.NotificationChannel[T]
	)
	if fsProvider, err := queries.GetFileSystemProvider(ctx); err == nil {
		p, err := fs.Connect[T](*fsProvider)
		if err == nil {
			providers = append(providers, p)
		}
	}
	if logProvider, err := queries.GetLogProvider(ctx); err == nil {
		providers = append(providers, log.Connect[T](*logProvider))
	}
	return providers
}
