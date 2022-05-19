/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: "ZITADEL Docs",
  trailingSlash: false,
  url: "https://docs.zitadel.com",
  baseUrl: "/",
  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "zitadel",
  projectName: "zitadel",
  scripts: [
    {
      src: "/proxy/js/script.js",
      async: true,
      defer: true,
      "data-domain": "docs.zitadel.com",
      "data-api": "/proxy/api/event",
    },
  ],
  themeConfig: {
    zoomSelector: ".markdown :not(em) > img",
    announcementBar: {
      id: 'documentation',
      content:
        'This page contains the documentation for ZITADEL version 2, if you are looking for version 1 please visit <a target="_blank" rel="noopener noreferrer" href="https://docs.zitadel.ch">https://docs.zitadel.ch</a>',
      backgroundColor: '#fafbfc',
      textColor: '#091E42',
      isCloseable: false,
    },
    navbar: {
      // title: 'ZITADEL',
      logo: {
        alt: "ZITADEL logo",
        src: "img/zitadel-logo-dark.svg",
        srcDark: "img/zitadel-logo-light.svg",
      },
      items: [
        {
          type: "doc",
          label: "Guides",
          docId: "guides/overview",
          position: "left",
        },
        {
          type: "doc",
          label: "Quickstarts",
          docId: "quickstarts/introduction",
          position: "left",
        },
        {
          type: "doc",
          label: "APIs",
          docId: "apis/introduction",
          position: "left",
        },
        {
          type: "doc",
          docId: "concepts/introduction",
          label: "Concepts",
          position: "left",
        },
        {
          type: "doc",
          docId: "manuals/introduction",
          label: "Help",
          position: "left",
        },
        {
          type: "doc",
          docId: "legal/introduction",
          label: "Legal",
          position: "left",
        },
        {
          href: "https://github.com/zitadel/zitadel",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      links: [
        {
          title: "Community",
          items: [
            {
              label: "GitHub Discussions",
              href: "https://github.com/zitadel/zitadel/discussions",
            },
            {
              label: "Twitter",
              href: "https://twitter.com/zitadel",
            },
            {
              label: "Linkedin",
              href: "https://www.linkedin.com/company/zitadel/",
            },
            {
              label: "Blog",
              href: "https://zitadel.com/blog",
            },
          ],
        },
        {
          title: "Company",
          items: [
            {
              label: "Team.",
              href: "https://zitadel.com/team",
            },
            {
              label: "Contact",
              href: "https://zitadel.com/contact/",
            },
            {
              label: "GitHub",
              href: "https://github.com/zitadel",
            },
            {
              label: "Status",
              href: "https://status.zitadel.ch/",
            },
            {
              label: "Terms and Conditions",
              href: "https://docs.zitadel.com/docs/legal/terms-of-service",
            },
            {
              label: "Privacy Policy",
              href: "https://docs.zitadel.com/docs/legal/privacy-policy",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} ZITADEL - Built with Docusaurus.`,
    },
    algolia: {
      appId: "8H6ZKXENLO",
      apiKey: "c3899716db098111f5e89c8987b9c427",
      indexName: "zitadel",
    },
    prism: {
      additionalLanguages: ["csharp", "dart", "groovy"],
    },
  },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/zitadel/zitadel/edit/main/docs/",
          remarkPlugins: [require("mdx-mermaid")],
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],
  plugins: [require.resolve("plugin-image-zoom")],
};
