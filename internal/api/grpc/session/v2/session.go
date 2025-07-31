package session

import (
	"context"
	"net"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/zitadel/internal/api/grpc/object/v2"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func (s *Server) CreateSession(ctx context.Context, req *connect.Request[session.CreateSessionRequest]) (*connect.Response[session.CreateSessionResponse], error) {
	checks, metadata, userAgent, lifetime, err := s.createSessionRequestToCommand(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	challengeResponse, cmds, err := s.challengesToCommand(req.Msg.GetChallenges(), checks)
	if err != nil {
		return nil, err
	}

	set, err := s.command.CreateSession(ctx, cmds, metadata, userAgent, lifetime)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&session.CreateSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionId:    set.ID,
		SessionToken: set.NewToken,
		Challenges:   challengeResponse,
	}), nil
}

func (s *Server) SetSession(ctx context.Context, req *connect.Request[session.SetSessionRequest]) (*connect.Response[session.SetSessionResponse], error) {
	checks, err := s.setSessionRequestToCommand(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	challengeResponse, cmds, err := s.challengesToCommand(req.Msg.GetChallenges(), checks)
	if err != nil {
		return nil, err
	}

	set, err := s.command.UpdateSession(ctx, req.Msg.GetSessionId(), req.Msg.GetSessionToken(), cmds, req.Msg.GetMetadata(), req.Msg.GetLifetime().AsDuration())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&session.SetSessionResponse{
		Details:      object.DomainToDetailsPb(set.ObjectDetails),
		SessionToken: set.NewToken,
		Challenges:   challengeResponse,
	}), nil
}

func (s *Server) DeleteSession(ctx context.Context, req *connect.Request[session.DeleteSessionRequest]) (*connect.Response[session.DeleteSessionResponse], error) {
	details, err := s.command.TerminateSession(ctx, req.Msg.GetSessionId(), req.Msg.GetSessionToken())
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(&session.DeleteSessionResponse{
		Details: object.DomainToDetailsPb(details),
	}), nil
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
	return q.GetUserByID(ctx, false, u.id)
}

type userSearchByLoginName struct {
	loginName string
}

func (u userSearchByLoginName) search(ctx context.Context, q *query.Queries) (*query.User, error) {
	return q.GetUserByLoginName(ctx, true, u.loginName)
}
