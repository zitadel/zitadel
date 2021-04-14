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
    'quickstarts',
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
      items: ['guides/introduction', 'guides/organizations', 'guides/projects'],
    },
  ],
  apis: [
    'apis',
    {
      type: 'category',
      label: 'Administration',
      items: ['guides/organizations', 'guides/projects'],
    },
    {
      type: 'category',
      label: 'Authentication',
      items: ['guides/organizations', 'guides/projects'],
    },
    {
      type: 'category',
      label: 'Management',
      items: ['guides/organizations', 'guides/projects'],
    },
    {
      type: 'category',
      label: 'OpenID Connect & OAuth',
      items: ['guides/organizations', 'guides/projects'],
    },
  ],
  architecture: [
    'architecture',
    'architecture/principles',
  ]
};