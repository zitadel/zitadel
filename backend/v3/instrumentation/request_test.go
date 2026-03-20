package instrumentation

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/rs/xid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testID = xid.New()

func TestSetInstanceID(t *testing.T) {
	originalXidWithTime := xidWithTime
	xidWithTime = func(time.Time) xid.ID {
		return testID
	}
	t.Cleanup(func() {
		xidWithTime = originalXidWithTime
	})

	tests := []struct {
		name        string
		ctx         context.Context
		instanceID  string
		wantDetails *requestDetails
	}{
		{
			name:        "no request details in context",
			ctx:         context.Background(),
			wantDetails: nil,
		},
		{
			name:       "request details in context",
			ctx:        WithRequestDetails(context.Background(), "instanceHost", "publicHost"),
			instanceID: "instanceID",
			wantDetails: &requestDetails{
				id:           testID,
				instanceID:   "instanceID",
				instanceHost: "instanceHost",
				publicHost:   "publicHost",
				userID:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetInstanceID(tt.ctx, tt.instanceID)
			got, ok := getRequestDetails(tt.ctx)
			require.Equal(t, tt.wantDetails != nil, ok)
			assert.Equal(t, tt.wantDetails, got)
		})
	}
}

func TestSetUserID(t *testing.T) {
	originalXidWithTime := xidWithTime
	xidWithTime = func(time.Time) xid.ID {
		return testID
	}
	t.Cleanup(func() {
		xidWithTime = originalXidWithTime
	})

	tests := []struct {
		name        string
		ctx         context.Context
		userID      string
		wantDetails *requestDetails
	}{
		{
			name:        "no request details in context",
			ctx:         context.Background(),
			wantDetails: nil,
		},
		{
			name:   "request details in context",
			ctx:    WithRequestDetails(context.Background(), "instanceHost", "publicHost"),
			userID: "userID",
			wantDetails: &requestDetails{
				id:           testID,
				instanceHost: "instanceHost",
				publicHost:   "publicHost",
				userID:       "userID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUserID(tt.ctx, tt.userID)
			got, ok := getRequestDetails(tt.ctx)
			require.Equal(t, tt.wantDetails != nil, ok)
			assert.Equal(t, tt.wantDetails, got)
		})
	}
}

func TestGetRequestID(t *testing.T) {
	originalXidWithTime := xidWithTime
	xidWithTime = func(time.Time) xid.ID {
		return testID
	}
	t.Cleanup(func() {
		xidWithTime = originalXidWithTime
	})

	reqCtx := WithRequestDetails(context.Background(), "instanceHost", "publicHost")
	tests := []struct {
		name   string
		ctx    context.Context
		wantID xid.ID
	}{
		{
			name:   "no request ID in context",
			ctx:    context.Background(),
			wantID: xid.NilID(),
		},
		{
			name:   "request ID in context",
			ctx:    reqCtx,
			wantID: testID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRequestID(tt.ctx)
			assert.Equal(t, tt.wantID, got)
		})
	}
}

func Test_requestDetails_slogAttributes(t *testing.T) {
	type fields struct {
		id           xid.ID
		instanceHost string
		publicHost   string
		instanceID   string
		userID       string
	}
	tests := []struct {
		name   string
		fields fields
		want   []any
	}{
		{
			name: "all fields set",
			fields: fields{
				id:           testID,
				instanceHost: "instanceHost",
				publicHost:   "publicHost",
				instanceID:   "instanceID",
				userID:       "userID",
			},
			want: []any{
				slog.String("id", testID.String()),
				slog.String("instance_host", "instanceHost"),
				slog.String("public_host", "publicHost"),
				slog.String("instance_id", "instanceID"),
				slog.String("user_id", "userID"),
			},
		},
		{
			name: "only required fields set",
			fields: fields{
				id:           testID,
				instanceHost: "instanceHost",
				publicHost:   "publicHost",
			},
			want: []any{
				slog.String("id", testID.String()),
				slog.String("instance_host", "instanceHost"),
				slog.String("public_host", "publicHost"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &requestDetails{
				id:           tt.fields.id,
				instanceHost: tt.fields.instanceHost,
				publicHost:   tt.fields.publicHost,
				instanceID:   tt.fields.instanceID,
				userID:       tt.fields.userID,
			}
			got := d.slogAttributes()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_requestDetailsExtractor(t *testing.T) {
	originalXidWithTime := xidWithTime
	xidWithTime = func(time.Time) xid.ID {
		return testID
	}
	t.Cleanup(func() {
		xidWithTime = originalXidWithTime
	})
	tests := []struct {
		name string
		ctx  context.Context
		want []slog.Attr
	}{
		{
			name: "no request details in context",
			ctx:  context.Background(),
			want: nil,
		},
		{
			name: "request details in context",
			ctx:  WithRequestDetails(context.Background(), "instanceHost", "publicHost"),
			want: []slog.Attr{
				slog.Group("request",
					slog.String("id", testID.String()),
					slog.String("instance_host", "instanceHost"),
					slog.String("public_host", "publicHost"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := requestDetailsExtractor(tt.ctx, time.Time{}, 0, "")
			assert.Equal(t, tt.want, got)
		})
	}
}
