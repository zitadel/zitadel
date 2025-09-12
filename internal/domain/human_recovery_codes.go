package domain

import (
	"github.com/google/uuid"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/id"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanRecoveryCodes struct {
	*ObjectDetails

	Codes []string
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

func GenerateRecoveryCodes(count int, format RecoveryCodeFormat, hasher *crypto.Hasher) ([]string, []string, error) {
	if count <= 0 {
		return nil, nil, zerrors.ThrowInvalidArgument(nil, "DOMAIN-7rp5j", "Errors.User.MFA.RecoveryCodes.InvalidCount")
	}

	hashedCodes, rawCodes := make([]string, count), make([]string, count)

	for i := 0; i < count; i++ {
		rawCode, err := makeRawCode(format)
		if err != nil {
			return nil, nil, err
		}
		hashedCode, err := hasher.Hash(rawCode)
		if err != nil {
			return nil, nil, err
		}
		hashedCodes[i] = hashedCode
		rawCodes[i] = rawCode
	}

	return hashedCodes, rawCodes, nil
}

func makeRawCode(format RecoveryCodeFormat) (string, error) {
	if format == RecoveryCodeFormatSonyFlake {
		return id.SonyFlakeGenerator().Next()
	}
	return uuid.New().String(), nil
}

func ValidateRecoveryCode(code string, recoveryCodes *HumanRecoveryCodes, hasher *crypto.Hasher) (string, error) {
	if code == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-9xrr0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
	}

	if recoveryCodes == nil {
		return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-17bgk", "Errors.User.MFA.RecoveryCodes.Missing")
	}

	for _, recoveryCode := range recoveryCodes.Codes {
		if _, verifyErr := hasher.Verify(recoveryCode, code); verifyErr != nil {
			continue
		}
		// return the hashed code to store in the success event data
		hashedCode, err := hasher.Hash(recoveryCode)
		if err != nil {
			return "", err
		}
		return hashedCode, nil
	}

	return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-6uvh0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
}
