import { registerOTel } from '@vercel/otel'
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { propagation } from "@opentelemetry/api";

export function register() {
  console.log('Registering OpenTelemetry for Login App');
  propagation.setGlobalPropagator(new W3CTraceContextPropagator());
  registerOTel({
     serviceName: 'login-app',
     propagators: ['tracecontext', 'baggage'],
     instrumentationConfig: {
      fetch: {
        // This URLs will have the tracing context propagated to them.
        propagateContextUrls: [
          'localhost:8080',
          'zitadel.com',
          'zitadel.localhost'
        ],
        // This URLs will not have the tracing context propagated to them.
        // dontPropagateContextUrls: [
        //   'some-third-party-service-domain.com',
        // ],
        // This URLs will be ignored and will not be traced.
        // ignoreUrls: ['my-internal-private-tool.com'],
      },
    },

    });
}