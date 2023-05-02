package session

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
)

func (s *Server) CreateSession(ctx context.Context, req *CreateSessionRequest) (*CreateSessionResponse, error) {
	checks, metadata, err := s.createSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	set, err := s.command.CreateSession(ctx, checks, metadata)
	if err != nil {
		return nil, err
	}
	return &CreateSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionId:    set.ID,
		SessionToken: set.Token,
	}, nil
}

func (s *Server) SetSession(ctx context.Context, req *SetSessionRequest) (*SetSessionResponse, error) {
	checks, err := s.setSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	set, err := s.command.UpdateSession(ctx, req.GetSessionId(), req.GetSessionToken(), checks, req.GetMetadata())
	if err != nil {
		return nil, err
	}
	return &SetSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionToken: set.Token,
	}, nil
}

func (s *Server) DeleteSession(ctx context.Context, req *DeleteSessionRequest) (*DeleteSessionResponse, error) {
	terminated, err := s.command.TerminateSession(ctx, req.GetSessionId(), req.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return &DeleteSessionResponse{
		Details: object.DomainToDetailsPb(terminated.ObjectDetails),
	}, nil
}

func (s *Server) createSessionRequestToCommand(ctx context.Context, req *CreateSessionRequest) (*command.SessionChecks, map[string][]byte, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, nil, err
	}
	return checks, req.GetMetadata(), nil
}

func (s *Server) setSessionRequestToCommand(ctx context.Context, req *SetSessionRequest) (*command.SessionChecks, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Server) checksToCommand(ctx context.Context, checks *Checks) (*command.SessionChecks, error) {
	checkUser, err := userCheck(checks.GetUser())
	if err != nil {
		return nil, err
	}
	sessionChecks := s.command.NewSessionChecks()
	if checkUser != nil {
		user, err := checkUser(ctx, s.query)
		if err != nil {
			return nil, err
		}
		sessionChecks.CheckUser(user.ID)
	}
	if password := checks.GetPassword(); password != nil {
		sessionChecks.CheckPassword(password.GetPassword())
	}
	return sessionChecks, nil
}

func userCheck(user *CheckUser) (userSearch, error) {
	if user == nil {
		return nil, nil
	}
	switch s := user.GetSearch().(type) {
	case *CheckUser_UserId:
		return userByID(s.UserId), nil
	case *CheckUser_LoginName:
		return userByLoginName(s.LoginName)
	default:
		return nil, caos_errs.ThrowUnimplementedf(nil, "SESSION-d3b4g0", "user search %T not implemented", s)
	}
}

type userSearch func(ctx context.Context, q *query.Queries) (*query.User, error)

func userByID(id string) userSearch {
	return func(ctx context.Context, q *query.Queries) (*query.User, error) {
		return q.GetUserByID(ctx, true, id, false)
	}
}

func userByLoginName(loginName string) (userSearch, error) {
	loginNameQuery, err := query.NewUserLoginNamesSearchQuery(loginName)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context, q *query.Queries) (*query.User, error) {
		return q.GetUser(ctx, true, false, loginNameQuery)
	}, nil
}
