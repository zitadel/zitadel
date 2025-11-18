import pino from 'pino';
import { context, trace } from '@opentelemetry/api';

const projectId = process.env.GOOGLE_CLOUD_PROJECT;

function getCurrentTraceFields() {
  const span = trace.getSpan(context.active());
  if (!span) return {};
  const spanContext = span.spanContext();
  if (!spanContext || !spanContext.traceId) return {};
  return {
    'logging.googleapis.com/trace': projectId
      ? `projects/${projectId}/traces/${spanContext.traceId}`
      : spanContext.traceId,
    'logging.googleapis.com/spanId': spanContext.spanId,
    'logging.googleapis.com/trace_sampled': spanContext.traceFlags === 1,
  };
}

export const logger = pino({
  level: process.env.NODE_ENV === 'production' ? 'info' : 'debug',
  formatters: {
    log(obj) {
      return { ...getCurrentTraceFields(), ...obj };
    },
  },
});

// Patch console to use pino (only in Node.js runtime, not Edge)
if (typeof window === 'undefined' && process.env.NEXT_RUNTIME === 'nodejs') {
  const enabled = process.env.LOG_PATCH_CONSOLE !== 'false';
  if (enabled) {
    (global as any).console = {
      ...console,
      log: (...args: any[]) => logger.info(args.length === 1 ? args[0] : args),
      info: (...args: any[]) => logger.info(args.length === 1 ? args[0] : args),
      warn: (...args: any[]) => logger.warn(args.length === 1 ? args[0] : args),
      error: (...args: any[]) => logger.error(args.length === 1 ? args[0] : args),
      debug: (...args: any[]) => logger.debug(args.length === 1 ? args[0] : args),
      trace: (...args: any[]) => logger.trace(args.length === 1 ? args[0] : args),
    };
  }
}
