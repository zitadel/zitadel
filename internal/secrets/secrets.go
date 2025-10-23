package secrets

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zitadel/logging"
)

// ProcessDockerSecretsIntoViper reads environment variables ending with _FILE
// and directly sets the corresponding values in Viper without exposing them
// as environment variables.
func ProcessDockerSecretsIntoViper(v *viper.Viper) error {
	logging.Info("Processing Docker secrets from _FILE environment variables")
	
	processedSecrets := make(map[string]string)
	
	for _, env := range os.Environ() {
		if !strings.Contains(env, "_FILE=") {
			continue
		}
		
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 || !strings.HasSuffix(parts[0], "_FILE") {
			continue
		}
		
		fileEnvVar := parts[0]
		filePath := parts[1]
		
		if filePath == "" {
			continue
		}
		
		content, err := os.ReadFile(filePath)
		if err != nil {
			logging.WithError(err).WithFields(logrus.Fields{
				"env_var":   fileEnvVar,
				"file_path": filePath,
			}).Error("Failed to read Docker secret file")
			continue
		}
		
		secretValue := strings.TrimRight(string(content), "\r\n")
		
		baseEnvVar := strings.TrimSuffix(fileEnvVar, "_FILE")
		
		viperKey := baseEnvVar
		if strings.HasPrefix(viperKey, "ZITADEL_") {
			viperKey = strings.TrimPrefix(viperKey, "ZITADEL_")
		}
		
		viperKey = strings.ToLower(strings.ReplaceAll(viperKey, "_", "."))
		
		v.Set(viperKey, secretValue)
		processedSecrets[viperKey] = filePath
		
		logging.WithFields(logrus.Fields{
			"viper_key": viperKey,
			"file_path": filePath,
		}).Info("Successfully loaded Docker secret")
	}
	
	if len(processedSecrets) > 0 {
		logging.WithFields(logrus.Fields{
			"count": len(processedSecrets),
		}).Info("Loaded Docker secrets")
	}
	
	return nil
}
