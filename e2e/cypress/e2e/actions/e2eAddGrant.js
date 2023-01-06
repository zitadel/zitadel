function e2eAddGrant(ctx, api) {
  api.userGrants.push({
    projectID: '<PROJECT_ID>',
    roles: ['<ROLE_KEY>'],
  });
}
