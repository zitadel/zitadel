package model

type possibleSteps []NextStep

func (p possibleSteps) add(step NextStep) {
	p = append(p, step)
}

type NextStep interface {
	Type() NextStepType
}

type NextStepType int32

const (
	NextStepUnspecified NextStepType = iota
	NextStepLogin
	NextStepUserSelection
	NextStepPassword
	NextStepChangePassword
	NextStepInitPassword
	NextStepVerifyEmail
	NextStepMfaPrompt
	NextStepMfaVerify
	NextStepRedirectToCallback
)

type UserSessionState int32

const (
	UserSessionStateActive UserSessionState = iota
	UserSessionStateTerminated
)

type MfaType int32

const (
	MfaTypeNone MfaType = iota
	MfaTypeOTP
	MFaTypeSMS
)

type LoginStep struct {
	ErrMsg string
}

func NewLoginStep(err error) *LoginStep {
	var msg string
	if err != nil {
		msg = err.Error()
	}
	return &LoginStep{ErrMsg: msg}
}

func (s *LoginStep) Type() NextStepType {
	return NextStepLogin
}

type SelectUserStep struct {
	Users []UserSelection
}

func (s *SelectUserStep) Type() NextStepType {
	return NextStepUserSelection
}

type UserSelection struct {
	UserID           string
	UserName         string
	UserSessionState UserSessionState
}

type PasswordStep struct {
	ErrMsg       string
	FailureCount uint16
}

func NewPasswordStep(failureCount uint16) *PasswordStep {
	step := &PasswordStep{FailureCount: failureCount}
	if failureCount > 0 {
		step.ErrMsg = "Password invalid"
	}
	return step
}

func (s *PasswordStep) Type() NextStepType {
	return NextStepPassword
}

type ChangePasswordStep struct {
}

func (s *ChangePasswordStep) Type() NextStepType {
	return NextStepChangePassword
}

type InitPasswordStep struct {
}

func (s *InitPasswordStep) Type() NextStepType {
	return NextStepInitPassword
}

type VerifyEMailStep struct {
}

func (s *VerifyEMailStep) Type() NextStepType {
	return NextStepVerifyEmail
}

type MfaPromptStep struct {
	Required     bool
	MfaProviders []MfaType
}

func (s *MfaPromptStep) Type() NextStepType {
	return NextStepMfaPrompt
}

type MfaVerificationStep struct {
	ErrMsg       string
	FailureCount uint16
	MfaProviders []MfaType
}

func NewMfaVerificationStep(failureCount uint16, mfaProviders []MfaType) *MfaVerificationStep {
	step := &MfaVerificationStep{
		FailureCount: failureCount,
		MfaProviders: mfaProviders,
	}
	if failureCount > 0 {
		step.ErrMsg = "Code invalid"
	}
	return step
}

func (s *MfaVerificationStep) Type() NextStepType {
	return NextStepMfaVerify
}

type RedirectToCallbackStep struct {
}

func (s *RedirectToCallbackStep) Type() NextStepType {
	return NextStepRedirectToCallback
}
