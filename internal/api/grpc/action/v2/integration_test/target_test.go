//go:build integration

package action_test

import (
	"context"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"

	target_domain "github.com/zitadel/zitadel/internal/execution/target"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/action/v2"
)

func TestServer_CreateTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type want struct {
		id           bool
		creationDate bool
		signingKey   bool
	}
	alreadyExistingTargetName := integration.TargetName()
	instance.CreateTarget(isolatedIAMOwnerCTX, t, alreadyExistingTargetName, "https://example.com", target_domain.TargetTypeAsync, false)
	tests := []struct {
		name string
		ctx  context.Context
		req  *action.CreateTargetRequest
		want
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.CreateTargetRequest{
				Name: integration.TargetName(),
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: "",
			},
			wantErr: true,
		},
		{
			name: "empty type",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:       integration.TargetName(),
				TargetType: nil,
			},
			wantErr: true,
		},
		{
			name: "empty webhook url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: integration.TargetName(),
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty request response url",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name: integration.TargetName(),
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty timeout",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{},
				},
				Timeout: nil,
			},
			wantErr: true,
		},
		{
			name: "async, already existing, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     alreadyExistingTargetName,
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestAsync{
					RestAsync: &action.RESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			wantErr: true,
		},
		{
			name: "async, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestAsync{
					RestAsync: &action.RESTAsync{},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: want{
				id:           true,
				creationDate: true,
				signingKey:   true,
			},
		},
		{
			name: "webhook, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: want{
				id:           true,
				creationDate: true,
				signingKey:   true,
			},
		},
		{
			name: "webhook, interrupt on error, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestWebhook{
					RestWebhook: &action.RESTWebhook{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: want{
				id:           true,
				creationDate: true,
				signingKey:   true,
			},
		},
		{
			name: "call, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: false,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: want{
				id:           true,
				creationDate: true,
				signingKey:   true,
			},
		},

		{
			name: "call, interruptOnError, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.CreateTargetRequest{
				Name:     integration.TargetName(),
				Endpoint: "https://example.com",
				TargetType: &action.CreateTargetRequest_RestCall{
					RestCall: &action.RESTCall{
						InterruptOnError: true,
					},
				},
				Timeout: durationpb.New(10 * time.Second),
			},
			want: want{
				id:           true,
				creationDate: true,
				signingKey:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			got, err := instance.Client.ActionV2.CreateTarget(tt.ctx, tt.req)
			changeDate := time.Now().UTC()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertCreateTargetResponse(t, creationDate, changeDate, tt.want.creationDate, tt.want.id, tt.want.signingKey, got)
		})
	}
}

func assertCreateTargetResponse(t *testing.T, creationDate, changeDate time.Time, expectedCreationDate, expectedID, expectedSigningKey bool, actualResp *action.CreateTargetResponse) {
	if expectedCreationDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetCreationDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.CreationDate)
	}

	if expectedID {
		assert.NotEmpty(t, actualResp.GetId())
	} else {
		assert.Nil(t, actualResp.Id)
	}

	if expectedSigningKey {
		assert.NotEmpty(t, actualResp.GetSigningKey())
	} else {
		assert.Nil(t, actualResp.SigningKey)
	}
}

