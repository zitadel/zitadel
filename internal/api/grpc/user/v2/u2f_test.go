package user

import (
	"io"
	"testing"
	"time"

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

func Test_u2fRegistrationDetailsToPb(t *testing.T) {
	type args struct {
		details *domain.WebAuthNRegistrationDetails
		err     error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterU2FResponse
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
			want: &user.RegisterU2FResponse{
				Details: &object.Details{
					Sequence: 22,
					ChangeDate: &timestamppb.Timestamp{
						Seconds: 3000,
						Nanos:   22,
					},
					ResourceOwner: "me",
				},
				U2FId: "123",
				PublicKeyCredentialCreationOptions: &structpb.Struct{
					Fields: map[string]*structpb.Value{"foo": {Kind: &structpb.Value_StringValue{StringValue: "bar"}}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u2fRegistrationDetailsToPb(tt.args.details, tt.args.err)
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
