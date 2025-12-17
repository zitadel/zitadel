package instrumentation

import (
	"cloud.google.com/go/profiler"

	"github.com/zitadel/zitadel/cmd/build"
)

type ProfileConfig struct {
	Exporter ExporterConfig
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
