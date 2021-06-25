package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	caos_errs "github.com/caos/zitadel/internal/errors"
)

func TestLabelPolicyPrimaryColorValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty primary, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyBackgroundColorValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty background, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyWarnColorValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty warn, valid",
			args: args{
				policy: &LabelPolicy{WarnColor: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{WarnColor: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{WarnColor: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{WarnColor: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{WarnColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{WarnColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyFontColorValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty font, valid",
			args: args{
				policy: &LabelPolicy{FontColor: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{FontColor: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{FontColor: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{FontColor: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{FontColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{FontColor: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyPrimaryColorDarkValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty primary dark, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{PrimaryColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyBackgroundColorDarkValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty background dark, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{BackgroundColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyWarnColorDarkValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty warn dark, valid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{WarnColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestLabelPolicyFontColorDarkValid(t *testing.T) {
	type args struct {
		policy *LabelPolicy
	}
	tests := []struct {
		name string
		args args
		err  func(error) bool
	}{
		{
			name: "empty font dark, valid",
			args: args{
				policy: &LabelPolicy{FontColorDark: ""},
			},
		},
		{
			name: "color code with 6 characters, valid",
			args: args{
				policy: &LabelPolicy{FontColorDark: "#ffffff"},
			},
		},
		{
			name: "color code with 3 characters, valid",
			args: args{
				policy: &LabelPolicy{FontColorDark: "#000"},
			},
		},
		{
			name: "color code with wrong characters, invalid",
			args: args{
				policy: &LabelPolicy{FontColorDark: "#0f9wfm"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with wrong count of characters, invalid",
			args: args{
				policy: &LabelPolicy{FontColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
		{
			name: "color code with no #, invalid",
			args: args{
				policy: &LabelPolicy{FontColorDark: "#00"},
			},
			err: caos_errs.IsErrorInvalidArgument,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.policy.IsValid()
			if tt.err == nil {
				assert.NoError(t, err)
			}
			if tt.err != nil && !tt.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
