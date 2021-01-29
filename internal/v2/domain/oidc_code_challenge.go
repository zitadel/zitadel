package domain

type OIDCCodeChallenge struct {
	Challenge string
	Method    OIDCCodeChallengeMethod
}

func (c *OIDCCodeChallenge) IsValid() bool {
	return c.Challenge != ""
}

type OIDCCodeChallengeMethod int32

const (
	CodeChallengeMethodPlain OIDCCodeChallengeMethod = iota
	CodeChallengeMethodS256
)
