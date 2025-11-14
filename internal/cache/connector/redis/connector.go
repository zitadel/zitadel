package redis

import (
	"crypto/tls"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Enabled bool

	// The network type, either tcp or unix.
	// Default is tcp.
	Network string
	// host:port address.
	Addr string
	// ClientName will execute the `CLIENT SETNAME ClientName` command for each conn.
	ClientName string
	// Use the specified Username to authenticate the current connection
	// with one of the connections defined in the ACL list when connecting
	// to a Redis 6.0 instance, or greater, that is using the Redis ACL system.
	Username string
	// Optional password. Must match the password specified in the
	// requirepass server configuration option (if connecting to a Redis 5.0 instance, or lower),
	// or the User Password when connecting to a Redis 6.0 instance, or greater,
	// that is using the Redis ACL system.
	Password string
	// Each ZITADEL cache uses an incremental DB namespace.
	// This option offsets the first DB so it doesn't conflict with other databases on the same server.
	// Note that ZITADEL uses FLUSHDB command to truncate a cache.
	// This can have destructive consequences when overlapping DB namespaces are used.
	DBOffset int

	// Maximum number of retries before giving up.
	// Default is 3 retries; -1 (not 0) disables retries.
	MaxRetries int
	// Minimum backoff between each retry.
	// Default is 8 milliseconds; -1 disables backoff.
	MinRetryBackoff time.Duration
	// Maximum backoff between each retry.
	// Default is 512 milliseconds; -1 disables backoff.
	MaxRetryBackoff time.Duration

	// Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout time.Duration
	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetReadDeadline calls completely.
	ReadTimeout time.Duration
	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.  Supported values:
	//   - `0` - default timeout (3 seconds).
	//   - `-1` - no timeout (block indefinitely).
	//   - `-2` - disables SetWriteDeadline calls completely.
	WriteTimeout time.Duration

	// Type of connection pool.
	// true for FIFO pool, false for LIFO pool.
	// Note that FIFO has slightly higher overhead compared to LIFO,
	// but it helps closing idle connections faster reducing the pool size.
	PoolFIFO bool
	// Base number of socket connections.
	// Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
	// If there is not enough connections in the pool, new connections will be allocated in excess of PoolSize,
	// you can limit it through MaxActiveConns
	PoolSize int
	// Amount of time client waits for connection if all connections
	// are busy before returning an error.
	// Default is ReadTimeout + 1 second.
	PoolTimeout time.Duration
	// Minimum number of idle connections which is useful when establishing
	// new connection is slow.
	// Default is 0. the idle connections are not closed by default.
	MinIdleConns int
	// Maximum number of idle connections.
	// Default is 0. the idle connections are not closed by default.
	MaxIdleConns int
	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	MaxActiveConns int
	// ConnMaxIdleTime is the maximum amount of time a connection may be idle.
	// Should be less than server's timeout.
	//
	// Expired connections may be closed lazily before reuse.
	// If d <= 0, connections are not closed due to a connection's idle time.
	//
	// Default is 30 minutes. -1 disables idle timeout check.
	ConnMaxIdleTime time.Duration
	// ConnMaxLifetime is the maximum amount of time a connection may be reused.
	//
	// Expired connections may be closed lazily before reuse.
	// If <= 0, connections are not closed due to a connection's age.
	//
	// Default is to not close idle connections.
	ConnMaxLifetime time.Duration

	EnableTLS bool

	// Disable set-lib on connect. Default is false.
	DisableIndentity bool

	// Add suffix to client name. Default is empty.
	IdentitySuffix string

	CircuitBreaker *CBConfig
}

type Connector struct {
	*redis.Client
	Config Config
}

func NewConnector(config Config) *Connector {
	if !config.Enabled {
		return nil
	}
	return &Connector{
		Client: redis.NewClient(optionsFromConfig(config)),
		Config: config,
	}
}

func optionsFromConfig(c Config) *redis.Options {
	opts := &redis.Options{
		Network:               c.Network,
		Addr:                  c.Addr,
		ClientName:            c.ClientName,
		Protocol:              3,
		Username:              c.Username,
		Password:              c.Password,
		MaxRetries:            c.MaxRetries,
		MinRetryBackoff:       c.MinRetryBackoff,
		MaxRetryBackoff:       c.MaxRetryBackoff,
		DialTimeout:           c.DialTimeout,
		ReadTimeout:           c.ReadTimeout,
		WriteTimeout:          c.WriteTimeout,
		ContextTimeoutEnabled: true,
		PoolFIFO:              c.PoolFIFO,
		PoolTimeout:           c.PoolTimeout,
		MinIdleConns:          c.MinIdleConns,
		MaxIdleConns:          c.MaxIdleConns,
		MaxActiveConns:        c.MaxActiveConns,
		ConnMaxIdleTime:       c.ConnMaxIdleTime,
		ConnMaxLifetime:       c.ConnMaxLifetime,
		DisableIndentity:      c.DisableIndentity,
		IdentitySuffix:        c.IdentitySuffix,
		Limiter:               newLimiter(c.CircuitBreaker, c.MaxActiveConns),
	}
	if c.EnableTLS {
		opts.TLSConfig = new(tls.Config)
	}
	return opts
}
