package domain

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// stableNow returns a timestamp which is not in the last second of a TOTP period.
// VerifyTOTP validates against the wall clock, so this guarantees the period
// does not tick over between generating a code and validating it,
// which would shift the expected validity window of the skew based test cases.
func stableNow(t *testing.T) time.Time {
	t.Helper()
	for {
		now := time.Now()
		if now.Unix()%TOTPPeriod != TOTPPeriod-1 {
			return now
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func TestVerifyTOTP(t *testing.T) {
	cryptoAlg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	key, err := NewTOTPKey("example.com", "user1")
	require.NoError(t, err)
	secret, err := crypto.Encrypt([]byte(key.Secret()), cryptoAlg)
	require.NoError(t, err)

	unknownKeySecret := &crypto.CryptoValue{
		CryptoType: secret.CryptoType,
		Algorithm:  secret.Algorithm,
		KeyID:      "unknown",
		Crypted:    secret.Crypted,
	}

	base := stableNow(t)
	codeAt := func(ts time.Time) string {
		t.Helper()
		code, err := totp.GenerateCode(key.Secret(), ts)
		require.NoError(t, err)
		return code
	}

	type args struct {
		code   string
		secret *crypto.CryptoValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "current period",
			args: args{
				code:   codeAt(base),
				secret: secret,
			},
		},
		{
			name: "previous period, within skew",
			args: args{
				code:   codeAt(base.Add(-TOTPPeriod * time.Second)),
				secret: secret,
			},
		},
		{
			name: "next period, within skew",
			args: args{
				code:   codeAt(base.Add(TOTPPeriod * time.Second)),
				secret: secret,
			},
		},
		{
			name: "two periods in the past, outside skew",
			args: args{
				code:   codeAt(base.Add(-2 * TOTPPeriod * time.Second)),
				secret: secret,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "two periods in the future, outside skew",
			args: args{
				code:   codeAt(base.Add(2 * TOTPPeriod * time.Second)),
				secret: secret,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "malformed code",
			args: args{
				code:   "foobar",
				secret: secret,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "empty code",
			args: args{
				code:   "",
				secret: secret,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "EVENT-8isk2", "Errors.User.MFA.OTP.InvalidCode"),
		},
		{
			name: "missing secret",
			args: args{
				code:   codeAt(base),
				secret: nil,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "CRYPT-mNsQwe", "input value cannot be nil"),
		},
		{
			name: "unknown encryption key",
			args: args{
				code:   codeAt(base),
				secret: unknownKeySecret,
			},
			wantErr: zerrors.ThrowInvalidArgument(nil, "CRYPT-Kq12vn", "value was encrypted with a different key"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyTOTP(tt.args.code, tt.args.secret, cryptoAlg)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTOTPHistory_AddRecent(t *testing.T) {
	// start of a window which is still open.
	start := time.Now().Add(-checkDuration)

	hash1 := crypto.NewHMACValue("123456")
	hash2 := crypto.NewHMACValue("654321")

	type args struct {
		ts    time.Time
		value *crypto.HMACValue
	}
	tests := []struct {
		name       string
		history    *TOTPHistory
		args       args
		wantValues []*crypto.HMACValue
	}{
		{
			name:    "zero start, timestamp inside the window",
			history: new(TOTPHistory),
			args: args{
				ts:    time.Now(),
				value: hash1,
			},
			wantValues: []*crypto.HMACValue{hash1},
		},
		{
			name:    "zero start, timestamp outside the window",
			history: new(TOTPHistory),
			args: args{
				ts:    time.Now().Add(-checkDuration - time.Minute),
				value: hash1,
			},
			wantValues: nil,
		},
		{
			name:    "timestamp exactly at the start of the window",
			history: &TOTPHistory{start: start},
			args: args{
				ts:    start,
				value: hash1,
			},
			wantValues: nil,
		},
		{
			name:    "appended after existing codes",
			history: &TOTPHistory{start: start, recentValues: []*crypto.HMACValue{hash1}},
			args: args{
				ts:    time.Now(),
				value: hash2,
			},
			wantValues: []*crypto.HMACValue{hash1, hash2},
		},
		{
			name:    "same value added twice",
			history: &TOTPHistory{start: start, recentValues: []*crypto.HMACValue{hash1}},
			args: args{
				ts:    time.Now(),
				value: hash1,
			},
			wantValues: []*crypto.HMACValue{hash1, hash1},
		},
		{
			name:    "nil value of a legacy event",
			history: &TOTPHistory{start: start},
			args: args{
				ts:    time.Now(),
				value: nil,
			},
			wantValues: nil,
		},
		{
			name:    "nil value of a legacy event, zero start",
			history: new(TOTPHistory),
			args: args{
				ts:    time.Now(),
				value: nil,
			},
			wantValues: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startWasZero := tt.history.start.IsZero()
			startBefore := tt.history.start

			tt.history.AddRecent(tt.args.ts, tt.args.value)

			assert.Equal(t, tt.wantValues, tt.history.recentValues)
			switch {
			case tt.args.value == nil:
				assert.Equal(t, startBefore, tt.history.start, "a nil value must not set the start")
			case startWasZero:
				assert.WithinDuration(t, time.Now().Add(-checkDuration), tt.history.start, time.Second)
			default:
				assert.Equal(t, startBefore, tt.history.start, "an already set start must not be recomputed")
			}
		})
	}
}

func TestTOTPHistory_CheckReuse(t *testing.T) {
	recentValues := []*crypto.HMACValue{
		crypto.NewHMACValue("123456"),
		crypto.NewHMACValue("654321"),
	}

	tests := []struct {
		name    string
		history *TOTPHistory
		code    string
		wantErr error
	}{
		{
			name:    "empty history",
			history: new(TOTPHistory),
			code:    "123456",
		},
		{
			name:    "code not used recently",
			history: &TOTPHistory{recentValues: recentValues},
			code:    "111111",
		},
		{
			name:    "first code reused",
			history: &TOTPHistory{recentValues: recentValues},
			code:    "123456",
			wantErr: zerrors.ThrowInvalidArgument(nil, "TOTP-Auw0a", "Errors.User.MFA.OTP.Reused"),
		},
		{
			name:    "last code reused",
			history: &TOTPHistory{recentValues: recentValues},
			code:    "654321",
			wantErr: zerrors.ThrowInvalidArgument(nil, "TOTP-Auw0a", "Errors.User.MFA.OTP.Reused"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.history.CheckReuse(tt.code)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestTOTPHistory_addAndCheck(t *testing.T) {
	history := new(TOTPHistory)
	history.AddRecent(time.Now(), crypto.NewHMACValue("123456"))
	history.AddRecent(time.Now().Add(-checkDuration-time.Minute), crypto.NewHMACValue("654321"))
	history.AddRecent(time.Now(), nil)

	require.Error(t, history.CheckReuse("123456"), "a code used inside the window must be reported as reused")
	require.NoError(t, history.CheckReuse("654321"), "a code used before the window must not be reported as reused")
	require.NoError(t, history.CheckReuse("111111"))
}
