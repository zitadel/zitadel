package session

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	user "github.com/zitadel/zitadel/internal/api/grpc/user/v2alpha"
)

func (s *Server) CreateSession(ctx context.Context, req *CreateSessionRequest) (*CreateSessionResponse, error) {
	return &CreateSessionResponse{
		SessionId: "hodor",
	}, nil
}
func (s *Server) SetSession(ctx context.Context, req *SetSessionRequest) (*SetSessionResponse, error) {
	test.id = req.GetSessionId()
	if user := req.Checks.GetUser(); user != nil {
		switch s := user.GetSearch().(type) {
		case *CheckUser_UserId:
			test.userID = s.UserId
		case *CheckUser_Username:
			if s.Username == "test" {
				test.userID = "id"
			}
		default:
			return nil, status.Error(codes.NotFound, "user not found")
		}
		test.userCheck = time.Now()
	}
	if pw := req.Checks.GetPassword(); pw != nil {
		if test.userID == "" {
			return nil, status.Error(codes.FailedPrecondition, "missing user")
		}
		if pw.GetPassword() != "correct" {
			return nil, status.Error(codes.InvalidArgument, "password invalid")
		}
		test.password = time.Now()
	}
	return &SetSessionResponse{}, nil
}

func (s *Server) GetSession(ctx context.Context, req *GetSessionRequest) (*GetSessionResponse, error) {
	if test.id != req.SessionId {
		return nil, status.Error(codes.NotFound, "invalid id")
	}
	factors := make(map[string]*Factor, 2)
	if test.userID != "" {
		factors["user"] = &Factor{
			Checked: timestamppb.New(test.userCheck),
			Factor:  &Factor_User{User: &user.User{Id: test.userID}},
		}
	}
	if !test.password.IsZero() {
		factors["password"] = &Factor{
			Checked: timestamppb.New(test.password),
		}
	}
	return &GetSessionResponse{
		Session: &Session{
			Id:      req.SessionId,
			Factors: factors,
		},
	}, nil
}

type session struct {
	id        string
	userID    string
	userCheck time.Time
	password  time.Time
}

var test session
