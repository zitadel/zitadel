package instrumentation

import (
	"cloud.google.com/go/profiler"

	"github.com/zitadel/zitadel/cmd/build"
)

type ProfileConfig struct {
	Exporter ExporterConfig
}

// TODO: remove for v5 release
type LegacyProfileConfig struct {
	Type      string
	ProjectID string
}

func (c *ProfileConfig) SetLegacyConfig(lc *LegacyProfileConfig) {
	typ := c.Exporter.Type
	if lc == nil || typ.isNone() {
		return
	}
	if lc.Type == "google" {
		c.Exporter.Type = ExporterTypeGoogle
		c.Exporter.GoogleProjectID = lc.ProjectID
	}
}

func startProfiler(cfg ProfileConfig, serviceName string) error {
	typ := cfg.Exporter.Type
	if typ.isNone() {
		return nil
	}
	if typ != ExporterTypeGoogle {
		return errExporterType(typ, "profiler")
	}
	pc := profiler.Config{
		Service:        serviceName,
		ServiceVersion: build.Version(),
		ProjectID:      cfg.Exporter.GoogleProjectID,
	}
	return profiler.Start(pc)
}
