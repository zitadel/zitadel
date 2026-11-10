import { GrpcTransportOptions } from "@connectrpc/connect-node";
import { createGrpcWebTransport } from "@connectrpc/connect-web";
import { NewAuthorizationBearerInterceptor } from "./interceptors.js";

/**
 * Create a client transport using grpc web with the given token and configuration options.
 * @param token
 * @param opts
 */
export function createClientTransport(token: string, opts: GrpcTransportOptions) {
  return createGrpcWebTransport({
    ...opts,
    interceptors: [...(opts.interceptors || []), NewAuthorizationBearerInterceptor(token)],
  });
}
