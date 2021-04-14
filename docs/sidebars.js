module.exports = {
  manuals: [
    'manuals',
    {
      type: 'category',
      label: 'User',
      items: ['manuals/user'],
    },
  ],
  quickstarts: [
    'quickstarts/quickstarts',
    {
      type: 'category',
      label: 'Single Page Applications',
      items: ['quickstarts/angular'],
    },
    {
      type: 'category',
      label: 'Identity Aware Proxy',
      items: ['quickstarts/oauth2-proxy'],
    }
  ],
  guides: [
    {
      type: 'category',
      label: 'Get to know ZITADEL',
      items: ['guides/introduction', 'guides/organizations', 'guides/projects', 'guides/serviceusers', 'guides/oauth-recommended-flows', 'guides/identity-brokering'],
    },
  ],
  apis: [
    'apis/apis',
    'apis/domains',
    'apis/authn',
    'apis/admin',
    'apis/mgmt',
    {
      type: 'category',
      label: 'OpenID Connect & OAuth',
      items: ['apis/openidoauth/endpoints', 'apis/openidoauth/scopes', 'apis/openidoauth/claims', 'apis/openidoauth/authn-methods', 'apis/openidoauth/grant-types'],
    },
  ],
  architecture: [
    'architecture',
    'architecture/principles',
  ]
};