package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"

	"github.com/caos/zitadel/internal/errors"
)

type Reader interface {
	Unmarshal(data []byte, o interface{}) error
}

type ValidatableConfiguration interface {
	Validate() error
}

type ReaderFunc func(data []byte, o interface{}) error

func (c ReaderFunc) Unmarshal(data []byte, o interface{}) error {
	return c(data, o)
}

var (
	JSONReader = ReaderFunc(json.Unmarshal)
	TOMLReader = ReaderFunc(toml.Unmarshal)
	YAMLReader = ReaderFunc(func(y []byte, o interface{}) error {
		return yaml.Unmarshal(y, o)
	})
)

// Read deserializes each config file to the target obj
// using a Reader (depending on file extension)
// env vars are replaced in the config file as well as the file path
func Read(obj interface{}, configFiles ...string) error {
	for _, cf := range configFiles {
		configReader, err := configReaderForFile(cf)
		if err != nil {
			return err
		}
		if err := readConfigFile(configReader, cf, obj); err != nil {
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

func readConfigFile(configReader Reader, configFile string, obj interface{}) error {
	configFile = os.ExpandEnv(configFile)

	if _, err := os.Stat(configFile); err != nil {
		return errors.ThrowNotFoundf(err, "CONFI-Hs93M", "config file %s does not exist", configFile)
	}

	configStr, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.ThrowInternalf(err, "CONFI-nJk2a", "failed to read config file %s", configFile)
	}

	configStr = []byte(os.ExpandEnv(string(configStr)))

	if err := configReader.Unmarshal(configStr, obj); err != nil {
		return errors.ThrowInternalf(err, "CONFI-2Mc3c", "error parse config file %s", configFile)
	}

	return nil
}

func configReaderForFile(configFile string) (Reader, error) {
	ext := filepath.Ext(configFile)
	switch ext {
	case ".yaml", ".yml":
		return YAMLReader, nil
	case ".json":
		return JSONReader, nil
	case ".toml":
		return TOMLReader, nil
	}
	return nil, errors.ThrowUnimplementedf(nil, "file extension (%s) not supported", ext)
}
