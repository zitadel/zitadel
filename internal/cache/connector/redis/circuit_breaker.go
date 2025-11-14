package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sony/gobreaker/v2"
	"github.com/zitadel/logging"
)

const defaultInflightSize = 100000

type CBConfig struct {
	// Interval when the counters are reset to 0.
	// 0 interval never resets the counters until the CB is opened.
	Interval time.Duration
	// Amount of consecutive failures permitted
	MaxConsecutiveFailures uint32
	// The ratio of failed requests out of total requests
	MaxFailureRatio float64
	// Timeout after opening of the CB, until the state is set to half-open.
	Timeout time.Duration
	// The allowed amount of requests that are allowed to pass when the CB is half-open.
	MaxRetryRequests uint32
}

func (config *CBConfig) readyToTrip(counts gobreaker.Counts) bool {
	if config.MaxConsecutiveFailures > 0 && counts.ConsecutiveFailures > config.MaxConsecutiveFailures {
		return true
	}
	if config.MaxFailureRatio > 0 && counts.Requests > 0 {
		failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
		return failureRatio > config.MaxFailureRatio
	}
	return false
}

// limiter implements [redis.Limiter] as a circuit breaker.
type limiter struct {
	inflight chan func(success bool)
	cb       *gobreaker.TwoStepCircuitBreaker[struct{}]
}

func newLimiter(config *CBConfig, maxActiveConns int) redis.Limiter {
	if config == nil {
		return nil
	}
	// The size of the inflight channel needs to be big enough for maxActiveConns to prevent blocking.
	// When that is 0 (no limit), we must set a sane default.
	if maxActiveConns <= 0 {
		maxActiveConns = defaultInflightSize
	}
	return &limiter{
		inflight: make(chan func(success bool), maxActiveConns),
		cb: gobreaker.NewTwoStepCircuitBreaker[struct{}](gobreaker.Settings{
			Name:        "redis cache",
			MaxRequests: config.MaxRetryRequests,
			Interval:    config.Interval,
			Timeout:     config.Timeout,
			ReadyToTrip: config.readyToTrip,
			OnStateChange: func(name string, from, to gobreaker.State) {
				logging.WithFields("name", name, "from", from, "to", to).Warn("circuit breaker state change")
			},
		}),
	}
}

// Allow implements [redis.Limiter].
func (l *limiter) Allow() error {
	done, err := l.cb.Allow()
	if err != nil {
		return err
	}
	l.inflight <- done
	return nil
}

// ReportResult implements [redis.Limiter].
//
// ReportResult checks the error returned by the Redis client.
// `nil`, [redis.Nil] and [context.Canceled] are not considered failures.
// Any other error, like connection or [context.DeadlineExceeded] is counted as a failure.
func (l *limiter) ReportResult(err error) {
	done := <-l.inflight
	done(err == nil ||
		errors.Is(err, redis.Nil) ||
		errors.Is(err, context.Canceled) ||
		redis.HasErrorPrefix(err, "NOSCRIPT"))
}
