package ratelimit

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPGRateLimiter_Check(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(`INSERT INTO signals.rate_limit_counters (key, count, window_start, window_secs)
		 VALUES ($1, 1, $2, $3)
		 ON CONFLICT (key) DO UPDATE SET
		     count = CASE
		         WHEN signals.rate_limit_counters.window_start + (signals.rate_limit_counters.window_secs || ' seconds')::INTERVAL <= $2
		         THEN 1
		         ELSE signals.rate_limit_counters.count + 1
		     END,
		     window_start = CASE
		         WHEN signals.rate_limit_counters.window_start + (signals.rate_limit_counters.window_secs || ' seconds')::INTERVAL <= $2
		         THEN $2
		         ELSE signals.rate_limit_counters.window_start
		     END,
		     window_secs = $3
		 RETURNING count`)
	mock.ExpectQuery(query).
		WithArgs("tenant:rule:key", sqlmock.AnyArg(), 300).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

	limiter := NewPGRateLimiter(db)
	count, allowed := limiter.Check(context.Background(), "tenant:rule:key", 5*time.Minute, 1)
	if count != 2 || allowed {
		t.Fatalf("Check() = (%d, %v), want (2, false)", count, allowed)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}

func TestPGRateLimiter_Prune(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}
	defer db.Close()

	query := regexp.QuoteMeta(`DELETE FROM signals.rate_limit_counters
		 WHERE window_start + (window_secs || ' seconds')::INTERVAL <= $1`)
	mock.ExpectExec(query).
		WithArgs(sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 3))

	limiter := NewPGRateLimiter(db)
	limiter.Prune(context.Background())

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("ExpectationsWereMet(): %v", err)
	}
}
