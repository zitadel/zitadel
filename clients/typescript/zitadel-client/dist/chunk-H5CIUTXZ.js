// src/interceptors.ts
function NewAuthorizationBearerInterceptor(token) {
  return (next) => (req) => {
    if (!req.header.get("Authorization")) {
      req.header.set("Authorization", `Bearer ${token}`);
    }
    return next(req);
  };
}

export {
  NewAuthorizationBearerInterceptor
};
//# sourceMappingURL=chunk-H5CIUTXZ.js.map