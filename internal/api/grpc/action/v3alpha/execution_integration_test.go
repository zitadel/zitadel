//go:build integration

package action_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	action "github.com/zitadel/zitadel/pkg/grpc/action/v3alpha"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
)

func TestServer_SetExecution_Request(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Request_Include(t *testing.T) {
	ensureFeatureEnabled(t)

	targetResp := Tester.CreateTarget(CTX, t)
	executionCond := "request"
	Tester.SetExecution(CTX, t,
		&action.Condition{
			ConditionType: &action.Condition_Request{
				Request: &action.RequestExecution{
					Condition: &action.RequestExecution_All{
						All: true,
					},
				},
			},
		},
		[]string{targetResp.GetId()},
		[]string{},
	)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "method, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Includes: []string{executionCond},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Includes: []string{executionCond},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{
								All: true,
							},
						},
					},
				},
				Includes: []string{executionCond},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteExecution_Request(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *action.DeleteExecutionRequest) error
		req     *action.DeleteExecutionRequest
		want    *action.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/NotExisting",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/GetSession",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_Service{
								Service: "zitadel.user.v2beta.UserService",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Request{
						Request: &action.RequestExecution{
							Condition: &action.RequestExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Response(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.NotExistingService/List",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/ListSessions",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "zitadel.session.v2beta.SessionService",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteExecution_Response(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *action.DeleteExecutionRequest) error
		req     *action.DeleteExecutionRequest
		want    *action.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/NotExisting",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "method, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Method{
								Method: "/zitadel.session.v2beta.SessionService/GetSession",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "service, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "NotExistingService",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "service, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_Service{
								Service: "zitadel.user.v2beta.UserService",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Event(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		/*
			//TODO event existing check

			{
				name: "event, not existing",
				ctx:  CTX,
				req: &action.SetExecutionRequest{
					Condition: &action.Condition{
						ConditionType: &action.Condition_Event{
							Event: &action.EventExecution{
								Condition: &action.EventExecution_Event{
									Event: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "event, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Event{
								Event: "xxx",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		/*
			// TODO:

			{
				name: "group, not existing",
				ctx:  CTX,
				req: &action.SetExecutionRequest{
					Condition: &action.Condition{
						ConditionType: &action.Condition_Event{
							Event: &action.EventExecution{
								Condition: &action.EventExecution_Group{
									Group: "xxx",
								},
							},
						},
					},
					Targets: []string{targetResp.GetId()},
				},
				wantErr: true,
			},
		*/
		{
			name: "group, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteExecution_Event(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *action.DeleteExecutionRequest) error
		req     *action.DeleteExecutionRequest
		want    *action.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{},
					},
				},
			},
			wantErr: true,
		},
		/*
			//TODO: add when check is implemented
			{
				name: "event, not existing",
				ctx:  CTX,
				req: &action.DeleteExecutionRequest{
					Condition: &action.Condition{
						ConditionType: &action.Condition_Event{
							Event: &action.EventExecution{
								Condition: &action.EventExecution_Event{
									Event: "xxx",
								},
							},
						},
					},
				},
				wantErr: true,
			},
		*/
		{
			name: "event, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Event{
								Event: "xxx",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "group, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "group, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_Group{
								Group: "xxx",
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
		{
			name: "all, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Event{
						Event: &action.EventExecution{
							Condition: &action.EventExecution_All{
								All: true,
							},
						},
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_SetExecution_Function(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.SetExecutionRequest
		want    *action.SetExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{All: true},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "function, not existing",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: "xxx",
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  CTX,
			req: &action.SetExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication",
					},
				},
				Targets: []string{targetResp.GetId()},
			},
			want: &action.SetExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.SetExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteExecution_Function(t *testing.T) {
	ensureFeatureEnabled(t)
	targetResp := Tester.CreateTarget(CTX, t)

	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, request *action.DeleteExecutionRequest) error
		req     *action.DeleteExecutionRequest
		want    *action.DeleteExecutionResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{
							Condition: &action.ResponseExecution_All{
								All: true,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no condition, error",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Response{
						Response: &action.ResponseExecution{},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "function, not existing",
			ctx:  CTX,
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: "xxx",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "function, ok",
			ctx:  CTX,
			dep: func(ctx context.Context, request *action.DeleteExecutionRequest) error {
				Tester.SetExecution(ctx, t, request.GetCondition(), []string{targetResp.GetId()}, []string{})
				return nil
			},
			req: &action.DeleteExecutionRequest{
				Condition: &action.Condition{
					ConditionType: &action.Condition_Function{
						Function: "Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication",
					},
				},
			},
			want: &action.DeleteExecutionResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Instance.InstanceID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				err := tt.dep(tt.ctx, tt.req)
				require.NoError(t, err)
			}

			got, err := Client.DeleteExecution(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}
