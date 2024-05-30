import { CallOptions, ClientMiddlewareCall, Metadata } from "nice-grpc-web";

export const authMiddleware = (token: string) =>
  async function* <Request, Response>(
    call: ClientMiddlewareCall<Request, Response>,
    options: CallOptions,
  ) {
    if (!options.metadata?.has("authorization")) {
      options.metadata ??= new Metadata();
      options.metadata?.set("authorization", `Bearer ${token}`);
    }

    return yield* call.next(call.request, options);
  };
