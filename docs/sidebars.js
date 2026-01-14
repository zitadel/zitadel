
const sidebar_api_auth = require("./docs/apis/resources/auth/sidebar.ts").default
const sidebar_api_mgmt = require("./docs/apis/resources/mgmt/sidebar.ts").default
const sidebar_api_admin = require("./docs/apis/resources/admin/sidebar.ts").default
const sidebar_api_system = require("./docs/apis/resources/system/sidebar.ts").default

const sidebar_api_user_service_v2 = require("./docs/apis/resources/user_service_v2/sidebar.ts").default
const sidebar_api_session_service_v2 = require("./docs/apis/resources/session_service_v2/sidebar.ts").default
const sidebar_api_oidc_service_v2 = require("./docs/apis/resources/oidc_service_v2/sidebar.ts").default
const sidebar_api_saml_service_v2 = require("./docs/apis/resources/saml_service_v2/sidebar.ts").default
const sidebar_api_settings_service_v2 = require("./docs/apis/resources/settings_service_v2/sidebar.ts").default
const sidebar_api_feature_service_v2 = require("./docs/apis/resources/feature_service_v2/sidebar.ts").default
const sidebar_api_org_service_v2 = require("./docs/apis/resources/org_service_v2/sidebar.ts").default
const sidebar_api_idp_service_v2 = require("./docs/apis/resources/idp_service_v2/sidebar.ts").default
const sidebar_api_actions_v2 = require("./docs/apis/resources/action_service_v2/sidebar.ts").default
const sidebar_api_project_service_v2 = require("./docs/apis/resources/project_service_v2/sidebar.ts").default
const sidebar_api_webkey_service_v2 = require("./docs/apis/resources/webkey_service_v2/sidebar.ts").default
const sidebar_api_instance_service_v2 = require("./docs/apis/resources/instance_service_v2/sidebar.ts").default
const sidebar_api_authorization_service_v2 = require("./docs/apis/resources/authorization_service_v2/sidebar.ts").default
const sidebar_api_internal_permission_service_v2 = require("./docs/apis/resources/internal_permission_service_v2/sidebar.ts").default
const sidebar_api_application_v2 = require("./docs/apis/resources/application_service_v2/sidebar.ts").default

