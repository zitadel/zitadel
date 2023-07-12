package crypto

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/passwap"
	"github.com/zitadel/passwap/argon2"
	"github.com/zitadel/passwap/bcrypt"
	"github.com/zitadel/passwap/md5"
	"github.com/zitadel/passwap/scrypt"
	"github.com/zitadel/passwap/verifier"

	"github.com/zitadel/zitadel/internal/errors"
)

type HashName string

const (
	HashNameArgon2   HashName = "argon2"   // used for the common argon2 verifier
	HashNameArgon2i  HashName = "argon2i"  // hash only
	HashNameArgon2id HashName = "argon2id" // hash only
	HashNameBcrypt   HashName = "bcrypt"   // hash and verify
	HashNameMd5      HashName = "md5"      // verify only, as hashing with md5 is insecure and deprecated
	HashNameScrypt   HashName = "scrypt"   // hash and verify
)

type PasswordHashConfig struct {
	Verifiers []HashName
	Hasher    HasherConfig
}

func (c *PasswordHashConfig) BuildSwapper() (*passwap.Swapper, error) {
	verifiers, err := c.buildVerifiers()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "CRYPT-sahW9", "password hash config invalid")
	}
	hasher, err := c.Hasher.buildHasher()
	if err != nil {
		return nil, errors.ThrowInvalidArgument(err, "CRYPT-Que4r", "password hash config invalid")
	}
	return passwap.NewSwapper(hasher, verifiers...), nil
}

// map HashNames to Verifier instances.
var knowVerifiers = map[HashName]verifier.Verifier{
	HashNameArgon2: argon2.Verifier,
	HashNameBcrypt: bcrypt.Verifier,
	HashNameMd5:    md5.Verifier,
	HashNameScrypt: scrypt.Verifier,
}

func (c *PasswordHashConfig) buildVerifiers() ([]verifier.Verifier, error) {
	verifiers := make([]verifier.Verifier, len(c.Verifiers))
	for i, name := range c.Verifiers {
		v, ok := knowVerifiers[name]
		if !ok {
			return nil, fmt.Errorf("invalid verifier %q", name)
		}
		verifiers[i] = v
	}
	return verifiers, nil
}

type HasherConfig struct {
	Algorithm HashName
	Params    map[string]any `mapstructure:",remain"`
}

func (c *HasherConfig) buildHasher() (hasher passwap.Hasher, err error) {
	switch c.Algorithm {
	case HashNameArgon2i:
		return c.argon2i()
	case HashNameArgon2id:
		return c.argon2id()
	case HashNameBcrypt:
		return c.bcrypt()
	case HashNameScrypt:
		return c.scrypt()
	case "":
		// return nil, fmt.Errorf("missing hasher name")
		return nil, nil
		// Discuss: the setup commands seems to run without taking defaults.yaml into acount.
		// That means if the Hasher is not configured in steps, migrations break.
		// On top of that, such failure seems to corrupt the database in a way
		// that migrations keep failing until the database recreated.
	case HashNameArgon2, HashNameMd5:
		fallthrough
	default:
		return nil, fmt.Errorf("invalid name %q", c.Algorithm)
	}
}

func (c *HasherConfig) decodeParams(dst any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		ErrorUnused: true,
		ErrorUnset:  true,
		Result:      dst,
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

func (c *HasherConfig) argon2i() (passwap.Hasher, error) {
	p, err := c.argon2Params(argon2.RecommendedIParams)
	if err != nil {
		return nil, err
	}
	return argon2.NewArgon2i(p), nil
}

func (c *HasherConfig) argon2id() (passwap.Hasher, error) {
	p, err := c.argon2Params(argon2.RecommendedIDParams)
	if err != nil {
		return nil, err
	}
	return argon2.NewArgon2id(p), nil
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

func (c *HasherConfig) bcrypt() (passwap.Hasher, error) {
	cost, err := c.bcryptCost()
	if err != nil {
		return nil, err
	}
	return bcrypt.New(cost), nil
}

func (c *HasherConfig) scryptParams() (scrypt.Params, error) {
	var dst = struct {
		Cost int `mapstructure:"Cost"`
	}{}
	if err := c.decodeParams(&dst); err != nil {
		return scrypt.Params{}, fmt.Errorf("decode scrypt params: %w", err)
	}
	p := scrypt.RecommendedParams // copy
	p.N = 1 << dst.Cost
	return p, nil
}

func (c *HasherConfig) scrypt() (passwap.Hasher, error) {
	p, err := c.scryptParams()
	if err != nil {
		return nil, err
	}
	return scrypt.New(p), nil
}
