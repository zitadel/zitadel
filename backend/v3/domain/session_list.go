package domain

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SessionFilter interface {
	sessionFilter()
}

// SessionIDsFilter filters sessions by one or more IDs.
type SessionIDsFilter struct{ IDs []string }

func (SessionIDsFilter) sessionFilter() {}

// SessionUserIDFilter filters sessions by the associated user's ID.
type SessionUserIDFilter struct{ UserID string }

func (SessionUserIDFilter) sessionFilter() {}

// SessionCreationDateFilter filters sessions by creation date.
type SessionCreationDateFilter struct {
	Op   database.NumberOperation
	Date time.Time
}

func (SessionCreationDateFilter) sessionFilter() {}

// SessionCreatorFilter filters sessions by the creator's ID.
type SessionCreatorFilter struct{ ID *string }

func (SessionCreatorFilter) sessionFilter() {}

// SessionUserAgentFilter filters sessions by user-agent fingerprint.
type SessionUserAgentFilter struct{ FingerprintID *string }

func (SessionUserAgentFilter) sessionFilter() {}

// SessionExpirationDateFilter filters sessions by expiration date.
type SessionExpirationDateFilter struct {
	Op   database.NumberOperation
	Date time.Time
}

func (SessionExpirationDateFilter) sessionFilter() {}

// SessionSortColumn determines the column used to sort the result set.
type SessionSortColumn int

const (
	SessionSortColumnUnspecified SessionSortColumn = iota
	SessionSortColumnCreationDate
)

// ListSessionsRequest holds all parameters for a ListSessions query.
type ListSessionsRequest struct {
	Filters    []SessionFilter
	SortColumn SessionSortColumn
	Ascending  bool
	Limit      uint32
	Offset     uint32
}

type ListSessionResponse struct {
	Session *Session
	User    *User
}

// ListSessionsQuery implements [Querier] for listing sessions.
type ListSessionsQuery struct {
	Request             *ListSessionsRequest
	listSessionResponse []ListSessionResponse
}

// Result implements [Querier].
func (l *ListSessionsQuery) Result() []ListSessionResponse { return l.listSessionResponse }

var _ Querier[[]ListSessionResponse] = (*ListSessionsQuery)(nil)

// NewListSessionsQuery returns a new [ListSessionsQuery] for the given request.
func NewListSessionsQuery(req *ListSessionsRequest) *ListSessionsQuery {
	return &ListSessionsQuery{Request: req}
}

// String implements [Querier].
func (*ListSessionsQuery) String() string { return "ListSessionsQuery" }

// Validate implements [Querier].
func (l *ListSessionsQuery) Validate(_ context.Context, _ *InvokeOpts) error {
	return nil
}

// Execute implements [Querier].
func (l *ListSessionsQuery) Execute(ctx context.Context, opts *InvokeOpts) error {
	sessionRepo := opts.sessionRepo

	conds, err := l.Conditions(ctx, sessionRepo)
	if err != nil {
		return err
	}

	sessions, err := sessionRepo.List(ctx, opts.DB(),
		database.WithCondition(conds),
		l.Sorting(sessionRepo),
		database.WithLimit(l.Request.Limit),
		database.WithOffset(l.Request.Offset),
	)
	if err != nil {
		return zerrors.ThrowInternal(err, "DOM-Yx8q2r", "Errors.Session.List")
	}

	if len(sessions) == 0 {
		return nil
	}

	userRepo := opts.userRepo
	sessionByUserID := make(map[string][]*Session, len(sessions))
	userIDs := make([]database.Condition, 0, len(sessions))
	l.listSessionResponse = make([]ListSessionResponse, 0, len(sessions))

	for _, session := range sessions {
		if session.UserID != "" {
			sessionByUserID[session.UserID] = append(sessionByUserID[session.UserID], session)
			userIDs = append(userIDs, userRepo.IDCondition(session.UserID))
			continue
		}
		// Sessions not associated with any user (e.g. newly created ones?) get added right away
		l.listSessionResponse = append(l.listSessionResponse, ListSessionResponse{Session: session})
	}

	users, err := userRepo.List(ctx, opts.DB(),
		database.WithCondition(database.And(
			userRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID()),
			database.Or(
				userIDs...,
			),
		)),
	)
	if err != nil {
		l.listSessionResponse = []ListSessionResponse{}
		return zerrors.ThrowInternal(err, "DOM-ILxIlM", "Errors.User.List")
	}

	if len(users) != len(userIDs) {
		// If there is a mismatch between the 2 slices, it means that a session has
		// a user ID that doesn't exist on DB. That should not happen.
		l.listSessionResponse = []ListSessionResponse{}
		return zerrors.ThrowNotFound(nil, "DOM-vh3SVZ", "Errors.User.Mismatch")
	}

	for _, user := range users {
		sessions, ok := sessionByUserID[user.ID]
		if !ok {
			l.listSessionResponse = []ListSessionResponse{}
			return zerrors.ThrowInternal(nil, "DOM-69ICaj", "Errors.SessionByUser.NotFound")
		}
		for _, s := range sessions {
			l.listSessionResponse = append(l.listSessionResponse, ListSessionResponse{Session: s, User: user})
		}
	}

	return nil
}

