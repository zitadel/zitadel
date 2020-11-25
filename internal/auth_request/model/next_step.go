package model

type NextStep interface {
	Type() NextStepType
}

type NextStepType int32

const (
	NextStepUnspecified NextStepType = iota
	NextStepLogin
	NextStepUserSelection
	NextStepInitUser
	NextStepPassword
	NextStepChangePassword
	NextStepInitPassword
	NextStepVerifyEmail
	NextStepMfaPrompt
	NextStepMfaVerify
	NextStepRedirectToCallback
	NextStepChangeUsername
	NextStepLinkUsers
	NextStepExternalNotFoundOption
	NextStepExternalLogin
	NextStepGrantRequired
	NextStepPasswordless
)

type UserSessionState int32

const (
	UserSessionStateActive UserSessionState = iota
	UserSessionStateTerminated
)

type LoginStep struct{}

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
	DisplayName      string
	LoginName        string
	UserSessionState UserSessionState
}

type InitUserStep struct {
	PasswordSet bool
}

type ExternalNotFoundOptionStep struct{}

func (s *ExternalNotFoundOptionStep) Type() NextStepType {
	return NextStepExternalNotFoundOption
}

func (s *InitUserStep) Type() NextStepType {
	return NextStepInitUser
}

type PasswordStep struct{}

func (s *PasswordStep) Type() NextStepType {
	return NextStepPassword
}

type ExternalLoginStep struct {
	SelectedIDPConfigID string
}

func (s *ExternalLoginStep) Type() NextStepType {
	return NextStepExternalLogin
}

type PasswordlessStep struct{}

func (s *PasswordlessStep) Type() NextStepType {
	return NextStepPasswordless
}

type ChangePasswordStep struct{}

func (s *ChangePasswordStep) Type() NextStepType {
	return NextStepChangePassword
}

type InitPasswordStep struct{}

func (s *InitPasswordStep) Type() NextStepType {
	return NextStepInitPassword
}

type ChangeUsernameStep struct{}

func (s *ChangeUsernameStep) Type() NextStepType {
	return NextStepChangeUsername
}

type VerifyEMailStep struct{}

func (s *VerifyEMailStep) Type() NextStepType {
	return NextStepVerifyEmail
}

type MfaPromptStep struct {
	Required     bool
	MfaProviders []MFAType
}

func (s *MfaPromptStep) Type() NextStepType {
	return NextStepMfaPrompt
}

type MfaVerificationStep struct {
	MfaProviders []MFAType
}

func (s *MfaVerificationStep) Type() NextStepType {
	return NextStepMfaVerify
}

type LinkUsersStep struct{}

func (s *LinkUsersStep) Type() NextStepType {
	return NextStepLinkUsers
}

type GrantRequiredStep struct{}

func (s *GrantRequiredStep) Type() NextStepType {
	return NextStepGrantRequired
}

type RedirectToCallbackStep struct{}

func (s *RedirectToCallbackStep) Type() NextStepType {
	return NextStepRedirectToCallback
}

type MFAType int

const (
	MFATypeOTP MFAType = iota
	MFATypeU2F
	MFATypeU2FUserVerification
)

type MFALevel int

const (
	MFALevelNotSetUp MFALevel = iota
	MFALevelSecondFactor
	MFALevelMultiFactor
	MFALevelMultiFactorCertified
)