module.exports = {
  guides: [
    {
      type: "category",
      label: "Get Started",
      collapsed: false,
      items: [
        "guides/overview",
        "guides/start/quickstart",
        {
          type: "category",
          label: "Key Concepts",
          link: {
            type: "generated-index",
            title: "Key Concepts",
            slug: "concepts",
            description:
              "This part of our documentation contains ZITADEL specific or general concepts required to understand the system or our guides.",
          },
          items: [
            "concepts/structure/instance",
            "guides/manage/console/organizations-overview",
            "guides/manage/console/projects-overview",
            "guides/manage/console/applications-overview",
            "guides/manage/console/users-overview",
            "concepts/structure/managers",
          ],
        },
        {
          type: "category",
          label: "Authenticate Users",
          link: {
            type: "generated-index",
            title: "Login users with ZITADEL",
            slug: "guides/integrate/login",
            description:
              "Sign-in users and application with ZITADEL. In this section you will find resources on how to authenticate your users by using the hosted login via OpenID Connect and SAML. Follow our dedicated guides to build your custom login user interface, if you want to customize the login behavior further.",
          },
          items: [
            "guides/integrate/login/login-users",
            "guides/integrate/login/hosted-login",
            "guides/integrate/login/oidc/logout",
          ],
        },
        {
          type: "category",
          label: "Example Applications",
          link: {
            type: "generated-index",
            title: "Example Applications",
            slug: "/examples/introduction",
            description:
              "Practical examples showing how to integrate ZITADEL authentication and secure APIs across different application types and frameworks.",
          },
          items: [
            {
              type: "category",
              label: "Frontend (SPA)",
              link: {
                type: "generated-index",
                title: "Frontend (SPA) Quickstart Guides",
                slug: "/examples/introduction/frontend",
                description:
                  "Quickstart guides for integrating ZITADEL authentication into frontend single-page applications.",
              },
              items: [
                {
                  type: "link",
                  label: "Vanilla-JS",
                  href: "https://github.com/zitadel/zitadel-vanilla-js",
                },
                "examples/login/react",
                "examples/login/angular",
                "examples/login/vue",
              ],
            },
            {
              type: "category",
              label: "Mobile & Native",
              link: {
                type: "generated-index",
                title: "Mobile & Native Quickstart Guides",
                slug: "/examples/introduction/mobile",
                description:
                  "Quickstart guides for integrating ZITADEL authentication into Mobile & Native applications.",
              },
              items: [
                "examples/login/flutter",
              ],
            },
            {
              type: "category",
              label: "Full-Stack / SSR",
              link: {
                type: "generated-index",
                title: "Full-Stack / SSR Quickstart Guides",
                slug: "/examples/introduction/fullstack",
                description:
                  "Quickstart guides for integrating ZITADEL authentication into Full-Stack / SSR applications.",
              },
              items: [
                "examples/login/nextjs",
                "examples/login/nextjs-b2b",
              ],
            },
            {
              type: "category",
              label: "Web App (Server-Side)",
              link: {
                type: "generated-index",
                title: "Web App (Server-Side) Quickstart Guides",
                slug: "/examples/introduction/webapp",
                description:
                  "Quickstart guides for integrating ZITADEL authentication into Web App (Server-Side) applications.",
              },
              items: [
                "examples/login/symfony",
                "examples/login/java-spring",
                "examples/login/python-django",
                "examples/login/go"
              ],
            },
            {
              type: "category",
              label: "APIs / Backend Services",
              link: {
                type: "generated-index",
                title: "APIs / Backend Quickstart Guides",
                slug: "/examples/introduction/backend",
                description:
                  "Quickstart guides for integrating ZITADEL authentication into APIs / Backend Services.",
              },
              items: [
                "examples/secure-api/go",
                "examples/secure-api/java-spring",
                "examples/secure-api/python-django",
                "examples/secure-api/python-flask",
                "examples/secure-api/nodejs-nestjs",
                "examples/secure-api/pylon",
              ],
            },
          ],
        },
        {
          type: "category",
          label: "Use Cases",
          link: {
            type: "generated-index",
            title: "Use Cases",
            slug: "guides/solution-scenarios/introduction",
            description:
              "Customers of an SaaS Identity and access management system usually have all distinct use cases and requirements. This guide attempts to explain real-world implementations and break them down into solution scenarios which aim to help you getting started with ZITADEL.",
          },
          items: [
            "guides/solution-scenarios/b2b",
            "guides/solution-scenarios/configurations",
            "guides/solution-scenarios/saas",
            "guides/solution-scenarios/b2c",
            "guides/solution-scenarios/frontend-calling-backend-API",
            {
              type: "category",
              label: "Machine-to-Machine (M2M)",
              items: [
                "guides/integrate/service-users/authenticate-service-users",
                "guides/integrate/service-users/private-key-jwt",
                "guides/integrate/service-users/client-credentials",
                "guides/integrate/service-users/personal-access-token",
              ],
            },
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
          label: "Branding & Customization",
          items: [
            "guides/manage/customize/branding",
            "guides/manage/customize/texts",
            "guides/manage/customize/restrictions",
            "guides/manage/customize/user-schema",
            "guides/manage/customize/user-metadata",
            "guides/manage/customize/notification-providers",
            "concepts/features/custom-domain",
          ]
        },
      ],
    },
    {
      type: "category",
      label: "Integrate & Authenticate",
      link: {
        type: "generated-index",
        title: "Integrate",
        slug: "guides/integrate",
        description:
          "Integrate your users and application with ZITADEL. In this section you will find resource on how to authenticate your users, configure external identity providers, access the ZITADEL APIs to manage resources, and integrate with third party services and tools.",
      },
      collapsed: true,
      items: [
        {
          type: "category",
          label: "OIDC & OAuth Flows",
          link: {
            type: "generated-index",
            title: "Authenticate users with OpenID Connect (OIDC)",
            slug: "guides/integrate/login/oidc",
            description:
              "This guide explains how to utilize ZITADEL for user authentication within your applications using OpenID Connect (OIDC). Here, we offer comprehensive guidance on seamlessly integrating ZITADEL's authentication features, ensuring both security and user experience excellence. Throughout this documentation, we'll cover the setup process for ZITADEL authentication, including the recommended OIDC flows tailored to different application types. Additionally, we'll provide clear instructions on securely signing out or logging out users from your application, ensuring data security and user privacy. With our guidance, you'll be equipped to leverage ZITADEL's authentication capabilities effectively, enhancing your application's security posture while delivering a seamless login experience for your users.",
          },
          items: [
            "guides/integrate/login/oidc/oauth-recommended-flows",
            "guides/integrate/login/oidc/login-users",
            "guides/integrate/login/oidc/device-authorization",
            "apis/openidoauth/endpoints",
            "apis/openidoauth/authn-methods",
            "apis/openidoauth/scopes",
            "apis/openidoauth/claims",
            "apis/openidoauth/grant-types",
            "apis/openidoauth/authrequest",
            "guides/integrate/login/oidc/webkeys",
            "guides/integrate/token-exchange",
          ],
        },
        {
          type: "category",
          label: "SAML",
          items: [
            "guides/integrate/login/saml",
            "apis/saml/endpoints",
          ],
        },
        {
          type: "category",
          label: "API Access",
          items: [
            "guides/integrate/zitadel-apis/access-zitadel-apis",
            "guides/integrate/zitadel-apis/access-zitadel-system-api",
            "guides/integrate/zitadel-apis/event-api",
            "guides/integrate/zitadel-apis/example-zitadel-api-with-go",
            "guides/integrate/zitadel-apis/example-zitadel-api-with-dot-net",
            {
              type: "link",
              label: "Zitadel APIs",
              href: "/docs/apis/introduction",
            },
          ],
        },
        {
          type: "category",
          label: "SDKs",
          items: [
            "sdk-examples/introduction",
            {
              type: "category",
              label: "Frontend (SPA)",
              items: [
                "sdk-examples/react",
                "sdk-examples/angular",
                "sdk-examples/vue",
              ],
            },
            {
              type: "category",
              label: "Mobile & Native",
              items: [
                {
                  type: "link",
                  label: "Dart / Flutter",
                  href: "https://github.com/smartive/zitadel-dart",
                },
                {
                  type: "link",
                  label: ".NET (MAUI/Xamarin)",
                  href: "https://github.com/smartive/zitadel-net",
                },
              ],
            },
            {
              type: "category",
              label: "Full-Stack / SSR",
              items: [
                "sdk-examples/nextjs",
                "sdk-examples/nuxtjs",
                "sdk-examples/sveltekit",
                "sdk-examples/qwik",
                "sdk-examples/solidstart",
                "sdk-examples/astro",
                {
                  type: "link",
                  label: "NextAuth",
                  href: "https://next-auth.js.org/providers/zitadel",
                },
              ],
            },
            {
              type: "category",
              label: "Backend & API",
              items: [
                {
                  type: "category",
                  label: "Node.js",
                  items: [
                    "sdk-examples/client-libraries/node",
                    "sdk-examples/expressjs",
                    "sdk-examples/fastify",
                    "sdk-examples/hono",
                    "sdk-examples/nestjs",
                    {
                      type: "link",
                      label: "Passport.js",
                      href: "https://github.com/buehler/node-passport-zitadel",
                    },
                    {
                      type: "link",
                      label: "Node.js (Community)",
                      href: "https://www.npmjs.com/package/@zitadel/node",
                    },
                  ],
                },
                {
                  type: "category",
                  label: "Python",
                  items: [
                    "sdk-examples/client-libraries/python",
                    "sdk-examples/flask",
                    "sdk-examples/django",
                    "sdk-examples/fastapi",
                  ],
                },
                {
                  type: "category",
                  label: "Go",
                  items: ["sdk-examples/go"],
                },
                {
                  type: "category",
                  label: "Java",
                  items: [
                    "sdk-examples/java",
                    "sdk-examples/spring",
                    "sdk-examples/client-libraries/java",
                  ],
                },
                {
                  type: "category",
                  label: "PHP",
                  items: [
                    "sdk-examples/symfony",
                    "sdk-examples/laravel",
                    "sdk-examples/client-libraries/php",
                  ],
                },
                {
                  type: "category",
                  label: "Other Languages",
                  items: [
                    "sdk-examples/client-libraries/ruby",
                    {
                      type: "link",
                      label: "Elixir",
                      href: "https://github.com/maennchen/zitadel_api",
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
              ],
            },
          ],
        },
        {
          type: "category",
          label: "SCIM",
          link: {
            type: "doc",
            id: "guides/manage/user/scim2",
          },
          items: ["guides/integrate/scim-okta-guide"]
        },
        {
          type: "category",
          label: "Token Introspection",
          link: {
            type: "doc",
            id: "guides/integrate/token-introspection/index",
          },
          items: [
            "guides/integrate/token-introspection/basic-auth",
            "guides/integrate/token-introspection/private-key-jwt",
          ]
        },
        "guides/integrate/back-channel-logout",
        {
          type: "category",
          label: "External Integrations",
          items: [
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
              items: [
                "guides/integrate/services/google-workspace",
                "guides/integrate/services/aws-saml",
                "guides/integrate/services/atlassian-saml",
                "guides/integrate/services/pingidentity-saml",
                "guides/integrate/services/google-cloud",
                "guides/integrate/services/cloudflare-oidc",
                "guides/integrate/services/gitlab-saml",
                "guides/integrate/services/auth0-saml",
                "guides/integrate/services/auth0-oidc",
                "guides/integrate/services/gitlab-self-hosted",
                {
                  type: "link",
                  label: "Bold BI",
                  href: "https://support.boldbi.com/kb/article/13708/how-to-configure-zitadel-oauth-login-in-bold-bi",
                },
                {
                  type: "link",
                  label: "Cloudflare Workers",
                  href: "https://zitadel.com/blog/increase-spa-security-with-cloudflare-workers",
                },
                {
                  type: "link",
                  label: "Firezone",
                  href: "https://www.firezone.dev/docs/authenticate/oidc/zitadel",
                },
                {
                  type: "link",
                  label: "Netbird",
                  href: "https://docs.netbird.io/selfhosted/identity-providers",
                },
                {
                  type: "link",
                  label: "Nextcloud",
                  href: "https://zitadel.com/blog/zitadel-as-sso-provider-for-selfhosting",
                },
                {
                  type: "link",
                  label: "Psono",
                  href: "https://doc.psono.com/admin/configuration/oidc-zitadel.html",
                },
                {
                  type: "link",
                  label: "Zoho Desk",
                  href: "https://help.zoho.com/portal/en/kb/desk/user-management-and-security/data-security/articles/setting-up-saml-single-signon-for-help-center#Zitadel_IDP",
                },
              ],
            },
            {
              type: "category",
              label: "Tools",
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
          ],
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
            "guides/integrate/login-ui/saml-standard",
            "guides/integrate/login-ui/device-auth",
            "guides/integrate/login-ui/login-app",
          ],
        },
        {
          type: "category",
          label: "Actions",
          items: [
            {
              type: "category",
              label: "V1",
              link: {
                type: "doc",
                id: "concepts/features/actions",
              },
              items: [
                "guides/manage/console/actions-overview",
                "apis/actions/modules",
                "apis/actions/code-examples",
                "apis/actions/internal-authentication",
                "apis/actions/external-authentication",
                "apis/actions/complement-token",
                "apis/actions/customize-samlresponse",
                "apis/actions/objects",
                "guides/manage/customize/behavior",
              ]
            },
            {
              type: "category",
              label: "V2",
              link: {
                type: "doc",
                id: "concepts/features/actions_v2",
              },
              items: [
                "guides/integrate/actions/usage",
                "guides/integrate/actions/testing-request",
                "guides/integrate/actions/testing-request-manipulation",
                "guides/integrate/actions/testing-response",
                "guides/integrate/actions/testing-response-manipulation",
                "guides/integrate/actions/testing-function",
                "guides/integrate/actions/testing-function-manipulation",
                "guides/integrate/actions/testing-event",
                "guides/integrate/actions/testing-request-signature",
                "guides/integrate/actions/migrate-from-v1",
                "guides/integrate/actions/webhook-site-setup",
              ],
            },
          ],
        },
        {
          type: "category",
          label: "Migrate",
          collapsed: true,
          items: [
            "guides/migrate/introduction",
            "guides/migrate/users",
            "guides/migrate/sources/zitadel",
            "guides/migrate/sources/auth0",
            "guides/migrate/sources/keycloak",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Configure Identity & Policies",
      collapsed: true,
      items: [
        "guides/manage/user/reg-create-user",
        "guides/manage/terraform-provider",
        {
          type: "category",
          label: "Identity Providers",
          link: {
            type: "doc",
            id: "guides/integrate/identity-providers/introduction",
          },
          items: [
            {
              type: "category",
              label: "Social Logins",
              collapsed: true,
              items: [
                "guides/integrate/identity-providers/google",
                "guides/integrate/identity-providers/apple",
                "guides/integrate/identity-providers/github",
                "guides/integrate/identity-providers/gitlab",
                "guides/integrate/identity-providers/linkedin-oauth",
              ],
            },
            {
              type: "category",
              label: "Enterprise (SAML & OIDC)",
              collapsed: true,
              items: [
                "guides/integrate/identity-providers/azure-ad-oidc",
                "guides/integrate/identity-providers/azure-ad-saml",
                "guides/integrate/identity-providers/okta-oidc",
                "guides/integrate/identity-providers/okta-saml",
                "guides/integrate/identity-providers/keycloak",
                "guides/integrate/identity-providers/onelogin-saml",
                "guides/integrate/identity-providers/pingfederate-saml",
              ],
            },
            {
              type: "category",
              label: "Legacy & Directory (LDAP)",
              items: [
                "guides/integrate/identity-providers/ldap",
                "guides/integrate/identity-providers/openldap",
              ],
            },
            {
              type: "category",
              label: "Custom & Generic",
              items: [
                "guides/integrate/identity-providers/generic-oidc",
                "guides/integrate/identity-providers/jwt_idp",
                "guides/integrate/identity-providers/mocksaml",
              ],
            },
            {
              type: "category",
              label: "Guides",
              items: [
                "guides/integrate/identity-providers/migrate",
                "guides/integrate/identity-providers/additional-information",
              ]
            }
          ],
        },
        {
          type: "category",
          label: "Policies",
          items: [
            "guides/manage/console/default-settings",
            "concepts/structure/policies",
          ],
        },
        {
          type: "category",
          label: "Roles & Permissions",
          items: [
            "guides/manage/console/roles",
            "guides/integrate/retrieve-user-roles",
            "concepts/structure/granted_projects",
            "guides/manage/console/managers",
          ],
        },
        {
          type: "category",
          label: "Compliance & Security",
          items: [
            "concepts/features/audit-trail",
            "guides/integrate/external-audit-log",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Test & Debug",
      collapsed: true,
      items: [
        {
          type: "link",
          label: "OIDC Playground",
          href: "https://zitadel.com/playgrounds/oidc",
        },
      ],
    },
    {
      type: "category",
      label: "Deploy & Operate",
      collapsed: true,
      items: [
        "guides/manage/console/console-overview",
        {
          type: "category",
          label: "Customer Portal",
          collapsed: true,
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
            "guides/manage/cloud/settings",
            "guides/manage/cloud/billing",
            "guides/manage/cloud/users",
          ],
        },
        {
          type: "category",
          label: "Self-Hosted",
          items: [
            "self-hosting/deploy/overview",
            "self-hosting/deploy/linux",
            "self-hosting/deploy/macos",
            "self-hosting/deploy/devcontainer",
            "self-hosting/deploy/compose",
            "self-hosting/deploy/kubernetes",
            {
              type: "category",
              label: "Manage",
              collapsed: true,
              items: [
                {
                  type: "category",
                  label: "Production & Operations",
                  collapsed: false,
                  items: [
                    "self-hosting/manage/production",
                    "self-hosting/manage/productionchecklist",
                    "self-hosting/manage/usage_control",
                  ],
                },
                {
                  type: "category",
                  label: "Configuration",
                  collapsed: true,
                  items: [
                    "self-hosting/manage/configure/configure",
                    "self-hosting/manage/custom-domain",
                    "self-hosting/manage/tls_modes",
                    "self-hosting/manage/http2",
                    "self-hosting/manage/login-client",
                  ],
                },
                {
                  type: "category",
                  collapsed: true,
                  label: "Reverse Proxy",
                  link: {
                    type: "doc",
                    id: "self-hosting/manage/reverseproxy/reverse_proxy",
                  },
                  items: [
                    "self-hosting/manage/reverseproxy/traefik/traefik",
                    "self-hosting/manage/reverseproxy/nginx/nginx",
                    "self-hosting/manage/reverseproxy/caddy/caddy",
                    "self-hosting/manage/reverseproxy/httpd/httpd",
                    "self-hosting/manage/reverseproxy/cloudflare/cloudflare",
                    "self-hosting/manage/reverseproxy/cloudflare_tunnel/cloudflare_tunnel",
                    "self-hosting/manage/reverseproxy/zitadel_cloud/zitadel_cloud",
                  ],
                },
                {
                  type: "category",
                  label: "Observability",
                  collapsed: true,
                  items: [
                    "self-hosting/manage/service_ping",
                    {
                      type: "category",
                      label: "Metrics",
                      collapsed: true,
                      link: {
                        type: "doc",
                        id: "self-hosting/manage/metrics/overview",
                      },
                      items: ["self-hosting/manage/metrics/prometheus"],
                    },
                  ],
                },
                {
                  type: "category",
                  label: "Tools",
                  collapsed: true,
                  items: [
                    {
                      type: "category",
                      label: "Command Line Interface",
                      collapsed: true,
                      link: {
                        type: "doc",
                        id: "self-hosting/manage/cli/overview",
                      },
                      items: ["self-hosting/manage/cli/mirror"],
                    },
                  ],
                },
              ],
            },
            "self-hosting/deploy/troubleshooting/troubleshooting",
          ],
        },
        {
          type: "category",
          label: "Scaling & Performance",
          items: [
            "self-hosting/manage/updating_scaling",
            "self-hosting/manage/database/database",
            "self-hosting/manage/cache",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Architecture & Concepts",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "System Architecture",
          items: [
            "concepts/architecture/solution",
            "concepts/architecture/software",
            "concepts/architecture/secrets",
          ],
        },
        "concepts/principles",
        {
          type: "category",
          label: "Advanced Topics",
          items: [
            {
              type: "category",
              label: "Event Store",
              items: [
                "concepts/eventstore/overview",
                "concepts/eventstore/implementation",
              ]
            },
            "concepts/features/selfservice",
            "concepts/features/account-linking",
            "concepts/features/external-user-grant",
            "concepts/features/identity-brokering",
            "concepts/features/passkeys",
            "concepts/knowledge/opaque-tokens",
            "guides/solution-scenarios/domain-discovery",
            "guides/solution-scenarios/restrict-console",
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Product, Releases & Support",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "Product Features",
          items: [
            "product/roadmap",
          ],
        },
        {
          type: "category",
          label: "Releases",
          items: [
            {
              type: "link",
              label: "Changelog",
              href: "https://zitadel.com/changelog",
            },
            "product/release-cycle",
          ],
        },
        {
          type: "category",
          label: "Support Resources",
          items: [
            "support/troubleshooting",
            "support/technical_advisory",
            {
              type: "link",
              label: "Support States",
              href: "https://help.zitadel.com/zitadel-support-states",
            },
            {
              type: "link",
              label: "Zitadel Release Cycle",
              href: "https://help.zitadel.com/zitadel-software-release-cycle",
            },
            {
              type: "link",
              label: "Knowledge Base",
              href: "https://help.zitadel.com",
            }
          ],
        },
        {
          type: "link",
          label: "Pricing",
          href: "https://zitadel.com/pricing",
        },
        {
          type: "link",
          label: "Contact Us",
          href: "https://zitadel.com/contact",
        },
        "guides/manage/cloud/support",
      ],
    },
  ],
  apis: [
    "apis/introduction",
    {
      type: "category",
      label: "Core Resources",
      collapsed: true,
      link: {
        type: "doc",
        id: "apis/apis/index",
      },
      items: [
        {
          type: "category",
          label: "V2",
          collapsed: true,
          link: {
            type: "doc",
            id: "apis/v2",
          },
          items: [
            {
              type: "category",
              label: "User",
              link: {
                type: "generated-index",
                title: "User Service API",
                slug: "/apis/resources/user_service_v2",
                description:
                  "This API is intended to manage users in a ZITADEL instance.\n",
              },
              items: sidebar_api_user_service_v2,
            },
            {
              type: "category",
              label: "Session",
              link: {
                type: "generated-index",
                title: "Session Service API",
                slug: "/apis/resources/session_service_v2",
                description:
                  "This API is intended to manage sessions in a ZITADEL instance.\n",
              },
              items: sidebar_api_session_service_v2,
            },
            {
              type: "category",
              label: "OIDC",
              link: {
                type: "generated-index",
                title: "OIDC Service API",
                slug: "/apis/resources/oidc_service_v2",
                description:
                  "Get OIDC Auth Request details and create callback URLs.\n",
              },
              items: sidebar_api_oidc_service_v2,
            },
            {
              type: "category",
              label: "SAML",
              link: {
                type: "generated-index",
                title: "SAML Service API",
                slug: "/apis/resources/saml_service_v2",
                description:
                  "Get SAML Request details and create responses.\n",
              },
              items: sidebar_api_saml_service_v2,
            },
            {
              type: "category",
              label: "Settings",
              link: {
                type: "generated-index",
                title: "Settings Service API",
                slug: "/apis/resources/settings_service_v2",
                description:
                  "This API is intended to manage settings in a ZITADEL instance.\n",
              },
              items: sidebar_api_settings_service_v2,
            },
            {
              type: "category",
              label: "Feature",
              link: {
                type: "generated-index",
                title: "Feature Service API",
                slug: "/apis/resources/feature_service_v2",
                description:
                  'This API is intended to manage features for ZITADEL. Feature settings that are available on multiple "levels", such as instance and organization. The higher level instance acts as a default for the lower level. When a feature is set on multiple levels, the lower level takes precedence. Features can be experimental where ZITADEL will assume a sane default, such as disabled. When over time confidence in such a feature grows, ZITADEL can default to enabling the feature. As a final step we might choose to always enable a feature and remove the setting from this API, reserving the proto field number. Such removal is not considered a breaking change. Setting a removed field will effectively result in a no-op.\n',
              },
              items: sidebar_api_feature_service_v2,
            },
            {
              type: "category",
              label: "Organization",
              link: {
                type: "generated-index",
                title: "Organization Service API",
                slug: "/apis/resources/org_service/v2",
                description:
                  "This API is intended to manage organizations for ZITADEL. \n",
              },
              items: sidebar_api_org_service_v2,
            },
            {
              type: "category",
              label: "Identity Provider",
              link: {
                type: "generated-index",
                title: "Identity Provider Service API",
                slug: "/apis/resources/idp_service_v2",
                description:
                  "This API is intended to manage identity providers (IdPs) for ZITADEL.\n",
              },
              items: sidebar_api_idp_service_v2,
            },
            {
              type: "category",
              label: "Web Key",
              link: {
                type: "generated-index",
                title: "Web Key Service API",
                slug: "/apis/resources/webkey_service_v2",
                description:
                  "This API is intended to manage web keys for a ZITADEL instance, used to sign and validate OIDC tokens.\n" +
                  "\n" +
                  "The public key endpoint (outside of this service) is used to retrieve the public keys of the active and inactive keys.\n",
              },
              items: sidebar_api_webkey_service_v2
            },
            {
              type: "category",
              label: "Action",
              link: {
                type: "generated-index",
                title: "Action Service API",
                slug: "/apis/resources/action_service_v2",
                description:
                  "This API is intended to manage custom executions and targets (previously known as actions) in a ZITADEL instance.\n" +
                  "\n" +
                  "The version 2 of actions provide much more options to customize ZITADELs behaviour than previous action versions.\n" +
                  "Also, v2 actions are available instance-wide, whereas previous actions had to be managed for each organization individually\n" +
                  "ZITADEL doesn't restrict the implementation languages, tooling and runtime for v2 action executions anymore.\n" +
                  "Instead, it calls external endpoints which are implemented and maintained by action v2 users."
              },
              items: sidebar_api_actions_v2,
            },
            {
              type: "category",
              label: "Instance",
              link: {
                type: "generated-index",
                title: "Instance Service API",
                slug: "/apis/resources/instance_service_v2",
                description:
                  "This API is intended to manage instances, custom domains and trusted domains in ZITADEL.\n" +
                  "\n" +
                  "This v2 of the API provides the same functionalities as the v1, but organised on a per resource basis.\n" +
                  "The whole functionality related to domains (custom and trusted) has been moved under this instance API."
                ,
              },
              items: sidebar_api_instance_service_v2,
            },
            {
              type: "category",
              label: "Project",
              link: {
                type: "generated-index",
                title: "Project Service API",
                slug: "/apis/resources/project_service_v2",
                description:
                  "This API is intended to manage projects and subresources for ZITADEL."
              },
              items: sidebar_api_project_service_v2,
            },
            {
              type: "category",
              label: "Application",
              link: {
                type: "generated-index",
                title: "Application Service API",
                slug: "/apis/resources/application_service_v2",
                description:
                  "This API lets you manage Zitadel applications (API, SAML, OIDC).\n" +
                  "\n" +
                  "The API offers generic endpoints that work for all app types (API, SAML, OIDC), "
              },
              items: sidebar_api_application_v2,
            },
            {
              type: "category",
              label: "Authorizations",
              link: {
                type: "generated-index",
                title: "Authorization Service API",
                slug: "/apis/resources/authorization_service_v2",
                description:
                  "AuthorizationService provides methods to manage authorizations for users within your projects and applications.\n" +
                  "\n" +
                  "For managing permissions and roles for ZITADEL internal resources, like organizations, projects,\n" +
                  "users, etc., please use the InternalPermissionService."
              },
              items: sidebar_api_authorization_service_v2,
            },
            {
              type: "category",
              label: "Internal Permissions",
              link: {
                type: "generated-index",
                title: "Internal Permission Service API",
                slug: "/apis/resources/internal_permission_service_v2",
                description:
                  "This API provides methods to manage permissions for resource and and their management in ZITADEL itself also known as \"administrators\"."
              },
              items: sidebar_api_internal_permission_service_v2,
            },
          ],
        },
        {
          type: "category",
          label: "V1",
          collapsed: true,
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
              items: sidebar_api_auth,
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
              items: sidebar_api_mgmt,
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
              items: sidebar_api_admin,
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
              items: sidebar_api_system,
            },
            "apis/migration_v1_to_v2"
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
      label: "Observability",
      collapsed: true,
      items: [
        "apis/observability/metrics",
        "apis/observability/health",
      ],
    },
    {
      type: "category",
      label: "Provision Users",
      collapsed: true,
      items: ["apis/scim2"],
    },
    {
      type: "doc",
      label: "gRPC Status Codes",
      id: "apis/statuscodes",
    },
    {
      type: "link",
      label: "Rate Limits (Cloud)",
      href: "/legal/policies/rate-limit-policy",
    },
    {
      type: "category",
      label: "Benchmarks",
      collapsed: true,
      link: {
        type: "doc",
        id: "apis/benchmarks/index",
      },
      items: [
        {
          type: "category",
          label: "v2.65.0",
          link: {
            title: "v2.65.0",
            slug: "/apis/benchmarks/v2.65.0",
            description: "Benchmark results of Zitadel v2.65.0\n",
          },
          items: ["apis/benchmarks/v2.65.0/machine_jwt_profile_grant/index"],
        },
        {
          type: "category",
          label: "v2.66.0",
          link: {
            title: "v2.66.0",
            slug: "/apis/benchmarks/v2.66.0",
            description: "Benchmark results of Zitadel v2.66.0\n",
          },
          items: ["apis/benchmarks/v2.66.0/machine_jwt_profile_grant/index"],
        },
        {
          type: "category",
          label: "v2.70.0",
          link: {
            title: "v2.70.0",
            slug: "/apis/benchmarks/v2.70.0",
            description: "Benchmark results of Zitadel v2.70.0\n",
          },
          items: [
            "apis/benchmarks/v2.70.0/machine_jwt_profile_grant/index",
            "apis/benchmarks/v2.70.0/oidc_session/index",
          ],
        },
        {
          type: "category",
          label: "v4",
          link: {
            title: "v4",
            slug: "/apis/benchmarks/v4",
            description: "Benchmark results of Zitadel v4\n",
          },
          items: [
            "apis/benchmarks/v4/add_session/index",
            "apis/benchmarks/v4/human_password_login/index",
            "apis/benchmarks/v4/introspect/index",
            "apis/benchmarks/v4/machine_client_credentials_login/index",
            "apis/benchmarks/v4/machine_jwt_profile_grant/index",
            "apis/benchmarks/v4/machine_pat_login/index",
            "apis/benchmarks/v4/manipulate_user/index",
            "apis/benchmarks/v4/oidc_session/index",
            "apis/benchmarks/v4/otp_session/index",
            "apis/benchmarks/v4/password_session/index",
            "apis/benchmarks/v4/user_info/index",
          ],
        },
      ],
    },
  ],
  selfHosting: [{
    type: "link",
    label: "Self-Hosting Home",
    href:
      "/docs/self-hosting/deploy/overview"
  }
  ],
  legal: [
    {
      type: "category",
      label: "Legal Agreements",
      collapsed: true,
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
          collapsed: true,
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
            {
              type: "link",
              label: "Billing",
              href: "https://help.zitadel.com/pricing-and-billing-of-zitadel-services"
            }
          ],
        },
        {
          type: "category",
          label: "Policies",
          collapsed: true,
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
