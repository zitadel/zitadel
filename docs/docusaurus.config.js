/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'ZITADEL Docs',
  url: 'https://docs.zitadel.ch',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'caos',
  projectName: 'zitadel',
  themeConfig: {
    navbar: {
      title: 'ZITADEL',
      logo: {
        alt: 'ZITADEL logo',
        src: 'img/zitadel-logo-solo-darkdesign.svg',
      },
      items: [
        {
          type: 'doc',
          docId: 'manuals/introduction',
          label: 'Manuals',
          position: 'left'
        },
        {
          type: 'doc',
          label: 'Quickstarts',
          docId: 'quickstarts/introduction',
          position: 'left'
        },
        {
          type: 'doc',
          label: 'Guides',
          docId: 'guides/introduction',
          position: 'left'
        },
        {
          type: 'doc',
          label: 'APIs',
          docId: 'apis/introduction',
          position: 'left'
        },
        {
          type: 'doc',
          docId: 'concepts/introduction',
          label: 'Concepts',
          position: 'left'
        },
        {
          href: 'https://github.com/caos/zitadel',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub Discussions',
              href: 'https://github.com/caos/zitadel/discussions',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/zitadel_ch',
            },
            {
              label: 'Linkedin',
              href: 'https://www.linkedin.com/company/caos-ag/',
            },
            {
              label: 'Blog',
              href: 'https://zitadel.ch/blog',
            },
          ],
        },
        {
          title: 'Company',
          items: [
            {
              label: 'About CAOS Ltd.',
              href: 'https://caos.ch/#about',
            },
            {
              label: 'Contact',
              href: 'https://zitadel.ch/contact/',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/caos',
            },
            {
              label: 'Status',
              href: 'https://status.zitadel.ch/',
            },
            {
              label: 'Terms and Conditions',
              href: 'https://zitadel.ch/pdf/tos.pdf',
            },
            {
              label: 'Privacy Policy',
              href: 'https://zitadel.ch/pdf/privacy.pdf',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} CAOS Ltd. Built with Docusaurus.`,
    },
    algolia: {
      apiKey: 'bff480bce03126c2d348345647854e91',
      indexName: 'zitadel'
    },
  },
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl:
            'https://github.com/caos/zitadel/edit/main/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
  plugins: [
    [
      'docusaurus-plugin-plausible',
      {
        domain: 'docs.zitadel.ch',
      },
    ]
  ],
};
