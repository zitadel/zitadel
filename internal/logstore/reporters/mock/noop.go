package mock

import (
	"context"

	"github.com/zitadel/zitadel/internal/logstore"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

var _ logstore.UsageReporter = (*inmemReporter)(nil)

type inmemReporter struct {
	quota     *query.Quota
	lastUsage uint64
}

func NewNoopReporter(quota *query.Quota) *inmemReporter {
	return &inmemReporter{quota: quota}
}

func (i *inmemReporter) GetQuota(_ context.Context, _ string, _ quota.Unit) (*query.Quota, error) {
	return i.quota, nil
}

func (i *inmemReporter) Report(_ context.Context, _ *query.Quota, _ uint64) error {
	return nil
}
