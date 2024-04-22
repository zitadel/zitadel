package template

import (
	"context"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/settings/template/v2"
)

func (s *Server) SetHTMLTemplate(_ context.Context, request *template.SetHTMLTemplateRequest) (*template.SetHTMLTemplateResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "TPL-NtK93", "Not implemented. Got: %v", request)
}

func (s *Server) ResolveHTMLTemplate(ctx context.Context, request *template.ResolveHTMLTemplateRequest) (*template.ResolveHTMLTemplateResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "TPL-4D8y6", "Not implemented. Got: %v", request)
}

func (s *Server) SetTextTemplate(ctx context.Context, request *template.SetTextTemplateRequest) (*template.SetTextTemplateResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "TPL-4YIBu", "Not implemented. Got: %v", request)
}

func (s *Server) ResolveTextTemplate(ctx context.Context, request *template.ResolveTextTemplateRequest) (*template.ResolveTextTemplateResponse, error) {
	return nil, zerrors.ThrowUnimplementedf(nil, "TPL-bzmxX", "Not implemented. Got: %v", request)
}
