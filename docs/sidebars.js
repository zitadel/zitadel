module.exports = {
  manuals: [
    'manuals/introduction',
    {
      type: 'category',
      label: 'User',
      items: ['manuals/user'],
    },
    {
      type: 'category',
      label: 'Administrator',
      items: ['manuals/admin-managers'],
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
      items: ['guides/get-started', 'guides/organizations', 'guides/projects', 'guides/oauth-recommended-flows', 'guides/serviceusers', 'guides/access-zitadel-apis', 'guides/identity-brokering'],
    },
  ],
  apis: [
    'apis/introduction',
    'apis/domains',
    'apis/apis',
    {
      type: 'category',
      label: 'Proto API Definition',
      items: [
        'apis/proto/auth',
        'apis/proto/management',
        'apis/proto/admin',
        'apis/proto/org',
        'apis/proto/user',
        'apis/proto/app',
        'apis/proto/policy',
        'apis/proto/auth_n_key',
        'apis/proto/change',
        'apis/proto/idp',
        'apis/proto/member',
        'apis/proto/message',
        'apis/proto/object',
        'apis/proto/options',
      ],
    },
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