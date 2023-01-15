function e2eAddAuthErrorMetadata(ctx, api) {
  const mgmt = api.v1.mgmt;
  const authRequest = ctx.v1.authRequest;
  const goCtx = ctx.v1.ctx;
  const authMethod = ctx.v1.authMethod;
  const authError = ctx.v1.authError;

  mgmt.setUserMetadata(goCtx, {
    id: authRequest.userID,
    key: `${authMethod} error`,
    value: authError,
  });
}
