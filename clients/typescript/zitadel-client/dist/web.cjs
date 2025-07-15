"use strict";Object.defineProperty(exports, "__esModule", {value: true});

var _chunkRFOVJJ5Mcjs = require('./chunk-RFOVJJ5M.cjs');

// src/web.ts
var _connectweb = require('@connectrpc/connect-web');
function createClientTransport(token, opts) {
  return _connectweb.createGrpcWebTransport.call(void 0, {
    ...opts,
    interceptors: [...opts.interceptors || [], _chunkRFOVJJ5Mcjs.NewAuthorizationBearerInterceptor.call(void 0, token)]
  });
}


exports.createClientTransport = createClientTransport;
//# sourceMappingURL=web.cjs.map