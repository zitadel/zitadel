module.exports = {
  guides: [
    "guides/overview",
    {
      type: "category",
      label: "Get started",
      collapsed: false,
      items: [
        "guides/start/quickstart",
      ],
    },
    "examples/sdks",
    {
      type: "category",
      label: "Quickstarts",
      items: [
        {
          type: "category",
          label: "Frontend",
          items: [
            "examples/login/angular",
            "examples/login/react",
            "examples/login/flutter",
            "examples/login/nextjs",
          ],
          collapsed: false,
        },
        {
          type: "category",
          label: "Backend",
          items: [
            "examples/secure-api/go", 
            "examples/secure-api/python-flask", 
            "examples/secure-api/dot-net"
          ],
          collapsed: false,
        },
      ],
      collapsed: true,
    },
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
            "guides/manage/customize/user-metadata",
          ],
        },
        {
          type: "category",
          label: "Terraform",
          items: ["guides/manage/terraform/basics"],
        },
        "guides/manage/user/reg-create-user",
      ],
    },
    {
      type: "category",
      label: "Integrate",
      collapsed: true,
      items: [
        {
          type: "category",
          label: "Authenticate Users",
          collapsed: true,
          items: [
            "guides/integrate/login-users",
            "guides/integrate/oauth-recommended-flows",
            "guides/integrate/identity-brokering",
            "guides/integrate/logout",
          ],
        },
        {
          type: "category",
          label: "Configure External IDPs",
          collapsed: true,
          items: [
            "guides/integrate/auth0-oidc",
            "guides/integrate/auth0-saml",
            "guides/integrate/azuread-oidc",
            "guides/integrate/pingidentity-saml",
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
          collapsed: false,
          items: [
            "guides/integrate/authenticated-mongodb-charts",
            "guides/integrate/gitlab-self-hosted",
            "guides/integrate/aws-saml",
            "guides/integrate/atlassian-saml",
            "guides/integrate/gitlab-saml",
          ],
        },
        {
          type: "category",
          label: "Infrastructure",
          collapsed: false,
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
        "guides/solution-scenarios/configurations",
      ],
    },
    {
      type: "category",
      label: "Concepts",
      collapsed: true,
      items: [
        "concepts/introduction",
        "concepts/principles",
        {
          type: "category",
          label: "Eventstore",
          collapsed: false,
          items: [
            "concepts/eventstore/overview",
            "concepts/eventstore/implementation",
          ],
        },
        {
          type: "category",
          label: "Architecture",
          collapsed: false,
          items: [
            "concepts/architecture/software",
            "concepts/architecture/solution",
            "concepts/architecture/secrets",
          ],
        },
        {
          type: "category",
          label: "Structure",
          collapsed: false,
          items: [
            "concepts/structure/overview",
            "concepts/structure/instance",
            "concepts/structure/organizations",
            "concepts/structure/projects",
            "concepts/structure/applications",
            "concepts/structure/granted_projects",
            "concepts/structure/users",
            "concepts/structure/managers",
            "concepts/structure/policies",
            "concepts/structure/jwt_idp",
          ],
        },
        {
          type: "category",
          label: "Use Cases",
          collapsed: false,
          items: ["concepts/usecases/saas"],
        },
        {
          type: "category",
          label: "Features",
          collapsed: false,
          items: [
            "concepts/features/actions",
            "concepts/features/selfservice"
          ],
        },
      ]
    },
  ],
  apis: [
    "apis/introduction",
    {
      type: "category",
      label: "API Definition",
      collapsed: false,
      items: [
        "apis/statuscodes",
        {
          type: "category",
          label: "Proto",
          collapsed: true,
          items: [
            "apis/proto/auth",
            "apis/proto/management",
            "apis/proto/admin",
            "apis/proto/system",
            "apis/proto/instance",
            "apis/proto/org",
            "apis/proto/user",
            "apis/proto/app",
            "apis/proto/policy",
            "apis/proto/auth_n_key",
            "apis/proto/change",
            "apis/proto/idp",
            "apis/proto/member",
            "apis/proto/metadata",
            "apis/proto/message",
            "apis/proto/text",
            "apis/proto/action",
            "apis/proto/object",
            "apis/proto/options",
          ],
        },
        {
          type: "category",
          label: "Assets API",
          collapsed: true,
          items: ["apis/assets/assets"],
        },
      ],
    },
    {
      type: "category",
      label: "OpenID Connect & OAuth",
      collapsed: false,
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
      label: "SAML",
      collapsed: false,
      items: ["apis/saml/endpoints"],
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
    "support/introduction",
    {
      type: "category",
      label: "Trainings",
      collapsed: true,
      items: [
        "guides/trainings/introduction",
        "guides/trainings/application",
        "guides/trainings/recurring",
        "guides/trainings/project",
      ],
    },
  ],
  concepts: [
    
  ],
  manuals: [
    "manuals/introduction",
    "manuals/user-profile",
    "manuals/user-login",
    "manuals/troubleshooting",
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
