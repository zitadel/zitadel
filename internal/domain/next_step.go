package domain

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
	NextStepPasswordlessRegistrationPrompt
	NextStepRegistration
	NextStepProjectRequired
	NextStepRedirectToExternalIDP
	NextStepLoginSucceeded
	NextStepVerifyInvite
)

type LoginStep struct{}

func (s *LoginStep) Type() NextStepType {
	return NextStepLogin
}

type RegistrationStep struct{}

func (s *RegistrationStep) Type() NextStepType {
	return NextStepRegistration
}

type SelectUserStep struct {
	Users []UserSelection
}

func (s *SelectUserStep) Type() NextStepType {
	return NextStepUserSelection
}

type UserSelection struct {
	UserID            string
	UserName          string
	DisplayName       string
	LoginName         string
	UserSessionState  UserSessionState
	SelectionPossible bool
	AvatarKey         string
	ResourceOwner     string
}

type UserSessionState int32

const (
	UserSessionStateActive UserSessionState = iota
	UserSessionStateTerminated
)

type RedirectToExternalIDPStep struct{}

func (s *RedirectToExternalIDPStep) Type() NextStepType {
	return NextStepRedirectToExternalIDP
}

type InitUserStep struct {
	PasswordSet bool
}

func (s *InitUserStep) Type() NextStepType {
	return NextStepInitUser
}

type ExternalNotFoundOptionStep struct{}

func (s *ExternalNotFoundOptionStep) Type() NextStepType {
	return NextStepExternalNotFoundOption
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

type PasswordlessStep struct {
	PasswordSet bool
}

func (s *PasswordlessStep) Type() NextStepType {
	return NextStepPasswordless
}

type PasswordlessRegistrationPromptStep struct{}

func (s *PasswordlessRegistrationPromptStep) Type() NextStepType {
	return NextStepPasswordlessRegistrationPrompt
}

type ChangePasswordStep struct {
	Expired bool
}

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

type VerifyEMailStep struct {
	InitPassword bool
}

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

type ProjectRequiredStep struct{}

func (s *ProjectRequiredStep) Type() NextStepType {
	return NextStepProjectRequired
}

type RedirectToCallbackStep struct{}

func (s *RedirectToCallbackStep) Type() NextStepType {
	return NextStepRedirectToCallback
}

type LoginSucceededStep struct{}

func (s *LoginSucceededStep) Type() NextStepType {
	return NextStepLoginSucceeded
}

type VerifyInviteStep struct{}

func (s *VerifyInviteStep) Type() NextStepType {
	return NextStepVerifyInvite
}
