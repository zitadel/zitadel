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
          to: 'docs/manuals',
          label: 'Manuals',
          position: 'left'
        },
        {
          to: 'docs/quickstarts',
          label: 'Quickstarts',
          position: 'left'
        },
        {
          to: 'docs/guides',
          label: 'Guides',
          position: 'left'
        },
        {
          to: 'docs/apis',
          label: 'APIs',
          position: 'left'
        },
        {
          to: 'docs/architecture',
          label: 'Architecture',
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
              label: 'GitHub',
              href: 'https://github.com/caos/zitadel/discussions',
            },
            {
              label: 'Twitter',
              href: 'https://twitter.com/zitadel_ch',
            },
          ],
        },
        {
          title: 'More',
          items: [
            {
              label: 'Blog',
              href: 'https://zitadel.ch/blog',
            },
            {
              label: 'GitHub',
              href: 'https://github.com/caos/zitadel',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} CAOS Built with Docusaurus.`,
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
