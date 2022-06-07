/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: "ZITADEL Docs",
  trailingSlash: false,
  url: "https://docs-v1.zitadel.com",
  baseUrl: "/",
  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "zitadel",
  projectName: "zitadel",
  scripts: [
    {
      src: "/proxy/js/script.js",
      async: true,
      defer: true,
      "data-domain": "docs.zitadel.ch",
      "data-api": "/proxy/api/event",
    },
  ],
  themeConfig: {
    zoomSelector: ".markdown :not(em) > img",
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
              href: "https://twitter.com/zitadel_ch",
            },
            {
              label: "Linkedin",
              href: "https://www.linkedin.com/company/zitadel/",
            },
            {
              label: "Blog",
              href: "https://zitadel.ch/blog",
            },
          ],
        },
        {
          title: "Company",
          items: [
            {
              label: "About CAOS Ltd.",
              href: "https://caos.ch/#about",
            },
            {
              label: "Contact",
              href: "https://zitadel.ch/contact/",
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
              href: "https://docs.zitadel.ch/docs/legal/terms-of-service",
            },
            {
              label: "Privacy Policy",
              href: "https://docs.zitadel.ch/docs/legal/privacy-policy",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} CAOS Ltd. Built with Docusaurus.`,
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
