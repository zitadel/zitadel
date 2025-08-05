package org

import (
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2beta/orgconnect"
)

var _ orgconnect.OrganizationServiceHandler = (*Server)(nil)

type Server struct {
	orgconnect.UnimplementedOrganizationServiceHandler

	logger logging.Logger
	tracer tracing.Tracer
}

func NewServer(logger logging.Logger, tracer tracing.Tracer) *Server {
	return &Server{
		logger: logger,
		tracer: tracer,
	}
}
