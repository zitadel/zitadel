import type { Interceptor } from "@zitadel/client";
import { trace, SpanKind, SpanStatusCode } from "@opentelemetry/api";

/**
 * OpenTelemetry interceptor for Connect RPC clients.
 * Automatically creates spans for all outgoing RPC calls with:
 * - RPC method name as span name
 * - Request/response metadata
 * - Error tracking
 */
export const createOTelConnectInterceptor = (tracerName = "connect-rpc-client"): Interceptor => {
  const tracer = trace.getTracer(tracerName);

  return (next) => async (req) => {
    // Extract service and method name from the request
    const serviceName = req.service.typeName || "unknown-service";
    const methodName = req.method.name || "unknown-method";
    const spanName = `${serviceName}/${methodName}`;

    // Start a new span for this RPC call
    return await tracer.startActiveSpan(
      spanName,
      {
        kind: SpanKind.CLIENT,
        attributes: {
          "rpc.system": "connect",
          "rpc.service": serviceName,
          "rpc.method": methodName,
          "net.peer.name": req.url,
        },
      },
      async (span) => {
        try {
          // Proceed with the request
          const response = await next(req);

          // Mark span as successful
          span.setStatus({ code: SpanStatusCode.OK });

          return response;
        } catch (error) {
          // Record error details
          span.setStatus({
            code: SpanStatusCode.ERROR,
            message: error instanceof Error ? error.message : String(error),
          });

          if (error instanceof Error) {
            span.recordException(error);

            // Add Connect-specific error attributes if available
            if ("code" in error) {
              span.setAttribute("rpc.grpc.status_code", String(error.code));
            }
            if ("message" in error) {
              span.setAttribute("rpc.error.message", error.message);
            }
          }

          // Re-throw the error
          throw error;
        } finally {
          // End the span
          span.end();
        }
      },
    );
  };
};