func TestServer_UpdateTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type args struct {
		ctx context.Context
		req *action.UpdateTargetRequest
	}
	type want struct {
		change     bool
		changeDate bool
		signingKey bool
	}
	tests := []struct {
		name    string
		prepare func(request *action.UpdateTargetRequest)
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(integration.TargetName()),
				},
			},
			wantErr: true,
		},
		{
			name: "not existing",
			prepare: func(request *action.UpdateTargetRequest) {
				request.Id = "notexisting"
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(integration.TargetName()),
				},
			},
			wantErr: true,
		},
		{
			name: "no change, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Endpoint: gu.Ptr("https://example.com"),
				},
			},
			want: want{
				change:     false,
				changeDate: true,
				signingKey: false,
			},
		},
		{
			name: "change name, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Name: gu.Ptr(integration.TargetName()),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: false,
			},
		},
		{
			name: "regenerate signingkey, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					ExpirationSigningKey: durationpb.New(0 * time.Second),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: true,
			},
		},
		{
			name: "change type, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestCall{
						RestCall: &action.RESTCall{
							InterruptOnError: true,
						},
					},
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: false,
			},
		},
		{
			name: "change url, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Endpoint: gu.Ptr("https://example.com/hooks/new"),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: false,
			},
		},
		{
			name: "change timeout, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					Timeout: durationpb.New(20 * time.Second),
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: false,
			},
		},
		{
			name: "change type async, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeAsync, false).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					TargetType: &action.UpdateTargetRequest_RestAsync{
						RestAsync: &action.RESTAsync{},
					},
				},
			},
			want: want{
				change:     true,
				changeDate: true,
				signingKey: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationDate := time.Now().UTC()
			tt.prepare(tt.args.req)

			got, err := instance.Client.ActionV2.UpdateTarget(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			changeDate := time.Time{}
			if tt.want.change {
				changeDate = time.Now().UTC()
			}
			assert.NoError(t, err)
			assertUpdateTargetResponse(t, creationDate, changeDate, tt.want.changeDate, tt.want.signingKey, got)
		})
	}
}

func assertUpdateTargetResponse(t *testing.T, creationDate, changeDate time.Time, expectedChangeDate, expectedSigningKey bool, actualResp *action.UpdateTargetResponse) {
	if expectedChangeDate {
		if !changeDate.IsZero() {
			assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), creationDate, changeDate)
		} else {
			assert.WithinRange(t, actualResp.GetChangeDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.ChangeDate)
	}

	if expectedSigningKey {
		assert.NotEmpty(t, actualResp.GetSigningKey())
	} else {
		assert.Nil(t, actualResp.SigningKey)
	}
}

func TestServer_DeleteTarget(t *testing.T) {
	instance := integration.NewInstance(CTX)
	iamOwnerCtx := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name             string
		ctx              context.Context
		prepare          func(request *action.DeleteTargetRequest) (time.Time, time.Time)
		req              *action.DeleteTargetRequest
		wantDeletionDate bool
		wantErr          bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.DeleteTargetRequest{
				Id: "notexisting",
			},
			wantErr: true,
		},
		{
			name: "empty id",
			ctx:  iamOwnerCtx,
			req: &action.DeleteTargetRequest{
				Id: "",
			},
			wantErr: true,
		},
		{
			name: "delete target, not existing",
			ctx:  iamOwnerCtx,
			req: &action.DeleteTargetRequest{
				Id: "notexisting",
			},
			wantDeletionDate: false,
		},
		{
			name: "delete target",
			ctx:  iamOwnerCtx,
			prepare: func(request *action.DeleteTargetRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				return creationDate, time.Time{}
			},
			req:              &action.DeleteTargetRequest{},
			wantDeletionDate: true,
		},
		{
			name: "delete target, already removed",
			ctx:  iamOwnerCtx,
			prepare: func(request *action.DeleteTargetRequest) (time.Time, time.Time) {
				creationDate := time.Now().UTC()
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", target_domain.TargetTypeWebhook, false).GetId()
				request.Id = targetID
				instance.DeleteTarget(iamOwnerCtx, t, targetID)
				return creationDate, time.Now().UTC()
			},
			req:              &action.DeleteTargetRequest{},
			wantDeletionDate: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var creationDate, deletionDate time.Time
			if tt.prepare != nil {
				creationDate, deletionDate = tt.prepare(tt.req)
			}
			got, err := instance.Client.ActionV2.DeleteTarget(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDeleteTargetResponse(t, creationDate, deletionDate, tt.wantDeletionDate, got)
		})
	}
}

func assertDeleteTargetResponse(t *testing.T, creationDate, deletionDate time.Time, expectedDeletionDate bool, actualResp *action.DeleteTargetResponse) {
	if expectedDeletionDate {
		if !deletionDate.IsZero() {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, deletionDate)
		} else {
			assert.WithinRange(t, actualResp.GetDeletionDate().AsTime(), creationDate, time.Now().UTC())
		}
	} else {
		assert.Nil(t, actualResp.DeletionDate)
	}
}
