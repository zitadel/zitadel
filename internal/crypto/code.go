package crypto

import (
	"crypto/rand"
	"time"

	"github.com/caos/zitadel/internal/errors"
)

var (
	LowerLetters = []rune("abcdefghijklmnopqrstuvwxyz")
	UpperLetters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	Digits       = []rune("0123456789")
	Symbols      = []rune("~!@#$^&*()_+`-={}|[]:<>?,./")
)

type Generator interface {
	Length() uint
	Expiry() time.Duration
	Alg() Crypto
	Runes() []rune
}

type EncryptionGenerator struct {
	length uint
	expiry time.Duration
	alg    EncryptionAlg
	runes  []rune
}

func (g *EncryptionGenerator) Length() uint {
	return g.length
}

func (g *EncryptionGenerator) Expiry() time.Duration {
	return g.expiry
}

func (g *EncryptionGenerator) Alg() Crypto {
	return g.alg
}

func (g *EncryptionGenerator) Runes() []rune {
	return g.runes
}

func NewEncryptionGenerator(length uint, expiry time.Duration, alg EncryptionAlg, runes []rune) *EncryptionGenerator {
	return &EncryptionGenerator{
		length: length,
		expiry: expiry,
		alg:    alg,
		runes:  runes,
	}
}

type HashGenerator struct {
	length uint
	expiry time.Duration
	alg    HashAlg
	runes  []rune
}

func (g *HashGenerator) Length() uint {
	return g.length
}

func (g *HashGenerator) Expiry() time.Duration {
	return g.expiry
}

func (g *HashGenerator) Alg() Crypto {
	return g.alg
}

func (g *HashGenerator) Runes() []rune {
	return g.runes
}

func NewHashGenerator(length uint, expiry time.Duration, alg HashAlg, runes []rune) *HashGenerator {
	return &HashGenerator{
		length: length,
		expiry: expiry,
		alg:    alg,
		runes:  runes,
	}
}

func NewCode(g Generator) (*CryptoValue, string, error) {
	code, err := generateRandomString(g.Length(), g.Runes())
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
	return creationDate.Add(expiry).Before(time.Now().UTC())
}

func VerifyCode(creationDate time.Time, expiry time.Duration, cryptoCode *CryptoValue, verificationCode string, g Generator) error {
	if IsCodeExpired(creationDate, expiry) {
		return errors.ThrowPreconditionFailed(nil, "CODE-QvUQ4P", "verification code is expired")
	}
	switch alg := g.Alg().(type) {
	case EncryptionAlg:
		return verifyEncryptedCode(cryptoCode, verificationCode, alg)
	case HashAlg:
		return verifyHashedCode(cryptoCode, verificationCode, alg)
	}
	return errors.ThrowInvalidArgument(nil, "CODE-fW2gNa", "generator alg is not supported")
}

func generateRandomString(length uint, chars []rune) (string, error) {
	if length == 0 {
		return "", nil
	}

	max := len(chars) - 1
	maxStr := int(length - 1)

	str := make([]rune, length)
	randBytes := make([]byte, length)
	if _, err := rand.Read(randBytes); err != nil {
		return "", err
	}
	for i, rb := range randBytes {
		str[i] = chars[int(rb)%max]
		if i == maxStr {
			return string(str), nil
		}
	}
	return "", nil
}

func verifyEncryptedCode(cryptoCode *CryptoValue, verificationCode string, alg EncryptionAlg) error {
	code, err := DecryptString(cryptoCode, alg)
	if err != nil {
		return err
	}

	if code != verificationCode {
		return errors.ThrowInvalidArgument(nil, "CODE-woT0xc", "verification code is invalid")
	}
	return nil
}

func verifyHashedCode(cryptoCode *CryptoValue, verificationCode string, alg HashAlg) error {
	return CompareHash(cryptoCode, []byte(verificationCode), alg)
}
