import { Code, ConnectError, Interceptor } from "@connectrpc/connect";

/**
 * Retries ZITADEL API calls that fail while connecting.
 *
 * A call is retried only when it fails with `Code.Unavailable` and the Node
 * error behind it came from opening the connection (`connect`) or resolving the
 * host (`getaddrinfo`), so the request never left the Login and repeating it is
 * safe for every method. Unavailable errors that arrive after the request was
 * sent, such as a proxy 503, and every other code pass through unchanged.
 *
 * Retries stop once RETRY_BUDGET_MS has passed since the first attempt, so a
 * slow connect failure is not repeated. The pause between attempts is jittered
 * so Login instances do not retry in step. Streaming calls are not retried:
 * the Login makes only unary calls, and a stream's input cannot be replayed.
 *
 * Place it first in the interceptor chain, before the tracing interceptor, so
 * each attempt gets its own span.
 */

export const MAX_RETRIES = 2;
/** Mean pause between attempts; each pause is uniform in [PAUSE_MS / 2, PAUSE_MS * 1.5). */
export const PAUSE_MS = 100;
export const RETRY_BUDGET_MS = 1000;
const PRE_SEND_SYSCALLS = new Set<string | undefined>(["connect", "getaddrinfo"]);

/** The Node error is on `cause`, or in an AggregateError when Node tried more than one address. */
function failedBeforeSend(reason: unknown): boolean {
  const causes = reason instanceof AggregateError ? reason.errors : [reason];
  return causes.some((cause) => PRE_SEND_SYSCALLS.has((cause as NodeJS.ErrnoException | undefined)?.syscall));
}

function isRetryable(err: unknown): boolean {
  return err instanceof ConnectError && err.code === Code.Unavailable && failedBeforeSend(err.cause);
}

/** Waits about PAUSE_MS, or less if the call is aborted meanwhile. */
function pause(signal: AbortSignal): Promise<void> {
  return new Promise((resolve) => {
    const timer = setTimeout(resolve, PAUSE_MS * (0.5 + Math.random()));
    signal.addEventListener(
      "abort",
      () => {
        clearTimeout(timer);
        resolve();
      },
      { once: true },
    );
  });
}

export const retryUnavailableInterceptor: Interceptor = (next) =>
  async function retriedCall(req) {
    if (req.stream) {
      return next(req);
    }
    const start = Date.now();
    for (let attempt = 0; ; attempt++) {
      try {
        return await next(req);
      } catch (err) {
        if (attempt >= MAX_RETRIES || !isRetryable(err) || req.signal.aborted || Date.now() - start >= RETRY_BUDGET_MS) {
          throw err;
        }
        await pause(req.signal);
        if (req.signal.aborted) {
          throw err;
        }
      }
    }
  };
