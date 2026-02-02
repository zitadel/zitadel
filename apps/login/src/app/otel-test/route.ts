import { NextResponse } from "next/server";
import { trace, metrics, SpanStatusCode } from "@opentelemetry/api";
import { createLogger } from "@/lib/logger";
import {
  recordAuthAttempt,
  recordAuthSuccess,
  recordAuthFailure,
  recordRequestStart,
  recordRequestEnd,
  recordSessionCreationDuration,
} from "@/lib/metrics";

const logger = createLogger("otel-test");

export async function GET() {
  const startTime = Date.now();

  // Test custom metrics from metrics.ts
  recordRequestStart();
  recordAuthAttempt("test", "test-org");

  // Test tracing
  const tracer = trace.getTracer("otel-test");
  const span = tracer.startSpan("test-span");

  try {
    // Test logging with trace context (via Winston)
    // The @opentelemetry/instrumentation-winston automatically:
    // 1. Injects trace_id, span_id into log entries
    // 2. Sends logs to OTEL LoggerProvider
    logger.info("OTEL test endpoint called", { test: true });

    // Test metrics
    const meter = metrics.getMeter("otel-test");
    const counter = meter.createCounter("otel_test_requests_total", {
      description: "Total number of OTEL test requests",
    });
    counter.add(1, { test: "true" });

    const histogram = meter.createHistogram("otel_test_duration_seconds", {
      description: "Duration of OTEL test requests",
      unit: "s",
    });

    // Set span attributes
    span.setAttribute("test.attribute", "test-value");
    span.setAttribute("test.timestamp", new Date().toISOString());

    // Simulate some work
    await new Promise((resolve) => setTimeout(resolve, 10));

    const duration = (Date.now() - startTime) / 1000;
    histogram.record(duration, { test: "true" });

    // Record custom metrics success
    recordAuthSuccess("test", "test-org");
    recordRequestEnd(startTime, 200, "/otel-test");
    recordSessionCreationDuration(startTime, true);

    span.setStatus({ code: SpanStatusCode.OK });
    span.end();

    logger.info("OTEL test completed successfully", {
      duration,
      traceId: span.spanContext().traceId,
    });

    return NextResponse.json({
      status: "ok",
      message: "OTEL test completed",
      traceId: span.spanContext().traceId,
      spanId: span.spanContext().spanId,
      duration: `${duration}s`,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    // Record custom metrics failure
    recordAuthFailure("test", "error", "test-org");
    recordRequestEnd(startTime, 500, "/otel-test");
    recordSessionCreationDuration(startTime, false);

    logger.error("OTEL test failed", { error });
    span.setStatus({ code: SpanStatusCode.ERROR, message: String(error) });
    span.end();

    return NextResponse.json(
      {
        status: "error",
        message: String(error),
        traceId: span.spanContext().traceId,
      },
      { status: 500 },
    );
  }
}
