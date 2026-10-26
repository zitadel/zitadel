import { Code, ConnectError, createRouterTransport, Interceptor } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-node";
import { SpanStatusCode } from "@opentelemetry/api";
import { resourceFromAttributes } from "@opentelemetry/resources";
import { InMemorySpanExporter, SimpleSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { createClientFor } from "@zitadel/client";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import type { AddressInfo } from "node:net";
import { createServer } from "node:net";
import { afterAll, afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { errorClassificationInterceptor, isClassifiedError } from "./error-classification";
import { otelGrpcInterceptor } from "./otel";
import { MAX_RETRIES, PAUSE_MS, retryUnavailableInterceptor } from "./retry";

const exporter = new InMemorySpanExporter();
const provider = new NodeTracerProvider({
  resource: resourceFromAttributes({ "service.name": "test-zitadel-client" }),
  spanProcessors: [new SimpleSpanProcessor(exporter)],
});
provider.register();

afterAll(() => provider.shutdown());

function nodeError(code: string, syscall: string) {
  return Object.assign(new Error(`${syscall} ${code}`), { code, syscall });
}

/**
 * Builds a client over the production chain order, with a fake network step
 * last. The step throws `failure` without calling `next` for the first
 * `failures` attempts (the request never left the Login), then calls `next`.
 */
function setup(opts: { failures: number; failure?: () => ConnectError; handlers?: Record<string, () => unknown> }) {
  const state = { attempts: 0 };
  const network: Interceptor = (next) => async (req) => {
    state.attempts++;
    if (state.attempts <= opts.failures) {
      throw (
        opts.failure?.() ??
        new ConnectError("connect ETIMEDOUT", Code.Unavailable, undefined, undefined, nodeError("ETIMEDOUT", "connect"))
      );
    }
    return next(req);
  };

  const transport = createRouterTransport(
    ({ service }) => {
      service(SessionService, {
        getSession: () => ({}),
        createSession: () => ({}),
        ...opts.handlers,
      });
    },
    {
      transport: {
        interceptors: [retryUnavailableInterceptor, otelGrpcInterceptor, errorClassificationInterceptor, network],
      },
    },
  );

  return { client: createClientFor(SessionService)(transport), state };
}

/** Starts a call and lets the fake clock run past every pause. */
async function settle<T>(call: Promise<T>) {
  const outcome = call.then(
    (value) => ({ value }),
    (error: unknown) => ({ error }),
  );
  await vi.advanceTimersByTimeAsync(PAUSE_MS * (MAX_RETRIES + 1));
  return outcome;
}

describe("retryUnavailableInterceptor", () => {
  beforeEach(() => {
    exporter.reset();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it.each([
    ["connect ETIMEDOUT", () => nodeError("ETIMEDOUT", "connect")],
    ["getaddrinfo ENOTFOUND", () => nodeError("ENOTFOUND", "getaddrinfo")],
    [
      "AggregateError of connect ECONNREFUSED",
      () => new AggregateError([nodeError("ECONNREFUSED", "connect"), nodeError("ECONNREFUSED", "connect")]),
    ],
  ])("returns the response after a %s failure before send", async (_, cause) => {
    const { client, state } = setup({
      failures: 1,
      failure: () => new ConnectError("unavailable", Code.Unavailable, undefined, undefined, cause()),
    });

    const result = await settle(client.getSession({}));

    expect(result).toHaveProperty("value");
    expect(state.attempts).toBe(2);
  });

  it("makes every attempt and then throws the classified Unavailable error", async () => {
    const { client, state } = setup({ failures: Infinity });

    const result = (await settle(client.getSession({}))) as { error: unknown };

    expect(state.attempts).toBe(MAX_RETRIES + 1);
    expect(isClassifiedError(result.error)).toBe(true);
    expect(result.error).toMatchObject({ code: Code.Unavailable, httpStatus: 503, isUserError: false });
  });

  it("waits for the pause before the next attempt", async () => {
    const { client, state } = setup({ failures: 1 });

    const call = client.getSession({});
    await vi.advanceTimersByTimeAsync(PAUSE_MS - 1);
    expect(state.attempts).toBe(1);

    await vi.advanceTimersByTimeAsync(1);
    expect(state.attempts).toBe(2);
    await call;
  });

  it("does not retry once the call is aborted", async () => {
    const { client, state } = setup({ failures: Infinity });
    const controller = new AbortController();

    const call = settle(client.getSession({}, { signal: controller.signal }));
    await vi.advanceTimersByTimeAsync(PAUSE_MS / 2);
    controller.abort();
    await call;

    expect(state.attempts).toBe(1);
  });

  it("does not retry an error with a code other than Unavailable", async () => {
    const { client, state } = setup({
      failures: Infinity,
      // ECONNRESET maps to Aborted in connect-node
      failure: () =>
        new ConnectError("read ECONNRESET", Code.Aborted, undefined, undefined, nodeError("ECONNRESET", "read")),
    });

    const result = (await settle(client.getSession({}))) as { error: unknown };

    expect(state.attempts).toBe(1);
    expect(result.error).toMatchObject({ code: Code.Aborted });
  });

  it("does not retry an Unavailable error that arrived after the request was sent", async () => {
    const { client, state } = setup({
      failures: 0,
      handlers: {
        createSession: () => {
          throw new ConnectError("upstream unavailable", Code.Unavailable);
        },
      },
    });

    const result = (await settle(client.createSession({}))) as { error: unknown };

    expect(state.attempts).toBe(1);
    expect(result.error).toMatchObject({ code: Code.Unavailable });
  });

  it("records one span per attempt, each with its own status code", async () => {
    const { client } = setup({ failures: 1 });

    await settle(client.getSession({}));

    const spans = exporter.getFinishedSpans().filter((s) => s.name.endsWith("/GetSession"));
    expect(spans).toHaveLength(2);
    expect(spans[0].status.code).toBe(SpanStatusCode.ERROR);
    expect(spans[0].attributes["rpc.grpc.status_code"]).toBe(14);
    expect(spans[1].status.code).toBe(SpanStatusCode.OK);
    expect(spans[1].attributes["rpc.grpc.status_code"]).toBe(0);
  });
});

describe("retryUnavailableInterceptor over connect-node", () => {
  it("retries a real connection refused error", async () => {
    const server = createServer();
    await new Promise<void>((resolve) => server.listen(0, "127.0.0.1", resolve));
    const { port } = server.address() as AddressInfo;
    await new Promise((resolve) => server.close(resolve));

    let attempts = 0;
    const countAttempts: Interceptor = (next) => (req) => {
      attempts++;
      return next(req);
    };
    const transport = createConnectTransport({
      httpVersion: "1.1",
      baseUrl: `http://127.0.0.1:${port}`,
      interceptors: [retryUnavailableInterceptor, countAttempts],
    });

    await expect(createClientFor(SessionService)(transport).getSession({})).rejects.toMatchObject({
      code: Code.Unavailable,
    });
    expect(attempts).toBe(MAX_RETRIES + 1);
  });
});
