package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSON writes the value as indented JSON to stdout.
func JSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling JSON: %w", err)
	}
	_, err = os.Stdout.Write(append(data, '\n'))
	return err
}
