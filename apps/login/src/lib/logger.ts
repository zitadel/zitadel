import "server-only";
import winston from "winston";

const isProduction = process.env.NODE_ENV === "production";

/**
 * Winston logger configuration for the login application.
 *
 * OpenTelemetry integration:
 * - @opentelemetry/instrumentation-winston automatically injects trace context
 *   (trace_id, span_id, trace_flags) into every log entry
 * - @opentelemetry/winston-transport automatically sends logs to OTEL LoggerProvider
 *   when the SDK is initialized with a LoggerProvider
 *
 * The trace context injection and OTEL log sending are handled by the instrumentation
 * registered in instrumentation.ts - no manual setup required here.
 */

// Custom format for development - colorized and readable
const devFormat = winston.format.combine(
  winston.format.colorize(),
  winston.format.timestamp({ format: "YYYY-MM-DD HH:mm:ss" }),
  winston.format.printf(({ level, message, timestamp, ...meta }) => {
    const metaStr = Object.keys(meta).length ? ` ${JSON.stringify(meta)}` : "";
    return `${timestamp} ${level}: ${message}${metaStr}`;
  }),
);

// Production format - JSON for structured logging
const prodFormat = winston.format.combine(
  winston.format.timestamp(),
  winston.format.errors({ stack: true }),
  winston.format.json(),
);

/**
 * The main logger instance.
 * In production, outputs JSON for structured logging.
 * In development, outputs colorized human-readable format.
 */
export const logger = winston.createLogger({
  level: process.env.OTEL_LOG_LEVEL || process.env.LOG_LEVEL || (isProduction ? "info" : "debug"),
  format: isProduction ? prodFormat : devFormat,
  defaultMeta: {
    service: process.env.OTEL_SERVICE_NAME || "zitadel-login",
  },
  transports: [
    new winston.transports.Console({
      // In production, use stderr for error level, stdout for others
      stderrLevels: isProduction ? ["error"] : [],
    }),
  ],
});

/**
 * Create a child logger with a specific context.
 * The context is added as metadata to all log entries from this logger.
 *
 * @param context - The context name for this logger (e.g., 'auth-flow', 'session')
 * @param metadata - Optional additional metadata to include in all log messages
 * @returns A child logger instance
 */
export function createLogger(
  context: string,
  metadata?: Record<string, unknown>,
): winston.Logger {
  return logger.child({ context, ...metadata });
}

// Re-export winston types for convenience
export type { Logger } from "winston";
