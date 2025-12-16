//go:build integration

package action_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/crypto"
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
	instance.CreateTarget(isolatedIAMOwnerCTX, t, alreadyExistingTargetName, "https://example.com", target_domain.TargetTypeAsync, false, action.PayloadType_PAYLOAD_TYPE_JSON)
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
				Timeout:     durationpb.New(10 * time.Second),
				PayloadType: action.PayloadType_PAYLOAD_TYPE_JSON,
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
				Timeout:     durationpb.New(10 * time.Second),
				PayloadType: action.PayloadType_PAYLOAD_TYPE_JSON,
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
				Timeout:     durationpb.New(10 * time.Second),
				PayloadType: action.PayloadType_PAYLOAD_TYPE_JSON,
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
				Timeout:     durationpb.New(10 * time.Second),
				PayloadType: action.PayloadType_PAYLOAD_TYPE_JSON,
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
				Timeout:     durationpb.New(10 * time.Second),
				PayloadType: action.PayloadType_PAYLOAD_TYPE_JSON,
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
			name: "change timeout, ok",
			prepare: func(request *action.UpdateTargetRequest) {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				request.Id = targetID
			},
			args: args{
				ctx: isolatedIAMOwnerCTX,
				req: &action.UpdateTargetRequest{
					PayloadType: action.PayloadType_PAYLOAD_TYPE_JWT,
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
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeAsync, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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
				targetID := instance.CreateTarget(iamOwnerCtx, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
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

var (
	publicKey = func() *ecdsa.PublicKey {
		privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		return &privateKey.PublicKey
	}()
	publicKeyBytes = func() []byte {
		pubBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
		return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})
	}()
	fingerprint = func() string {
		block, _ := x509.MarshalPKIXPublicKey(publicKey)
		hash := sha256.Sum256(block)
		return fmt.Sprintf("SHA256:%s", base64.RawStdEncoding.EncodeToString(hash[:]))
	}()

	rsaPublicKey = func() *rsa.PublicKey {
		privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		return &privateKey.PublicKey
	}()
	rsaPublicKeyBytes = func() []byte {
		data, _ := crypto.PublicKeyToBytes(rsaPublicKey)
		return data
	}()
	rsaFingerprint = func() string {
		block, _ := x509.MarshalPKIXPublicKey(rsaPublicKey)
		hash := sha256.Sum256(block)
		return fmt.Sprintf("SHA256:%s", base64.RawStdEncoding.EncodeToString(hash[:]))
	}()
)

func TestServer_AddPublicKey(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	type want struct {
		id           bool
		creationDate bool
	}
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.AddPublicKeyRequest
		want    want
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.AddPublicKeyRequest{
				TargetId:       "targetID",
				PublicKey:      publicKeyBytes,
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "empty target id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.AddPublicKeyRequest{
				TargetId:       "",
				PublicKey:      publicKeyBytes,
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "empty public key",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.AddPublicKeyRequest{
				TargetId:       "targetID",
				PublicKey:      nil,
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "invalid public key",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.AddPublicKeyRequest{
				TargetId:       "targetID",
				PublicKey:      []byte("invalid"),
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "expiration date in the past",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.AddPublicKeyRequest{
				TargetId:       "targetID",
				PublicKey:      publicKeyBytes,
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(-1 * time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "target not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.AddPublicKeyRequest{
				TargetId:       "notexisting",
				PublicKey:      publicKeyBytes,
				ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
			},
			wantErr: true,
		},
		{
			name: "add public key, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.AddPublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				return &action.AddPublicKeyRequest{
					TargetId:       targetID,
					PublicKey:      publicKeyBytes,
					ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
				}
			}(),
			want: want{
				id:           true,
				creationDate: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			got, err := instance.Client.ActionV2.AddPublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, got.GetKeyId())
			assertDateWithinRangeClockSkew(t, got.GetCreationDate().AsTime(), start, time.Now())
		})
	}
}

