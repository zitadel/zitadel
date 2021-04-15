module.exports = {
  manuals: [
    'manuals/introduction',
    {
      type: 'category',
      label: 'User',
      items: ['manuals/user'],
    },
  ],
  quickstarts: [
    'quickstarts/introduction',
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
    'guides/introduction',
    {
      type: 'category',
      label: 'Get to know ZITADEL',
      items: ['guides/get-started', 'guides/organizations', 'guides/projects', 'guides/serviceusers', 'guides/oauth-recommended-flows', 'guides/identity-brokering'],
    },
  ],
  apis: [
    'apis/introduction',
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
  concepts: [
    'concepts/introduction',
    'concepts/architecture',
    'concepts/principles',
  ]
};