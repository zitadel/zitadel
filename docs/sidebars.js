module.exports = {
  manuals: [
    {
      type: 'category',
      label: 'User',
      items: [
          'manuals/user',
      ],
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
      items: [
          'guides/usage/organizations',
          'guides/usage/projects',
          'guides/usage/serviceusers',
          'guides/usage/oauth-recommended-flows',
          'guides/usage/identity-brokering'
      ],
    },
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
