package cache

type Config interface {
	NewCache() (Cache, error)
}
