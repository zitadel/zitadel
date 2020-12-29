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
	NextStepMFAPrompt
	NextStepMFAVerify
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
	UserSessionStateInitiated
	UserSessionStateDeleted //TODO: implement delete => special case of terminated
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
	UserID            string
	DisplayName       string
	LoginName         string
	UserSessionState  UserSessionState
	SelectionPossible bool
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

type MFAPromptStep struct {
	Required     bool
	MFAProviders []MFAType
}

func (s *MFAPromptStep) Type() NextStepType {
	return NextStepMFAPrompt
}

type MFAVerificationStep struct {
	MFAProviders []MFAType
}

func (s *MFAVerificationStep) Type() NextStepType {
	return NextStepMFAVerify
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
