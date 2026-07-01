package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/zitadel/zitadel/internal/api/action/otp"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/execution"
	exec_repo "github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// preOTPSMSCodeHookFromTargets is the default implementation of the preotpsmscode
// hook invocation. It queries registered execution targets and calls them with a
// context snapshot. Returns (nil, nil) when no targets are registered or the
// response is unusable, letting the caller fall back to the standard generation path.
func (c *Commands) preOTPSMSCodeHookFromTargets(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPSMSCodeResponse, error) {
	fnID := exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreOTPSMSCode.LocalizationKey())
	targets := execution.QueryExecutionTargetsForFunction(ctx, fnID)
	if len(targets) == 0 {
		return nil, nil
	}

	phoneWM, err := c.phoneWriteModelByID(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	// OTP SMS enrollment requires a verified phone (AddHumanOTPSMS precondition),
	// so phoneWM.Phone is the verified number at this point.
	ctxInfo := &otp.PreOTPSMSCodeContext{
		FunctionName:         domain.ActionFunctionPreOTPSMSCode.LocalizationKey(),
		RecipientPhoneNumber: string(phoneWM.Phone),
		GeneratorConfig:      publicGeneratorConfigFrom(effectiveConfig),
	}

	resp, err := execution.CallTargets(ctx, targets, ctxInfo, c.targetEncryption, c.GetActiveSigningWebKey, c.ActionsV2DenyList)
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*otp.PreOTPSMSCodeResponse)
	if !ok {
		return nil, nil
	}
	return response, nil
}

// newPhoneCodeWithHook wraps newPhoneCode so that, when AllowOTPCodeOverride is
// enabled and no Twilio Verify provider owns the code lifecycle, the
// preotpsmscode hook can supply an override code or adjust generator parameters
// before local code generation.
func (c *Commands) newPhoneCodeWithHook(userID, resourceOwner string) encryptedCodeGeneratorWithDefaultFunc {
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, secretGeneratorType domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		if !authz.GetFeatures(ctx).AllowOTPCodeOverride {
			slog.DebugContext(ctx, "preotpsmscode hook skipped: AllowOTPCodeOverride disabled", "userID", userID)
			return c.newPhoneCode(ctx, filter, secretGeneratorType, alg, defaultConfig)
		}
		if c.preOTPSMSCodeHook == nil {
			slog.DebugContext(ctx, "preotpsmscode hook skipped: no hook implementation", "userID", userID)
			return c.newPhoneCode(ctx, filter, secretGeneratorType, alg, defaultConfig)
		}

		externalID, err := c.activeSMSProvider(ctx)
		if err != nil {
			return nil, "", err
		}
		if externalID != "" {
			// Twilio Verify owns the code lifecycle when configured; skip the hook.
			slog.DebugContext(ctx, "preotpsmscode hook skipped: external SMS provider owns code", "userID", userID, "provider", externalID)
			return nil, externalID, nil
		}

		effectiveConfig, err := cryptoGeneratorConfigWithDefault(ctx, filter, secretGeneratorType, defaultConfig)
		if err != nil {
			return nil, "", err
		}

		hookResp, err := c.preOTPSMSCodeHook(ctx, userID, resourceOwner, effectiveConfig)
		if err != nil {
			return nil, "", err
		}

		if hookResp != nil && hookResp.Code != nil {
			code, err := encryptOverriddenOTPCode(*hookResp.Code, alg, overrideExpiry(effectiveConfig.Expiry, hookResp.Expiry))
			if err != nil {
				return nil, "", err
			}
			return code, "", nil
		}

		codeConfig := effectiveConfig
		if hookResp != nil && (hookResp.Generate != nil || hookResp.Expiry != nil) {
			codeConfig = applyGenerationOverrides(effectiveConfig, hookResp.Generate, hookResp.Expiry)
			if err := validateGenerationConfig(codeConfig); err != nil {
				return nil, "", err
			}
		}
		crypted, plain, err := crypto.NewCode(crypto.NewEncryptionGenerator(*codeConfig, alg))
		if err != nil {
			return nil, "", err
		}
		return &EncryptedCode{Crypted: crypted, Plain: plain, Expiry: codeConfig.Expiry}, "", nil
	}
}

// validateGenerationConfig rejects configs that would cause crypto.NewCode to
// produce an empty code or panic: non-zero length with no character class set.
// This guards the new override path — pre-existing instance configs are trusted.
func validateGenerationConfig(cfg *crypto.GeneratorConfig) error {
	if cfg.Length == 0 {
		return nil
	}
	if !cfg.IncludeLowerLetters && !cfg.IncludeUpperLetters && !cfg.IncludeDigits && !cfg.IncludeSymbols {
		return zerrors.ThrowPreconditionFailed(nil, "ACTION-w4n9p", "Errors.Execution.Invalid")
	}
	return nil
}

