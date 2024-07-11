package session

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/muhlemmer/gu"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/authz"
	object "github.com/zitadel/zitadel/internal/api/grpc/object/v2beta"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	objpb "github.com/zitadel/zitadel/pkg/grpc/object"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
)

var (
	timestampComparisons = map[objpb.TimestampQueryMethod]query.TimestampComparison{
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_EQUALS:            query.TimestampEquals,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER:           query.TimestampGreater,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_GREATER_OR_EQUALS: query.TimestampGreaterOrEquals,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS:              query.TimestampLess,
		objpb.TimestampQueryMethod_TIMESTAMP_QUERY_METHOD_LESS_OR_EQUALS:    query.TimestampLessOrEquals,
	}
)

func (s *Server) GetSession(ctx context.Context, req *session.GetSessionRequest) (*session.GetSessionResponse, error) {
	res, err := s.query.SessionByID(ctx, true, req.GetSessionId(), req.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return &session.GetSessionResponse{
		Session: sessionToPb(res),
	}, nil
}

func (s *Server) ListSessions(ctx context.Context, req *session.ListSessionsRequest) (*session.ListSessionsResponse, error) {
	queries, err := listSessionsRequestToQuery(ctx, req)
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

func (s *Server) CreateSession(ctx context.Context, req *session.CreateSessionRequest) (*session.CreateSessionResponse, error) {
	checks, metadata, userAgent, lifetime, err := s.createSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	challengeResponse, cmds, err := s.challengesToCommand(req.GetChallenges(), checks)
	if err != nil {
		return nil, err
	}

	set, err := s.command.CreateSession(ctx, cmds, metadata, userAgent, lifetime)
	if err != nil {
		return nil, err
	}

	return &session.CreateSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionId:    set.ID,
		SessionToken: set.NewToken,
		Challenges:   challengeResponse,
	}, nil
}

func (s *Server) SetSession(ctx context.Context, req *session.SetSessionRequest) (*session.SetSessionResponse, error) {
	checks, err := s.setSessionRequestToCommand(ctx, req)
	if err != nil {
		return nil, err
	}
	challengeResponse, cmds, err := s.challengesToCommand(req.GetChallenges(), checks)
	if err != nil {
		return nil, err
	}

	set, err := s.command.UpdateSession(ctx, req.GetSessionId(), cmds, req.GetMetadata(), req.GetLifetime().AsDuration())
	if err != nil {
		return nil, err
	}
	return &session.SetSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionToken: set.NewToken,
		Challenges:   challengeResponse,
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

func sessionsToPb(sessions []*query.Session) []*session.Session {
	s := make([]*session.Session, len(sessions))
	for i, session := range sessions {
		s[i] = sessionToPb(session)
	}
	return s
}

func sessionToPb(s *query.Session) *session.Session {
	return &session.Session{
		Id:             s.ID,
		CreationDate:   timestamppb.New(s.CreationDate),
		ChangeDate:     timestamppb.New(s.ChangeDate),
		Sequence:       s.Sequence,
		Factors:        factorsToPb(s),
		Metadata:       s.Metadata,
		UserAgent:      userAgentToPb(s.UserAgent),
		ExpirationDate: expirationToPb(s.Expiration),
	}
}

func userAgentToPb(ua domain.UserAgent) *session.UserAgent {
	if ua.IsEmpty() {
		return nil
	}

	out := &session.UserAgent{
		FingerprintId: ua.FingerprintID,
		Description:   ua.Description,
	}
	if ua.IP != nil {
		out.Ip = gu.Ptr(ua.IP.String())
	}
	if ua.Header == nil {
		return out
	}
	out.Header = make(map[string]*session.UserAgent_HeaderValues, len(ua.Header))
	for k, v := range ua.Header {
		out.Header[k] = &session.UserAgent_HeaderValues{
			Values: v,
		}
	}
	return out
}

func expirationToPb(expiration time.Time) *timestamppb.Timestamp {
	if expiration.IsZero() {
		return nil
	}
	return timestamppb.New(expiration)
}

func factorsToPb(s *query.Session) *session.Factors {
	user := userFactorToPb(s.UserFactor)
	if user == nil {
		return nil
	}
	return &session.Factors{
		User:     user,
		Password: passwordFactorToPb(s.PasswordFactor),
		WebAuthN: webAuthNFactorToPb(s.WebAuthNFactor),
		Intent:   intentFactorToPb(s.IntentFactor),
		Totp:     totpFactorToPb(s.TOTPFactor),
		OtpSms:   otpFactorToPb(s.OTPSMSFactor),
		OtpEmail: otpFactorToPb(s.OTPEmailFactor),
	}
}

func passwordFactorToPb(factor query.SessionPasswordFactor) *session.PasswordFactor {
	if factor.PasswordCheckedAt.IsZero() {
		return nil
	}
	return &session.PasswordFactor{
		VerifiedAt: timestamppb.New(factor.PasswordCheckedAt),
	}
}

func intentFactorToPb(factor query.SessionIntentFactor) *session.IntentFactor {
	if factor.IntentCheckedAt.IsZero() {
		return nil
	}
	return &session.IntentFactor{
		VerifiedAt: timestamppb.New(factor.IntentCheckedAt),
	}
}

func webAuthNFactorToPb(factor query.SessionWebAuthNFactor) *session.WebAuthNFactor {
	if factor.WebAuthNCheckedAt.IsZero() {
		return nil
	}
	return &session.WebAuthNFactor{
		VerifiedAt:   timestamppb.New(factor.WebAuthNCheckedAt),
		UserVerified: factor.UserVerified,
	}
}

func totpFactorToPb(factor query.SessionTOTPFactor) *session.TOTPFactor {
	if factor.TOTPCheckedAt.IsZero() {
		return nil
	}
	return &session.TOTPFactor{
		VerifiedAt: timestamppb.New(factor.TOTPCheckedAt),
	}
}

func otpFactorToPb(factor query.SessionOTPFactor) *session.OTPFactor {
	if factor.OTPCheckedAt.IsZero() {
		return nil
	}
	return &session.OTPFactor{
		VerifiedAt: timestamppb.New(factor.OTPCheckedAt),
	}
}

func userFactorToPb(factor query.SessionUserFactor) *session.UserFactor {
	if factor.UserID == "" || factor.UserCheckedAt.IsZero() {
		return nil
	}
	return &session.UserFactor{
		VerifiedAt:     timestamppb.New(factor.UserCheckedAt),
		Id:             factor.UserID,
		LoginName:      factor.LoginName,
		DisplayName:    factor.DisplayName,
		OrganizationId: factor.ResourceOwner,
	}
}

func listSessionsRequestToQuery(ctx context.Context, req *session.ListSessionsRequest) (*query.SessionsSearchQueries, error) {
	offset, limit, asc := object.ListQueryToQuery(req.Query)
	queries, err := sessionQueriesToQuery(ctx, req.GetQueries())
	if err != nil {
		return nil, err
	}
	return &query.SessionsSearchQueries{
		SearchRequest: query.SearchRequest{
			Offset:        offset,
			Limit:         limit,
			Asc:           asc,
			SortingColumn: fieldNameToSessionColumn(req.GetSortingColumn()),
		},
		Queries: queries,
	}, nil
}

func sessionQueriesToQuery(ctx context.Context, queries []*session.SearchQuery) (_ []query.SearchQuery, err error) {
	q := make([]query.SearchQuery, len(queries)+1)
	for i, v := range queries {
		q[i], err = sessionQueryToQuery(v)
		if err != nil {
			return nil, err
		}
	}
	creatorQuery, err := query.NewSessionCreatorSearchQuery(authz.GetCtxData(ctx).UserID)
	if err != nil {
		return nil, err
	}
	q[len(queries)] = creatorQuery
	return q, nil
}

func sessionQueryToQuery(sq *session.SearchQuery) (query.SearchQuery, error) {
	switch q := sq.Query.(type) {
	case *session.SearchQuery_IdsQuery:
		return idsQueryToQuery(q.IdsQuery)
	case *session.SearchQuery_UserIdQuery:
		return query.NewUserIDSearchQuery(q.UserIdQuery.GetId())
	case *session.SearchQuery_CreationDateQuery:
		return creationDateQueryToQuery(q.CreationDateQuery)
	default:
		return nil, zerrors.ThrowInvalidArgument(nil, "GRPC-Sfefs", "List.Query.Invalid")
	}
}

func idsQueryToQuery(q *session.IDsQuery) (query.SearchQuery, error) {
	return query.NewSessionIDsSearchQuery(q.Ids)
}

func creationDateQueryToQuery(q *session.CreationDateQuery) (query.SearchQuery, error) {
	comparison := timestampComparisons[q.GetMethod()]
	return query.NewCreationDateQuery(q.GetCreationDate().AsTime(), comparison)
}

func fieldNameToSessionColumn(field session.SessionFieldName) query.Column {
	switch field {
	case session.SessionFieldName_SESSION_FIELD_NAME_CREATION_DATE:
		return query.SessionColumnCreationDate
	default:
		return query.Column{}
	}
}

func (s *Server) createSessionRequestToCommand(ctx context.Context, req *session.CreateSessionRequest) ([]command.SessionCommand, map[string][]byte, *domain.UserAgent, time.Duration, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, nil, nil, 0, err
	}
	return checks, req.GetMetadata(), userAgentToCommand(req.GetUserAgent()), req.GetLifetime().AsDuration(), nil
}

func userAgentToCommand(userAgent *session.UserAgent) *domain.UserAgent {
	if userAgent == nil {
		return nil
	}
	out := &domain.UserAgent{
		FingerprintID: userAgent.FingerprintId,
		IP:            net.ParseIP(userAgent.GetIp()),
		Description:   userAgent.Description,
	}
	if len(userAgent.Header) > 0 {
		out.Header = make(http.Header, len(userAgent.Header))
		for k, values := range userAgent.Header {
			out.Header[k] = values.GetValues()
		}
	}
	return out
}

func (s *Server) setSessionRequestToCommand(ctx context.Context, req *session.SetSessionRequest) ([]command.SessionCommand, error) {
	checks, err := s.checksToCommand(ctx, req.Checks)
	if err != nil {
		return nil, err
	}
	return checks, nil
}

func (s *Server) checksToCommand(ctx context.Context, checks *session.Checks) ([]command.SessionCommand, error) {
	checkUser, err := userCheck(checks.GetUser())
	if err != nil {
		return nil, err
	}
	sessionChecks := make([]command.SessionCommand, 0, 7)
	if checkUser != nil {
		user, err := checkUser.search(ctx, s.query)
		if err != nil {
			return nil, err
		}
		if !user.State.IsEnabled() {
			return nil, zerrors.ThrowPreconditionFailed(nil, "SESSION-Gj4ko", "Errors.User.NotActive")
		}

		var preferredLanguage *language.Tag
		if user.Human != nil && !user.Human.PreferredLanguage.IsRoot() {
			preferredLanguage = &user.Human.PreferredLanguage
		}
		sessionChecks = append(sessionChecks, command.CheckUser(user.ID, user.ResourceOwner, preferredLanguage))
	}
	if password := checks.GetPassword(); password != nil {
		sessionChecks = append(sessionChecks, command.CheckPassword(password.GetPassword()))
	}
	if intent := checks.GetIdpIntent(); intent != nil {
		sessionChecks = append(sessionChecks, command.CheckIntent(intent.GetIdpIntentId(), intent.GetIdpIntentToken()))
	}
	if passkey := checks.GetWebAuthN(); passkey != nil {
		sessionChecks = append(sessionChecks, s.command.CheckWebAuthN(passkey.GetCredentialAssertionData()))
	}
	if totp := checks.GetTotp(); totp != nil {
		sessionChecks = append(sessionChecks, command.CheckTOTP(totp.GetCode()))
	}
	if otp := checks.GetOtpSms(); otp != nil {
		sessionChecks = append(sessionChecks, command.CheckOTPSMS(otp.GetCode()))
	}
	if otp := checks.GetOtpEmail(); otp != nil {
		sessionChecks = append(sessionChecks, command.CheckOTPEmail(otp.GetCode()))
	}
	return sessionChecks, nil
}

func (s *Server) challengesToCommand(challenges *session.RequestChallenges, cmds []command.SessionCommand) (*session.Challenges, []command.SessionCommand, error) {
	if challenges == nil {
		return nil, cmds, nil
	}
	resp := new(session.Challenges)
	if req := challenges.GetWebAuthN(); req != nil {
		challenge, cmd := s.createWebAuthNChallengeCommand(req)
		resp.WebAuthN = challenge
		cmds = append(cmds, cmd)
	}
	if req := challenges.GetOtpSms(); req != nil {
		challenge, cmd := s.createOTPSMSChallengeCommand(req)
		resp.OtpSms = challenge
		cmds = append(cmds, cmd)
	}
	if req := challenges.GetOtpEmail(); req != nil {
		challenge, cmd, err := s.createOTPEmailChallengeCommand(req)
		if err != nil {
			return nil, nil, err
		}
		resp.OtpEmail = challenge
		cmds = append(cmds, cmd)
	}
	return resp, cmds, nil
}

func (s *Server) createWebAuthNChallengeCommand(req *session.RequestChallenges_WebAuthN) (*session.Challenges_WebAuthN, command.SessionCommand) {
	challenge := &session.Challenges_WebAuthN{
		PublicKeyCredentialRequestOptions: new(structpb.Struct),
	}
	userVerification := userVerificationRequirementToDomain(req.GetUserVerificationRequirement())
	return challenge, s.command.CreateWebAuthNChallenge(userVerification, req.GetDomain(), challenge.PublicKeyCredentialRequestOptions)
}

func userVerificationRequirementToDomain(req session.UserVerificationRequirement) domain.UserVerificationRequirement {
	switch req {
	case session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_UNSPECIFIED:
		return domain.UserVerificationRequirementUnspecified
	case session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_REQUIRED:
		return domain.UserVerificationRequirementRequired
	case session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_PREFERRED:
		return domain.UserVerificationRequirementPreferred
	case session.UserVerificationRequirement_USER_VERIFICATION_REQUIREMENT_DISCOURAGED:
		return domain.UserVerificationRequirementDiscouraged
	default:
		return domain.UserVerificationRequirementUnspecified
	}
}

func (s *Server) createOTPSMSChallengeCommand(req *session.RequestChallenges_OTPSMS) (*string, command.SessionCommand) {
	if req.GetReturnCode() {
		challenge := new(string)
		return challenge, s.command.CreateOTPSMSChallengeReturnCode(challenge)
	}

	return nil, s.command.CreateOTPSMSChallenge()

}

func (s *Server) createOTPEmailChallengeCommand(req *session.RequestChallenges_OTPEmail) (*string, command.SessionCommand, error) {
	switch t := req.GetDeliveryType().(type) {
	case *session.RequestChallenges_OTPEmail_SendCode_:
		cmd, err := s.command.CreateOTPEmailChallengeURLTemplate(t.SendCode.GetUrlTemplate())
		if err != nil {
			return nil, nil, err
		}
		return nil, cmd, nil
	case *session.RequestChallenges_OTPEmail_ReturnCode_:
		challenge := new(string)
		return challenge, s.command.CreateOTPEmailChallengeReturnCode(challenge), nil
	case nil:
		return nil, s.command.CreateOTPEmailChallenge(), nil
	default:
		return nil, nil, zerrors.ThrowUnimplementedf(nil, "SESSION-k3ng0", "delivery_type oneOf %T in OTPEmailChallenge not implemented", t)
	}
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
		return nil, zerrors.ThrowUnimplementedf(nil, "SESSION-d3b4g0", "user search %T not implemented", s)
	}
}

type userSearch interface {
	search(ctx context.Context, q *query.Queries) (*query.User, error)
}

func userByID(userID string) userSearch {
	return userSearchByID{userID}
}

func userByLoginName(loginName string) (userSearch, error) {
	return userSearchByLoginName{loginName}, nil
}

type userSearchByID struct {
	id string
}

func (u userSearchByID) search(ctx context.Context, q *query.Queries) (*query.User, error) {
	return q.GetUserByID(ctx, true, u.id)
}

type userSearchByLoginName struct {
	loginName string
}

func (u userSearchByLoginName) search(ctx context.Context, q *query.Queries) (*query.User, error) {
	return q.GetUserByLoginName(ctx, true, u.loginName)
}
