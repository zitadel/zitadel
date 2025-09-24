package domain

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanRecoveryCodes struct {
	*ObjectDetails

	Codes []string
}

type ImportHumanRecoveryCode struct {
	RawCode    string
	HashedCode string
}

func HashRecoveryCodesIfNeeded(codes []ImportHumanRecoveryCode, hasher *crypto.Hasher) ([]string, error) {
	hashedCodes := make([]string, len(codes))
	for i, code := range codes {
		hashed, err := HashRecoveryCodeIfNeeded(code, hasher)
		if err != nil {
			return nil, err
		}
		hashedCodes[i] = hashed
	}
	return hashedCodes, nil
}

func HashRecoveryCodeIfNeeded(code ImportHumanRecoveryCode, hasher *crypto.Hasher) (string, error) {
	if code.RawCode != "" {
		return hasher.Hash(code.RawCode)
	}
	if !hasher.EncodingSupported(code.HashedCode) {
		return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-JDk4t", "Errors.Hash.NotSupported")
	}
	return code.HashedCode, nil
}

func GenerateRecoveryCodes(count int, config RecoveryCodesConfig, hasher *crypto.Hasher) ([]string, []string, error) {
	hashedCodes, rawCodes := make([]string, count), make([]string, count)

	for i := range count {
		rawCode, err := makeRawCode(config)
		if err != nil {
			return nil, nil, err
		}
		hashedCode, err := hasher.Hash(rawCode)
		if err != nil {
			return nil, nil, err
		}
		hashedCodes[i], rawCodes[i] = hashedCode, rawCode
	}

	return hashedCodes, rawCodes, nil
}

func makeRawCode(config RecoveryCodesConfig) (string, error) {
	switch config.Format {
	case RecoveryCodeFormatAlphanumeric:
		return generateAlphanumericCode(config.Length, config.WithHyphen)
	case RecoveryCodeFormatUUID:
		code := uuid.New().String()
		if !config.WithHyphen {
			code = strings.ReplaceAll(code, "-", "")
		}
		return code, nil
	default:
		return "", zerrors.ThrowError(nil, "DOMAIN-g52kn", "Errors.User.MFA.RecoveryCodes.ConfigInvalid")
	}
}

func generateAlphanumericCode(length int, withHyphen bool) (string, error) {
	if length <= 0 {
		return "", zerrors.ThrowError(nil, "DOMAIN-68mvq", "Errors.User.MFA.RecoveryCodes.ConfigInvalid")
	}

	// lower-cased base32 character set https://www.crockford.com/base32.html
	chars := []rune("0123456789abcdefghjkmnpqrstvwxyz")

	code, err := crypto.GenerateRandomString(uint(length), chars)
	if err != nil {
		return "", err
	}

	if withHyphen && length > 2 {
		mid := length / 2
		return code[:mid] + "-" + code[mid:], nil
	}

	return code, nil
}

func ValidateRecoveryCode(ctx context.Context, code string, recoveryCodes *HumanRecoveryCodes, hasher *crypto.Hasher) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if code == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-9xrr0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
	}

	for _, recoveryCode := range recoveryCodes.Codes {
		_, spanCodeComparison := tracing.NewNamedSpan(ctx, "passwap.Verify")
		_, verifyErr := hasher.Verify(recoveryCode, code)
		spanCodeComparison.EndWithError(verifyErr)
		if verifyErr != nil {
			continue
		}
		return recoveryCode, nil
	}

	return "", zerrors.ThrowInvalidArgument(nil, "DOMAIN-6uvh0", "Errors.User.MFA.RecoveryCodes.InvalidCode")
}
