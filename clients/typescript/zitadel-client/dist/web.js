import {
  NewAuthorizationBearerInterceptor
} from "./chunk-H5CIUTXZ.js";

// src/web.ts
import { createGrpcWebTransport } from "@connectrpc/connect-web";
function createClientTransport(token, opts) {
  return createGrpcWebTransport({
    ...opts,
    interceptors: [...opts.interceptors || [], NewAuthorizationBearerInterceptor(token)]
  });
}
export {
  createClientTransport
};
//# sourceMappingURL=web.js.map