import { Code, ConnectError, createRouterTransport } from "@connectrpc/connect";
import { SpanKind, SpanStatusCode } from "@opentelemetry/api";
import { resourceFromAttributes } from "@opentelemetry/resources";
import { InMemorySpanExporter, SimpleSpanProcessor } from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { createClientFor } from "@zitadel/client";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { afterAll, beforeEach, describe, expect, it } from "vitest";
import { errorClassificationInterceptor } from "./error-classification";
import { otelGrpcInterceptor } from "./otel";

const exporter = new InMemorySpanExporter();
const provider = new NodeTracerProvider({
  resource: resourceFromAttributes({ "service.name": "test-zitadel-client" }),
  spanProcessors: [new SimpleSpanProcessor(exporter)],
});
provider.register();

afterAll(() => provider.shutdown());

describe("otelGrpcInterceptor", () => {
  beforeEach(() => exporter.reset());

  it("injects traceparent header into gRPC requests", async () => {
    let capturedHeaders: Headers | undefined;

    const mockTransport = createRouterTransport(() => {}, {
      transport: {
        interceptors: [
          otelGrpcInterceptor,
          (next) => (req) => {
            capturedHeaders = req.header;
            return next(req);
          },
        ],
      },
    });

    const client = createClientFor(SessionService)(mockTransport);
    await expect(client.createSession({})).rejects.toThrow();
    expect(capturedHeaders!.get("traceparent")!).toMatch(/^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$/);
  });

  it("creates client spans with RPC attributes", async () => {
    const mockTransport = createRouterTransport(() => {}, {
      transport: {
        interceptors: [otelGrpcInterceptor],
      },
    });

    const client = createClientFor(UserService)(mockTransport);
    await expect(client.getUserByID({})).rejects.toThrow();

    const spans = exporter.getFinishedSpans();

    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.user.v2.UserService/GetUserByID",
      kind: SpanKind.CLIENT,
      attributes: {
        "rpc.system": "grpc",
        "rpc.service": "zitadel.user.v2.UserService",
        "rpc.method": "GetUserByID",
        "server.address": "in-memory",
      },
    });
  });

  it("records error spans with gRPC status codes and classification", async () => {
    const mockTransport = createRouterTransport(
      ({ service }) => {
        service(SessionService, {
          deleteSession: () => {
            throw new ConnectError("Session not found", Code.NotFound);
          },
        });
      },
      {
        transport: {
          interceptors: [otelGrpcInterceptor, errorClassificationInterceptor],
        },
      },
    );

    const client = createClientFor(SessionService)(mockTransport);
    await expect(client.deleteSession({})).rejects.toThrow();

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.session.v2.SessionService/DeleteSession",
      status: {
        code: SpanStatusCode.ERROR,
        // https://github.com/connectrpc/connect-es/blob/beb1864a45786f01c8ecb76c93a08ac14483885a/packages/connect/src/connect-error.ts#L32
        // code prefix gets added by connect
        message: "[not_found] Session not found",
      },
      attributes: {
        "rpc.grpc.status_code": 5,
        "error.is_user_error": true,
        "http.status_code": 404,
      },
    });
    expect(spans[0].events).toHaveLength(1);
    expect(spans[0].events[0].name).toBe("exception");
  });

  it("returns the response from the next handler", async () => {
    const expectedResponse = {
      details: {},
      sessionId: "test-session-id",
      sessionToken: "test-token",
    };

    const mockTransport = createRouterTransport(
      ({ service }) => {
        service(SessionService, {
          createSession: () => expectedResponse,
        });
      },
      {
        transport: {
          interceptors: [otelGrpcInterceptor],
        },
      },
    );

    const client = createClientFor(SessionService)(mockTransport);
    const result = await client.createSession({});
    expect(result).toMatchObject(expectedResponse);
  });

  it("handles in-memory transport URL gracefully", async () => {
    const mockTransport = createRouterTransport(() => {}, {
      transport: {
        interceptors: [otelGrpcInterceptor],
      },
    });

    const client = createClientFor(SessionService)(mockTransport);
    await expect(client.createSession({})).rejects.toThrow();

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].attributes["server.address"]).toBeDefined();
  });

  it("sets correct service and method names from protobuf definitions", async () => {
    const mockTransport = createRouterTransport(() => {}, {
      transport: {
        interceptors: [otelGrpcInterceptor],
      },
    });

    const client = createClientFor(SessionService)(mockTransport);
    await expect(client.deleteSession({})).rejects.toThrow();

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.session.v2.SessionService/DeleteSession",
      attributes: {
        "rpc.service": "zitadel.session.v2.SessionService",
        "rpc.method": "DeleteSession",
      },
    });
  });

  it("handles non-Error exceptions", async () => {
    const mockTransport = createRouterTransport(
      ({ service }) => {
        service(SessionService, {
          createSession: () => {
            throw "string error";
          },
        });
      },
      {
        transport: {
          interceptors: [otelGrpcInterceptor],
        },
      },
    );

    const client = createClientFor(SessionService)(mockTransport);
    await expect(client.createSession({})).rejects.toThrow();

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].status.code).toBe(SpanStatusCode.ERROR);
  });
});
