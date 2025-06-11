package app

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/app/v2beta/convert"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func (s *Server) GetApplication(ctx context.Context, req *app.GetApplicationRequest) (*app.GetApplicationResponse, error) {
	res, err := s.query.AppByID(ctx, req.GetApplicationId(), false)
	if err != nil {
		return nil, err
	}

	return &app.GetApplicationResponse{
		App: convert.AppToPb(res),
	}, nil
}
