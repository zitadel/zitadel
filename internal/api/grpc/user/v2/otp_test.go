package user

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zitadel/zitadel/internal/domain"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_otpDetailsToPb(t *testing.T) {
	type args struct {
		otp *domain.OTPv2
		err error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RegisterOTPResponse
		wantErr error
	}{
		{
			name: "error",
			args: args{
				err: io.ErrClosedPipe,
			},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "success",
			args: args{
				otp: &domain.OTPv2{
					ObjectDetails: &domain.ObjectDetails{
						Sequence:      123,
						EventDate:     time.Unix(456, 789),
						ResourceOwner: "me",
					},
					Secret: "secret",
					URI:    "URI",
				},
			},
			want: &user.RegisterOTPResponse{
				Details: &object.Details{
					Sequence: 123,
					ChangeDate: &timestamppb.Timestamp{
						Seconds: 456,
						Nanos:   789,
					},
					ResourceOwner: "me",
				},
				Secret: "secret",
				Uri:    "URI",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := otpDetailsToPb(tt.args.otp, tt.args.err)
			require.ErrorIs(t, err, tt.wantErr)
			if !proto.Equal(tt.want, got) {
				t.Errorf("RegisterOTPResponse =\n%v\nwant\n%v", got, tt.want)
			}
		})
	}
}
