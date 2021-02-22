package domain

import (
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
)

var (
	//most of us won't survive until 12-31-9999 23:59:59, maybe ZITADEL does
	defaultExpDate = time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)
)

type AuthNKey interface {
}

type authNKey interface {
	setPublicKey([]byte)
	setPrivateKey([]byte)
	expirationDate() time.Time
	setExpirationDate(time.Time)
}

type AuthNKeyType int32

const (
	AuthNKeyTypeNONE = iota
	AuthNKeyTypeJSON

	keyCount
)

func (k AuthNKeyType) Valid() bool {
	return k >= 0 && k < keyCount
}

func (key *MachineKey) GenerateNewMachineKeyPair(keySize int) error {
	privateKey, publicKey, err := crypto.GenerateKeyPair(keySize)
	if err != nil {
		return err
	}
	key.PublicKey, err = crypto.PublicKeyToBytes(publicKey)
	if err != nil {
		return err
	}
	key.PrivateKey = crypto.PrivateKeyToBytes(privateKey)
	return nil
}

func EnsureValidExpirationDate(key authNKey) error {
	if key.expirationDate().IsZero() {
		key.setExpirationDate(defaultExpDate)
	}
	if key.expirationDate().Before(time.Now()) {
		return errors.ThrowInvalidArgument(nil, "AUTHN-dv3t5", "Errors.AuthNKey.ExpireBeforeNow")
	}
	return nil
}

func SetNewAuthNKeyPair(key authNKey, keySize int) error {
	privateKey, publicKey, err := NewAuthNKeyPair(keySize)
	if err != nil {
		return err
	}
	key.setPrivateKey(privateKey)
	key.setPublicKey(publicKey)
	return nil
}

func NewAuthNKeyPair(keySize int) (privateKey, publicKey []byte, err error) {
	private, public, err := crypto.GenerateKeyPair(keySize)
	if err != nil {
		logging.Log("AUTHN-Ud51I").WithError(err).Error("unable to create authn key pair")
		return nil, nil, errors.ThrowInternal(err, "AUTHN-gdg2l", "Errors.Project.CouldNotGenerateClientSecret")
	}
	publicKey, err = crypto.PublicKeyToBytes(public)
	if err != nil {
		logging.Log("AUTHN-Dbb35").WithError(err).Error("unable to convert public key")
		return nil, nil, errors.ThrowInternal(err, "AUTHN-Bne3f", "Errors.Project.CouldNotGenerateClientSecret")
	}
	privateKey = crypto.PrivateKeyToBytes(private)
	return privateKey, publicKey, nil
}
