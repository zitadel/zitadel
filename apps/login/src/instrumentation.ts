import { registerOTel } from '@vercel/otel'
import { W3CTraceContextPropagator } from "@opentelemetry/core";
import { propagation } from "@opentelemetry/api";

export function register() {
  propagation.setGlobalPropagator(new W3CTraceContextPropagator());
  registerOTel({ serviceName: 'login-app', propagators: ['tracecontext', 'baggage'] });
}