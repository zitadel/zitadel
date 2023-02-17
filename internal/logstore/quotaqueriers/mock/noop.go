package mock

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.QuotaQuerier = (*inmemReporter)(nil)

type inmemReporter struct {
	config      *quota.AddedEvent
	startPeriod time.Time
}

func NewNoopQuerier(quota *quota.AddedEvent, startPeriod time.Time) *inmemReporter {
	return &inmemReporter{config: quota, startPeriod: startPeriod}
}

func (i *inmemReporter) GetCurrentQuotaPeriod(context.Context, string, quota.Unit) (*quota.AddedEvent, time.Time, error) {
	return i.config, i.startPeriod, nil
}

func (*inmemReporter) GetDueQuotaNotifications(context.Context, *quota.AddedEvent, time.Time, uint64) ([]*quota.NotifiedEvent, error) {
	return nil, nil
}
