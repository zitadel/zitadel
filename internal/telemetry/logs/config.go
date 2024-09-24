package logs

import (
	"fmt"
	"github.com/zitadel/zitadel/pkg/streams"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"
	"github.com/sirupsen/logrus"
	"github.com/zitadel/logging"
	"github.com/zitadel/logging/otel"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/otel/log"
)

type Hook string

const (
	GCPLoggingOtelExporter Hook = "GCPLoggingOtelExporter"
)

type Config struct {
	Log   logging.Config `mapstructure:",squash"`
	Hooks map[string]map[string]interface{}
}

type GCPExporterConfig struct {
	Attributes map[string]string
}

func (g *GCPExporterConfig) ToAttributes() (attributes []log.KeyValue) {
	for k, v := range g.Attributes {
		attributes = append(attributes, log.KeyValue{Key: k, Value: log.StringValue(fmt.Sprintf("%v", v))})
	}
	return attributes
}

func (c *Config) SetLogger() (err error) {
	var hooks []logrus.Hook
	for name, rawCfg := range c.Hooks {
		switch name {
		case strings.ToLower(string(GCPLoggingOtelExporter)):
			var hook *otel.GcpLoggingExporterHook
			addedAttributes := &GCPExporterConfig{}
			if err = decodeRawConfig(rawCfg, addedAttributes); err != nil {
				return err
			}
			hook, err = otel.NewGCPLoggingExporterHook(
				otel.WithExporterConfig(func(cfg *googlecloudexporter.Config) {
					cfg.LogConfig.DefaultLogName = "zitadel"
					cfg.LogConfig.ServiceResourceLabels = false
					err = decodeRawConfig(rawCfg, cfg)
				}),
				otel.WithOtelSettings(func(cfg *exporter.Settings) {
					err = decodeRawConfig(rawCfg, cfg)
				}),
				otel.WithInclude(func(entry *logrus.Entry) bool {
					stream, ok := entry.Data["stream"].(streams.Stream)
					return ok && stream == streams.LogFieldValueStreamActivity
				}),
				otel.WithLevels([]logrus.Level{logrus.InfoLevel}),
				otel.WithAttributes(addedAttributes.ToAttributes()),
			)
			if err != nil {
				return err
			}
			if err = hook.Start(); err != nil {
				return err
			}
			hooks = append(hooks, hook)
		default:
			return fmt.Errorf("unknown hook: %s", name)
		}
	}
	return c.Log.SetLogger(
		logging.AddHooks(hooks...),
	)
}

func decodeRawConfig(rawConfig map[string]interface{}, typedConfig any) (err error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		MatchName: func(mapKey, fieldName string) bool {
			return strings.ToLower(mapKey) == strings.ToLower(fieldName)
		},
		WeaklyTypedInput: true,
		Result:           typedConfig,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(rawConfig)
}
