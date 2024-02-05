package crypto

import "context"

type KeyStorage interface {
	ReadKeys() (Keys, error)
	ReadKey(id string) (*Key, error)
	CreateKeys(context.Context, ...*Key) error
}
