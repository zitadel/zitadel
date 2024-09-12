//go:build !integration

package sink

// StartServer and its returned close function are a no-op
// when the `integration` build tag is disabled.
func StartServer() (close func()) {
	return func() {}
}
