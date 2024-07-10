package user

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	object "github.com/zitadel/zitadel/pkg/grpc/object/v2"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func Test_passkeyAuthenticatorToDomain(t *testing.T) {
	tests := []struct {
		pa   user.PasskeyAuthenticator
		want domain.AuthenticatorAttachment
	}{
		{
			pa:   user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_UNSPECIFIED,
			want: domain.AuthenticatorAttachmentUnspecified,
		},
		{
			pa:   user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_PLATFORM,
			want: domain.AuthenticatorAttachmentPlattform,
		},
		{
			pa:   user.PasskeyAuthenticator_PASSKEY_AUTHENTICATOR_CROSS_PLATFORM,
			want: domain.AuthenticatorAttachmentCrossPlattform,
		},
		{
			pa:   999,
			want: domain.AuthenticatorAttachmentUnspecified,
		},
	}
	for _, tt := range tests {
		t.Run(tt.pa.String(), func(t *testing.T) {
			got := passkeyAuthenticatorToDomain(tt.pa)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_passkeyRegistrationDetailsToPb(t *testing.T) {
	type args struct {
		details *domain.WebAuthNRegistrationDetails
		err     error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterPasskeyResponse
		wantErr error
	}{
		{
			name: "an error",
			args: args{
				details: nil,
				err:     io.ErrClosedPipe,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "unmarshall error",
			args: args{
				details: &domain.WebAuthNRegistrationDetails{
					ObjectDetails: &domain.ObjectDetails{
						Sequence:      22,
						EventDate:     time.Unix(3000, 22),
						ResourceOwner: "me",
					},
					ID:                                 "123",
					PublicKeyCredentialCreationOptions: []byte(`\\`),
				},
				err: nil,
			},
			wantErr: zerrors.ThrowInternal(nil, "USERv2-Dohr6", "Errors.Internal"),
		},
		{
			name: "ok",
			args: args{
				details: &domain.WebAuthNRegistrationDetails{
					ObjectDetails: &domain.ObjectDetails{
						Sequence:      22,
						EventDate:     time.Unix(3000, 22),
						ResourceOwner: "me",
					},
					ID:                                 "123",
					PublicKeyCredentialCreationOptions: []byte(`{"foo": "bar"}`),
				},
				err: nil,
			},
			want: &user.RegisterPasskeyResponse{
				Details: &object.Details{
					Sequence: 22,
					ChangeDate: &timestamppb.Timestamp{
						Seconds: 3000,
						Nanos:   22,
					},
					ResourceOwner: "me",
				},
				PasskeyId: "123",
				PublicKeyCredentialCreationOptions: &structpb.Struct{
					Fields: map[string]*structpb.Value{"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := passkeyRegistrationDetailsToPb(tt.args.details, tt.args.err)
			require.ErrorIs(t, err, tt.wantErr)
			if !proto.Equal(tt.want, got) {
				t.Errorf("Not equal:\nExpected\n%s\nActual:%s", tt.want, got)
			}
			if tt.want != nil {
				grpc.AllFieldsSet(t, got.ProtoReflect())
			}
		})
	}
}

func Test_passkeyDetailsToPb(t *testing.T) {
	type args struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name string
		args args
		want *user.CreatePasskeyRegistrationLinkResponse
	}{
		{
			name: "an error",
			args: args{
				details: nil,
				err:     io.ErrClosedPipe,
			},
		},
		{
			name: "ok",
			args: args{
				details: &domain.ObjectDetails{
					Sequence:      22,
					EventDate:     time.Unix(3000, 22),
					ResourceOwner: "me",
				},
				err: nil,
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					Sequence: 22,
					ChangeDate: &timestamppb.Timestamp{
						Seconds: 3000,
						Nanos:   22,
					},
					ResourceOwner: "me",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := passkeyDetailsToPb(tt.args.details, tt.args.err)
			require.ErrorIs(t, err, tt.args.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_passkeyCodeDetailsToPb(t *testing.T) {
	type args struct {
		details *domain.PasskeyCodeDetails
		err     error
	}
	tests := []struct {
		name string
		args args
		want *user.CreatePasskeyRegistrationLinkResponse
	}{
		{
			name: "an error",
			args: args{
				details: nil,
				err:     io.ErrClosedPipe,
			},
		},
		{
			name: "ok",
			args: args{
				details: &domain.PasskeyCodeDetails{
					ObjectDetails: &domain.ObjectDetails{
						Sequence:      22,
						EventDate:     time.Unix(3000, 22),
						ResourceOwner: "me",
					},
					CodeID: "123",
					Code:   "456",
				},
				err: nil,
			},
			want: &user.CreatePasskeyRegistrationLinkResponse{
				Details: &object.Details{
					Sequence: 22,
					ChangeDate: &timestamppb.Timestamp{
						Seconds: 3000,
						Nanos:   22,
					},
					ResourceOwner: "me",
				},
				Code: &user.PasskeyRegistrationCode{
					Id:   "123",
					Code: "456",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := passkeyCodeDetailsToPb(tt.args.details, tt.args.err)
			require.ErrorIs(t, err, tt.args.err)
			assert.Equal(t, tt.want, got)
			if tt.want != nil {
				grpc.AllFieldsSet(t, got.ProtoReflect())
			}
		})
	}
}
