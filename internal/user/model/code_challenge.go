package model

type OIDCCodeChallenge struct {
	Challenge string
	Method    OIDCCodeChallengeMethod
}

type OIDCCodeChallengeMethod int32

const (
	CodeChallengeMethodPlain OIDCCodeChallengeMethod = iota
	CodeChallengeMethodS256
)
