/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: 'ZITADEL Docs',
  url: 'https://docs.zitadel.ch',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'caos', // Usually your GitHub org/user name.
  projectName: 'zitadel', // Usually your repo name.
  themeConfig: {
    navbar: {
      title: 'Docs',
      logo: {
        alt: 'ZITADEL logo',
        src: 'img/zitadel-logo-light.svg',
      },
      items: [
        {
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Docs',
          position: 'left',
        },
        {to: 'blog', label: 'Blog', position: 'left'},
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
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: 'docs/',
            },
          ],
        },
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
          // Please change this to your repo.
          editUrl:
            'https://github.com/caos/zitadel/edit/main/docs/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
