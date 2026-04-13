package domain

import "time"

type SessionCommandOption interface {
	CreateSessionOption
	SetSessionOption
}

func WithSessionMetadata(metadata ...*SessionMetadata) SessionCommandOption {
	return sessionMetadataOption{metadata: metadata}
}

type sessionMetadataOption struct {
	metadata []*SessionMetadata
}

// ApplyOnSetSessionCommand implements [SetSessionOption].
func (opt sessionMetadataOption) ApplyOnSetSessionCommand(cmd *SetSessionCommand) {
	cmd.metadata = append(cmd.metadata, opt.metadata...)
}

// ApplyOnCreateSessionCommand implements [CreateSessionOption].
func (opt sessionMetadataOption) ApplyOnCreateSessionCommand(cmd *CreateSessionCommand) {
	cmd.session.Metadata = append(cmd.session.Metadata, opt.metadata...)
}

var (
	_ CreateSessionOption = sessionMetadataOption{}
	_ SetSessionOption    = sessionMetadataOption{}
)

func WithSessionUserAgent(userAgent *SessionUserAgent) CreateSessionOption {
	return sessionUserAgentOption{userAgent: userAgent}
}

type sessionUserAgentOption struct {
	userAgent *SessionUserAgent
}

func (opt sessionUserAgentOption) ApplyOnCreateSessionCommand(cmd *CreateSessionCommand) {
	cmd.session.UserAgent = opt.userAgent
}

var _ CreateSessionOption = sessionUserAgentOption{}

func WithSessionLifetime(lifetime time.Duration) SessionCommandOption {
	return sessionLifetimeOption{lifetime: lifetime}
}

type sessionLifetimeOption struct {
	lifetime time.Duration
}

// ApplyOnSetSessionCommand implements [SetSessionOption].
func (opt sessionLifetimeOption) ApplyOnSetSessionCommand(cmd *SetSessionCommand) {
	cmd.lifetime = &opt.lifetime
}

// ApplyOnCreateSessionCommand implements [CreateSessionOption].
func (opt sessionLifetimeOption) ApplyOnCreateSessionCommand(cmd *CreateSessionCommand) {
	cmd.session.Lifetime = opt.lifetime
}

var (
	_ CreateSessionOption = sessionLifetimeOption{}
	_ SetSessionOption    = sessionLifetimeOption{}
)
