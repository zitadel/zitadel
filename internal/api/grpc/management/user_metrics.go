package management

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/telemetry/metrics"
	mgmt_pb "github.com/zitadel/zitadel/pkg/grpc/management"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
	"go.opentelemetry.io/otel/metric"
)

const (
	TotalHumanUsersGauge       = "zitadel_human_users_total"
	TotalHumanUsersGaugeDesc   = "The total number of human users"
	TotalMachineUsersGauge     = "zitadel_machine_users_total"
	TotalMachineUsersGaugeDesc = "The total number of machine users"
)

var (
	HumanUserCount   uint64
	MachineUserCount uint64
)

func updateUserGaugeMetric(count uint64, req *mgmt_pb.ListUsersRequest) {
	var userType string

	userType = ""
	for _, query := range req.Queries {
		if typeQuery := query.GetTypeQuery(); typeQuery != nil {
			userType = typeQuery.Type.String()
			break
		}
	}

	if len(req.Queries) > 1 || userType == "" {
		return
	}

	switch userType {
	case user_pb.Type_TYPE_HUMAN.String():
		HumanUserCount = count
		callback := func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(int64(HumanUserCount))
			return nil
		}
		if err := metrics.RegisterValueObserver(
			TotalHumanUsersGauge,
			TotalHumanUsersGaugeDesc,
			callback,
		); err != nil {
			logging.WithError(err).Warn("failed to register", TotalHumanUsersGauge, "total observer")
		}
	case user_pb.Type_TYPE_MACHINE.String():
		MachineUserCount = count
		callback := func(ctx context.Context, observer metric.Int64Observer) error {
			observer.Observe(int64(MachineUserCount))
			return nil
		}
		if err := metrics.RegisterValueObserver(
			TotalMachineUsersGauge,
			TotalMachineUsersGaugeDesc,
			callback,
		); err != nil {
			logging.WithError(err).Warn("failed to register", TotalMachineUsersGauge, "total observer")
		}
	}
}
