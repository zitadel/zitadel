package dialect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectionConfig_takeRatio(t *testing.T) {
	type fields struct {
		MaxOpenConns uint32
		MaxIdleConns uint32
	}
	tests := []struct {
		name    string
		fields  fields
		ratio   float64
		wantOut *ConnectionConfig
		wantErr error
	}{
		{
			name:    "ratio less than 0 error",
			ratio:   -0.1,
			wantErr: ErrNegativeRatio,
		},
		{
			name: "zero values",
			fields: fields{
				MaxOpenConns: 0,
				MaxIdleConns: 0,
			},
			ratio: 0,
			wantOut: &ConnectionConfig{
				MaxOpenConns: 0,
				MaxIdleConns: 0,
			},
		},
		{
			name: "max conns, ratio 0",
			fields: fields{
				MaxOpenConns: 10,
				MaxIdleConns: 5,
			},
			ratio: 0,
			wantOut: &ConnectionConfig{
				MaxOpenConns: 0,
				MaxIdleConns: 0,
			},
		},
		{
			name: "half ratio",
			fields: fields{
				MaxOpenConns: 10,
				MaxIdleConns: 5,
			},
			ratio: 0.5,
			wantOut: &ConnectionConfig{
				MaxOpenConns: 5,
				MaxIdleConns: 2,
			},
		},
		{
			name: "minimal 1",
			fields: fields{
				MaxOpenConns: 2,
				MaxIdleConns: 2,
			},
			ratio: 0.1,
			wantOut: &ConnectionConfig{
				MaxOpenConns: 1,
				MaxIdleConns: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := &ConnectionConfig{
				MaxOpenConns: tt.fields.MaxOpenConns,
				MaxIdleConns: tt.fields.MaxIdleConns,
			}
			got, err := in.takeRatio(tt.ratio)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestNewConnectionConfig(t *testing.T) {
	type args struct {
		openConns       uint32
		idleConns       uint32
		pusherRatio     float64
		projectionRatio float64
		purpose         DBPurpose
	}
	tests := []struct {
		name    string
		args    args
		want    *ConnectionConfig
		wantErr error
	}{
		{
			name: "illegal open conns error",
			args: args{
				openConns: 2,
				idleConns: 3,
			},
			wantErr: ErrIllegalMaxOpenConns,
		},
		{
			name: "illegal idle conns error",
			args: args{
				openConns: 3,
				idleConns: 2,
			},
			wantErr: ErrIllegalMaxIdleConns,
		},
		{
			name: "high ration sum error",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.5,
				projectionRatio: 0.5,
			},
			wantErr: ErrHighSumRatio,
		},
		{
			name: "illegal pusher ratio error",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     -0.1,
				projectionRatio: 0.5,
			},
			wantErr: ErrNegativeRatio,
		},
		{
			name: "illegal projection ratio error",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.5,
				projectionRatio: -0.1,
			},
			wantErr: ErrNegativeRatio,
		},
		{
			name: "invalid purpose error",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.4,
				projectionRatio: 0.4,
				purpose:         99,
			},
			wantErr: ErrInvalidPurpose,
		},
		{
			name: "min values, query purpose",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeQuery,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 1,
				MaxIdleConns: 1,
			},
		},
		{
			name: "min values, pusher purpose",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeEventPusher,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 1,
				MaxIdleConns: 1,
			},
		},
		{
			name: "min values, projection purpose",
			args: args{
				openConns:       3,
				idleConns:       3,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeProjectionSpooler,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 1,
				MaxIdleConns: 1,
			},
		},
		{
			name: "high values, query purpose",
			args: args{
				openConns:       10,
				idleConns:       5,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeQuery,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 6,
				MaxIdleConns: 3,
			},
		},
		{
			name: "high values, pusher purpose",
			args: args{
				openConns:       10,
				idleConns:       5,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeEventPusher,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 2,
				MaxIdleConns: 1,
			},
		},
		{
			name: "high values, projection purpose",
			args: args{
				openConns:       10,
				idleConns:       5,
				pusherRatio:     0.2,
				projectionRatio: 0.2,
				purpose:         DBPurposeProjectionSpooler,
			},
			want: &ConnectionConfig{
				MaxOpenConns: 2,
				MaxIdleConns: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConnectionConfig(tt.args.openConns, tt.args.idleConns, tt.args.pusherRatio, tt.args.projectionRatio, tt.args.purpose)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
