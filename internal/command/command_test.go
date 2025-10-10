package command

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
	"github.com/zitadel/zitadel/internal/repository/user"
)

var (
	SupportedLanguages   = []language.Tag{language.English, language.German}
	OnlyAllowedLanguages = []language.Tag{language.English}
	AllowedLanguage      = language.English
	DisallowedLanguage   = language.German
	UnsupportedLanguage  = language.Spanish
)

func TestMain(m *testing.M) {
	i18n.SupportLanguages(SupportedLanguages...)
	os.Exit(m.Run())
}

func TestCommands_asyncPush(t *testing.T) {
	// make sure the test terminates on deadlock
	background := context.Background()
	agg := user.NewAggregate("userID", "orgID")
	cmd := user.NewMachineSecretHashUpdatedEvent(background, &agg.Aggregate, "updatedSecret")

	tests := []struct {
		name         string
		pushCtx      func() (context.Context, context.CancelFunc)
		eventstore   func(*testing.T) *eventstore.Eventstore
		closeCtx     func() (context.Context, context.CancelFunc)
		wantCloseErr bool
	}{
		{
			name: "push error",
			pushCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(background)
			},
			eventstore: expectEventstore(
				expectPushFailed(io.ErrClosedPipe, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second)
			},
			wantCloseErr: false,
		},
		{
			name: "success",
			pushCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(background)
			},
			eventstore: expectEventstore(
				expectPushSlow(time.Second/10, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second)
			},
			wantCloseErr: false,
		},
		{
			name: "success after push context cancels",
			pushCtx: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(background)
				cancel()
				return ctx, cancel
			},
			eventstore: expectEventstore(
				expectPushSlow(time.Second/10, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second)
			},
			wantCloseErr: false,
		},
		{
			name: "success after push context timeout",
			pushCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second/100)
			},
			eventstore: expectEventstore(
				expectPushSlow(time.Second/10, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second)
			},
			wantCloseErr: false,
		},
		{
			name: "success after push context timeout",
			pushCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second/100)
			},
			eventstore: expectEventstore(
				expectPushSlow(time.Second/10, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second)
			},
			wantCloseErr: false,
		},
		{
			name: "close timeout error",
			pushCtx: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(background)
			},
			eventstore: expectEventstore(
				expectPushSlow(time.Second/10, cmd),
			),
			closeCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(background, time.Second/100)
			},
			wantCloseErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			c.eventstore.PushTimeout = 10 * time.Second
			pushCtx, cancel := tt.pushCtx()
			c.asyncPush(pushCtx, cmd)
			cancel()

			closeCtx, cancel := tt.closeCtx()
			defer cancel()
			err := c.Close(closeCtx)
			if tt.wantCloseErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
