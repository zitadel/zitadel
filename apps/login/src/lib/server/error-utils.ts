import { isClassifiedError } from "@/lib/grpc/interceptors/error-classification";
import { createLogger } from "@/lib/logger";

const logger = createLogger("action-error");

/**
 * Converts a thrown service-call error into the `{ error }` result shape server
 * actions return to their client components.
 *
 * A server action that throws makes the whole POST respond with HTTP 500 — even
 * when the client catches the rejected promise — which pollutes availability
 * SLOs with user-caused failures. User errors (gRPC 4xx equivalents: wrong
 * code, session gone, precondition failed, ...) therefore resolve to the given
 * translated message so the request completes with 200 and the form renders the
 * error. Genuine server failures (Internal, Unavailable, ...) are rethrown so
 * they still surface as 500s and stay visible to alerting.
 */
export function catchUserError(error: unknown, message: string, context?: Record<string, unknown>): { error: string } {
  // Some session helpers (passwordAttemptsHandler) throw a plain `{ error }`
  // object with a preformatted user message; pass it through unchanged.
  if (error !== null && typeof error === "object" && "error" in error && typeof error.error === "string") {
    return { error: error.error };
  }

  if (isClassifiedError(error) && error.isUserError) {
    logger.warn("Service call rejected with user error", {
      grpcCode: error.code,
      httpStatus: error.httpStatus,
      message: error.rawMessage,
      ...context,
    });
    return { error: message };
  }

  throw error;
}