func (l *ListSessionsQuery) permissionCondition(ctx context.Context, sessionRepo SessionRepository) database.Condition {
	ctxData := authz.GetCtxData(ctx)
	return database.Or(
		sessionRepo.UserIDCondition(ctxData.UserID),
		sessionRepo.UserAgentIDCondition(ctxData.AgentID),
		sessionRepo.CreatorIDCondition(ctxData.UserID),
		database.PermissionCheck(SessionReadPermission, false),
	)
}

func (l *ListSessionsQuery) Conditions(ctx context.Context, sessionRepo SessionRepository) (database.Condition, error) {
	conds := make([]database.Condition, 2, 2+len(l.Request.Filters))
	conds[0] = sessionRepo.InstanceIDCondition(authz.GetInstance(ctx).InstanceID())
	conds[1] = l.permissionCondition(ctx, sessionRepo)

	for _, f := range l.Request.Filters {
		cond, err := l.sessionFilterToCondition(ctx, sessionRepo, f)
		if err != nil {
			return nil, err
		}
		conds = append(conds, cond)
	}
	return database.And(conds...), nil
}

func (l *ListSessionsQuery) sessionFilterToCondition(ctx context.Context, repo SessionRepository, filter SessionFilter) (database.Condition, error) {
	switch typedFilter := filter.(type) {
	case SessionIDsFilter:
		idConds := make([]database.Condition, len(typedFilter.IDs))
		for i, id := range typedFilter.IDs {
			idConds[i] = repo.IDCondition(id)
		}
		return database.Or(idConds...), nil

	case SessionUserIDFilter:
		return repo.UserIDCondition(typedFilter.UserID), nil

	case SessionCreationDateFilter:
		return repo.CreatedAtCondition(typedFilter.Op, typedFilter.Date), nil

	case SessionCreatorFilter:
		creatorID := authz.GetCtxData(ctx).UserID
		if typedFilter.ID != nil {
			creatorID = *typedFilter.ID
		}
		return repo.CreatorIDCondition(creatorID), nil

	case SessionUserAgentFilter:
		fp := authz.GetCtxData(ctx).AgentID
		if typedFilter.FingerprintID != nil {
			fp = *typedFilter.FingerprintID
		}
		return repo.UserAgentIDCondition(fp), nil

	case SessionExpirationDateFilter:
		// Unlike ES side, we are also including NOT EQUAL. This is an impossible case though since gRPC has no value for it
		expCond := repo.ExpirationCondition(typedFilter.Op, typedFilter.Date)
		if typedFilter.Op == database.NumberOperationGreaterThan || typedFilter.Op == database.NumberOperationGreaterThanOrEqual {
			// Include also nil expirations (i.e. non-expiring sessions)
			return database.Or(expCond, database.IsNull(repo.ExpirationColumn())), nil
		}
		return expCond, nil

	default:
		return nil, zerrors.ThrowInvalidArgument(NewUnexpectedQueryTypeError(typedFilter), "DOM-Cz5s3t", "List.Query.Invalid")
	}
}

// Sorting returns the QueryOption for ordering results.
func (l *ListSessionsQuery) Sorting(sessionRepo SessionRepository) database.QueryOption {
	if l.Request.SortColumn != SessionSortColumnCreationDate {
		return func(*database.QueryOpts) {}
	}
	direction := database.OrderDirectionDesc
	if l.Request.Ascending {
		direction = database.OrderDirectionAsc
	}
	return database.WithOrderBy(direction, sessionRepo.CreatedAtColumn())
}
