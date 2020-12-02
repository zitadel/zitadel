package server

import (
	"context"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/proto"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type ValidationFunction func(ctx context.Context) error

type Validator struct {
	validations map[string]ValidationFunction
}

func NewValidator(validations map[string]ValidationFunction) *Validator {
	return &Validator{validations: validations}
}

func (v *Validator) Healthz(_ context.Context, e *empty.Empty) (*empty.Empty, error) {
	return e, nil
}

func (v *Validator) Ready(ctx context.Context, e *empty.Empty) (*empty.Empty, error) {
	if len(validate(ctx, v.validations)) == 0 {
		return e, nil
	}
	return nil, errors.ThrowInternal(nil, "API-2jD9a", "not ready")
}

func (v *Validator) Validate(ctx context.Context, _ *empty.Empty) (*structpb.Struct, error) {
	validations := validate(ctx, v.validations)
	return proto.ToPBStruct(validations)
}

func validate(ctx context.Context, validations map[string]ValidationFunction) map[string]error {
	errors := make(map[string]error)
	for id, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Error("validation failed")
			errors[id] = err
		}
	}
	return errors
}
