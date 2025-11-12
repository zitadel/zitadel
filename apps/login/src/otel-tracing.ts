import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';

// Cloud Run automatically provides credentials via Workload Identity
// The exporter will use cloudtrace.googleapis.com when GOOGLE_CLOUD_PROJECT is set
const exporter = new OTLPTraceExporter({
  url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'cloudtrace.googleapis.com:443',
  headers: {},
});

const spanProcessor = new BatchSpanProcessor(exporter);

const provider = new NodeTracerProvider({
  spanProcessors: [spanProcessor],
});

provider.register();

console.log('OpenTelemetry tracing initialized', {
  endpoint: process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'cloudtrace.googleapis.com:443',
  projectId: process.env.GOOGLE_CLOUD_PROJECT,
});
