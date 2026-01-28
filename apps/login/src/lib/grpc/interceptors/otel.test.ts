import { SpanKind, SpanStatusCode } from "@opentelemetry/api";
import { Resource } from "@opentelemetry/resources";
import {
  InMemorySpanExporter,
  SimpleSpanProcessor,
} from "@opentelemetry/sdk-trace-base";
import { NodeTracerProvider } from "@opentelemetry/sdk-trace-node";
import { afterAll, beforeEach, describe, expect, it } from "vitest";
import { otelGrpcInterceptor } from "./otel";

const exporter = new InMemorySpanExporter();
const provider = new NodeTracerProvider({
  resource: new Resource({ "service.name": "test-zitadel-client" }),
  spanProcessors: [new SimpleSpanProcessor(exporter)],
});
provider.register();

afterAll(() => provider.shutdown());

describe("otelGrpcInterceptor", () => {
  beforeEach(() => exporter.reset());

  it("injects traceparent header into gRPC requests", async () => {
    const capturedHeaders: Record<string, string> = {};

    await otelGrpcInterceptor(async () => ({ success: true }))({
      service: { typeName: "zitadel.session.v2.SessionService" },
      method: { name: "CreateSession" },
      url: "https://api.zitadel.example.com/zitadel.session.v2/CreateSession",
      header: {
        set: (key: string, value: string) => {
          capturedHeaders[key.toLowerCase()] = value;
        },
      },
    });

    expect(capturedHeaders["traceparent"]).toMatch(
      /^[0-9a-f]{2}-[0-9a-f]{32}-[0-9a-f]{16}-[0-9a-f]{2}$/,
    );
  });

  it("creates client spans with RPC attributes", async () => {
    await otelGrpcInterceptor(async () => ({ user: { id: "123" } }))({
      service: { typeName: "zitadel.user.v2.UserService" },
      method: { name: "GetUserByID" },
      url: "https://api.zitadel.example.com/zitadel.user.v2/GetUserByID",
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.user.v2.UserService/GetUserByID",
      kind: SpanKind.CLIENT,
      attributes: {
        "rpc.system": "grpc",
        "rpc.service": "zitadel.user.v2.UserService",
        "rpc.method": "GetUserByID",
        "server.address": "api.zitadel.example.com",
      },
    });
  });

  it("records error spans with gRPC status codes", async () => {
    const error = Object.assign(new Error("Session not found"), { code: 5 });

    await expect(
      otelGrpcInterceptor(async () => {
        throw error;
      })({
        service: { typeName: "zitadel.session.v2.SessionService" },
        method: { name: "DeleteSession" },
        url: "https://api.zitadel.example.com",
        header: { set: () => {} },
      }),
    ).rejects.toThrow("Session not found");

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.session.v2.SessionService/DeleteSession",
      status: { code: SpanStatusCode.ERROR, message: "Session not found" },
      attributes: { "rpc.grpc.status_code": 5 },
    });
    expect(spans[0].events).toHaveLength(1);
    expect(spans[0].events[0].name).toBe("exception");
  });

  it("handles missing service/method gracefully", async () => {
    await otelGrpcInterceptor(async () => ({ ok: true }))({
      url: "https://api.zitadel.example.com",
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "unknown/unknown",
      attributes: { "rpc.service": "unknown", "rpc.method": "unknown" },
    });
  });

  it("returns the response from the next handler", async () => {
    const expectedResponse = { data: "test-data", items: [1, 2, 3] };

    const result = await otelGrpcInterceptor(async () => expectedResponse)({
      service: { typeName: "TestService" },
      method: { name: "TestMethod" },
      url: "https://test.example.com",
      header: { set: () => {} },
    });

    expect(result).toEqual(expectedResponse);
  });

  it("records error spans without gRPC status code", async () => {
    const error = new Error("Network failure");

    await expect(
      otelGrpcInterceptor(async () => {
        throw error;
      })({
        service: { typeName: "zitadel.user.v2.UserService" },
        method: { name: "ListUsers" },
        url: "https://api.zitadel.example.com",
        header: { set: () => {} },
      }),
    ).rejects.toThrow("Network failure");

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "zitadel.user.v2.UserService/ListUsers",
      status: { code: SpanStatusCode.ERROR, message: "Network failure" },
    });
    expect(spans[0].attributes).not.toHaveProperty("rpc.grpc.status_code");
  });

  it("handles malformed URL gracefully", async () => {
    await otelGrpcInterceptor(async () => ({ ok: true }))({
      service: { typeName: "TestService" },
      method: { name: "TestMethod" },
      url: "not-a-valid-url",
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].attributes["server.address"]).toBe("unknown");
  });

  it("handles missing URL gracefully", async () => {
    await otelGrpcInterceptor(async () => ({ ok: true }))({
      service: { typeName: "TestService" },
      method: { name: "TestMethod" },
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].attributes["server.address"]).toBe("unknown");
  });

  it("handles empty string URL gracefully", async () => {
    await otelGrpcInterceptor(async () => ({ ok: true }))({
      service: { typeName: "TestService" },
      method: { name: "TestMethod" },
      url: "",
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].attributes["server.address"]).toBe("unknown");
  });

  it("uses empty strings as-is for service/method names", async () => {
    await otelGrpcInterceptor(async () => ({ ok: true }))({
      service: { typeName: "" },
      method: { name: "" },
      url: "https://api.example.com",
      header: { set: () => {} },
    });

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0]).toMatchObject({
      name: "/",
      attributes: { "rpc.service": "", "rpc.method": "" },
    });
  });

  it("handles non-Error exceptions", async () => {
    await expect(
      otelGrpcInterceptor(async () => {
        throw "string error";
      })({
        service: { typeName: "TestService" },
        method: { name: "TestMethod" },
        url: "https://api.example.com",
        header: { set: () => {} },
      }),
    ).rejects.toBe("string error");

    const spans = exporter.getFinishedSpans();
    expect(spans).toHaveLength(1);
    expect(spans[0].status.code).toBe(SpanStatusCode.ERROR);
  });
});
