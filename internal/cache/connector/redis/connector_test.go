package redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionsFromURL(t *testing.T) {
	t.Run("basic URL", func(t *testing.T) {
		c := Config{URL: "redis://localhost:6379/0"}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, "localhost:6379", opts.Addr)
		assert.Equal(t, 0, opts.DB)
		assert.Equal(t, 3, opts.Protocol)
		assert.True(t, opts.ContextTimeoutEnabled)
	})

	t.Run("with auth", func(t *testing.T) {
		c := Config{URL: "redis://myuser:mypass@host:6380/2"}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, "host:6380", opts.Addr)
		assert.Equal(t, "myuser", opts.Username)
		assert.Equal(t, "mypass", opts.Password)
		assert.Equal(t, 2, opts.DB)
	})

	t.Run("TLS via rediss scheme", func(t *testing.T) {
		c := Config{URL: "rediss://host:6380/0"}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.NotNil(t, opts.TLSConfig)
	})

	t.Run("password only auth", func(t *testing.T) {
		c := Config{URL: "redis://:secret@host:6379/0"}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, "", opts.Username)
		assert.Equal(t, "secret", opts.Password)
	})

	t.Run("URL without DB number", func(t *testing.T) {
		c := Config{URL: "redis://localhost:6379"}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, "localhost:6379", opts.Addr)
		assert.Equal(t, 0, opts.DB)
	})

	t.Run("overlays pool settings", func(t *testing.T) {
		c := Config{
			URL:            "redis://localhost:6379/0",
			PoolSize:       42,
			MaxActiveConns: 100,
			MinIdleConns:   5,
		}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, 42, opts.PoolSize)
		assert.Equal(t, 100, opts.MaxActiveConns)
		assert.Equal(t, 5, opts.MinIdleConns)
	})

	t.Run("overlays timeout settings", func(t *testing.T) {
		c := Config{
			URL:         "redis://localhost:6379/0",
			DialTimeout: 10 * time.Second,
			ReadTimeout: 5 * time.Second,
		}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, 10*time.Second, opts.DialTimeout)
		assert.Equal(t, 5*time.Second, opts.ReadTimeout)
	})

	t.Run("overlays retry settings", func(t *testing.T) {
		c := Config{
			URL:        "redis://localhost:6379/0",
			MaxRetries: 5,
		}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, 5, opts.MaxRetries)
	})

	t.Run("overlays client name", func(t *testing.T) {
		c := Config{
			URL:        "redis://localhost:6379/0",
			ClientName: "zitadel",
		}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.Equal(t, "zitadel", opts.ClientName)
	})

	t.Run("overlays circuit breaker", func(t *testing.T) {
		c := Config{
			URL: "redis://localhost:6379/0",
			CircuitBreaker: &CBConfig{
				MaxConsecutiveFailures: 10,
			},
			MaxActiveConns: 50,
		}
		opts, err := optionsFromURL(c)
		require.NoError(t, err)
		assert.NotNil(t, opts.Limiter)
	})

	t.Run("invalid URL returns error", func(t *testing.T) {
		c := Config{URL: "garbage://???"}
		_, err := optionsFromURL(c)
		require.Error(t, err)
	})
}

func TestOptionsFromConfig(t *testing.T) {
	t.Run("maps all fields", func(t *testing.T) {
		c := Config{
			Addr:       "myhost:6380",
			Username:   "user",
			Password:   "pass",
			PoolSize:   42,
			MaxRetries: 5,
		}
		opts := optionsFromConfig(c)
		assert.Equal(t, "myhost:6380", opts.Addr)
		assert.Equal(t, "user", opts.Username)
		assert.Equal(t, "pass", opts.Password)
		assert.Equal(t, 42, opts.PoolSize)
		assert.Equal(t, 5, opts.MaxRetries)
		assert.Equal(t, 3, opts.Protocol)
		assert.True(t, opts.ContextTimeoutEnabled)
	})

	t.Run("EnableTLS sets TLSConfig", func(t *testing.T) {
		c := Config{Addr: "host:6379", EnableTLS: true}
		opts := optionsFromConfig(c)
		assert.NotNil(t, opts.TLSConfig)
	})
}

func TestNewConnector(t *testing.T) {
	t.Run("disabled returns nil without error", func(t *testing.T) {
		c, err := NewConnector(Config{Enabled: false, URL: "redis://localhost:6379/0"})
		require.NoError(t, err)
		assert.Nil(t, c)
	})

	t.Run("URL mode connects to miniredis", func(t *testing.T) {
		server := miniredis.RunT(t)
		url := "redis://" + server.Addr() + "/0"
		c, err := NewConnector(Config{
			Enabled:          true,
			URL:              url,
			DisableIndentity: true,
		})
		require.NoError(t, err)
		require.NotNil(t, c)
		t.Cleanup(func() { c.Close() })
	})

	t.Run("legacy mode connects to miniredis", func(t *testing.T) {
		server := miniredis.RunT(t)
		c, err := NewConnector(Config{
			Enabled:          true,
			Addr:             server.Addr(),
			DisableIndentity: true,
		})
		require.NoError(t, err)
		require.NotNil(t, c)
		t.Cleanup(func() { c.Close() })
	})

	t.Run("invalid URL errors", func(t *testing.T) {
		_, err := NewConnector(Config{
			Enabled: true,
			URL:     "garbage://???",
		})
		require.Error(t, err)
	})
}
