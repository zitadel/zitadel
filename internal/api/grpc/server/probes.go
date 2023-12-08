package server

import (
	"context"

	"github.com/zitadel/logging"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type ValidationFunction func(ctx context.Context) error

type Validator struct {
	validations map[string]ValidationFunction
}

func NewValidator(validations map[string]ValidationFunction) *Validator {
	return &Validator{validations: validations}
}

func (v *Validator) Healthz(_ context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	return e, nil
}

func (v *Validator) Ready(ctx context.Context, e *emptypb.Empty) (*emptypb.Empty, error) {
	if len(validate(ctx, v.validations)) == 0 {
		return e, nil
	}
	return nil, zerrors.ThrowInternal(nil, "API-2jD9a", "not ready")
}

func (v *Validator) Validate(ctx context.Context, _ *emptypb.Empty) (*structpb.Struct, error) {
	return structpb.NewStruct(validate(ctx, v.validations))
}

func validate(ctx context.Context, validations map[string]ValidationFunction) map[string]any {
	errors := make(map[string]any)
	for id, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Error("validation failed")
			errors[id] = err
		}
	}
	return errors
}
