//go:build integration

package cmd

// DefaultConfig returns the embedded defaults.yaml content.
// This is used by integration test orchestrators to configure
// ZITADEL without going through cobra/viper.
func DefaultConfig() []byte {
	return defaultConfig
}
