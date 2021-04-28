module.exports = {
  manuals: [
    'manuals/introduction',
    {
      type: 'category',
      label: 'User',
      items: ['manuals/user-register', 'manuals/user-login', 'manuals/user-password', 'manuals/user-factors', 'manuals/user-email', 'manuals/user-phone', 'manuals/user-social-login', ],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Administrator',
      items: ['manuals/admin-managers'],
      collapsed: false,
    },
  ],
  quickstarts: [
    'quickstarts/introduction',
    {
      type: 'category',
      label: 'Single Page Applications',
      items: ['quickstarts/angular'],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Backends',
      items: ['quickstarts/go', 'quickstarts/dot-net'],
      collapsed: false,
    },
    {
      type: 'category',
      label: 'Identity Aware Proxy',
      items: ['quickstarts/oauth2-proxy'],
      collapsed: false,
    }
  ],
  guides: [
    'guides/introduction',
    {
      type: 'category',
      label: 'Get to know ZITADEL',
      collapsed: false,
      items: [
          'guides/usage/get-started',
          'guides/usage/organizations',
          'guides/usage/projects',
          'guides/usage/oauth-recommended-flows',
          'guides/usage/serviceusers',
          'guides/usage/access-zitadel-apis',
          'guides/usage/identity-brokering',
      ],
    },
      {
          type: 'category',
          label: 'Installation',
          collapsed: false,
          items: [
              {
                  type: 'category',
                  label: 'CAOS Managed',
                  collapsed: true,
                  items: [
                      'guides/installation/shared-cloud',
                      'guides/installation/managed-dedicated-instance'
                  ],
              },
              {
                  type: 'category',
                  label: 'Self Managed',
                  collapsed: true,
                  items: [
                      'guides/installation/crd',
                      'guides/installation/gitops',
                      'guides/installation/orbos'
                  ],
              },
          ],
      }
  ],
  apis: [
    'apis/introduction',
    'apis/domains',
    'apis/apis',
    {
      type: 'category',
      label: 'Proto API Definition',
      collapsed: false,
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
      collapsed: false,
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
