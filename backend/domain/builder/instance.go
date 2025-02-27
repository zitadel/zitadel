package builder

import (
	"context"
	"log/slog"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/cache"
	"github.com/zitadel/zitadel/backend/telemetry/logging"
	"github.com/zitadel/zitadel/backend/telemetry/tracing"
)

type InstanceBuilder struct {
	*Builder[*InstanceBuilder]

	tracer tracing.Tracer
	logger *logging.Logger
	cache  *cache.Instance
	db     repository.InstanceRepository
}

func NewInstanceBuilder() *InstanceBuilder {
	return &InstanceBuilder{
		Builder: NewBuilder(func() *InstanceBuilder {
			return new(InstanceBuilder)
		}),
		cache:  new(cache.Instance),
		tracer: tracing.NewTracer("instance"),
		logger: &logging.Logger{Logger: slog.Default().With("service", "instance")},
	}
}

var _ builder = (*InstanceBuilder)(nil)

func (b *InstanceBuilder) reset() {
	b.db = nil
}

var _ repository.InstanceRepository = (*InstanceBuilder)(nil)

// ByDomain implements repository.InstanceRepository.
func (b *InstanceBuilder) ByDomain(ctx context.Context, domain string) (*repository.Instance, error) {
	panic("unimplemented")
}

// ByID implements repository.InstanceRepository.
func (b *InstanceBuilder) ByID(ctx context.Context, id string) (*repository.Instance, error) {
	panic("unimplemented")
}

// SetUp implements repository.InstanceRepository.
func (b *InstanceBuilder) SetUp(ctx context.Context, instance *repository.Instance) error {
	panic("unimplemented")
}
