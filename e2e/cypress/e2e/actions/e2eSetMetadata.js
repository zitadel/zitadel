let console = require('zitadel/log')

function <IDENTIFIER>(ctx, api) {
  const mgmt = api.v1.mgmt;
  const authRequest = ctx.v1.authRequest;
  const goCtx = ctx.v1.ctx;

  // TODO: Why is this undefined?
  console.log('mgmt', mgmt.setUserMetadata)

  mgmt.setUserMetadata(goCtx, {
    id: authRequest.userID,
    key: `<KEY>`,
    value: `<VALUE>`,
  });
}
