/**
 * OpenTelemetry gRPC/Connect Interceptor
 *
 * Provides distributed tracing for Connect/gRPC calls. Creates client spans
 * with RPC semantic conventions, propagates W3C trace context headers
 * (traceparent/tracestate) and baggage, and records errors with gRPC status
 * codes.
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
  SpanKind,
  SpanStatusCode,
  trace,
} from "@opentelemetry/api";
import { GrpcError, Interceptor } from "./types";

const TRACER_NAME = "zitadel-login-grpc" as const;

/**
 * Extracts the hostname from a URL string. Returns "unknown" if the URL is
 * undefined or cannot be parsed.
 */
function parseHostname(url: string | undefined): string {
  if (!url) {
    return "unknown";
  }
  try {
    return new URL(url).hostname;
  } catch {
    return "unknown";
  }
}

/**
 * OpenTelemetry interceptor for Connect/gRPC calls. Wraps each RPC call in a
 * client span, injects trace context headers for distributed tracing, and
 * records any errors that occur during the call.
 */
export function otelGrpcInterceptor(next: Parameters<Interceptor>[0]): ReturnType<Interceptor> {
  return async function tracedCall(req) {
    const serviceName: string = req.service?.typeName ?? "unknown";
    const methodName: string = req.method?.name ?? "unknown";

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
        const ctx = trace.setSpan(context.active(), span);
        propagation.inject(ctx, req.header, {
          set: (carrier, key, value) => carrier.set(key, value),
        });

        try {
          const response = await next(req);
          span.setStatus({ code: SpanStatusCode.OK });
          return response;
        } catch (err: unknown) {
          const exception =
            err instanceof Error ? err : new Error(String(err));
          span.recordException(exception);
          const grpcError = err as GrpcError;
          if (grpcError.code !== undefined) {
            span.setAttribute("rpc.grpc.status_code", grpcError.code);
          }
          span.setStatus({
            code: SpanStatusCode.ERROR,
            message: exception.message,
          });
          throw err;
        } finally {
          span.end();
        }
      },
    );
  };
}
