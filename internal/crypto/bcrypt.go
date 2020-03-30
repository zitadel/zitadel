package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

var _ HashAlgorithm = (*BCrypt)(nil)

type BCrypt struct {
	cost int
}

func NewBCrypt(cost int) *BCrypt {
	return &BCrypt{cost: cost}
}

func (b *BCrypt) Algorithm() string {
	return "bcrypt"
}

func (b *BCrypt) Hash(value []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(value, b.cost)
}

func (b *BCrypt) CompareHash(hashed, value []byte) error {
	return bcrypt.CompareHashAndPassword(hashed, value)
}
