/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: "ZITADEL Docs",
  trailingSlash: false,
  url: "https://zitadel.com",
  baseUrl: "/docs",
  onBrokenLinks: "warn",
  onBrokenAnchors: "warn",
  onBrokenMarkdownLinks: "throw",
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
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} ZITADEL Docs - Built with Docusaurus.`,
    },
    prism: {
      additionalLanguages: [
        "csharp",
        "dart",
        "groovy",
        "regex",
        "java",
        "php",
        "python",
        "protobuf",
        "json",
        "bash",
      ],
    },
    colorMode: {
      defaultMode: "dark",
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
    codeblock: {
      showGithubLink: true,
      githubLinkLabel: "View on GitHub",
      showRunmeLink: false,
      runmeLinkLabel: "Checkout via Runme",
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
          remarkPlugins: [require("mdx-mermaid")],
          docItemComponent: "@theme/ApiItem",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      }),
    ],
  ],
  plugins: [
    [
      "@inkeep/cxkit-docusaurus",
      {
        SearchBar: {
          baseSettings: {
            apiKey: process.env.INKEEP_API_KEY,
            primaryBrandColor: "#ff2069",
            organizationDisplayName: "ZITADEL",
          },
        },
        SearchSettings: {
          tabs: ["All", "Docs", "GitHub", "Forums", "Discord"],
        },
      },
    ],
    [
      "@signalwire/docusaurus-plugin-llms-txt",
      {
        depth: 3,
        logLevel: 1,
        content: {
          excludeRoutes: ["/search"],
          enableMarkdownFiles: true,
        },
      },
    ],
    [
      "docusaurus-plugin-openapi-docs",
      {
        id: "apiDocs",
        docsPluginId: "classic",
        config: {
          auth: {
            specPath: ".artifacts/openapi/zitadel/auth.swagger.json",
            outputDir: "docs/apis/resources/auth",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          mgmt: {
            specPath: ".artifacts/openapi/zitadel/management.swagger.json",
            outputDir: "docs/apis/resources/mgmt",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          admin: {
            specPath: ".artifacts/openapi/zitadel/admin.swagger.json",
            outputDir: "docs/apis/resources/admin",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          system: {
            specPath: ".artifacts/openapi/zitadel/system.swagger.json",
            outputDir: "docs/apis/resources/system",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          user_v2: {
            specPath:
              ".artifacts/openapi/zitadel/user/v2/user_service.swagger.json",
            outputDir: "docs/apis/resources/user_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          session_v2: {
            specPath:
              ".artifacts/openapi/zitadel/session/v2/session_service.swagger.json",
            outputDir: "docs/apis/resources/session_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          oidc_v2: {
            specPath:
              ".artifacts/openapi/zitadel/oidc/v2/oidc_service.swagger.json",
            outputDir: "docs/apis/resources/oidc_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          saml_v2: {
            specPath:
              ".artifacts/openapi/zitadel/saml/v2/saml_service.swagger.json",
            outputDir: "docs/apis/resources/saml_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          settings_v2: {
            specPath:
              ".artifacts/openapi/zitadel/settings/v2/settings_service.swagger.json",
            outputDir: "docs/apis/resources/settings_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          action_v2: {
            specPath:
              ".artifacts/openapi/zitadel/action/v2beta/action_service.swagger.json",
            outputDir: "docs/apis/resources/action_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          webkey_v2: {
            specPath:
              ".artifacts/openapi3/zitadel/webkey/v2/webkey_service.openapi.yaml",
            outputDir: "docs/apis/resources/webkey_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          feature_v2: {
            specPath:
              ".artifacts/openapi/zitadel/feature/v2/feature_service.swagger.json",
            outputDir: "docs/apis/resources/feature_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          org_v2: {
            specPath:
              ".artifacts/openapi/zitadel/org/v2/org_service.swagger.json",
            outputDir: "docs/apis/resources/org_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          idp_v2: {
            specPath:
              ".artifacts/openapi/zitadel/idp/v2/idp_service.swagger.json",
            outputDir: "docs/apis/resources/idp_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          org_v2beta: {
            specPath:
              ".artifacts/openapi3/zitadel/org/v2beta/org_service.openapi.yaml",
            outputDir: "docs/apis/resources/org_service_v2beta",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          project_v2beta: {
            specPath:
              ".artifacts/openapi3/zitadel/project/v2beta/project_service.openapi.yaml",
            outputDir: "docs/apis/resources/project_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          application_v2: {
            specPath: ".artifacts/openapi3/zitadel/app/v2beta/app_service.openapi.yaml",
            outputDir: "docs/apis/resources/application_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          instance_v2: {
            specPath:
              ".artifacts/openapi3/zitadel/instance/v2beta/instance_service.openapi.yaml",
            outputDir: "docs/apis/resources/instance_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          authorization_v2: {
            specPath:
              ".artifacts/openapi3/zitadel/authorization/v2beta/authorization_service.openapi.yaml",
            outputDir: "docs/apis/resources/authorization_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
          },
          internal_permission_v2: {
            specPath:
              ".artifacts/openapi3/zitadel/internal_permission/v2beta/internal_permission_service.openapi.yaml",
            outputDir: "docs/apis/resources/internal_permission_service_v2",
            sidebarOptions: {
              groupPathsBy: "tag",
              categoryLinkSource: "auto",
            },
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
  markdown: {
    mermaid:true,
  },
  themes: [
    "docusaurus-theme-github-codeblock",
    "docusaurus-theme-openapi-docs",
    '@docusaurus/theme-mermaid',
  ],
  future: {
    v4: false, // Disabled because of some problems related to https://github.com/facebook/docusaurus/issues/11040
    experimental_faster: {
      swcJsLoader: false, // Disabled because of memory usage > 8GB which is a problem on vercel default runners
      swcJsMinimizer: true,
      swcHtmlMinimizer: true,
      lightningCssMinimizer: true,
      mdxCrossCompilerCache: true,
      ssgWorkerThreads: false, // Disabled because of some problems related to https://github.com/facebook/docusaurus/issues/11040
      rspackBundler: true,
      rspackPersistentCache: true,
    },
  },
};
