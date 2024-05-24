import {
  CallOptions,
  ClientMiddleware,
  ClientMiddlewareCall,
  Metadata,
} from "nice-grpc";

export function authMiddleware(token: string): ClientMiddleware {
  return async function* <Request, Response>(
    call: ClientMiddlewareCall<Request, Response>,
    options: CallOptions,
  ) {
    if (!options.metadata?.has("authorization")) {
      options.metadata ??= new Metadata();
      options.metadata?.set("authorization", `Bearer ${token}`);
    }

    return yield* call.next(call.request, options);
  };
}

export const orgMetadata = (orgId: string) =>
  new Metadata({ "x-zitadel-orgid": orgId });
