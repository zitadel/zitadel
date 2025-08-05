package connect_middleware

import (
	"context"

	"connectrpc.com/connect"
	// import to make sure go.mod does not lose it
	// because dependency is only needed for generated code
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
)

func ValidationHandler() connect.UnaryInterceptorFunc {
	return func(handler connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			return validate(ctx, req, handler)
		}
	}
}

// validator interface needed for github.com/envoyproxy/protoc-gen-validate
// (it does not expose an interface itself)
type validator interface {
	Validate() error
}

func validate(ctx context.Context, req connect.AnyRequest, handler connect.UnaryFunc) (connect.AnyResponse, error) {
	validate, ok := req.Any().(validator)
	if !ok {
		return handler(ctx, req)
	}
	err := validate.Validate()
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return handler(ctx, req)
}
