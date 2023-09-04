package query

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_calculateThreshold(t *testing.T) {
	type args struct {
		usedRel             uint16
		notificationPercent uint16
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "80 - below configuration",
			args: args{
				usedRel:             70,
				notificationPercent: 80,
			},
			want: 0,
		},
		{
			name: "80 - below 100 percent use",
			args: args{
				usedRel:             90,
				notificationPercent: 80,
			},
			want: 80,
		},
		{
			name: "80 - above 100 percent use",
			args: args{
				usedRel:             120,
				notificationPercent: 80,
			},
			want: 80,
		},
		{
			name: "80 - more than twice the use",
			args: args{
				usedRel:             190,
				notificationPercent: 80,
			},
			want: 180,
		},
		{
			name: "100 - below 100 percent use",
			args: args{
				usedRel:             90,
				notificationPercent: 100,
			},
			want: 0,
		},
		{
			name: "100 - above 100 percent use",
			args: args{
				usedRel:             120,
				notificationPercent: 100,
			},
			want: 100,
		},
		{
			name: "100 - more than twice the use",
			args: args{
				usedRel:             210,
				notificationPercent: 100,
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateThreshold(tt.args.usedRel, tt.args.notificationPercent)
			assert.Equal(t, int(tt.want), int(got))
		})
	}
}
