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
      items: ['quickstarts/vue', 'quickstarts/angular'],
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
      items: ['guides/introduction', 'guides/organizations', 'guides/projects', 'guides/serviceusers'],
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
      items: ['apis/openidoauth/endpoints', 'apis/openidoauth/scopes', 'apis/openidoauth/claims'],
    },
  ],
  architecture: [
    'architecture',
    'architecture/principles',
  ]
};