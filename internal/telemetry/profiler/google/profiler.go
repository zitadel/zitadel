package google

import (
	"cloud.google.com/go/profiler"

	"github.com/zitadel/zitadel/cmd/build"
)

type Config struct {
	ProjectID string
}

func NewProfiler(rawConfig map[string]interface{}) (err error) {
	c := new(Config)
	c.ProjectID, _ = rawConfig["projectid"].(string)
	return c.NewProfiler()
}

func (c *Config) NewProfiler() (err error) {
	cfg := profiler.Config{
		Service:        "zitadel",
		ServiceVersion: build.Version(),
		ProjectID:      c.ProjectID,
	}
	return profiler.Start(cfg)
}
