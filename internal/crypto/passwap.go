package crypto

import (
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/argon2"
	"github.com/zitadel/passwap/bcrypt"
	"github.com/zitadel/passwap/drupal7"
	"github.com/zitadel/passwap/md5"
	"github.com/zitadel/passwap/md5plain"
	"github.com/zitadel/passwap/md5salted"
	"github.com/zitadel/passwap/pbkdf2"
	"github.com/zitadel/passwap/phpass"
	"github.com/zitadel/passwap/scrypt"
	"github.com/zitadel/passwap/sha2"
	"github.com/zitadel/passwap/verifier"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type Hasher struct {
	*passwap.Swapper
	Prefixes     []string
	HexSupported bool
}

func (h *Hasher) EncodingSupported(encodedHash string) bool {
	for _, prefix := range h.Prefixes {
		if strings.HasPrefix(encodedHash, prefix) {
			return true
		}
	}
	if h.HexSupported {
		_, err := hex.DecodeString(encodedHash)
		if err == nil {
			return true
		}
	}
	return false
}

// ValidateEncodedHash checks that encoded is parseable by a configured verifier
// and that its cost parameters are within configured bounds.
// Use when accepting pre-encoded hashes from untrusted sources (import, create with hash).
// It intentionally returns password-scoped errors for compatibility with current callers.
func (h *Hasher) ValidateEncodedHash(encoded string) error {
	if encoded == "" {
		return nil
	}
	err := h.Validate(encoded)
	if err == nil {
		return nil
	}
	if errors.Is(err, passwap.ErrNoVerifier) {
		return zerrors.ThrowInvalidArgument(err, "CRYPT-2xK7d", "Errors.Hash.NotSupported")
	}
	var bounds *verifier.BoundsError
	if errors.As(err, &bounds) {
		return zerrors.ThrowInvalidArgument(err, "CRYPT-5uV9n", "Errors.User.Password.Invalid")
	}
	return zerrors.ThrowInvalidArgument(err, "CRYPT-r8Qm2", "Errors.User.Password.Invalid")
}

type HashName string

const (
	HashNameArgon2    HashName = "argon2"    // used for the common argon2 verifier
	HashNameArgon2i   HashName = "argon2i"   // hash only
	HashNameArgon2id  HashName = "argon2id"  // hash only
	HashNameBcrypt    HashName = "bcrypt"    // hash and verify
	HashNameMd5       HashName = "md5"       // verify only, as hashing with md5 is insecure and deprecated
	HashNameMd5Plain  HashName = "md5plain"  // verify only, as hashing with md5 is insecure and deprecated
	HashNameMd5Salted HashName = "md5salted" // verify only, as hashing with md5 is insecure and deprecated
	HashNamePHPass    HashName = "phpass"    // verify only, as hashing with md5 is insecure and deprecated
	HashNameSha2      HashName = "sha2"      // hash and verify
	HashNameScrypt    HashName = "scrypt"    // hash and verify
	HashNamePBKDF2    HashName = "pbkdf2"    // hash and verify
	HashNameDrupal7   HashName = "drupal7"   // verify only, Drupal 7 legacy hashes
)

type HashMode string

// HashMode defines a underlying [hash.Hash] implementation
// for algorithms like pbkdf2
const (
	HashModeSHA1   HashMode = "sha1"
	HashModeSHA224 HashMode = "sha224"
	HashModeSHA256 HashMode = "sha256"
	HashModeSHA384 HashMode = "sha384"
	HashModeSHA512 HashMode = "sha512"
)

type HashConfig struct {
	Verifiers []HashName
	Hasher    HasherConfig
	Limits    HashLimitsConfig
}

type HashLimitsConfig struct {
	Bcrypt  BcryptLimitsConfig
	Argon2  Argon2LimitsConfig
	Scrypt  ScryptLimitsConfig
	PBKDF2  PBKDF2LimitsConfig
	Sha2    Sha2LimitsConfig
	PHPass  PHPassLimitsConfig
	Drupal7 Drupal7LimitsConfig
}

type BcryptLimitsConfig struct {
	MinCost int
	MaxCost int
}

func (l BcryptLimitsConfig) validationOpts() *bcrypt.ValidationOpts {
	return &bcrypt.ValidationOpts{
		MinCost: l.MinCost,
		MaxCost: l.MaxCost,
	}
}

type Argon2LimitsConfig struct {
	MinTime    uint32
	MaxTime    uint32
	MinMemory  uint32
	MaxMemory  uint32
	MinThreads uint8
	MaxThreads uint8
}

func (l Argon2LimitsConfig) validationOpts() *argon2.ValidationOpts {
	return &argon2.ValidationOpts{
		MinTime:    l.MinTime,
		MaxTime:    l.MaxTime,
		MinMemory:  l.MinMemory,
		MaxMemory:  l.MaxMemory,
		MinThreads: l.MinThreads,
		MaxThreads: l.MaxThreads,
	}
}

type ScryptLimitsConfig struct {
	MinLN int
	MaxLN int
	MinR  int
	MaxR  int
	MinP  int
	MaxP  int
}

func (l ScryptLimitsConfig) validationOpts() *scrypt.ValidationOpts {
	return &scrypt.ValidationOpts{
		MinLN: l.MinLN,
		MaxLN: l.MaxLN,
		MinR:  l.MinR,
		MaxR:  l.MaxR,
		MinP:  l.MinP,
		MaxP:  l.MaxP,
	}
}

type PBKDF2LimitsConfig struct {
	MinRounds uint32
	MaxRounds uint32
}

func (l PBKDF2LimitsConfig) validationOpts() *pbkdf2.ValidationOpts {
	return &pbkdf2.ValidationOpts{
		MinRounds: l.MinRounds,
		MaxRounds: l.MaxRounds,
	}
}

type Sha2LimitsConfig struct {
	MinSha256Rounds int
	MaxSha256Rounds int
	MinSha512Rounds int
	MaxSha512Rounds int
}

func (l Sha2LimitsConfig) validationOpts() *sha2.ValidationOpts {
	return &sha2.ValidationOpts{
		MinSha256Rounds: l.MinSha256Rounds,
		MaxSha256Rounds: l.MaxSha256Rounds,
		MinSha512Rounds: l.MinSha512Rounds,
		MaxSha512Rounds: l.MaxSha512Rounds,
	}
}

type PHPassLimitsConfig struct {
	MinRounds int
	MaxRounds int
}

func (l PHPassLimitsConfig) validationOpts() *phpass.ValidationOpts {
	return &phpass.ValidationOpts{
		MinRounds: l.MinRounds,
		MaxRounds: l.MaxRounds,
	}
}

type Drupal7LimitsConfig struct {
	MinIterations int
	MaxIterations int
}

func (l Drupal7LimitsConfig) validationOpts() *drupal7.ValidationOpts {
	return &drupal7.ValidationOpts{
		MinIterations: l.MinIterations,
		MaxIterations: l.MaxIterations,
	}
}

func (c *HashConfig) NewHasher() (*Hasher, error) {
	verifiers, vPrefixes, err := c.buildVerifiers()
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "CRYPT-sahW9", "password hash config invalid")
	}
	hasher, hPrefixes, err := c.Hasher.buildHasher(c.Limits)
	if err != nil {
		return nil, zerrors.ThrowInvalidArgument(err, "CRYPT-Que4r", "password hash config invalid")
	}
	return &Hasher{
		Swapper:      passwap.NewSwapper(hasher, verifiers...),
		Prefixes:     append(hPrefixes, vPrefixes...),
		HexSupported: slices.Contains(c.Verifiers, HashNameMd5Plain),
	}, nil
}

