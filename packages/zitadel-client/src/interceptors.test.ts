import { Int32Value } from "@bufbuild/protobuf/wkt";
import { compileService } from "@bufbuild/protocompile";
import { createRouterTransport, HandlerContext } from "@connectrpc/connect";
import { describe, expect, test, vitest } from "vitest";
import { NewAuthorizationBearerInterceptor } from "./interceptors.js";

const TestService = compileService(`
  syntax = "proto3";
  package handwritten;
  service TestService {
    rpc Unary(Int32Value) returns (StringValue);
  }
  message Int32Value {
    int32 value = 1;
  }
  message StringValue {
    string value = 1;
  }
`);

describe("NewAuthorizationBearerInterceptor", () => {
  const transport = {
    interceptors: [NewAuthorizationBearerInterceptor("mytoken")],
  };

  test("injects the authorization token", async () => {
    const handler = vitest.fn((request: Int32Value, _context: HandlerContext) => {
      return { value: request.value.toString() };
    });

    const service = createRouterTransport(
      ({ rpc }) => {
        rpc(TestService.method.unary, handler);
      },
      { transport },
    );

    await service.unary(TestService.method.unary, undefined, undefined, {}, { value: 9001 });

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe("Bearer mytoken");
  });

  test("do not overwrite the previous authorization token", async () => {
    const handler = vitest.fn((request: Int32Value, _context: HandlerContext) => {
      return { value: request.value.toString() };
    });

    const service = createRouterTransport(
      ({ rpc }) => {
        rpc(TestService.method.unary, handler);
      },
      { transport },
    );

    await service.unary(
      TestService.method.unary,
      undefined,
      undefined,
      { Authorization: "Bearer somethingelse" },
      { value: 9001 },
    );

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe("Bearer somethingelse");
  });
});
