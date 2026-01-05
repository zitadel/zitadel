
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
            "guides/manage/console/users",
            "concepts/structure/users",
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
            "guides/integrate/login-ui/logout",
            "guides/integrate/login/oidc/logout",
          ],
        },
        {
          type: "category",
          label: "Examples & SDKs",
          items: [
            {
              type: "category",
              label: "SPA / Frontend",
              items: [
                "sdk-examples/react",
                "sdk-examples/angular",
                "sdk-examples/vue",
                "sdk-examples/nextjs",
                "sdk-examples/nuxtjs",
                "sdk-examples/sveltekit",
                "sdk-examples/qwik",
                "sdk-examples/solidstart",
                "sdk-examples/astro",
                "examples/login/react",
                "examples/login/angular",
                "examples/login/vue",
                "examples/login/nextjs",
                "examples/login/flutter",
              ],
            },
            {
              type: "category",
              label: "Web Applications",
              items: [
                "sdk-examples/expressjs",
                "sdk-examples/fastify",
                "sdk-examples/hono",
                "sdk-examples/nestjs",
                "sdk-examples/symfony",
                "examples/login/symfony",
                "examples/login/java-spring",
                "examples/login/python-django",
              ],
            },
            {
              type: "category",
              label: "APIs / Backend Services",
              items: [
                "sdk-examples/go",
                "sdk-examples/java",
                "sdk-examples/python-flask",
                "sdk-examples/python-django",
                "examples/secure-api/go",
                "examples/secure-api/java-spring",
                "examples/secure-api/python-django",
                "examples/secure-api/python-flask",
                "examples/secure-api/nodejs-nestjs",
                "examples/secure-api/pylon",
              ],
            },
            {
              type: "category",
              label: "Hybrid / Full-Stack",
              items: [
                "examples/login/nextjs-b2b",
              ],
            },
            {
              type: "category",
              label: "Sample Applications",
              items: [
                "sdk-examples/introduction",
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
                  label: "FastAPI",
                  href: "https://github.com/cleanenergyexchange/fastapi-zitadel-auth",
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
                {
                  type: "link",
                  label: "Vanilla-JS",
                  href: "https://github.com/zitadel/zitadel-vanilla-js",
                },
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
      ],
    },
    {
      type: "category",
      label: "Configure Identity & Policies",
      collapsed: true,
      items: [
        "guides/manage/console/overview",
        "guides/manage/user/reg-create-user",
        "guides/manage/user/scim2",
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
            "guides/manage/customize/behavior",
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
            "apis/openidoauth/endpoints",
            "apis/openidoauth/scopes",
            "apis/openidoauth/claims",
            "apis/openidoauth/authn-methods",
            "apis/openidoauth/grant-types",
            "apis/openidoauth/authrequest",
            "guides/integrate/login/oidc/oauth-recommended-flows",
            "guides/integrate/login/oidc/device-authorization",
            "guides/integrate/login/oidc/webkeys",
            "guides/integrate/token-exchange",
          ],
        },
        {
          type: "category",
          label: "API Access",
          items: [
            {
              type: "category",
              label: "gRPC APIs",
              items: [
                "apis/introduction",
                "apis/v2",
                "apis/apis/index",
                "apis/migration_v1_to_v2",
                "guides/integrate/zitadel-apis/access-zitadel-apis",
                "guides/integrate/zitadel-apis/access-zitadel-system-api",
                "guides/integrate/zitadel-apis/event-api",
                "guides/integrate/zitadel-apis/example-zitadel-api-with-go",
                "guides/integrate/zitadel-apis/example-zitadel-api-with-dot-net",
              ],
            },
            {
              type: "category",
              label: "REST APIs",
              items: [
                {
                  type: "link",
                  label: "Zitadel APIs",
                  href: "/docs/apis/introduction",
                }
              ],
            },
          ],
        },
        {
          type: "category",
          label: "SDKs",
          items: [
            {
              type: "category",
              label: "Go",
              items: [
                "sdk-examples/go",
                "sdk-examples/client-libraries/node",
              ],
            },
            {
              type: "category",
              label: "JavaScript/TypeScript",
              items: [
                "sdk-examples/react",
                "sdk-examples/nextjs",
                "sdk-examples/expressjs",
                "sdk-examples/client-libraries/node",
              ],
            },
            {
              type: "category",
              label: "Python",
              items: [
                "sdk-examples/python-flask",
                "sdk-examples/python-django",
                "sdk-examples/client-libraries/python",
              ],
            },
            {
              type: "category",
              label: "Others",
              items: [
                "sdk-examples/java",
                "sdk-examples/client-libraries/java",
                "sdk-examples/client-libraries/php",
                "sdk-examples/client-libraries/ruby",
                {
                  type: "link",
                  label: "Dart",
                  href: "https://github.com/smartive/zitadel-dart",
                },
                {
                  type: "link",
                  label: ".Net",
                  href: "https://github.com/smartive/zitadel-net",
                },
                {
                  type: "link",
                  label: "Rust",
                  href: "https://github.com/smartive/zitadel-rust",
                },
              ],
            },
          ],
        },
        {
          type: "category",
          label: "External Integrations",
          items: [
            "apis/saml/endpoints",
            "guides/integrate/login/saml",
            "guides/integrate/identity-providers/ldap",
            "guides/integrate/identity-providers/openldap",
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
            "guides/integrate/scim-okta-guide",
            "guides/integrate/token-introspection/index",
            "guides/integrate/token-introspection/basic-auth",
            "guides/integrate/token-introspection/private-key-jwt",
            "guides/integrate/back-channel-logout",
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
              ],
            },
            {
              type: "category",
              label: "Tools",
              items: [
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
            // TODO: File missing - Token Errors specific doc
            // TODO: File missing - Integration Failures specific doc
          ],
        }
      ],
    },
    {
      type: "category",
      label: "Customer Portal",
      collapsed: true,
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
      label: "Deploy & Operate",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "Deployment Options",
          items: [
            {
              type: "category",
              label: "Self-Hosted",
              items: [
                "self-hosting/deploy/overview",
                "self-hosting/deploy/linux",
                "self-hosting/deploy/macos",
              ],
            },
            {
              type: "category",
              label: "Docker Compose",
              items: [
                "self-hosting/deploy/compose",
              ],
            },
            {
              type: "category",
              label: "Kubernetes",
              items: [
                "self-hosting/deploy/kubernetes",
              ],
            },
            "self-hosting/deploy/devcontainer",
            "self-hosting/deploy/troubleshooting/troubleshooting",
          ],
        },
        {
          type: "category",
          label: "Configuration",
          items: [
            "self-hosting/manage/configure/configure",
            "self-hosting/manage/production",
            "self-hosting/manage/productionchecklist",
            "self-hosting/manage/custom-domain",
            "self-hosting/manage/tls_modes",
            "self-hosting/manage/http2",
            "self-hosting/manage/login-client",
            // TODO: File missing - Environment Variables specific doc
            // TODO: File missing - Advanced Settings specific doc
          ],
        },
        {
          type: "category",
          label: "Scaling & Performance",
          items: [
            "self-hosting/manage/updating_scaling",
            "self-hosting/manage/database/database",
            "self-hosting/manage/cache",
            // TODO: File missing - High Availability specific doc
            // TODO: File missing - Backup & Restore specific doc
          ],
        },
        {
          type: "category",
          label: "Monitoring & Maintenance",
          items: [
            "apis/observability/metrics",
            "apis/observability/health",
            "self-hosting/manage/service_ping",
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
        {
          type: "category",
          label: "Principles",
          items: [
            "concepts/principles",
            "guides/solution-scenarios/saas",
            "apis/openidoauth/authrequest",
          ],
        },
        {
          type: "category",
          label: "Advanced Topics",
          items: [
            "guides/manage/customize/branding",
            "guides/manage/customize/texts",
            "guides/manage/customize/restrictions",
            "guides/manage/customize/behavior",
            "guides/manage/customize/user-schema",
            "guides/manage/customize/user-metadata",
            "guides/manage/customize/notification-providers",
            "concepts/features/actions",
            "concepts/features/actions_v2",
            "concepts/features/selfservice",
            "concepts/features/custom-domain",
            "concepts/features/account-linking",
            "concepts/features/console",
            "concepts/features/external-user-grant",
            "concepts/features/identity-brokering",
            "concepts/features/passkeys",
            "concepts/knowledge/opaque-tokens",
            "concepts/eventstore/overview",
            "concepts/eventstore/implementation",
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
      collapsed: false,
      link: {
        type: "doc",
        id: "apis/apis/index",
      },
      items: [
        {
          type: "category",
          label: "V2",
          collapsed: false,
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
      label: "Provision Users",
      collapsed: true,
      items: ["apis/scim2"],
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
      collapsed: false,
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
        "self-hosting/deploy/kubernetes",
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
        "self-hosting/manage/login-client",
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
            "self-hosting/manage/reverseproxy/httpd/httpd",
            "self-hosting/manage/reverseproxy/cloudflare/cloudflare",
            "self-hosting/manage/reverseproxy/cloudflare_tunnel/cloudflare_tunnel",
            "self-hosting/manage/reverseproxy/zitadel_cloud/zitadel_cloud",
          ],
        },
        "self-hosting/manage/custom-domain",
        "self-hosting/manage/http2",
        "self-hosting/manage/tls_modes",
        "self-hosting/manage/database/database",
        "self-hosting/manage/cache",
        "self-hosting/manage/service_ping",
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
