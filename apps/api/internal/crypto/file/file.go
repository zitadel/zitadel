package file

import (
	"fmt"
	"os"

	"github.com/zitadel/zitadel/internal/config"
	"github.com/zitadel/zitadel/internal/crypto"
)

const (
	ZitadelKeyPath = "ZITADEL_KEY_PATH"
)

type Storage struct{}

func (d *Storage) ReadKeys() (crypto.Keys, error) {
	path := os.Getenv(ZitadelKeyPath)
	if path == "" {
		return nil, fmt.Errorf("no path set, %s is empty", ZitadelKeyPath)
	}
	keys := new(crypto.Keys)
	err := config.Read(keys, path)
	return *keys, err
}

func (d *Storage) ReadKey(id string) (*crypto.Key, error) {
	keys, err := d.ReadKeys()
	if err != nil {
		return nil, err
	}
	key, ok := keys[id]
	if !ok {
		return nil, fmt.Errorf("key no found")
	}
	return &crypto.Key{
		ID:    id,
		Value: key,
	}, nil
}

func (d *Storage) CreateKeys(keys ...*crypto.Key) error {
	return fmt.Errorf("this provider is not able to store new keys")
}
