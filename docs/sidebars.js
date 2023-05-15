module.exports = {
  guides: [
    "guides/overview",
    {
      type: "category",
      label: "Get started",
      collapsed: false,
      items: [
        "guides/start/quickstart",
        {
          type: "category",
          label: "Frontend",
          items: [
            "examples/login/angular",
            "examples/login/react",
            "examples/login/flutter",
            "examples/login/nextjs",
          ],
          collapsed: true,
        },
        {
          type: "category",
          label: "Backend",
          items: [
            "examples/secure-api/go",
            "examples/secure-api/python-flask",
            "examples/secure-api/dot-net"
          ],
          collapsed: true,
        },
      ],
    },
    "examples/sdks",
    {
      type: "category",
      label: "Example Applications",
      items: [
        "examples/introduction",
        {
          type: 'link',
          label: 'Frontend', // The link label
          href: '/examples/introduction#frontend', // The internal path
        },
        {
          type: 'link',
          label: 'Backend', // The link label
          href: '/examples/introduction#backend', // The internal path
        }
      ],
      collapsed: true,
    },
    {
      type: "category",
      label: "Manage",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "Cloud",
          items: [
            "guides/manage/cloud/overview",
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
          items: [
            "guides/manage/console/overview",
            "guides/manage/console/instance-settings",
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
          ],
        },
        {
          type: "category",
          label: "Terraform",
          items: ["guides/manage/terraform/basics"],
        },
        {
          type: "category",
          label: "Users",
          items: [
            "guides/manage/user/reg-create-user",
            "guides/manage/customize/user-metadata",
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
        {
          type: "category",
          label: "Sources",
          collapsed: true,
          items: [
            "guides/migrate/sources/zitadel",
            "guides/migrate/sources/auth0",
          ]
        },
      ]
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
          label: "Authenticate Users",
          collapsed: true,
          items: [
            "guides/integrate/login-users",
            "guides/integrate/oauth-recommended-flows",
            "guides/integrate/logout",
          ],
        },
        {
          type: "category",
          label: "Configure Identity Providers",
          link: {
            type: "generated-index",
            title: "Let users login with their preferred identity provider",
            slug: "/guides/integrate/identity-providers",
            description:
              "In the following guides you will learn how to configure and setup your preferred external identity provider in ZITADEL.",

          },
          collapsed: true,
          items: [
            "guides/integrate/identity-providers/google",
            "guides/integrate/identity-providers/azure-ad",
            "guides/integrate/identity-providers/github",
            "guides/integrate/identity-providers/gitlab",
            "guides/integrate/identity-providers/ldap",
            "guides/integrate/identity-providers/openldap",
            "guides/integrate/identity-providers/google-oidc",
            "guides/integrate/identity-providers/azuread-oidc",
          ],
        },
        {
          type: "category",
          label: "Access ZITADEL APIs",
          collapsed: true,
          items: [
            {
              type: "category",
              label: "Authenticate Service Users",
              link: {
                type: "generated-index",
                title: "Authenticate Service Users",
                slug: "/guides/integrate/serviceusers",
                description:
                  "How to authenticate service users",
              },
              collapsed: true,
              items: [
                "guides/integrate/private-key-jwt",
                "guides/integrate/client-credentials",
                "guides/integrate/pat",
              ],
            },
            "guides/integrate/access-zitadel-apis",
            "guides/integrate/access-zitadel-system-api",
            "guides/integrate/event-api",
            {
              type: "category",
              label: "Example Code",
              items: [
                "examples/call-zitadel-api/go",
                "examples/call-zitadel-api/dot-net",
              ],
              collapsed: true,
            },
          ],
        },
        {
          type: "category",
          label: "Services",
          link: {
            type: "generated-index",
            title: "Integrate ZITADEL with your favorite services",
            slug: "/guides/integrate/services",
            description:
              "With the guides in this section you will learn how to integrate ZITADEL with your services.",

          },
          collapsed: true,
          items: [
            "guides/integrate/services/gitlab-self-hosted",
            "guides/integrate/services/aws-saml",
            "guides/integrate/services/google-cloud",
            "guides/integrate/services/atlassian-saml",
            "guides/integrate/services/gitlab-saml",
            "guides/integrate/services/auth0-oidc",
            "guides/integrate/services/auth0-saml",
            "guides/integrate/services/pingidentity-saml",
          ],
        },
        {
          type: "category",
          label: "Tools",
          link: {
            type: "generated-index",
            title: "Integrate ZITADEL with your tools",
            slug: "/guides/integrate/tools",
            description:
              "With the guides in this section you will learn how to integrate ZITADEL with your favorite tools.",

          },
          collapsed: true,
          items: [
            "guides/integrate/authenticated-mongodb-charts",
            "examples/identity-proxy/oauth2-proxy"
          ],
        },
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
          "Customers of an SaaS Identity and Access Management System usually have all distinct use cases and requirements. This guide attempts to explain real-world implementations and break them down into Solution Scenarios which aim to help you getting started with ZITADEL.",
      },
      collapsed: true,
      items: [
        "guides/solution-scenarios/b2c",
        "guides/solution-scenarios/b2b",
        "guides/solution-scenarios/saas",
        "guides/solution-scenarios/domain-discovery",
        "guides/solution-scenarios/configurations",
      ],
    },
    {
      type: "category",
      label: "Concepts",
      collapsed: true,
      items: [
        "concepts/introduction",
        "concepts/structure/instance",
        "concepts/structure/organizations",
        "concepts/structure/projects",
        "concepts/structure/applications",
        "concepts/structure/granted_projects",
        "concepts/structure/users",
        "concepts/structure/managers",
        "concepts/structure/policies",
        "concepts/features/identity-brokering",
        "concepts/structure/jwt_idp",
        "concepts/features/actions",
        "concepts/features/audit-trail",
        "concepts/features/selfservice",
      ]
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
          label: "Eventstore",
          collapsed: true,
          items: [
            "concepts/eventstore/overview",
            "concepts/eventstore/implementation",
          ],
        },
      ]
    },
    {
      type: "category",
      label: "Support",
      collapsed: true,
      items: [
        "support/troubleshooting",
        {
          type: 'category',
          label: "Technical Advisory",
          link: {
            type: 'doc',
            id: 'support/technical_advisory',
          },
          collapsed: true,
          items: [
              {
                type: 'autogenerated',
                dirName: 'support/advisory',
              },
          ],
        },
        {
          type: "category",
          label: "Trainings",
          collapsed: true,
          items: [
            "support/trainings/introduction",
            "support/trainings/application",
            "support/trainings/recurring",
            "support/trainings/project",
          ],
        },
      ]
    },
  ],
  apis: [
    "apis/introduction",
    {
      type: "category",
      label: "Core Resources",
      collapsed: false,
      items: [
        {
          type: "category",
          label: "Authenticated User",
          link: {
            type: "generated-index",
            title: "Auth API",
            slug: "/apis/auth",
            description:
              "The authentication API (aka Auth API) is used for all operations on the currently logged in user. The user id is taken from the sub claim in the token.",

          },
          items: require("./docs/apis/auth/sidebar.js"),
        },
        {
          type: "category",
          label: "Organization Objects",
          link: {
            type: "generated-index",
            title: "Management API",
            slug: "/apis/mgmt",
            description:
              "The management API is as the name states the interface where systems can mutate IAM objects like, organizations, projects, clients, users and so on if they have the necessary access rights. To identify the current organization you can send a header x-zitadel-orgid or if no header is set, the organization of the authenticated user is set.",
          },
          items: require("./docs/apis/mgmt/sidebar.js"),
        },
        {
          type: "category",
          label: "Instance Objects",
          link: {
            type: "generated-index",
            title: "Admin API",
            slug: "/apis/admin",
            description:
              "This API is intended to configure and manage one ZITADEL instance itself.",
          },
          items: require("./docs/apis/admin/sidebar.js"),
        },
        {
          type: "category",
          label: "Instance Lifecycle",
          link: {
            type: "generated-index",
            title: "System API",
            slug: "/apis/system",
            description:
              "This API is intended to manage the different ZITADEL instances within the system.\n" +
              "\n" +
              "Checkout the guide how to access the ZITADEL System API.",
          },
          items: require("./docs/apis/system/sidebar.js"),
        },
        {
          type: "category",
          label: "User Lifecycle (Alpha)",
          link: {
            type: "generated-index",
            title: "User Service API (Alpha)",
            slug: "/apis/user_service",
            description:
              "This API is intended to manage users in a ZITADEL instance.\n"+
              "\n"+
              "This project is in alpha state. It can AND will continue breaking until the services provide the same functionality as the current login.",
          },
          items: require("./docs/apis/user_service/sidebar.js"),
        },
        {
          type: "category",
          label: "Session Lifecycle (Alpha)",
          link: {
            type: "generated-index",
            title: "Session Service API (Alpha)",
            slug: "/apis/session_service",
            description:
              "This API is intended to manage sessions in a ZITADEL instance.\n"+
              "\n"+
              "This project is in alpha state. It can AND will continue breaking until the services provide the same functionality as the current login.",
          },
          items: require("./docs/apis/session_service/sidebar.js"),
        },
        {
          type: "category",
          label: "Settings Lifecycle (Alpha)",
          link: {
            type: "generated-index",
            title: "Settings Service API (Alpha)",
            slug: "/apis/settings_service",
            description:
              "This API is intended to manage settings in a ZITADEL instance.\n"+
              "\n"+
              "This project is in alpha state. It can AND will continue breaking until the services provide the same functionality as the current login.",
          },
          items: require("./docs/apis/settings_service/sidebar.js"),
        },
        {
          type: "category",
          label: "Assets",
          collapsed: true,
          items: ["apis/assets/assets"],
        },
      ]
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
        "apis/actions/internal-authentication",
        "apis/actions/external-authentication",
        "apis/actions/complement-token",
        "apis/actions/objects",
      ]
    },
    {
      type: "doc",
      label: "gRPC Status Codes",
      id: "apis/statuscodes"
    },
    {
      type: "category",
      label: "Observability",
      collapsed: false,
      items: ["apis/observability/metrics", "apis/observability/health"],
    },
    {
      type: 'link',
      label: 'Rate Limits (Cloud)', // The link label
      href: '/legal/rate-limit-policy', // The internal path
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
        "self-hosting/deploy/knative",
        "self-hosting/deploy/kubernetes",
        "self-hosting/deploy/loadbalancing-example/loadbalancing-example"
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
        "self-hosting/manage/reverseproxy/reverse_proxy",
        "self-hosting/manage/custom-domain",
        "self-hosting/manage/http2",
        "self-hosting/manage/tls_modes",
        "self-hosting/manage/database/database",
        "self-hosting/manage/updating_scaling",
        "self-hosting/manage/quotas"
      ],
    },
  ],
  support: [
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
        {
          type: "category",
          label: "Service Description",
          collapsed: false,
          items: [
            "legal/cloud-service-description",
            "legal/service-level-description",
            "legal/support-services",
          ],
        },
        {
          type: "category",
          label: "Support Program",
          collapsed: true,
          items: [
            "legal/terms-support-service",
          ],
        },
        {
          type: "category",
          label: "Policies",
          collapsed: false,
          items: [
            "legal/privacy-policy",
            "legal/acceptable-use-policy",
            "legal/rate-limit-policy",
            "legal/vulnerability-disclosure-policy",
          ],
        },
      ]
    },
  ],
};