type prefixVerifier struct {
	prefixes []string
	verifier verifier.Verifier
}

type verifierFactory func(limits HashLimitsConfig) prefixVerifier

var knowVerifiers = map[HashName]verifierFactory{
	HashNameArgon2: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			// verifier for both argon2i and argon2id.
			prefixes: []string{argon2.Prefix},
			verifier: argon2.NewVerifier(l.Argon2.validationOpts()),
		}
	},
	HashNameBcrypt: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{bcrypt.Prefix},
			verifier: bcrypt.NewVerifier(l.Bcrypt.validationOpts()),
		}
	},
	HashNameMd5: func(_ HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{md5.Prefix},
			verifier: md5.NewVerifier(),
		}
	},
	HashNameMd5Plain: func(_ HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			// hex encoded without identifier or prefix.
			prefixes: nil,
			verifier: md5plain.NewVerifier(),
		}
	},
	HashNameScrypt: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{scrypt.Prefix, scrypt.Prefix_Linux},
			verifier: scrypt.NewVerifier(l.Scrypt.validationOpts()),
		}
	},
	HashNamePBKDF2: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{pbkdf2.Prefix},
			verifier: pbkdf2.NewVerifier(l.PBKDF2.validationOpts()),
		}
	},
	HashNameMd5Salted: func(_ HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{md5salted.Prefix},
			verifier: md5salted.NewVerifier(),
		}
	},
	HashNameSha2: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{sha2.Sha256Identifier, sha2.Sha512Identifier},
			verifier: sha2.NewVerifier(l.Sha2.validationOpts()),
		}
	},
	HashNamePHPass: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{phpass.IdentifierP, phpass.IdentifierH},
			verifier: phpass.NewVerifier(l.PHPass.validationOpts()),
		}
	},
	HashNameDrupal7: func(l HashLimitsConfig) prefixVerifier {
		return prefixVerifier{
			prefixes: []string{drupal7.Identifier},
			verifier: drupal7.NewVerifier(l.Drupal7.validationOpts()),
		}
	},
}

