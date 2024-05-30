import { describe, expect, test, vitest } from "vitest";
import { Int32Value, MethodKind, StringValue } from "@bufbuild/protobuf";
import { createRouterTransport, HandlerContext } from "@connectrpc/connect";
import { NewAuthorizationBearerInterceptor } from "./interceptors";

const TestService = {
  typeName: "handwritten.TestService",
  methods: {
    unary: {
      name: "Unary",
      I: Int32Value,
      O: StringValue,
      kind: MethodKind.Unary,
    },
  },
} as const;

describe("NewAuthorizationBearerInterceptor", () => {
  const transport = {
    interceptors: [NewAuthorizationBearerInterceptor("mytoken")],
  };

  test("injects the authorization token", async () => {
    const handler = vitest.fn(
      (request: Int32Value, context: HandlerContext) => {
        return { value: request.value.toString() };
      },
    );

    const service = createRouterTransport(
      ({ service }) => {
        service(TestService, { unary: handler });
      },
      { transport },
    );

    await service.unary(
      TestService,
      TestService.methods.unary,
      undefined,
      undefined,
      {},
      { value: 9001 },
    );

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe(
      "Bearer mytoken",
    );
  });

  test("do not overwrite the previous authorization token", async () => {
    const handler = vitest.fn(
      (request: Int32Value, context: HandlerContext) => {
        return { value: request.value.toString() };
      },
    );

    const service = createRouterTransport(
      ({ service }) => {
        service(TestService, { unary: handler });
      },
      { transport },
    );

    await service.unary(
      TestService,
      TestService.methods.unary,
      undefined,
      undefined,
      { Authorization: "Bearer somethingelse" },
      { value: 9001 },
    );

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe(
      "Bearer somethingelse",
    );
  });
});
