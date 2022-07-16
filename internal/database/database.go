package database

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"

	_ "github.com/zitadel/zitadel/internal/database/cockroach"
	"github.com/zitadel/zitadel/internal/database/dialect"
	_ "github.com/zitadel/zitadel/internal/database/postgres"
	"github.com/zitadel/zitadel/internal/errors"
)

type Config struct {
	Dialects map[string]interface{} `mapstructure:",remain"`
	dialect  dialect.Matcher
}

type User struct {
	Username string
	Password string
	SSL      SSL
}

type SSL struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

func Connect(config Config) (*sql.DB, error) {
	// if len(config) > 2 {
	// 	return nil, errors.ThrowInvalidArgument(nil, "DATAB-RbwLj", "too many dialects")
	// }
	// if len(config) == 1 {
	// 	return config["default"].Connect()
	// }
	// for key, conn := range connectors {
	// 	dialectConfig, ok := config[key]
	// 	if !ok {
	// 		continue
	// 	}
	// 	return conn(dialectConfig)
	// }

	return nil, errors.ThrowNotFound(nil, "DATAB-rhOQ6", "dialect not found")
}

func DecodeHook(from, to reflect.Value) (interface{}, error) {
	if to.Type() != reflect.TypeOf(Config{}) {
		return from.Interface(), nil
	}

	confgi := &Config{}
	asdf := mapstructure.Decode(from.Interface(), confgi)
	fmt.Println(asdf, confgi)

	configuredDialects, ok := from.Interface().(map[string]interface{})
	if !ok {
		return from.Interface(), nil
	}

	configuredDialect := dialect.SelectByConfig(configuredDialects)
	configs := make([]interface{}, 0, len(configuredDialects)-1)

	for name, dialectConfig := range configuredDialects {
		if !configuredDialect.Matcher.MatchName(name) {
			continue
		}

		configs = append(configs, dialectConfig)
	}

	return configuredDialect.Matcher.Decode(configs)
}
