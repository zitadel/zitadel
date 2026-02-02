/**
 * OpenTelemetry gRPC/Connect Interceptor
 *
 * Provides distributed tracing for Connect/gRPC calls by:
 * - Creating client spans with RPC semantic conventions
 * - Propagating W3C trace context (traceparent/tracestate) headers
 * - Propagating W3C baggage headers
 * - Recording errors with gRPC status codes
 *
 * @see https://opentelemetry.io/docs/specs/semconv/rpc/rpc-spans/
 *
 * @example
 * ```typescript
 * import { otelGrpcInterceptor } from "@/lib/grpc/interceptors/otel";
 *
 * const transport = createTransport({
 *   baseUrl: "https://api.example.com",
 *   interceptors: [otelGrpcInterceptor],
 * });
 * ```
 */

import {
  context,
  propagation,
  Span,
  SpanKind,
  SpanStatusCode,
  trace,
} from "@opentelemetry/api";
import { GrpcError, Interceptor } from "./types";

const TRACER_NAME = "zitadel-login-grpc" as const;

const parseHostname = (url: string | undefined): string => {
  if (!url) return "unknown";
  try {
    return new URL(url).hostname;
  } catch {
    return "unknown";
  }
};

const isGrpcError = (err: unknown): err is GrpcError => {
  return err instanceof Error && "code" in err;
};

const handleError = (span: Span, err: unknown): never => {
  const errorMessage = err instanceof Error ? err.message : String(err);
  if (err instanceof Error) {
    span.recordException(err);
  }
  if (isGrpcError(err) && err.code !== undefined) {
    span.setAttribute("rpc.grpc.status_code", err.code);
  }
  span.setStatus({ code: SpanStatusCode.ERROR, message: errorMessage });
  throw err;
};

export const otelGrpcInterceptor: Interceptor = (next) => async (req) => {
  const serviceName = req.service?.typeName ?? "unknown";
  const methodName = req.method?.name ?? "unknown";

  return trace.getTracer(TRACER_NAME).startActiveSpan(
    `${serviceName}/${methodName}`,
    {
      kind: SpanKind.CLIENT,
      attributes: {
        "rpc.system": "grpc",
        "rpc.service": serviceName,
        "rpc.method": methodName,
        "server.address": parseHostname(req.url),
      },
    },
    async (span) => {
      // Explicitly set the span in context before injection to ensure
      // traceparent headers are propagated even if async context is lost
      const ctx = trace.setSpan(context.active(), span);
      propagation.inject(ctx, req.header, {
        set: (carrier, key, value) => carrier.set(key, value),
      });

      try {
        const response = await next(req);
        span.setStatus({ code: SpanStatusCode.OK });
        return response;
      } catch (err) {
        handleError(span, err);
      } finally {
        span.end();
      }
    },
  );
};
