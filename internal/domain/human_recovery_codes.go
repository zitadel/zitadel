package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zitadel/passwap"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type RecoveryCode struct {
	HashedCode string
	CheckAt    time.Time
}

type HumanRecoveryCodes struct {
	*ObjectDetails

	Codes []RecoveryCode
}

func RecoveryCodesFromRaw(codes []string, hasher *crypto.Hasher) ([]string, error) {
	if len(codes) == 0 {
		return nil, zerrors.ThrowInvalidArgument(nil, "DOMAIN-vee93", "Errors.User.MFA.RecoveryCodes.InvalidCount")
	}

	hashedCodes := make([]string, len(codes))
	for i, code := range codes {
		hashedCode, err := hasher.Hash(code)
		if err != nil {
			return nil, err
		}
		hashedCodes[i] = hashedCode
	}

	return hashedCodes, nil
}

func GenerateRecoveryCodes(count int, hasher *crypto.Hasher) ([]string, []string, error) {
	if count <= 0 {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "DOMAIN-7rp5j", "Errors.User.MFA.RecoveryCodes.InvalidCount")
	}

	hashedCodes, rawCodes := make([]string, count), make([]string, count)

	for i := 0; i < count; i++ {
		rawCode := uuid.New().String()
		hashedCode, err := hasher.Hash(rawCode)
		if err != nil {
			return nil, nil, err
		}
		hashedCodes[i] = hashedCode
		rawCodes[i] = rawCode
	}

	return hashedCodes, rawCodes, nil
}

func ValidateRecoveryCode(code string, recoveryCodes *HumanRecoveryCodes, hasher *crypto.Hasher) (valid bool, index int, err error) {
	index = -1

	if code == "" {
		return false, index, zerrors.ThrowInvalidArgument(nil, "DOMAIN-9xrr0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
	}

	if recoveryCodes == nil {
		return false, index, zerrors.ThrowInvalidArgument(nil, "DOMAIN-17bgk", "Errors.User.MFA.RecoveryCodes.InvalidCode")
	}

	for i, recoveryCode := range recoveryCodes.Codes {
		if !recoveryCode.CheckAt.IsZero() {
			continue
		}

		hashedCode, err := hasher.Hash(code)
		if err != nil {
			return false, index, err
		}

		// Ignoring the updated hash value, if any, as the code can only be checked once regardless
		_, verifyErr := hasher.Verify(recoveryCode.HashedCode, hashedCode)
		if verifyErr != nil {
			if errors.Is(verifyErr, passwap.ErrPasswordMismatch) {
				continue
			} else {
				return false, index, zerrors.ThrowInvalidArgument(verifyErr, "DOMAIN-ecn95", "Errors.User.MFA.RecoveryCodes.InvalidCode")
			}
		}
		return true, i, nil
	}

	return false, index, zerrors.ThrowInvalidArgument(nil, "DOMAIN-6uvh0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
}
