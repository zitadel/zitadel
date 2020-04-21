package model

import "github.com/caos/zitadel/internal/user_agent/model"

type OIDCCodeChallenge struct {
	Challenge string
	Method    int32
}

func OIDCCodeChallengeFromModel(challenge *model.OIDCCodeChallenge) *OIDCCodeChallenge {
	return &OIDCCodeChallenge{
		Challenge: challenge.Challenge,
		Method:    int32(challenge.Method),
	}
}

func OIDCCodeChallengeToModel(challenge *OIDCCodeChallenge) *model.OIDCCodeChallenge {
	return &model.OIDCCodeChallenge{
		Challenge: challenge.Challenge,
		Method:    model.OIDCCodeChallengeMethod(challenge.Method),
	}
}
