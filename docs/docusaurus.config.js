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
      copyright: `Copyright © ${new Date().getFullYear()} ZITADEL Docs - Built with Docusaurus.`,
    },
    algolia: {
      appId: "1FD5EU44EF",
      apiKey: "2ddc6c35930ec608baf832bf07b5e150",
      indexName: "netlify_13929a6d-0c07-4205-af9b-42121333ec5d_v2-alpha-readme_all",
      contextualSearch: false,
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
          editUrl: "https://github.com/zitadel/zitadel/edit/v2-alpha/docs/",
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
