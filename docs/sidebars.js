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
      label: "Integrate",
      collapsed: true,
      link: {
        type: 'generated-index',
        title: 'Overview',
        slug: 'guides/integrate',
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
          collapsed: true,
          items: [
            "guides/integrate/identity-providers/introduction",
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
              collapsed: true,
              items: [
                "guides/integrate/serviceusers",
                "guides/integrate/client-credentials",
                "guides/integrate/pat",
              ],
            },
            "guides/integrate/access-zitadel-apis",
            "guides/integrate/access-zitadel-system-api",
            "guides/integrate/event-api",
            "guides/integrate/export-and-import",
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
          collapsed: true,
          items: [
            "guides/integrate/services/gitlab-self-hosted",
            "guides/integrate/services/aws-saml",
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
      collapsed: true,
      items: [
        "guides/solution-scenarios/introduction",
        "guides/solution-scenarios/b2c",
        "guides/solution-scenarios/b2b",
        "concepts/usecases/saas",
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
        "concepts/structure/jwt_idp",
        "concepts/features/actions",
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
      type: "category",
      label: "Rate Limits",
      collapsed: false,
      items: ["apis/ratelimits/ratelimits", "legal/rate-limit-policy"],
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
    "legal/introduction",
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
      label: "Additional terms",
      collapsed: true,
      items: [
        "legal/terms-support-service",
        "legal/terms-of-service-dedicated",
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
      ],
    },
  ],
};