// preOTPEmailCodeHookFromTargets is the default implementation of the preotpemailcode
// hook invocation. Mirrors preOTPSMSCodeHookFromTargets for the email channel.
func (c *Commands) preOTPEmailCodeHookFromTargets(ctx context.Context, userID, resourceOwner string, effectiveConfig *crypto.GeneratorConfig) (*otp.PreOTPEmailCodeResponse, error) {
	fnID := exec_repo.ID(domain.ExecutionTypeFunction, domain.ActionFunctionPreOTPEmailCode.LocalizationKey())
	targets := execution.QueryExecutionTargetsForFunction(ctx, fnID)
	if len(targets) == 0 {
		return nil, nil
	}

	emailWM, err := c.emailWriteModel(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	// OTP Email enrollment requires a verified email (AddHumanOTPEmail precondition),
	// so emailWM.Email is the verified address at this point.
	ctxInfo := &otp.PreOTPEmailCodeContext{
		FunctionName:          domain.ActionFunctionPreOTPEmailCode.LocalizationKey(),
		RecipientEmailAddress: string(emailWM.Email),
		GeneratorConfig:       publicGeneratorConfigFrom(effectiveConfig),
	}

	resp, err := execution.CallTargets(ctx, targets, ctxInfo, c.targetEncryption, c.GetActiveSigningWebKey, c.ActionsV2DenyList)
	if err != nil {
		return nil, err
	}
	response, ok := resp.(*otp.PreOTPEmailCodeResponse)
	if !ok {
		return nil, nil
	}
	return response, nil
}

// newEmailCodeWithHook wraps the default email code generator so that, when
// AllowOTPCodeOverride is enabled, the preotpemailcode hook can supply an
// override code or adjust generator parameters before local code generation.
// Unlike the SMS wrapper there is no external-provider bypass — email has no
// provider-level CodeGenerator equivalent to Twilio Verify.
func (c *Commands) newEmailCodeWithHook(userID, resourceOwner string) encryptedCodeGeneratorWithDefaultFunc {
	base := func(ctx context.Context, filter preparation.FilterToQueryReducer, secretGeneratorType domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		code, err := c.newEncryptedCodeWithDefault(ctx, filter, secretGeneratorType, alg, defaultConfig)
		return code, "", err
	}
	return func(ctx context.Context, filter preparation.FilterToQueryReducer, secretGeneratorType domain.SecretGeneratorType, alg crypto.EncryptionAlgorithm, defaultConfig *crypto.GeneratorConfig) (*EncryptedCode, string, error) {
		if !authz.GetFeatures(ctx).AllowOTPCodeOverride {
			slog.DebugContext(ctx, "preotpemailcode hook skipped: AllowOTPCodeOverride disabled", "userID", userID)
			return base(ctx, filter, secretGeneratorType, alg, defaultConfig)
		}
		if c.preOTPEmailCodeHook == nil {
			slog.DebugContext(ctx, "preotpemailcode hook skipped: no hook implementation", "userID", userID)
			return base(ctx, filter, secretGeneratorType, alg, defaultConfig)
		}

		effectiveConfig, err := cryptoGeneratorConfigWithDefault(ctx, filter, secretGeneratorType, defaultConfig)
		if err != nil {
			return nil, "", err
		}

		hookResp, err := c.preOTPEmailCodeHook(ctx, userID, resourceOwner, effectiveConfig)
		if err != nil {
			return nil, "", err
		}

		if hookResp != nil && hookResp.Code != nil {
			code, err := encryptOverriddenOTPCode(*hookResp.Code, alg, overrideExpiry(effectiveConfig.Expiry, hookResp.Expiry))
			if err != nil {
				return nil, "", err
			}
			return code, "", nil
		}

		codeConfig := effectiveConfig
		if hookResp != nil && (hookResp.Generate != nil || hookResp.Expiry != nil) {
			codeConfig = applyGenerationOverrides(effectiveConfig, hookResp.Generate, hookResp.Expiry)
			if err := validateGenerationConfig(codeConfig); err != nil {
				return nil, "", err
			}
		}
		crypted, plain, err := crypto.NewCode(crypto.NewEncryptionGenerator(*codeConfig, alg))
		if err != nil {
			return nil, "", err
		}
		return &EncryptedCode{Crypted: crypted, Plain: plain, Expiry: codeConfig.Expiry}, "", nil
	}
}

func publicGeneratorConfigFrom(cfg *crypto.GeneratorConfig) *otp.PublicGeneratorConfig {
	if cfg == nil {
		return nil
	}
	return &otp.PublicGeneratorConfig{
		Length:              uint32(cfg.Length),
		Expiry:              otp.Duration(cfg.Expiry),
		IncludeLowerLetters: cfg.IncludeLowerLetters,
		IncludeUpperLetters: cfg.IncludeUpperLetters,
		IncludeDigits:       cfg.IncludeDigits,
		IncludeSymbols:      cfg.IncludeSymbols,
	}
}

func applyGenerationOverrides(base *crypto.GeneratorConfig, gen *otp.GenerationOverrides, expiry *otp.Duration) *crypto.GeneratorConfig {
	out := *base
	if gen != nil {
		if gen.Length != nil {
			out.Length = uint(*gen.Length)
		}
		if gen.IncludeLowerLetters != nil {
			out.IncludeLowerLetters = *gen.IncludeLowerLetters
		}
		if gen.IncludeUpperLetters != nil {
			out.IncludeUpperLetters = *gen.IncludeUpperLetters
		}
		if gen.IncludeDigits != nil {
			out.IncludeDigits = *gen.IncludeDigits
		}
		if gen.IncludeSymbols != nil {
			out.IncludeSymbols = *gen.IncludeSymbols
		}
	}
	if expiry != nil {
		out.Expiry = time.Duration(*expiry)
	}
	return &out
}

func overrideExpiry(base time.Duration, override *otp.Duration) time.Duration {
	if override == nil {
		return base
	}
	return time.Duration(*override)
}

func encryptOverriddenOTPCode(plain string, alg crypto.EncryptionAlgorithm, expiry time.Duration) (*EncryptedCode, error) {
	crypted, err := crypto.Encrypt([]byte(plain), alg)
	if err != nil {
		return nil, err
	}
	return &EncryptedCode{Crypted: crypted, Plain: plain, Expiry: expiry}, nil
}
