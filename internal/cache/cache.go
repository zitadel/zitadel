package cache

type Cache interface {
	Set(key string, object interface{}) error
	Get(key string, ptrToObject interface{}) error
	Delete(key string) error
}
