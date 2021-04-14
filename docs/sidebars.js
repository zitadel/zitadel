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
    'guides',
    {
      type: 'category',
      label: 'Get to know ZITADEL',
      items: ['guides/organizations', 'guides/projects'],
    },
  ],
  apis: [
    'apis'
  ],
  architecture: [
    'architecture',
    'architecture/principles',
  ]
};