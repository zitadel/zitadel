/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: "ZITADEL Docs",
  trailingSlash: false,
  url: "https://zitadel.com",
  baseUrl: "/docs",
  onBrokenLinks: "warn",
  onBrokenAnchors: "warn",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "zitadel",
  projectName: "zitadel",
  scripts: [
    {
      src: "/docs/proxy/js/script.js",
      async: true,
      defer: true,
      "data-domain": "zitadel.com",
      "data-api": "/docs/proxy/api/event",
    },
  ],
  customFields: {
    description:
      "Documentation for ZITADEL - Identity infrastructure, simplified for you.",
  },
  themeConfig: {
    metadata: [
      {
        name: "keywords",
        content:
          "zitadel, documentation, jwt, saml, oauth2, authentication, serverless, login, auth, authorization, sso, openid-connect, oidc, mfa, 2fa, passkeys, fido2, docker",
      },
      {
        property: "og:type",
        content: "website",
      },
      { property: "og:url", content: "https://www.zitadel.com/docs" },
      {
        property: "og:image",
        content: "https://www.zitadel.com/docs/img/preview.png",
      },
      { property: "twitter:card", content: "summary_large_image" },
      { property: "twitter:url", content: "https://www.zitadel.com/docs" },
      { property: "twitter:title", content: "ZITADEL Docs" },
      {
        property: "twitter:image",
        content: "https://www.zitadel.com/docs/img/preview.png",
      },
    ],
    zoom: {
      selector: ".markdown :not(em) > img",
      background: {
        light: "rgb(243, 244, 246)",
        dark: "rgb(55, 59, 82)",
      },
      // options you can specify via https://github.com/francoischalifour/medium-zoom#usage
      config: {},
    },
    navbar: {
      // title: 'ZITADEL',
      logo: {
        alt: "ZITADEL logo",
        src: "img/zitadel-logo-dark.svg",
        srcDark: "img/zitadel-logo-light.svg",
        href: "https://zitadel.com",
        target: "_blank",
      },
      items: [
        {
          type: "doc",
          label: "🚀 Quick Start",
          docId: "guides/start/quickstart",
          position: "left",
        }, 
        {
          type: "doc",
          label: "Documentation",
          docId: "guides/overview",
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
          label: "Self-Hosting",
          docId: "self-hosting/deploy/overview",
          position: "left",
        },
        {
          to: "/redocusaurus/v2",
          label: "/redocusaurus/v2",
          position: "left",
        },
        {
          type: "doc",
          docId: "legal",
          label: "Legal",
          position: "right",
        },
        {
          type: "html",
          position: "right",
          value:
            '<a href="https://github.com/zitadel/zitadel/discussions" style="text-decoration: none; width: 20px; height: 24px; display: flex"><i class="las la-comments"></i></a>',
        },
        {
          type: "html",
          position: "right",
          value:
            '<a href="https://github.com/zitadel/zitadel" style="text-decoration: none; width: 20px; height: 24px; display: flex"><i class="lab la-github"></i></a>',
        },
        {
          type: "html",
          position: "right",
          value:
            '<a href="https://zitadel.com/chat" style="text-decoration: none; width: 20px; height: 24px; display: flex; margin: 0 .5rem 0 0"><i class="lab la-discord"></i></a>',
        },
      ],
    },
    footer: {
      links: [
        {
          title: "Community",
          items: [
            {
              label: "Chat",
              href: "https://zitadel.com/chat",
            },
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
          title: "Legal",
          items: [
            {
              label: "Terms and Conditions",
              href: "/legal/terms-of-service",
            },
            {
              label: "Privacy Policy",
              href: "/legal/policies/privacy-policy",
            },
          ],
        },
        {
          title: "About",
          items: [
            {
              label: "Website",
              href: "https://zitadel.com",
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
              href: "https://status.zitadel.com/",
            }
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} ZITADEL Docs - Built with Docusaurus.`,
    },
    algolia: {
      appId: "8H6ZKXENLO",
      apiKey: "124fe1c102a184bc6fc70c75dc84f96f",
      indexName: "zitadel",
      selector: "div#",
    },
    prism: {
      additionalLanguages: ["bash", "csharp", "java", "php", "ruby", "scala"],
    },
    colorMode: {
      defaultMode: "dark",
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    codeblock: {
      showGithubLink: true,
      githubLinkLabel: 'View on GitHub',
      showRunmeLink: false,
      runmeLinkLabel: 'Checkout via Runme'
    },
  },
  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          routeBasePath: "/",
          sidebarPath: require.resolve("./sidebars.js"),
          showLastUpdateAuthor: true,
          showLastUpdateTime: true,
          editUrl: "https://github.com/zitadel/zitadel/edit/main/docs/",
          
          docItemComponent:  '@theme/ApiItem'
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      })
    ],
    [
      'redocusaurus',
      {
        specs: [
          { 
            id: 'v2',
            spec: './static/zitadel.swagger.yaml',
            route: '/redocusaurus/v2'
          },
        ],
        theme: {},
        }
    ]
  ],
  plugins: [
    [
    '@scalar/docusaurus',
    {
      label: '/scalar/v2',
      route: '/docs/scalar/v2',
      configuration: {
        spec: {
          // Put the URL to your OpenAPI document here:
          url: '/docs/zitadel.swagger.yaml',
        },
      },
    },
    ],
    require.resolve("docusaurus-plugin-image-zoom"),
    async function myPlugin(context, options) {
      return {
        name: "docusaurus-tailwindcss",
        configurePostCss(postcssOptions) {
          // Appends TailwindCSS and AutoPrefixer.
          postcssOptions.plugins.push(require("tailwindcss"));
          postcssOptions.plugins.push(require("autoprefixer"));
          return postcssOptions;
        },
      };
    },
  ],
  themes: [ "docusaurus-theme-github-codeblock", "docusaurus-theme-openapi-docs"],
  future: {
    // See https://docusaurus.io/blog/releases/3.6
    experimental_faster: true,
  },
};
