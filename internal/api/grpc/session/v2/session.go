package session

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

func (s *Server) GetSession(ctx context.Context, req *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	res, err := s.query.SessionByID(ctx, req.GetSessionId())
	if err != nil {
		return nil, err
	}
	return &session.GetSessionResponse{
		//Details: object.DomainToDetailsPb(), //TODO: ?
		Session: sessionToPb(res),
	}, nil
}

func (s *Server) ListSessions(ctx context.Context, req *session.ListSessionsRequest) (*session.ListSessionsResponse, error) {
	queries, err := listSessionsRequestToQuery(req)
	if err != nil {
		return nil, err
	}
	sessions, err := s.query.SearchSessions(ctx, queries)
	if err != nil {
		return nil, err
	}
	return &session.ListSessionsResponse{
		Details:  object.ToListDetails(sessions.SearchResponse),
		Sessions: sessionsToPb(sessions.Sessions),
	}, nil
}

func sessionsToPb(sessions []*query.Session) []*session.Session {
	s := make([]*session.Session, len(sessions))
	for i, session := range sessions {
		s[i] = sessionToPb(session)
	}
	return s
}

func sessionToPb(s *query.Session) *session.Session {
	return &session.Session{
		Id:       s.ID,
		Factors:  factorsToPb(s),
		Metadata: s.Metadata,
	}
}

func factorsToPb(s *query.Session) *session.Factors {
	factors := new(session.Factors)
	if s.UserID != "" && !s.UserCheckedAt.IsZero() {
		factors.User = &session.UserFactor{
			VerifiedAt: timestamppb.New(s.UserCheckedAt),
			User:       &user.User{Id: s.UserID},
		}
	}
	if !s.PasswordCheckedAt.IsZero() {
		factors.Password = &session.PasswordFactor{
			VerifiedAt: timestamppb.New(s.PasswordCheckedAt),
		}
	}
	return factors
}

func (s *Server) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	checks, metadata, err := s.createSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	set, err := s.command.CreateSession(ctx, checks, metadata)
	if err != nil {
		return nil, err
	}
	return &session.CreateSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionId:    set.ID,
		SessionToken: set.NewToken,
	}, nil
}

func (s *Server) SetSession(ctx context.Context, req *session.SetSessionRequest) (*session.SetSessionResponse, error) {
	checks, err := s.setSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	set, err := s.command.UpdateSession(ctx, req.GetSessionId(), req.GetSessionToken(), checks, req.GetMetadata())
	if err != nil {
		return nil, err
	}
	// if there's no new token, just return the current
	if set.NewToken == "" {
		set.NewToken = req.GetSessionToken()
	}
	return &session.SetSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionToken: set.NewToken,
	}, nil
}

func (s *Server) DeleteSession(ctx context.Context, req *session.DeleteSessionRequest) (*session.DeleteSessionResponse, error) {
	details, err := s.command.TerminateSession(ctx, req.GetSessionId(), req.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return &session.DeleteSessionResponse{
		Details: object.DomainToDetailsPb(details),
	}, nil
}

func listSessionsRequestToQuery(req *session.ListSessionsRequest) (*query.SessionsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := sessionQueriesToQuery(req.GetQueries())
	if err != nil {
		return nil, err
	}
	return &query.SessionsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset: offset,
			Limit:  limit,
			Asc:    asc,
		},
		Queries: queries,
	}, nil
}

func sessionQueriesToQuery(queries []*session.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries))
	for i, query := range queries {
		q[i], err = sessionQueryToQuery(query)
		if err != nil {
			return nil, err
		}
	}
	return q, nil
}

func sessionQueryToQuery(query *session.SearchQuery) (query.SearchQuery, error) {
	switch q := query.Query.(type) {
	case *session.SearchQuery_IdsQuery:
		return idsQueryToQuery(q.IdsQuery)
	//case *session.SearchQuery_TokensQuery:
	//	return TokensQueryToQuery(q.TokensQuery)
	default:
		return nil, caos_errs.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid")
	}
}

func idsQueryToQuery(q *session.IDsQuery) (query.SearchQuery, error) {
	return query.NewSessionIDsSearchQuery(q.Ids)
}

func (s *Server) createSessionRequestToCommand(ctx context.Context, req *session.CreateSessionRequest) ([]command.SessionCheck, map[string][]byte, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, nil, err
	}
	return checks, req.GetMetadata(), nil
}

func (s *Server) setSessionRequestToCommand(ctx context.Context, req *session.SetSessionRequest) ([]command.SessionCheck, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Server) checksToCommand(ctx context.Context, checks *session.Checks) ([]command.SessionCheck, error) {
	checkUser, err := userCheck(checks.GetUser())
	if err != nil {
		return nil, err
	}
	sessionChecks := make([]command.SessionCheck, 0, 2)
	if checkUser != nil {
		user, err := checkUser(ctx, s.query)
		if err != nil {
			return nil, err
		}
		sessionChecks = append(sessionChecks, command.CheckUser(user.ID))
	}
	if password := checks.GetPassword(); password != nil {
		sessionChecks = append(sessionChecks, command.CheckPassword(password.GetPassword()))
	}
	return sessionChecks, nil
}

func userCheck(user *session.CheckUser) (userSearch, error) {
	if user == nil {
		return nil, nil
	}
	switch s := user.GetSearch().(type) {
	case *session.CheckUser_UserId:
		return userByID(s.UserId), nil
	case *session.CheckUser_LoginName:
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
