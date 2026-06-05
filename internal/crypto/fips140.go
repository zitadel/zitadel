package crypto

import (
	"crypto/fips140"
	"fmt"
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/instrumentation/logging"
)

// fipsPBKDF2MinIterations is the minimum PBKDF2 iteration count per NIST SP 800-132 §5.2.
const fipsPBKDF2MinIterations uint32 = 1000

// IsFIPSCompliant reports whether the algorithm is approved for password
// hashing/verification under FIPS 140-3 (NIST SP 800-132).
func (n HashName) IsFIPSCompliant() bool {
	return n == HashNamePBKDF2
}

// IsFIPSCompliant reports whether the hash mode is an approved PBKDF2 PRF
// under FIPS 140-3 (SHA-2 family; SHA-1/SHA-224 rejected for new use).
func (m HashMode) IsFIPSCompliant() bool {
	return m == HashModeSHA256 || m == HashModeSHA384 || m == HashModeSHA512
}

func nonFIPSVerifiersConfigured(verifiers []HashName) []HashName {
	var found []HashName
	for _, v := range verifiers {
		if !v.IsFIPSCompliant() {
			found = append(found, v)
		}
	}
	return found
}

func validateFIPSPBKDF2Hasher(c HasherConfig) error {
	p, hashMode, err := c.pbkdf2Params()
	if err != nil {
		return fmt.Errorf("decode pbkdf2 hasher for FIPS validation: %w", err)
	}
	if !hashMode.IsFIPSCompliant() {
		return fmt.Errorf(
			"application cannot start in uncertified cryptographic state: pbkdf2 hash mode %q is not FIPS 140-3 compliant while FIPS mode is enabled",
			hashMode,
		)
	}
	if p.Rounds < fipsPBKDF2MinIterations {
		return fmt.Errorf(
			"application cannot start in uncertified cryptographic state: pbkdf2 iteration count %d is below the FIPS minimum of %d while FIPS mode is enabled",
			p.Rounds, fipsPBKDF2MinIterations,
		)
	}
	return nil
}

func (c *HashConfig) validateFIPS140() error {
	if !fips140.Enabled() {
		return nil
	}

	alg := c.Hasher.Algorithm
	if !alg.IsFIPSCompliant() {
		return fmt.Errorf(
			"application cannot start in uncertified cryptographic state: password hasher algorithm %q is not FIPS 140-3 compliant while FIPS mode is enabled",
			alg,
		)
	}

	if alg == HashNamePBKDF2 {
		if err := validateFIPSPBKDF2Hasher(c.Hasher); err != nil {
			return err
		}
	}

	if legacy := nonFIPSVerifiersConfigured(c.Verifiers); len(legacy) > 0 {
		logging.New(logging.StreamRuntime).Warn(
			"Non-FIPS compliant password verifiers are active for migration. This instance is temporarily non-compliant until these verifiers are disabled",
			slog.Any("verifiers", legacy),
		)
	}

	if c.Limits.PBKDF2.MinRounds < fipsPBKDF2MinIterations {
		logging.New(logging.StreamRuntime).Warn(
			"PBKDF2 MinRounds is below the FIPS 140-3 minimum iteration count; imported hashes may use non-compliant cost parameters until limits are raised",
			slog.Uint64("min_rounds", uint64(c.Limits.PBKDF2.MinRounds)),
			slog.Uint64("fips_minimum", uint64(fipsPBKDF2MinIterations)),
		)
	}

	return nil
}
