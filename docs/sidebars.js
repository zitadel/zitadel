module.exports = {
  guides: [
    "guides/overview",
    {
      type: "category",
      label: "Get Started",
      collapsed: false,
      link: {type: "doc", id: "guides/start/quickstart"},
      items: [
        "guides/start/quickstart",
        {
          type: "category",
          label: "Frontend",
          items: [
            "examples/login/angular",
            "examples/login/flutter",
            "examples/login/go",
            "examples/login/java-spring",
            "examples/login/nextjs",
            "examples/login/python-django",
            "examples/login/react",
            "examples/login/symfony",
            "examples/login/vue",
            {
              type: "link",
              label: ".Net",
              href: "https://github.com/smartive/zitadel-net",
            },
          ],
          collapsed: true,
        },
        {
          type: "category",
          label: "Backend",
          items: [
            "examples/secure-api/go",
            "examples/secure-api/java-spring",
            "examples/secure-api/python-django",
            "examples/secure-api/python-flask",
            "examples/secure-api/nodejs-nestjs",
            "examples/secure-api/pylon",
            {
              type: "link",
              label: ".Net",
              href: "https://github.com/smartive/zitadel-net",
            },
          ],
          collapsed: true,
        },
      ],
    },
    {
      type: "category",
      label: "Examples & SDKs",
      link: {type: "doc", id: "sdk-examples/introduction"},
      items: [
        "sdk-examples/introduction",
        "sdk-examples/angular",
        "sdk-examples/flutter",
        "sdk-examples/go",
        "sdk-examples/java",
        "sdk-examples/nestjs",
        "sdk-examples/nextjs",
        "sdk-examples/python-flask",
        "sdk-examples/python-django",
        "sdk-examples/react",
        "sdk-examples/symfony",
        "sdk-examples/vue",
        {
          type: "link",
          label: "Dart",
          href: "https://github.com/smartive/zitadel-dart",
        },
        {
          type: "link",
          label: "Elixir",
          href: "https://github.com/maennchen/zitadel_api",
        },
        {
          type: "link",
          label: "NextAuth",
          href: "https://next-auth.js.org/providers/zitadel",
        },
        {
          type: "link",
          label: "Node.js",
          href: "https://www.npmjs.com/package/@zitadel/node",
        },
        {
          type: "link",
          label: ".Net",
          href: "https://github.com/smartive/zitadel-net",
        },
        {
          type: "link",
          label: "Passport.js",
          href: "https://github.com/buehler/node-passport-zitadel",
        },
        {
          type: "link",
          label: "Rust",
          href: "https://github.com/smartive/zitadel-rust",
        },
        {
          type: "link",
          label: "Pylon",
          href: "https://github.com/getcronit/pylon",
        },
      ],
    },
    {
      type: "category",
      label: "Manage",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "Cloud",
          link: {
            type: "generated-index",
            title: "Overview",
            slug: "guides/manage/cloud/overview",
            description:
              "Our customer portal is used to manage all your ZITADEL instances. You can also manage your subscriptions, billing, newsletters and support requests.",
          },
          items: [
            "guides/manage/cloud/start",
            "guides/manage/cloud/instances",
            "guides/manage/cloud/billing",
            "guides/manage/cloud/users",
            "guides/manage/cloud/support",
          ],
        },
        {
          type: "category",
          label: "Console",
          link: {
            type: "doc",
            id: "guides/manage/console/overview",
          },
          items: [
            "guides/manage/console/overview",
            "guides/manage/console/default-settings",
            "guides/manage/console/organizations",
            "guides/manage/console/projects",
            "guides/manage/console/roles",
            "guides/manage/console/applications",
            "guides/manage/console/users",
            "guides/manage/console/managers",
            "guides/manage/console/actions",
          ],
        },
        {
          type: "category",
          label: "Customize",
          items: [
            "guides/manage/customize/branding",
            "guides/manage/customize/texts",
            "guides/manage/customize/behavior",
            "guides/manage/customize/restrictions",
          ],
        },
        {
          type: "category",
          label: "Users",
          items: [
            "guides/manage/user/reg-create-user",
            "guides/manage/customize/user-metadata",
            "guides/manage/customize/user-schema",
          ],
        },
        "guides/manage/terraform-provider",
      ],
    },
    {
      type: "category",
      label: "Migrate",
      collapsed: true,
      items: [
        "guides/migrate/introduction",
        "guides/migrate/users",
        {
          type: "category",
          label: "Sources",
          collapsed: true,
          items: [
            "guides/migrate/sources/zitadel",
            "guides/migrate/sources/auth0",
            "guides/migrate/sources/keycloak",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Integrate",
      link: {
        type: "generated-index",
        title: "Integrate",
        slug: "guides/integrate",
        description:
          "Integrate your users and application with ZITADEL. In this section you will find resource on how to authenticate your users, configure external identity providers, access the ZITADEL APIs to manage resources, and integrate with third party services and tools.",
      },
      items: [
        {
          type: "category",
          label: "Login Users",
          link: {
            type: "generated-index",
            title: "Login users with ZITADEL",
            slug: "guides/integrate/login",
            description:
              "Sign-in users and application with ZITADEL. In this section you will find resources on how to authenticate your users by using the hosted login via OpenID Connect and SAML. Follow our dedicated guides to build your custom login user interface, if you want to customize the login behavior further.",
          },
          items: [
            "guides/integrate/login/login-users",
            {
              type: "category",
              label: "OpenID Connect",
              collapsed: true,
              link: {
                type: "generated-index",
                title: "Authenticate users with OpenID Connect (OIDC)",
                slug: "guides/integrate/login/oidc",
                description:
                  "This guide explains how to utilize ZITADEL for user authentication within your applications using OpenID Connect (OIDC). Here, we offer comprehensive guidance on seamlessly integrating ZITADEL's authentication features, ensuring both security and user experience excellence. Throughout this documentation, we'll cover the setup process for ZITADEL authentication, including the recommended OIDC flows tailored to different application types. Additionally, we'll provide clear instructions on securely signing out or logging out users from your application, ensuring data security and user privacy. With our guidance, you'll be equipped to leverage ZITADEL's authentication capabilities effectively, enhancing your application's security posture while delivering a seamless login experience for your users.",
              },
              items: [
                "guides/integrate/login/oidc/login-users",
                "guides/integrate/login/oidc/oauth-recommended-flows",
                "guides/integrate/login/oidc/device-authorization",
                "guides/integrate/login/oidc/logout",
              ],
            },
            "guides/integrate/login/saml",
          ],
        },
        {
          type: "category",
          label: "Onboard Customers and Users",
          link: {
            type: "generated-index",
            title: "Onboard Customers and Users",
            slug: "/guides/integrate/onboarding",
            description:
              "When building your own application, one of the first questions you have to face, is 'How do my customers onboard to my application?'\n" +
              "These guides will explain the built-in solution for onboarding new tenants, customers, and users and how you can handle more advanced onboarding use cases. ",
          },
          collapsed: true,
          items: [
            "guides/integrate/onboarding/b2b",
            "guides/integrate/onboarding/end-users",
          ],
        },
        {
          type: "category",
          label: "Token Introspection",
          collapsed: true,
          items: [
            {
              type: "autogenerated",
              dirName: "guides/integrate/token-introspection",
            },
          ],
        },
        "guides/integrate/token-exchange",
        {
          type: "category",
          label: "Service Users",
          link: {
            type: "doc",
            id: "guides/integrate/service-users/authenticate-service-users",
          },
          collapsed: true,
          items: [
            {
              type: "autogenerated",
              dirName: "guides/integrate/service-users",
            },
          ],
        },
        {
          type: "category",
          label: "Role Management",
          collapsed: true,
          items: ["guides/integrate/retrieve-user-roles"],
        },
        {
          type: "category",
          label: "Build your own Login UI",
          link: {
            type: "generated-index",
            title: "Build your own Login UI",
            slug: "/guides/integrate/login-ui",
            description:
              "In the following guides you will learn how to create your own login UI with our APIs. The different scenarios like username/password, external identity provider, etc. will be shown.",
          },
          collapsed: true,
          items: [
            "guides/integrate/login-ui/session-validation",
            "guides/integrate/login-ui/username-password",
            "guides/integrate/login-ui/external-login",
            "guides/integrate/login-ui/passkey",
            "guides/integrate/login-ui/mfa",
            "guides/integrate/login-ui/select-account",
            "guides/integrate/login-ui/password-reset",
            "guides/integrate/login-ui/logout",
            "guides/integrate/login-ui/oidc-standard",
            "guides/integrate/login-ui/typescript-repo",
          ],
        },
        {
          type: "category",
          label: "Login users with SSO",
          link: {
            type: "doc",
            id: "guides/integrate/identity-providers/introduction",
          },
          collapsed: true,
          items: [
            "guides/integrate/identity-providers/google",
            "guides/integrate/identity-providers/azure-ad-oidc",
            "guides/integrate/identity-providers/azure-ad-saml",
            "guides/integrate/identity-providers/github",
            "guides/integrate/identity-providers/gitlab",
            "guides/integrate/identity-providers/apple",
            "guides/integrate/identity-providers/ldap",
            "guides/integrate/identity-providers/openldap",
            "guides/integrate/identity-providers/okta-oidc",
            "guides/integrate/identity-providers/okta-saml",
            "guides/integrate/identity-providers/keycloak",
            "guides/integrate/identity-providers/linkedin-oauth",
            "guides/integrate/identity-providers/mocksaml",
            "guides/integrate/identity-providers/jwt_idp",
            "guides/integrate/identity-providers/migrate",
            "guides/integrate/identity-providers/additional-information",
          ],
        },
        {
          type: "category",
          label: "ZITADEL APIs",
          link: {
            type: "doc",
            id: "guides/integrate/zitadel-apis/access-zitadel-apis",
          },
          collapsed: true,
          items: [
            {
              type: "autogenerated",
              dirName: "guides/integrate/zitadel-apis",
            },
          ],
        },
        {
          type: "category",
          label: "Services",
          link: {
            type: "generated-index",
            title: "Integrate ZITADEL with your Favorite Services",
            slug: "/guides/integrate/services",
            description:
              "With the guides in this section you will learn how to integrate ZITADEL with your services.",
          },
          collapsed: true,
          items: [
            {
              type: "autogenerated",
              dirName: "guides/integrate/services",
            },
            {
              type: "link",
              label: "Bold BI (boldbi.com)",
              href: "https://support.boldbi.com/kb/article/13708/how-to-configure-zitadel-oauth-login-in-bold-bi",
            },
            {
              type: "link",
              label: "Cloudflare workers",
              href: "https://zitadel.com/blog/increase-spa-security-with-cloudflare-workers",
            },
            {
              type: "link",
              label: "Firezone (firezone.dev)",
              href: "https://www.firezone.dev/docs/authenticate/oidc/zitadel",
            },
            {
              type: "link",
              label: "Nextcloud",
              href: "https://zitadel.com/blog/zitadel-as-sso-provider-for-selfhosting",
            },
            {
              type: "link",
              label: "Netbird (netbird.io)",
              href: "https://docs.netbird.io/selfhosted/identity-providers",
            },
            {
              type: "link",
              label: "Psono (psono.com)",
              href: "https://doc.psono.com/admin/configuration/oidc-zitadel.html",
            },
            {
              type: "link",
              label: "Zoho Desk (zoho.com)",
              href: "https://help.zoho.com/portal/en/kb/desk/user-management-and-security/data-security/articles/setting-up-saml-single-signon-for-help-center#Zitadel_IDP",
            },
          ],
        },
        {
          type: "category",
          label: "Tools",
          link: {
            type: "generated-index",
            title: "Integrate ZITADEL with your Tools",
            slug: "/guides/integrate/tools",
            description:
              "With the guides in this section you will learn how to integrate ZITADEL with your favorite tools.",
          },
          collapsed: true,
          items: [
            {
              type: "link",
              label: "Argo CD",
              href: "https://argo-cd.readthedocs.io/en/latest/operator-manual/user-management/zitadel/",
            },
            "guides/integrate/tools/apache2",
            "guides/integrate/authenticated-mongodb-charts",
            "examples/identity-proxy/oauth2-proxy",
          ],
        },
        "guides/integrate/external-audit-log",
      ],
    },
    {
      type: "category",
      label: "Solution Scenarios",
      link: {
        type: "generated-index",
        title: "Solution Scenarios",
        slug: "guides/solution-scenarios/introduction",
        description:
          "Customers of an SaaS Identity and access management system usually have all distinct use cases and requirements. This guide attempts to explain real-world implementations and break them down into solution scenarios which aim to help you getting started with ZITADEL.",
      },
      collapsed: true,
      items: [
        "guides/solution-scenarios/b2c",
        "guides/solution-scenarios/b2b",
        "guides/solution-scenarios/saas",
        "guides/solution-scenarios/domain-discovery",
        "guides/solution-scenarios/configurations",
        "guides/solution-scenarios/frontend-calling-backend-API",
        "guides/solution-scenarios/restrict-console",
      ],
    },
    {
      type: "category",
      label: "Concepts",
      collapsed: true,
      link: {
        type: "generated-index",
        title: "Concepts and Features",
        slug: "concepts",
        description:
          "This part of our documentation contains ZITADEL specific or general concepts required to understand the system or our guides.",
      },
      items: [
        {
          type: "category",
          label: "Resources",
          collapsed: false,
          items: [
            {
              type: "autogenerated",
              dirName: "concepts/structure",
            },
          ],
        },
        {
          type: "category",
          label: "Features",
          collapsed: false,
          items: [
            {
              type: "autogenerated",
              dirName: "concepts/features",
            },
          ],
        },
        {
          type: "autogenerated",
          dirName: "concepts/knowledge",
        },
      ],
    },
    {
      type: "category",
      label: "Architecture",
      collapsed: true,
      items: [
        "concepts/architecture/software",
        "concepts/architecture/solution",
        "concepts/architecture/secrets",
        "concepts/principles",
        {
          type: "category",
          label: "Event Store",
          collapsed: true,
          items: [
            "concepts/eventstore/overview",
            "concepts/eventstore/implementation",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Support",
      collapsed: true,
      items: [
        "support/software-release-cycles-support",
        "support/troubleshooting",
        {
          type: "category",
          label: "Technical Advisory",
          link: {
            type: "doc",
            id: "support/technical_advisory",
          },
          collapsed: true,
          items: [
            {
              type: "autogenerated",
              dirName: "support/advisory",
            },
          ],
        },
      ],
    },
  ],
  apis: [
    "apis/introduction",
    {
      type: "category",
      label: "Core Resources",
      collapsed: false,
      link: {
        type: "doc",
        id: "apis/apis/index",
      },
      items: [
        {
          type: "category",
          label: "V1 (General Available)",
          collapsed: false,
          link: {
            type: "generated-index",
            title: "APIs V1 (GA)",
            slug: "/apis/services/",
            description:
              "APIs V1 organize access by context (authenticated user, organisation, instance, system), unlike resource-specific V2 APIs.",
          },
          items: [
            {
              type: "category",
              label: "Authenticated User",
              link: {
                type: "generated-index",
                title: "Auth API",
                slug: "/apis/resources/auth",
                description:
                  "The authentication API (aka Auth API) is used for all operations on the currently logged in user. The user id is taken from the sub claim in the token.",
              },
              items: require("./docs/apis/resources/auth/sidebar.ts"),
            },
            {
              type: "category",
              label: "Organization Objects",
              link: {
                type: "generated-index",
                title: "Management API",
                slug: "/apis/resources/mgmt",
                description:
                  "The management API is as the name states the interface where systems can mutate IAM objects like, organizations, projects, clients, users and so on if they have the necessary access rights. To identify the current organization you can send a header x-zitadel-orgid or if no header is set, the organization of the authenticated user is set.",
              },
              items: require("./docs/apis/resources/mgmt/sidebar.ts"),
            },
            {
              type: "category",
              label: "Instance Objects",
              link: {
                type: "generated-index",
                title: "Admin API",
                slug: "/apis/resources/admin",
                description:
                  "This API is intended to configure and manage one ZITADEL instance itself.",
              },
              items: require("./docs/apis/resources/admin/sidebar.ts"),
            },
            {
              type: "category",
              label: "Instance Lifecycle",
              link: {
                type: "generated-index",
                title: "System API",
                slug: "/apis/resources/system",
                description:
                  "This API is intended to manage the different ZITADEL instances within the system.\n" +
                  "\n" +
                  "Checkout the guide how to access the ZITADEL System API.",
              },
              items: require("./docs/apis/resources/system/sidebar.ts"),
            },
          ],
        },
        {
          type: "category",
          label: "V2 (Beta)",
          collapsed: false,
          link: {
            type: "doc",
            id: "apis/v2beta",
          },
          items: [
            {
              type: "category",
              label: "User Lifecycle (Beta)",
              link: {
                type: "generated-index",
                title: "User Service API (Beta)",
                slug: "/apis/resources/user_service_v2beta",
                description:
                  "This API is intended to manage users in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in beta state. It can AND will continue breaking until the services provide the same functionality as the current login.",
              },
              items: require("./docs/apis/resources/user_service_v2beta/sidebar.ts"),
            },
            {
              type: "category",
              label: "Session Lifecycle (Beta)",
              link: {
                type: "generated-index",
                title: "Session Service API (Beta)",
                slug: "/apis/resources/session_service_v2beta",
                description:
                  "This API is intended to manage sessions in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in beta state. It can AND will continue breaking until the services provide the same functionality as the current login.",
              },
              items: require("./docs/apis/resources/session_service_v2beta/sidebar.ts"),
            },
            {
              type: "category",
              label: "OIDC Lifecycle (Beta)",
              link: {
                type: "generated-index",
                title: "OIDC Service API (Beta)",
                slug: "/apis/resources/oidc_service_v2beta",
                description:
                  "Get OIDC Auth Request details and create callback URLs.\n" +
                  "\n" +
                  "This project is in beta state. It can AND will continue breaking until the services provide the same functionality as the current login.",
              },
              items: require("./docs/apis/resources/oidc_service_v2beta/sidebar.ts"),
            },
            {
              type: "category",
              label: "Settings Lifecycle (Beta)",
              link: {
                type: "generated-index",
                title: "Settings Service API (Beta)",
                slug: "/apis/resources/settings_service_v2beta",
                description:
                  "This API is intended to manage settings in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in beta state. It can AND will continue to break until the services provide the same functionality as the current login.",
              },
              items: require("./docs/apis/resources/settings_service_v2beta/sidebar.ts"),
            },
            {
              type: "category",
              label: "Feature Lifecycle (Beta)",
              link: {
                type: "generated-index",
                title: "Feature Service API (Beta)",
                slug: "/apis/resources/feature_service_v2beta",
                description:
                  'This API is intended to manage features for ZITADEL. Feature settings that are available on multiple "levels", such as instance and organization. The higher level instance acts as a default for the lower level. When a feature is set on multiple levels, the lower level takes precedence. Features can be experimental where ZITADEL will assume a sane default, such as disabled. When over time confidence in such a feature grows, ZITADEL can default to enabling the feature. As a final step we might choose to always enable a feature and remove the setting from this API, reserving the proto field number. Such removal is not considered a breaking change. Setting a removed field will effectively result in a no-op.\n' +
                  "\n" +
                  "This project is in beta state. It can AND will continue breaking until a stable version is released.",
              },
              items: require("./docs/apis/resources/feature_service_v2beta/sidebar.ts"),
            },
          ],
        },
        {
          type: "category",
          label: "V2 (General Available)",
          collapsed: false,
          link: {
            type: "doc",
            id: "apis/v2",
          },
          items: [
            {
              type: "category",
              label: "User Lifecycle",
              link: {
                type: "generated-index",
                title: "User Service API",
                slug: "/apis/resources/user_service_v2",
                description:
                  "This API is intended to manage users in a ZITADEL instance.\n"
              },
              items: require("./docs/apis/resources/user_service_v2/sidebar.ts"),
            },
            {
              type: "category",
              label: "Session Lifecycle",
              link: {
                type: "generated-index",
                title: "Session Service API",
                slug: "/apis/resources/session_service_v2",
                description:
                  "This API is intended to manage sessions in a ZITADEL instance.\n"
              },
              items: require("./docs/apis/resources/session_service_v2/sidebar.ts"),
            },
            {
              type: "category",
              label: "OIDC Lifecycle",
              link: {
                type: "generated-index",
                title: "OIDC Service API",
                slug: "/apis/resources/oidc_service_v2",
                description:
                  "Get OIDC Auth Request details and create callback URLs.\n"
              },
              items: require("./docs/apis/resources/oidc_service_v2/sidebar.ts"),
            },
            {
              type: "category",
              label: "Settings Lifecycle",
              link: {
                type: "generated-index",
                title: "Settings Service API",
                slug: "/apis/resources/settings_service_v2",
                description:
                  "This API is intended to manage settings in a ZITADEL instance.\n"
              },
              items: require("./docs/apis/resources/settings_service_v2/sidebar.ts"),
            },
            {
              type: "category",
              label: "Feature Lifecycle",
              link: {
                type: "generated-index",
                title: "Feature Service API",
                slug: "/apis/resources/feature_service/v2",
                description:
                  'This API is intended to manage features for ZITADEL. Feature settings that are available on multiple "levels", such as instance and organization. The higher level instance acts as a default for the lower level. When a feature is set on multiple levels, the lower level takes precedence. Features can be experimental where ZITADEL will assume a sane default, such as disabled. When over time confidence in such a feature grows, ZITADEL can default to enabling the feature. As a final step we might choose to always enable a feature and remove the setting from this API, reserving the proto field number. Such removal is not considered a breaking change. Setting a removed field will effectively result in a no-op.\n'
              },
              items: require("./docs/apis/resources/feature_service_v2/sidebar.ts"),
            },
          ],
        },
        {
          type: "category",
          label: "V3 (Preview)",
          collapsed: false,
          link: {
            type: "doc",
            id: "apis/v3",
          },
          items: [
            {
              type: "category",
              label: "User Schema Lifecycle (Preview)",
              link: {
                type: "generated-index",
                title: "User Schema Service API (Preview)",
                slug: "/apis/resources/user_schema_service",
                description:
                  "This API is intended to manage data schemas for users in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in Preview state. It can AND will continue breaking until the service provides the same functionality as the v1 and v2 user services.",
              },
              items: require("./docs/apis/resources/user_schema_service_v3/sidebar.ts"),
            },
            {
              type: "category",
              label: "User Lifecycle (Preview)",
              link: {
                type: "generated-index",
                title: "User Service API (Preview)",
                slug: "/apis/resources/user_service_v3",
                description:
                  "This API is intended to manage users with your own data schema in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in Preview state. It can AND will continue breaking until the service provides the same functionality as the v1 and v2 user services.",
              },
              items: require("./docs/apis/resources/user_service_v3/sidebar.ts"),
            },
            {
              type: "category",
              label: "Action Lifecycle (Preview)",
              link: {
                type: "generated-index",
                title: "Action Service API (Preview)",
                slug: "/apis/resources/action_service_v3",
                description:
                  "This API is intended to manage custom executions and targets (previously known as actions) in a ZITADEL instance.\n" +
                  "\n" +
                  "This project is in Preview state. It can AND will continue breaking until the services provide the same functionality as the current actions.",
              },
              items: require("./docs/apis/resources/action_service_v3/sidebar.ts"),
            },
          ],
        },
        {
          type: "category",
          label: "Assets",
          collapsed: true,
          items: ["apis/assets/assets"],
        },
      ],
    },
    {
      type: "category",
      label: "Sign In Users ",
      collapsed: false,
      items: [
        {
          type: "category",
          label: "OpenID Connect & OAuth",
          collapsed: true,
          items: [
            "apis/openidoauth/endpoints",
            "apis/openidoauth/authrequest",
            "apis/openidoauth/scopes",
            "apis/openidoauth/claims",
            "apis/openidoauth/authn-methods",
            "apis/openidoauth/grant-types",
          ],
        },
        {
          type: "category",
          label: "SAML 2.0",
          collapsed: true,
          items: ["apis/saml/endpoints"],
        },
      ],
    },
    {
      type: "category",
      label: "Actions",
      collapsed: false,
      items: [
        "apis/actions/introduction",
        "apis/actions/modules",
        "apis/actions/code-examples",
        "apis/actions/internal-authentication",
        "apis/actions/external-authentication",
        "apis/actions/complement-token",
        "apis/actions/customize-samlresponse",
        "apis/actions/objects",
      ],
    },
    {
      type: "category",
      label: "Actions V2",
      collapsed: false,
      items: ["apis/actionsv2/introduction", "apis/actionsv2/execution-local"],
    },
    {
      type: "doc",
      label: "gRPC Status Codes",
      id: "apis/statuscodes",
    },
    {
      type: "category",
      label: "Observability",
      collapsed: false,
      items: ["apis/observability/metrics", "apis/observability/health"],
    },
    {
      type: "link",
      label: "Rate Limits (Cloud)", // The link label
      href: "/legal/policies/rate-limit-policy", // The internal path
    },
  ],
  selfHosting: [
    {
      type: "category",
      label: "Deploy",
      collapsed: false,
      items: [
        "self-hosting/deploy/overview",
        "self-hosting/deploy/linux",
        "self-hosting/deploy/macos",
        "self-hosting/deploy/compose",
        "self-hosting/deploy/devcontainer",
        "self-hosting/deploy/knative",
        "self-hosting/deploy/kubernetes",
        "self-hosting/deploy/loadbalancing-example/loadbalancing-example",
        "self-hosting/deploy/troubleshooting/troubleshooting",
      ],
    },
    {
      type: "category",
      label: "Manage",
      collapsed: false,
      items: [
        "self-hosting/manage/production",
        "self-hosting/manage/productionchecklist",
        "self-hosting/manage/configure/configure",
        {
          type: "category",
          collapsed: false,
          label: "Reverse Proxy",
          link: {
            type: "doc",
            id: "self-hosting/manage/reverseproxy/reverse_proxy",
          },
          items: [
            "self-hosting/manage/reverseproxy/traefik/traefik",
            "self-hosting/manage/reverseproxy/nginx/nginx",
            "self-hosting/manage/reverseproxy/caddy/caddy",
            //            "self-hosting/manage/reverseproxy/httpd/httpd", grpc NOT WORKING
            "self-hosting/manage/reverseproxy/cloudflare/cloudflare",
            "self-hosting/manage/reverseproxy/cloudflare_tunnel/cloudflare_tunnel",
            "self-hosting/manage/reverseproxy/zitadel_cloud/zitadel_cloud",
          ],
        },
        "self-hosting/manage/custom-domain",
        "self-hosting/manage/http2",
        "self-hosting/manage/tls_modes",
        "self-hosting/manage/database/database",
        "self-hosting/manage/updating_scaling",
        "self-hosting/manage/usage_control",
        {
          type: "category",
          label: "Command Line Interface",
          collapsed: false,
          link: {
            type: "doc",
            id: "self-hosting/manage/cli/overview",
          },
          items: ["self-hosting/manage/cli/mirror"],
        },
      ],
    },
  ],
  legal: [
    {
      type: "category",
      label: "Legal Agreements",
      collapsed: false,
      link: {
        type: "generated-index",
        title: "Legal Agreements",
        slug: "legal",
        description:
          "This section contains important agreements, policies and appendices relevant for users of our websites and services. All documents will be provided in English language.",
      },
      items: [
        "legal/terms-of-service",
        "legal/data-processing-agreement",
        "legal/subprocessors",
        "legal/annex-support-services",
        {
          type: "category",
          label: "Service Description",
          collapsed: false,
          link: {
            type: "generated-index",
            title: "Service description",
            slug: "/legal/service-description",
            description:
              "Description of services and service levels for ZITADEL Cloud and Enterprise subscriptions.",
          },
          items: [
            {
              type: "autogenerated",
              dirName: "legal/service-description",
            },
          ],
        },
        {
          type: "category",
          label: "Policies",
          collapsed: false,
          link: {
            type: "generated-index",
            title: "Policies",
            slug: "/legal/policies",
            description:
              "Policies and guidelines in addition to our terms of services.",
          },
          items: [
            {
              type: "autogenerated",
              dirName: "legal/policies",
            },
          ],
        },
      ],
    },
  ],
};
