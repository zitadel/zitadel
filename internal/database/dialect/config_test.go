package dialect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBPurpose_AppName(t *testing.T) {
	tests := []struct {
		p    DBPurpose
		want string
	}{
		{
			p:    DBPurposeQuery,
			want: QueryAppName,
		},
		{
			p:    DBPurposeEventPusher,
			want: EventstorePusherAppName,
		},
		{
			p:    DBPurposeProjectionSpooler,
			want: ProjectionSpoolerAppName,
		},
		{
			p:    99,
			want: defaultAppName,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.AppName())
		})
	}
}