func (c *HashConfig) buildVerifiers() (verifiers []verifier.Verifier, prefixes []string, err error) {
	verifiers = make([]verifier.Verifier, len(c.Verifiers))
	prefixes = make([]string, 0, len(c.Verifiers)+1)
	for i, name := range c.Verifiers {
		factory, ok := knowVerifiers[name]
		if !ok {
			return nil, nil, fmt.Errorf("invalid verifier %q", name)
		}
		v := factory(c.Limits)
		verifiers[i] = v.verifier
		prefixes = append(prefixes, v.prefixes...)
	}
	return verifiers, prefixes, nil
}

type HasherConfig struct {
	Algorithm HashName
	Params    map[string]any `mapstructure:",remain"`
}

func (c *HasherConfig) buildHasher(limits HashLimitsConfig) (hasher passwap.Hasher, prefixes []string, err error) {
	switch c.Algorithm {
	case HashNameArgon2i:
		return c.argon2i(limits)
	case HashNameArgon2id:
		return c.argon2id(limits)
	case HashNameBcrypt:
		return c.bcrypt(limits)
	case HashNameScrypt:
		return c.scrypt(limits)
	case HashNamePBKDF2:
		return c.pbkdf2(limits)
	case HashNameSha2:
		return c.sha2(limits)
	case "":
		return nil, nil, fmt.Errorf("missing hasher algorithm")
	case HashNameArgon2, HashNameMd5, HashNameMd5Plain, HashNameMd5Salted, HashNamePHPass, HashNameDrupal7:
		fallthrough
	default:
		return nil, nil, fmt.Errorf("invalid algorithm %q", c.Algorithm)
	}
}

// decodeParams uses a mapstructure decoder from the Params map to dst.
// The decoder fails when there are unused fields in dst.
// It uses weak input typing, to allow conversion of env strings to ints.
func (c *HasherConfig) decodeParams(dst any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused:      false,
		ErrorUnset:       true,
		WeaklyTypedInput: true,
		Result:           dst,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(c.Params)
}

// argon2Params decodes [HasherConfig.Params] into a [argon2.Params] used as defaults.
// p is passed a copy and therfore will not be modified.
func (c *HasherConfig) argon2Params(p argon2.Params) (argon2.Params, error) {
	var dst struct {
		Time    uint32 `mapstructure:"Time"`
		Memory  uint32 `mapstructure:"Memory"`
		Threads uint8  `mapstructure:"Threads"`
	}
	if err := c.decodeParams(&dst); err != nil {
		return argon2.Params{}, fmt.Errorf("decode argon2i params: %w", err)
	}
	p.Time = dst.Time
	p.Memory = dst.Memory
	p.Threads = dst.Threads
	return p, nil
}

func (c *HasherConfig) argon2i(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	p, err := c.argon2Params(argon2.RecommendedIParams)
	if err != nil {
		return nil, nil, err
	}
	return argon2.NewArgon2i(p, limits.Argon2.validationOpts()), []string{argon2.Prefix}, nil
}

func (c *HasherConfig) argon2id(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	p, err := c.argon2Params(argon2.RecommendedIDParams)
	if err != nil {
		return nil, nil, err
	}
	return argon2.NewArgon2id(p, limits.Argon2.validationOpts()), []string{argon2.Prefix}, nil
}

func (c *HasherConfig) bcryptCost() (int, error) {
	var dst = struct {
		Cost int `mapstructure:"Cost"`
	}{}
	if err := c.decodeParams(&dst); err != nil {
		return 0, fmt.Errorf("decode bcrypt params: %w", err)
	}
	return dst.Cost, nil
}

