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
      label: 'Installation',
      items: [
          {
              type: 'category',
              label: 'CAOS Managed',
              items: [
                  'guides/installation/shared-cloud',
                  'guides/installation/managed-dedicated-instance'
              ],
          },
          {
              type: 'category',
              label: 'Self Managed',
              items: [
                  'guides/installation/crd',
                  'guides/installation/gitops',
                  'guides/installation/orbos'
              ],
          },
      ],
    },
    {
      type: 'category',
      label: 'Get to know ZITADEL',
      items: [
          'guides/usage/get-started',
          'guides/usage/organizations',
          'guides/usage/projects',
          'guides/usage/oauth-recommended-flows',
          'guides/usage/serviceusers',
          'guides/usage/access-zitadel-apis',
          'guides/usage/identity-brokering',
      ]
    }
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
      items: [
          'apis/openidoauth/endpoints',
          'apis/openidoauth/scopes',
          'apis/openidoauth/claims',
          'apis/openidoauth/authn-methods',
          'apis/openidoauth/grant-types'
      ],
    },
  ],
  concepts: [
    'concepts/introduction',
    'concepts/architecture',
    'concepts/principles',
  ]
};