func TestServer_ActivatePublicKey(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.ActivatePublicKeyRequest
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.ActivatePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty target id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.ActivatePublicKeyRequest{
				TargetId: "",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty key id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.ActivatePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "",
			},
			wantErr: true,
		},
		{
			name: "target not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.ActivatePublicKeyRequest{
				TargetId: "notexisting",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "key not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.ActivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				return &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    "notexisting",
				}
			}(),
			wantErr: true,
		},
		{
			name: "already expired key",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.ActivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:       targetID,
					PublicKey:      publicKeyBytes,
					ExpirationDate: timestamppb.New(time.Now().Add(1 * time.Second)),
				})
				require.NoError(t, err)
				time.Sleep(2 * time.Second)
				return &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: true,
		},
		{
			name: "already activated key",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.ActivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:  targetID,
					PublicKey: publicKeyBytes,
				})
				require.NoError(t, err)
				_, err = instance.Client.ActionV2.ActivatePublicKey(isolatedIAMOwnerCTX, &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				})
				require.NoError(t, err)
				return &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: false,
		},
		{
			name: "activate key, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.ActivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:       targetID,
					PublicKey:      publicKeyBytes,
					ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
				})
				require.NoError(t, err)
				return &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			got, err := instance.Client.ActionV2.ActivatePublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDateWithinRangeClockSkew(t, got.GetChangeDate().AsTime(), start, time.Now())
		})
	}
}

func TestServer_DeactivatePublicKey(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.DeactivatePublicKeyRequest
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.DeactivatePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty target id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.DeactivatePublicKeyRequest{
				TargetId: "",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty key id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.DeactivatePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "",
			},
			wantErr: true,
		},
		{
			name: "target not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.DeactivatePublicKeyRequest{
				TargetId: "notexisting",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "key not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.DeactivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				return &action.DeactivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    "notexisting",
				}
			}(),
			wantErr: true,
		},
		{
			name: "already deactivated key",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.DeactivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:       targetID,
					PublicKey:      publicKeyBytes,
					ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
				})
				require.NoError(t, err)
				return &action.DeactivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: false,
		},
		{
			name: "deactivate key, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.DeactivatePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:       targetID,
					PublicKey:      publicKeyBytes,
					ExpirationDate: timestamppb.New(time.Now().UTC().Add(time.Hour)),
				})
				require.NoError(t, err)
				_, err = instance.Client.ActionV2.ActivatePublicKey(isolatedIAMOwnerCTX, &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				})
				require.NoError(t, err)
				return &action.DeactivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			got, err := instance.Client.ActionV2.DeactivatePublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assertDateWithinRangeClockSkew(t, got.GetChangeDate().AsTime(), start, time.Now())
		})
	}
}

func TestServer_RemovePublicKey(t *testing.T) {
	instance := integration.NewInstance(CTX)
	isolatedIAMOwnerCTX := instance.WithAuthorizationToken(CTX, integration.UserTypeIAMOwner)
	tests := []struct {
		name    string
		ctx     context.Context
		req     *action.RemovePublicKeyRequest
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  instance.WithAuthorizationToken(context.Background(), integration.UserTypeOrgOwner),
			req: &action.RemovePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty target id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.RemovePublicKeyRequest{
				TargetId: "",
				KeyId:    "keyID",
			},
			wantErr: true,
		},
		{
			name: "empty key id",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.RemovePublicKeyRequest{
				TargetId: "targetID",
				KeyId:    "",
			},
			wantErr: true,
		},
		{
			name: "target not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: &action.RemovePublicKeyRequest{
				TargetId: "notexisting",
				KeyId:    "keyID",
			},
			wantErr: false,
		},
		{
			name: "key not existing",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.RemovePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				return &action.RemovePublicKeyRequest{
					TargetId: targetID,
					KeyId:    "notexisting",
				}
			}(),
			wantErr: false,
		},
		{
			name: "remove activated key",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.RemovePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:  targetID,
					PublicKey: publicKeyBytes,
				})
				require.NoError(t, err)
				_, err = instance.Client.ActionV2.ActivatePublicKey(isolatedIAMOwnerCTX, &action.ActivatePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				})
				require.NoError(t, err)
				return &action.RemovePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: true,
		},
		{
			name: "remove key, ok",
			ctx:  isolatedIAMOwnerCTX,
			req: func() *action.RemovePublicKeyRequest {
				targetID := instance.CreateTarget(isolatedIAMOwnerCTX, t, "", "https://example.com", target_domain.TargetTypeWebhook, false, action.PayloadType_PAYLOAD_TYPE_JSON).GetId()
				resp, err := instance.Client.ActionV2.AddPublicKey(isolatedIAMOwnerCTX, &action.AddPublicKeyRequest{
					TargetId:  targetID,
					PublicKey: publicKeyBytes,
				})
				require.NoError(t, err)
				return &action.RemovePublicKeyRequest{
					TargetId: targetID,
					KeyId:    resp.GetKeyId(),
				}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := instance.Client.ActionV2.RemovePublicKey(tt.ctx, tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func assertDateWithinRangeClockSkew(t assert.TestingT, actual time.Time, start, end time.Time) {
	assert.WithinRange(t, actual, start.Add(-time.Second), end.Add(time.Second))
}
