"use strict";Object.defineProperty(exports, "__esModule", {value: true});// src/interceptors.ts
function NewAuthorizationBearerInterceptor(token) {
  return (next) => (req) => {
    if (!req.header.get("Authorization")) {
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return next(req);
  };
}



exports.NewAuthorizationBearerInterceptor = NewAuthorizationBearerInterceptor;
//# sourceMappingURL=chunk-RFOVJJ5M.cjs.map