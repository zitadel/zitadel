/**
 * Error Classification Interceptor
 *
 * Enriches ConnectError instances with HTTP status code equivalents and
 * client/server classification. This interceptor runs at the transport level,
 * so every gRPC/Connect call automatically gets classified errors so we can properly catch them in client code.
 *
 * Purpose:
 * - Prevent client-side gRPC errors (4xx equivalents) from being surfaced as HTTP 500s
 * - Provide correct HTTP status codes for route handler error responses
 *
 * @see https://cloud.google.com/apis/design/errors#handling_errors
 */

import { Code, ConnectError, Interceptor } from "@connectrpc/connect";

/** Unique brand symbol for ClassifiedConnectError type guard detection */
const CLASSIFIED_BRAND = Symbol.for("ClassifiedConnectError");

/** Canonical gRPC → HTTP status code mapping */
const GRPC_TO_HTTP: Readonly<Record<number, number>> = {
  [Code.InvalidArgument]: 400,
  [Code.FailedPrecondition]: 400,
  [Code.OutOfRange]: 400,
  [Code.Unauthenticated]: 401,
  [Code.PermissionDenied]: 403,
  [Code.NotFound]: 404,
  [Code.AlreadyExists]: 409,
  [Code.Aborted]: 409,
  [Code.ResourceExhausted]: 429,
  [Code.Canceled]: 499,
  [Code.Unimplemented]: 501,
  [Code.Unavailable]: 503,
  [Code.DeadlineExceeded]: 504,
  [Code.DataLoss]: 500,
  [Code.Internal]: 500,
  [Code.Unknown]: 500,
};

/** gRPC codes that represent user input errors (not genuine server failures) */
const CLIENT_ERROR_CODES: ReadonlySet<Code> = new Set([
  Code.InvalidArgument,
  Code.FailedPrecondition,
  Code.OutOfRange,
  Code.Unauthenticated,
  Code.PermissionDenied,
  Code.NotFound,
  Code.AlreadyExists,
  Code.Aborted,
  Code.ResourceExhausted,
  Code.Canceled,
]);

/**
 * A ConnectError enriched with HTTP status classification.
 *
 * All ConnectErrors thrown by service calls through the classified transport
 * will be instances of this class, allowing callers to inspect `httpStatus`
 * and `isUserError` without manual mapping.
 */
export class ClassifiedConnectError extends ConnectError {
  /** The equivalent HTTP status code for this gRPC error */
  readonly httpStatus: number;

  /** Whether this error represents a user input error (true) or a server failure (false) */
  readonly isUserError: boolean;

  /** @internal Brand symbol for type guard detection */
  readonly [CLASSIFIED_BRAND] = true as const;

  constructor(source: ConnectError) {
    super(source.rawMessage, source.code, source.metadata, undefined, source.cause);
    // ConnectError's constructor resets the prototype chain via Object.setPrototypeOf.
    // We must restore it so instanceof ClassifiedConnectError works correctly.
    Object.setPrototypeOf(this, ClassifiedConnectError.prototype);
    this.name = "ClassifiedConnectError";
    // Preserve the original stack trace so debugging/alert triage can see the RPC call site.
    if (source.stack) {
      this.stack = source.stack;
    }
    this.httpStatus = GRPC_TO_HTTP[source.code] ?? 500;
    this.isUserError = CLIENT_ERROR_CODES.has(source.code);

    // Copy details from the source error (avoids OutgoingDetail/IncomingDetail type mismatch)
    if (source.details.length > 0) {
      Object.defineProperty(this, "details", { value: source.details, writable: false });
    }

    // Preserve the raw message from the original error
    if ("rawMessage" in source) {
      Object.defineProperty(this, "rawMessage", { value: source.rawMessage, writable: false });
    }
  }
}

/**
 * Type guard for ClassifiedConnectError.
 * Use this in catch blocks to safely access httpStatus/isUserError.
 */
export function isClassifiedError(error: unknown): error is ClassifiedConnectError {
  return error !== null && typeof error === "object" && CLASSIFIED_BRAND in error;
}

/**
 * Maps a gRPC Code to its HTTP status equivalent.
 * Useful when you have the code but not a full ClassifiedConnectError instance.
 */
export function grpcCodeToHttpStatus(code: Code): number {
  return GRPC_TO_HTTP[code] ?? 500;
}

/**
 * Transport-level interceptor that catches ConnectError and re-throws
 * it as a ClassifiedConnectError with httpStatus and isUserError metadata.
 *
 * Plug this into the transport's interceptor chain to automatically classify
 * every error from every service call.
 */
export const errorClassificationInterceptor: Interceptor = (next) =>
  async function classifiedCall(req) {
    try {
      return await next(req);
    } catch (err) {
      if (err instanceof ConnectError) {
        throw new ClassifiedConnectError(err);
      }
      throw err;
    }
  };
