import { NodeTracerProvider } from '@opentelemetry/sdk-trace-node';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-grpc';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';

const exporter = new OTLPTraceExporter({
  url: process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'http://localhost:4318/v1/traces',
  headers: process.env.GOOGLE_CLOUD_PROJECT
    ? { 'x-goog-project-id': process.env.GOOGLE_CLOUD_PROJECT }
    : {},
});

const spanProcessor = new BatchSpanProcessor(exporter);

const provider = new NodeTracerProvider({
  spanProcessors: [spanProcessor],
});

provider.register();
