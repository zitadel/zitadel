package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"sigs.k8s.io/yaml"

	"github.com/zitadel/zitadel/internal/zerrors"
)

type ValidatableConfiguration interface {
	Validate() error
}

type ReaderFunc func(data []byte, o interface{}) error

var (
	JSONReader = json.Unmarshal
	TOMLReader = toml.Unmarshal
	YAMLReader = func(data []byte, o interface{}) error { return yaml.Unmarshal(data, o) }
)

// Read deserializes each config file to the target obj
// using a Reader (depending on file extension)
// env vars are replaced in the config file as well as the file path
func Read(obj interface{}, configFiles ...string) error {
	for _, cf := range configFiles {
		readerFunc, err := readerFuncForFile(cf)
		if err != nil {
			return err
		}
		if err := readConfigFile(readerFunc, cf, obj); err != nil {
			return err
		}
	}

	if validatable, ok := obj.(ValidatableConfiguration); ok {
		if err := validatable.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func readConfigFile(readerFunc ReaderFunc, configFile string, obj interface{}) error {
	configFile = os.ExpandEnv(configFile)

	configStr, err := os.ReadFile(configFile)
	if err != nil {
		return zerrors.ThrowInternalf(err, "CONFI-nJk2a", "failed to read config file %s", configFile)
	}

	configStr = []byte(os.ExpandEnv(string(configStr)))

	if err := readerFunc(configStr, obj); err != nil {
		return zerrors.ThrowInternalf(err, "CONFI-2Mc3c", "error parse config file %s", configFile)
	}

	return nil
}

func readerFuncForFile(configFile string) (ReaderFunc, error) {
	ext := filepath.Ext(configFile)
	switch ext {
	case ".yaml", ".yml":
		return YAMLReader, nil
	case ".json":
		return JSONReader, nil
	case ".toml":
		return TOMLReader, nil
	}
	return nil, zerrors.ThrowUnimplementedf(nil, "CONFI-ZLk4u", "file extension (%s) not supported", ext)
}
