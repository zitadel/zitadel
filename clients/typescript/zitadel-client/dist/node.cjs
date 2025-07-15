"use strict";Object.defineProperty(exports, "__esModule", {value: true}); function _nullishCoalesce(lhs, rhsFn) { if (lhs != null) { return lhs; } else { return rhsFn(); } }

var _chunkRFOVJJ5Mcjs = require('./chunk-RFOVJJ5M.cjs');

// src/node.ts
var _connectnode = require('@connectrpc/connect-node');
var _jose = require('jose');
function createServerTransport(token, opts) {
  return _connectnode.createGrpcTransport.call(void 0, {
    ...opts,
    interceptors: [...opts.interceptors || [], _chunkRFOVJJ5Mcjs.NewAuthorizationBearerInterceptor.call(void 0, token)]
  });
}
async function newSystemToken({
  audience,
  subject,
  key,
  expirationTime
}) {
  return await new (0, _jose.SignJWT)({}).setProtectedHeader({ alg: "RS256" }).setIssuedAt().setExpirationTime(_nullishCoalesce(expirationTime, () => ( "1h"))).setIssuer(subject).setSubject(subject).setAudience(audience).sign(await _jose.importPKCS8.call(void 0, key, "RS256"));
}



exports.createServerTransport = createServerTransport; exports.newSystemToken = newSystemToken;
//# sourceMappingURL=node.cjs.map