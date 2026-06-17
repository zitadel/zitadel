package crypto

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	lowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
	upperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	digits       = []rune("0123456789")
	symbols      = []rune("~!@#$^&*()_+`-={}|[]:<>?,./")
)

type GeneratorConfig struct {
	Length              uint
	Expiry              time.Duration
	IncludeLowerLetters bool
	IncludeUpperLetters bool
	IncludeDigits       bool
	IncludeSymbols      bool
}

//go:generate mockgen -source code.go -destination ./code_mock.go -package crypto

type Generator interface {
	Length() uint
	Expiry() time.Duration
	Alg() EncryptionAlgorithm
	Runes() []rune
}

type generator struct {
	length uint
	expiry time.Duration
	runes  []rune
}

func (g *generator) Length() uint {
	return g.length
}

func (g *generator) Expiry() time.Duration {
	return g.expiry
}

func (g *generator) Runes() []rune {
	return g.runes
}

type encryptionGenerator struct {
	generator
	alg EncryptionAlgorithm
}

func (g *encryptionGenerator) Alg() EncryptionAlgorithm {
	return g.alg
}

func NewEncryptionGenerator(config GeneratorConfig, algorithm EncryptionAlgorithm) Generator {
	return &encryptionGenerator{
		newGenerator(config),
		algorithm,
	}
}

type HashGenerator struct {
	generator
	hasher *Hasher
}

func NewHashGenerator(config GeneratorConfig, hasher *Hasher) *HashGenerator {
	return &HashGenerator{
		newGenerator(config),
		hasher,
	}
}

func (g *HashGenerator) NewCode() (encoded, plain string, err error) {
	plain, err = GenerateRandomString(g.Length(), g.Runes())
	if err != nil {
		return "", "", err
	}
	encoded, err = g.hasher.Hash(plain)
	if err != nil {
		return "", "", err
	}
	return encoded, plain, nil
}

func newGenerator(config GeneratorConfig) generator {
	var runes []rune
	if config.IncludeLowerLetters {
		runes = append(runes, lowerLetters...)
	}
	if config.IncludeUpperLetters {
		runes = append(runes, upperLetters...)
	}
	if config.IncludeDigits {
		runes = append(runes, digits...)
	}
	if config.IncludeSymbols {
		runes = append(runes, symbols...)
	}
	return generator{
		length: config.Length,
		expiry: config.Expiry,
		runes:  runes,
	}
}

func NewCode(g Generator) (*CryptoValue, string, error) {
	code, err := GenerateRandomString(g.Length(), g.Runes())
	if err != nil {
		return nil, "", err
	}
	crypto, err := Crypt([]byte(code), g.Alg())
	if err != nil {
		return nil, "", err
	}
	return crypto, code, nil
}

func IsCodeExpired(creationDate time.Time, expiry time.Duration) bool {
	if expiry == 0 {
		return false
	}
	return creationDate.Add(expiry).Before(time.Now().UTC())
}

func VerifyCode(creationDate time.Time, expiry time.Duration, cryptoCode *CryptoValue, verificationCode string, algorithm EncryptionAlgorithm) error {
	if IsCodeExpired(creationDate, expiry) {
		return zerrors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "Errors.User.Code.Expired")
	}
	return verifyEncryptedCode(cryptoCode, verificationCode, algorithm)
}

func GenerateRandomString(length uint, chars []rune) (string, error) {
	if length == 0 {
		return "", nil
	}
	if len(chars) == 0 {
		return "", zerrors.ThrowInvalidArgument(nil, "CODE-aa1wf", "chars must not be empty")
	}

	str := make([]rune, length)
	max := big.NewInt(int64(len(chars)))
	for i := range str {
		idx, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		str[i] = chars[int(idx.Int64())]
	}
	return string(str), nil
}

func verifyEncryptedCode(cryptoCode *CryptoValue, verificationCode string, alg EncryptionAlgorithm) error {
	if cryptoCode == nil {
		return zerrors.ThrowInvalidArgument(nil, "CRYPT-aqrFV", "Errors.User.Code.CryptoCodeNil")
	}
	code, err := DecryptString(cryptoCode, alg)
	if err != nil {
		return err
	}

	if code != verificationCode {
		return zerrors.ThrowInvalidArgument(nil, "CODE-woT0xc", "Errors.User.Code.Invalid")
	}
	return nil
}
