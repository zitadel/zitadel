function e2eSetLastUsernameMD(ctx, api) {
  const mgmt = api.v1.mgmt;
  const authRequest = ctx.v1.authRequest;
  const goCtx = ctx.v1.ctx;

  mgmt.setUserMetadata(goCtx, {
    id: authRequest.userID,
    key: 'last username used',
    value: authRequest.userName,
  });
}
