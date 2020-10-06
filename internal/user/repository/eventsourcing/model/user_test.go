package model

import (
	"testing"

	"github.com/caos/zitadel/internal/user/model"
)

func TestAppendDeactivatedEvent(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append deactivate event",
			args: args{
				user: &User{},
			},
			result: &User{State: int32(model.UserStateInactive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.appendDeactivatedEvent()
			if tt.args.user.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendReactivatedEvent(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append reactivate event",
			args: args{
				user: &User{},
			},
			result: &User{State: int32(model.UserStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.appendReactivatedEvent()
			if tt.args.user.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendLockEvent(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append lock event",
			args: args{
				user: &User{},
			},
			result: &User{State: int32(model.UserStateLocked)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.appendLockedEvent()
			if tt.args.user.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}

func TestAppendUnlockEvent(t *testing.T) {
	type args struct {
		user *User
	}
	tests := []struct {
		name   string
		args   args
		result *User
	}{
		{
			name: "append unlock event",
			args: args{
				user: &User{},
			},
			result: &User{State: int32(model.UserStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.user.appendUnlockedEvent()
			if tt.args.user.State != tt.result.State {
				t.Errorf("got wrong result: expected: %v, actual: %v ", tt.result, tt.args.user)
			}
		})
	}
}
