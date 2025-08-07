package info

import (
	"context"

	"google.golang.org/grpc/codes"
)

type activityInfoKey struct{}

type ActivityInfo struct {
	Method        string
	Path          string
	RequestMethod string
	GRPCStatus    codes.Code
	HTTPStatus    int
}

func (a *ActivityInfo) IntoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, activityInfoKey{}, a)
}

func ActivityInfoFromContext(ctx context.Context) *ActivityInfo {
	m := ctx.Value(activityInfoKey{})
	if m == nil {
		return &ActivityInfo{}
	}
	ai, ok := m.(*ActivityInfo)
	if !ok {
		return &ActivityInfo{}
	}
	return ai
}

func (a *ActivityInfo) SetMethod(method string) *ActivityInfo {
	a.Method = method
	return a
}

func (a *ActivityInfo) SetPath(path string) *ActivityInfo {
	a.Path = path
	return a
}

func (a *ActivityInfo) SetRequestMethod(method string) *ActivityInfo {
	a.RequestMethod = method
	return a
}

func (a *ActivityInfo) SetGRPCStatus(status codes.Code) *ActivityInfo {
	a.GRPCStatus = status
	return a
}

func (a *ActivityInfo) SetHTTPStatus(status int) *ActivityInfo {
	a.HTTPStatus = status
	return a
}
