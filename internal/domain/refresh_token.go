package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func NewRefreshToken(userID, tokenID string, algorithm crypto.AuthAlgorithm) (string, error) {
	return RefreshToken(userID, tokenID, tokenID, algorithm)
}

func RefreshToken(userID, tokenID, token string, algorithm crypto.AuthAlgorithm) (string, error) {
	return algorithm.EncryptToken(userID + ":" + tokenID + ":" + token)
}

func FromRefreshToken(refreshToken string, algorithm crypto.AuthAlgorithm) (userID, tokenID, token string, err error) {
	decrypted, err := algorithm.DecryptToken(refreshToken)
	if err != nil {
		return "", "", "", zerrors.ThrowInvalidArgument(err, "DOMAIN-rie9A", "Errors.User.RefreshToken.Invalid")
	}
	split := strings.Split(decrypted, ":")
	if len(split) != 3 {
		return "", "", "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-Se8oh", "Errors.User.RefreshToken.Invalid")
	}
	return split[0], split[1], split[2], nil
}
