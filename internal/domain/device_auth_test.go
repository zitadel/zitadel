package domain

import (
	"testing"
)

func TestDeviceAuthState_Exists(t *testing.T) {
	tests := []struct {
		s    DeviceAuthState
		want bool
	}{
		{
			s:    DeviceAuthStateUndefined,
			want: false,
		},
		{
			s:    DeviceAuthStateInitiated,
			want: true,
		},
		{
			s:    DeviceAuthStateApproved,
			want: true,
		},
		{
			s:    DeviceAuthStateDenied,
			want: true,
		},
		{
			s:    DeviceAuthStateExpired,
			want: true,
		},
		{
			s:    DeviceAuthStateRemoved,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s.String(), func(t *testing.T) {
			if got := tt.s.Exists(); got != tt.want {
				t.Errorf("DeviceAuthState.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceAuthState_Done(t *testing.T) {
	tests := []struct {
		s    DeviceAuthState
		want bool
	}{
		{
			s:    DeviceAuthStateUndefined,
			want: false,
		},
		{
			s:    DeviceAuthStateInitiated,
			want: false,
		},
		{
			s:    DeviceAuthStateApproved,
			want: true,
		},
		{
			s:    DeviceAuthStateDenied,
			want: false,
		},
		{
			s:    DeviceAuthStateExpired,
			want: false,
		},
		{
			s:    DeviceAuthStateRemoved,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.s.String(), func(t *testing.T) {
			if got := tt.s.Done(); got != tt.want {
				t.Errorf("DeviceAuthState.Done() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceAuthState_Denied(t *testing.T) {
	tests := []struct {
		name string
		s    DeviceAuthState
		want bool
	}{
		{
			s:    DeviceAuthStateUndefined,
			want: false,
		},
		{
			s:    DeviceAuthStateInitiated,
			want: false,
		},
		{
			s:    DeviceAuthStateApproved,
			want: false,
		},
		{
			s:    DeviceAuthStateDenied,
			want: true,
		},
		{
			s:    DeviceAuthStateExpired,
			want: true,
		},
		{
			s:    DeviceAuthStateRemoved,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Denied(); got != tt.want {
				t.Errorf("DeviceAuthState.Denied() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeviceAuthCanceled_State(t *testing.T) {
	tests := []struct {
		name string
		c    DeviceAuthCanceled
		want DeviceAuthState
	}{
		{
			name: "empty",
			want: DeviceAuthStateUndefined,
		},
		{
			name: "invalid",
			c:    "foo",
			want: DeviceAuthStateUndefined,
		},
		{
			name: "denied",
			c:    DeviceAuthCanceledDenied,
			want: DeviceAuthStateDenied,
		},
		{
			name: "expired",
			c:    DeviceAuthCanceledExpired,
			want: DeviceAuthStateExpired,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.State(); got != tt.want {
				t.Errorf("DeviceAuthCanceled.State() = %v, want %v", got, tt.want)
			}
		})
	}
}
