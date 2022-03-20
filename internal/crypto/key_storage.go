package crypto

type KeyStorage interface {
	ReadKeys() (Keys, error)
	ReadKey(id string) (*Key, error)
	CreateKeys(...*Key) error
}
