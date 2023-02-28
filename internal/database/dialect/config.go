package dialect

import (
	"database/sql"
	"sync"
	"time"
)

type Config struct {
	Dialects map[string]interface{} `mapstructure:",remain"`
	Dialect  Matcher
}

type Dialect struct {
	Matcher   Matcher
	Config    Connector
	IsDefault bool
}

var (
	dialects       []*Dialect
	defaultDialect *Dialect
	dialectsMu     sync.Mutex
)

type Matcher interface {
	MatchName(string) bool
	Decode([]interface{}) (Connector, error)
}

type Connector interface {
	Connect(useAdmin bool) (*sql.DB, error)
	Password() string
	Database
}

type Database interface {
	DatabaseName() string
	Username() string
	Type() string
	Timetravel(time.Duration) string
}

func Register(matcher Matcher, config Connector, isDefault bool) {
	dialectsMu.Lock()
	defer dialectsMu.Unlock()

	d := &Dialect{Matcher: matcher, Config: config}

	if isDefault {
		defaultDialect = d
		return
	}

	dialects = append(dialects, d)
}

func SelectByConfig(config map[string]interface{}) *Dialect {
	for key := range config {
		for _, d := range dialects {
			if d.Matcher.MatchName(key) {
				return d
			}
		}
	}

	return defaultDialect
}
