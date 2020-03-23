package grpc

import (
	"context"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/proto"
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
	return e, ready(ctx, v.validations)
}

func (v *Validator) Validate(ctx context.Context, _ *empty.Empty) (*structpb.Struct, error) {
	validations := validate(ctx, v.validations)
	return proto.ToPBStruct(validations)
}

func ready(ctx context.Context, validations map[string]ValidationFunction) error {
	if len(validate(ctx, validations)) == 0 {
		return nil
	}
	return errors.ThrowInternal(nil, "API-2jD9a", "not ready")
}

func validate(ctx context.Context, validations map[string]ValidationFunction) map[string]error {
	errors := make(map[string]error)
	for id, validation := range validations {
		if err := validation(ctx); err != nil {
			logging.Log("API-vf823").WithError(err).Error("validation failed")
			errors[id] = err
		}
	}
	return errors
}
