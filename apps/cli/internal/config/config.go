package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config is the top-level CLI configuration.
type Config struct {
	ActiveContext string             `toml:"active-context" mapstructure:"active-context"`
	Contexts      map[string]Context `toml:"contexts"       mapstructure:"contexts"`
}

// Context represents one configured ZITADEL instance.
type Context struct {
	Instance   string `toml:"instance"    mapstructure:"instance"`
	AuthMethod string `toml:"auth-method" mapstructure:"auth-method"`
	PAT        string `toml:"pat"         mapstructure:"pat"`
	ClientID   string `toml:"client-id"   mapstructure:"client-id"`
	Token      string `toml:"token"       mapstructure:"token"`
}

// Path returns the config file path, respecting XDG_CONFIG_HOME.
func Path() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, _ := os.UserHomeDir()
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "zitadel", "config.toml")
}

// Load reads the config from disk. Returns an empty config if the file doesn't exist.
// It also loads a .env file from the current directory if one exists, populating
// env vars that the ZITADEL_TOKEN / ZITADEL_INSTANCE overrides will pick up.
func Load() (*Config, error) {
	// Silently ignore missing .env — only fail on parse errors
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	cfg := &Config{Contexts: make(map[string]Context)}
	path := Path()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Contexts == nil {
		cfg.Contexts = make(map[string]Context)
	}
	return cfg, nil
}

// Save writes the config to disk, creating parent directories as needed.
func Save(cfg *Config) error {
	path := Path()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.Set("active-context", cfg.ActiveContext)
	v.Set("contexts", cfg.Contexts)
	return v.WriteConfigAs(path)
}

// ActiveCtx returns the active context. It checks env overrides first.
func ActiveCtx(cfg *Config) (*Context, string, error) {
	name := cfg.ActiveContext
	if name == "" && len(cfg.Contexts) == 0 {
		// No config file — check env vars for ad-hoc usage
		inst := os.Getenv("ZITADEL_INSTANCE")
		token := os.Getenv("ZITADEL_TOKEN")
		if inst != "" || token != "" {
			ctx := &Context{
				Instance:   inst,
				AuthMethod: "pat",
				PAT:        token,
			}
			return ctx, "", nil
		}
		return nil, "", fmt.Errorf("no active context configured; set one with 'zitadel context use <name>' or set ZITADEL_INSTANCE")
	}

	ctx, ok := cfg.Contexts[name]
	if !ok {
		return nil, "", fmt.Errorf("context %q not found in config", name)
	}

	// Env overrides
	if inst := os.Getenv("ZITADEL_INSTANCE"); inst != "" {
		ctx.Instance = inst
	}
	if token := os.Getenv("ZITADEL_TOKEN"); token != "" {
		ctx.PAT = token
		ctx.AuthMethod = "pat"
	}

	return &ctx, name, nil
}
