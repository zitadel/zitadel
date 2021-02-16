package common

import (
	"bytes"
	"gopkg.in/yaml.v3"
)

func MarshalYAML(sth interface{}) []byte {
	if sth == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buf)

	defer func() {
		encoder.Close()
		buf.Truncate(0)
	}()

	encoder.SetIndent(2)
	if err := encoder.Encode(sth); err != nil {
		panic(err)
	}
	return buf.Bytes()
}