func (c *HasherConfig) bcrypt(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	cost, err := c.bcryptCost()
	if err != nil {
		return nil, nil, err
	}
	return bcrypt.New(cost, limits.Bcrypt.validationOpts()), []string{bcrypt.Prefix}, nil
}

func (c *HasherConfig) scryptParams() (scrypt.Params, error) {
	var dst = struct {
		Cost int `mapstructure:"Cost"`
	}{}
	if err := c.decodeParams(&dst); err != nil {
		return scrypt.Params{}, fmt.Errorf("decode scrypt params: %w", err)
	}
	p := scrypt.RecommendedParams // copy
	p.LN = dst.Cost
	return p, nil
}

func (c *HasherConfig) scrypt(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	p, err := c.scryptParams()
	if err != nil {
		return nil, nil, err
	}
	return scrypt.New(p, limits.Scrypt.validationOpts()), []string{scrypt.Prefix, scrypt.Prefix_Linux}, nil
}

func (c *HasherConfig) pbkdf2Params() (p pbkdf2.Params, _ HashMode, _ error) {
	var dst = struct {
		Rounds uint32   `mapstructure:"Rounds"`
		Hash   HashMode `mapstructure:"Hash"`
	}{}
	if err := c.decodeParams(&dst); err != nil {
		return p, "", fmt.Errorf("decode pbkdf2 params: %w", err)
	}
	switch dst.Hash {
	case HashModeSHA1:
		p = pbkdf2.RecommendedSHA1Params
	case HashModeSHA224:
		p = pbkdf2.RecommendedSHA224Params
	case HashModeSHA256:
		p = pbkdf2.RecommendedSHA256Params
	case HashModeSHA384:
		p = pbkdf2.RecommendedSHA384Params
	case HashModeSHA512:
		p = pbkdf2.RecommendedSHA512Params
	}
	p.Rounds = dst.Rounds
	return p, dst.Hash, nil
}

func (c *HasherConfig) pbkdf2(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	p, hash, err := c.pbkdf2Params()
	if err != nil {
		return nil, nil, err
	}
	opts := limits.PBKDF2.validationOpts()
	prefix := []string{pbkdf2.Prefix}
	switch hash {
	case HashModeSHA1:
		return pbkdf2.NewSHA1(p, opts), prefix, nil
	case HashModeSHA224:
		return pbkdf2.NewSHA224(p, opts), prefix, nil
	case HashModeSHA256:
		return pbkdf2.NewSHA256(p, opts), prefix, nil
	case HashModeSHA384:
		return pbkdf2.NewSHA384(p, opts), prefix, nil
	case HashModeSHA512:
		return pbkdf2.NewSHA512(p, opts), prefix, nil
	default:
		return nil, nil, fmt.Errorf("unsupported pbkdf2 hash mode: %s", hash)
	}
}

func (c *HasherConfig) sha2Params() (use512 bool, rounds int, err error) {
	var dst = struct {
		Rounds uint32   `mapstructure:"Rounds"`
		Hash   HashMode `mapstructure:"Hash"`
	}{}
	if err := c.decodeParams(&dst); err != nil {
		return false, 0, fmt.Errorf("decode sha2 params: %w", err)
	}
	switch dst.Hash {
	case HashModeSHA256:
		use512 = false
	case HashModeSHA512:
		use512 = true
	case HashModeSHA1, HashModeSHA224, HashModeSHA384:
		fallthrough
	default:
		return false, 0, fmt.Errorf("cannot use %s with sha2", dst.Hash)
	}
	if dst.Rounds > sha2.RoundsMax {
		return false, 0, fmt.Errorf("rounds with sha2 cannot be larger than %d", sha2.RoundsMax)
	} else {
		rounds = int(dst.Rounds)
	}
	return use512, rounds, nil
}

func (c *HasherConfig) sha2(limits HashLimitsConfig) (passwap.Hasher, []string, error) {
	use512, rounds, err := c.sha2Params()
	if err != nil {
		return nil, nil, err
	}
	opts := limits.Sha2.validationOpts()
	if use512 {
		return sha2.New512(rounds, opts), []string{sha2.Sha256Identifier, sha2.Sha512Identifier}, nil
	}
	return sha2.New256(rounds, opts), []string{sha2.Sha256Identifier, sha2.Sha512Identifier}, nil
}
