
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
const sidebar_api_org_service_v2beta = require("./docs/apis/resources/org_service_v2beta/sidebar.ts").default
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
          items: [
            "concepts/structure/instance",
            "concepts/structure/organizations",
            "guides/manage/console/organizations",
            "concepts/structure/projects",
            "guides/manage/console/projects",
            "concepts/structure/applications",
            "guides/manage/console/applications",
            "concepts/structure/users",
            "guides/manage/console/users",
            "concepts/structure/managers",
          ],
        },
        {
          type: "category",
          label: "Authenticate Users",
          items: [
            "guides/integrate/login/login-users",
            "guides/integrate/login/oidc/login-users",
            "guides/integrate/login/hosted-login",
            "guides/integrate/login/oidc/logout",
          ],
        },
        {
          type: "category",
          label: "Example Applications",
          items: [
            {
              type: "category",
              label: "Frontend (SPA)",
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
              items: [
                "examples/login/flutter",
              ],
            },
            {
              type: "category",
              label: "Full-Stack / SSR",
              items: [
                "examples/login/nextjs",
                "examples/login/nextjs-b2b",
              ],
            },
            {
              type: "category",
              label: "Web App (Server-Side)",
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
          items: [
            "guides/solution-scenarios/b2b",
            "guides/solution-scenarios/b2c",
            "guides/solution-scenarios/saas",
            "guides/solution-scenarios/frontend-calling-backend-API",
            "guides/solution-scenarios/configurations",
            "guides/integrate/onboarding/b2b",
            "guides/integrate/onboarding/end-users",
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
      label: "Configure Identity & Policies",
      collapsed: true,
      items: [
        "guides/manage/user/reg-create-user",
        "guides/manage/terraform-provider",
        {
          type: "category",
          label: "Identity Providers",
          items: [
            {
              type: "category",
              label: "External IDPs",
              items: [
                "guides/integrate/identity-providers/introduction",
                "guides/integrate/identity-providers/google",
                "guides/integrate/identity-providers/azure-ad-oidc",
                "guides/integrate/identity-providers/azure-ad-saml",
                "guides/integrate/identity-providers/github",
                "guides/integrate/identity-providers/gitlab",
                "guides/integrate/identity-providers/apple",
                "guides/integrate/identity-providers/okta-oidc",
                "guides/integrate/identity-providers/okta-saml",
                "guides/integrate/identity-providers/keycloak",
                "guides/integrate/identity-providers/linkedin-oauth",
                "guides/integrate/identity-providers/onelogin-saml",
                "guides/integrate/identity-providers/pingfederate-saml",
              ],
            },
            {
              type: "category",
              label: "Custom Providers",
              items: [
                "guides/integrate/identity-providers/generic-oidc",
                "guides/integrate/identity-providers/jwt_idp",
                "guides/integrate/identity-providers/ldap",
                "guides/integrate/identity-providers/openldap",
                "guides/integrate/identity-providers/mocksaml",
                "guides/integrate/identity-providers/migrate",
                "guides/integrate/identity-providers/additional-information",
              ],
            },
          ],
        },
        {
          type: "category",
          label: "Policies",
          items: [
            "guides/manage/console/default-settings",
            "concepts/structure/policies",
            // TODO: File missing - Password Policies specific doc
            // TODO: File missing - Lockout Policies specific doc
            // TODO: File missing - Privacy Policies specific doc
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
            // TODO: File missing - Certifications specific doc
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Integrate & Authenticate",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "OIDC & OAuth Flows",
          items: [
            "guides/integrate/login/oidc/oauth-recommended-flows",
            "apis/openidoauth/authn-methods",
            "apis/openidoauth/endpoints",
            "apis/openidoauth/scopes",
            "apis/openidoauth/claims",
            "apis/openidoauth/grant-types",
            "apis/openidoauth/authrequest",
            "guides/integrate/login/oidc/device-authorization",
            "guides/integrate/login/oidc/webkeys",
            "guides/integrate/token-exchange",
          ],
        },
        {
          type: "category",
          label: "SAML",
          items: [
            "apis/saml/endpoints",
            "guides/integrate/login/saml",
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
                    "sdk-examples/client-libraries/node", // Generic Node SDK
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
                    "sdk-examples/python-flask",
                    "sdk-examples/python-django",
                    {
                      type: "link",
                      label: "FastAPI",
                      href: "https://github.com/cleanenergyexchange/fastapi-zitadel-auth",
                    },
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
                    "sdk-examples/client-libraries/java",
                  ],
                },
                {
                  type: "category",
                  label: "PHP",
                  items: [
                    "sdk-examples/symfony",
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
          label: "Build Custom Login UI",
          items: [
            "guides/integrate/login-ui/login-app",
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
                "guides/manage/console/actions",
                "apis/actions/introduction",
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
        {
          type: "category",
          label: "Common Issues",
          items: [
            "guides/solution-scenarios/domain-discovery",
            "guides/solution-scenarios/restrict-console",
          ],
        }
      ],
    },
    {
      type: "category",
      label: "Deploy & Operate",
      collapsed: true,
      items: [
        "concepts/features/console",
        "guides/manage/console/overview",
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
            "guides/manage/cloud/support",
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
            // TODO: File missing - Data Flow Diagrams specific doc
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
            // TODO: File missing - Pricing & Tiers specific doc
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
              type: "autogenerated",
              dirName: "support/advisory",
            },
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
            // TODO: File missing - Community links doc
            // TODO: File missing - FAQ doc
            // TODO: File missing - Contact Support doc
            // TODO: File missing - Contribute doc
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
              label: "Organization (Beta)",
              link: {
                type: "generated-index",
                title: "Organization Service Beta API",
                slug: "/apis/resources/org_service/v2beta",
                description:
                  "This beta API is intended to manage organizations for ZITADEL. Expect breaking changes to occur. Please use the v2 version for a stable API. \n",
              },
              items: sidebar_api_org_service_v2beta,
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
