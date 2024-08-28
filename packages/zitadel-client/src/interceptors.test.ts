import { Int32Value, Int32ValueSchema, StringValueSchema } from "@bufbuild/protobuf/wkt";
import { createRouterTransport, HandlerContext } from "@connectrpc/connect";
import { describe, expect, test, vitest } from "vitest";
import { NewAuthorizationBearerInterceptor } from "./interceptors";

const TestService = {
  typeName: "handwritten.TestService",
  methods: {
    unary: {
      input: Int32ValueSchema,
      output: StringValueSchema,
      methodKind: "unary",
    },
  },
} as const;

describe.skip("NewAuthorizationBearerInterceptor", () => {
  const transport = {
    interceptors: [NewAuthorizationBearerInterceptor("mytoken")],
  };

  test("injects the authorization token", async () => {
    const handler = vitest.fn((request: Int32Value, context: HandlerContext) => {
      return { value: request.value.toString() };
    });

    const service = createRouterTransport(
      ({ service }) => {
        service(TestService, { unary: handler });
      },
      { transport },
    );

    await service.unary(TestService, TestService.methods.unary, undefined, undefined, {}, { value: 9001 });

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe("Bearer mytoken");
  });

  test("do not overwrite the previous authorization token", async () => {
    const handler = vitest.fn((request: Int32Value, context: HandlerContext) => {
      return { value: request.value.toString() };
    });

    const service = createRouterTransport(
      ({ service }) => {
        service(TestService, { unary: handler });
      },
      { transport },
    );

    await service.unary(
      TestService.methods.unary,
      undefined,
      undefined,
      { Authorization: "Bearer somethingelse" },
      { value: 9001 },
    );

    expect(handler).toBeCalled();
    expect(handler.mock.calls[0][1].requestHeader.get("Authorization")).toBe("Bearer somethingelse");
  });
});
