import json

# Docusaurus GUIDES structure (simplified for conversion)
guides_sidebar = [
    {
      "type": "category",
      "label": "Get Started",
      "collapsed": False,
      "items": [
        "guides/overview",
        "guides/start/quickstart",
        {
          "type": "category",
          "label": "Key Concepts",
          "items": [
            "concepts/structure/instance",
            "guides/manage/console/organizations-overview",
            "guides/manage/console/projects-overview",
            "guides/manage/console/applications-overview",
            "guides/manage/console/users-overview",
            "concepts/structure/managers",
          ],
        },
        {
          "type": "category",
          "label": "Authenticate Users",
          "items": [
            "guides/integrate/login/login-users",
            "guides/integrate/login/hosted-login",
            "guides/integrate/login/oidc/logout",
          ],
        },
        {
          "type": "category",
          "label": "Example Applications",
          "items": [
            {
              "type": "category",
              "label": "Frontend (SPA)",
              "items": [
                {
                  "type": "link",
                  "label": "Vanilla-JS",
                  "href": "https://github.com/zitadel/zitadel-vanilla-js",
                },
                "examples/login/react",
                "examples/login/angular",
                "examples/login/vue",
              ],
            },
            {
              "type": "category",
              "label": "Mobile & Native",
              "items": [
                "examples/login/flutter",
              ],
            },
            {
              "type": "category",
              "label": "Full-Stack / SSR",
              "items": [
                "examples/login/nextjs",
                "examples/login/nextjs-b2b",
              ],
            },
            {
              "type": "category",
              "label": "Web App (Server-Side)",
              "items": [
                "examples/login/symfony",
                "examples/login/java-spring",
                "examples/login/python-django",
                "examples/login/go"
              ],
            },
            {
              "type": "category",
              "label": "APIs / Backend Services",
              "items": [
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
          "type": "category",
          "label": "Use Cases",
          "items": [
            "guides/solution-scenarios/b2b",
            "guides/solution-scenarios/configurations",
            "guides/solution-scenarios/saas",
            "guides/solution-scenarios/b2c",
            "guides/solution-scenarios/frontend-calling-backend-API",
            {
              "type": "category",
              "label": "Machine-to-Machine (M2M)",
              "items": [
                "guides/integrate/service-users/authenticate-service-users",
                "guides/integrate/service-users/private-key-jwt",
                "guides/integrate/service-users/client-credentials",
                "guides/integrate/service-users/personal-access-token",
              ],
            },
          ],
        },
        {
          "type": "category",
          "label": "Onboard Customers and Users",
          "items": [
            "guides/integrate/onboarding/b2b",
            "guides/integrate/onboarding/end-users",
          ],
        },
        {
          "type": "category",
          "label": "Branding & Customization",
          "items": [
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
      "type": "category",
      "label": "Integrate & Authenticate",
      "items": [
        {
          "type": "category",
          "label": "OIDC & OAuth Flows",
          "items": [
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
          "type": "category",
          "label": "SAML",
          "items": [
            "guides/integrate/login/saml",
            "apis/saml/endpoints",
          ],
        },
        {
          "type": "category",
          "label": "API Access",
          "items": [
            "guides/integrate/zitadel-apis/access-zitadel-apis",
            "guides/integrate/zitadel-apis/access-zitadel-system-api",
            "guides/integrate/zitadel-apis/event-api",
            "guides/integrate/zitadel-apis/example-zitadel-api-with-go",
            "guides/integrate/zitadel-apis/example-zitadel-api-with-dot-net",
            {
              "type": "link",
              "label": "Zitadel APIs",
              "href": "/docs/apis/introduction",
            },
          ],
        },
        {
          "type": "category",
          "label": "SDKs",
          "items": [
            "sdk-examples/introduction",
            {
              "type": "category",
              "label": "Frontend (SPA)",
              "items": [
                "sdk-examples/react",
                "sdk-examples/angular",
                "sdk-examples/vue",
              ],
            },
            {
              "type": "category",
              "label": "Mobile & Native",
              "items": [
                {
                  "type": "link",
                  "label": "Dart / Flutter",
                  "href": "https://github.com/smartive/zitadel-dart",
                },
                {
                  "type": "link",
                  "label": ".NET (MAUI/Xamarin)",
                  "href": "https://github.com/smartive/zitadel-net",
                },
              ],
            },
            {
              "type": "category",
              "label": "Full-Stack / SSR",
              "items": [
                "sdk-examples/nextjs",
                "sdk-examples/nuxtjs",
                "sdk-examples/sveltekit",
                "sdk-examples/qwik",
                "sdk-examples/solidstart",
                "sdk-examples/astro",
                {
                  "type": "link",
                  "label": "NextAuth",
                  "href": "https://next-auth.js.org/providers/zitadel",
                },
              ],
            },
            {
              "type": "category",
              "label": "Backend & API",
              "items": [
                {
                  "type": "category",
                  "label": "Node.js",
                  "items": [
                    "sdk-examples/client-libraries/node",
                    "sdk-examples/expressjs",
                    "sdk-examples/fastify",
                    "sdk-examples/hono",
                    "sdk-examples/nestjs",
                    {
                      "type": "link",
                      "label": "Passport.js",
                      "href": "https://github.com/buehler/node-passport-zitadel",
                    },
                    {
                      "type": "link",
                      "label": "Node.js (Community)",
                      "href": "https://www.npmjs.com/package/@zitadel/node",
                    },
                  ],
                },
                {
                  "type": "category",
                  "label": "Python",
                  "items": [
                    "sdk-examples/client-libraries/python",
                    "sdk-examples/flask",
                    "sdk-examples/django",
                    "sdk-examples/fastapi",
                  ],
                },
                {
                  "type": "category",
                  "label": "Go",
                  "items": ["sdk-examples/go"],
                },
                {
                  "type": "category",
                  "label": "Java",
                  "items": [
                    "sdk-examples/java",
                    "sdk-examples/spring",
                    "sdk-examples/client-libraries/java",
                  ],
                },
                {
                  "type": "category",
                  "label": "PHP",
                  "items": [
                    "sdk-examples/symfony",
                    "sdk-examples/laravel",
                    "sdk-examples/client-libraries/php",
                  ],
                },
                {
                  "type": "category",
                  "label": "Other Languages",
                  "items": [
                    "sdk-examples/client-libraries/ruby",
                    {
                      "type": "link",
                      "label": "Elixir",
                      "href": "https://github.com/maennchen/zitadel_api",
                    },
                    {
                      "type": "link",
                      "label": "Rust",
                      "href": "https://github.com/smartive/zitadel-rust",
                    },
                    {
                      "type": "link",
                      "label": "Pylon",
                      "href": "https://github.com/getcronit/pylon",
                    },
                  ],
                },
              ],
            },
          ],
        },
        {
          "type": "category",
          "label": "SCIM",
          "items": ["guides/integrate/scim-okta-guide"]
        },
        {
          "type": "category",
          "label": "Token Introspection",
          "items": [
            "guides/integrate/token-introspection/basic-auth",
            "guides/integrate/token-introspection/private-key-jwt",
          ]
        },
        "guides/integrate/back-channel-logout",
        {
          "type": "category",
          "label": "External Integrations",
          "items": [
            {
              "type": "category",
              "label": "Services",
              "items": [
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
                  "type": "link",
                  "label": "Bold BI",
                  "href": "https://support.boldbi.com/kb/article/13708/how-to-configure-zitadel-oauth-login-in-bold-bi",
                },
                {
                  "type": "link",
                  "label": "Cloudflare Workers",
                  "href": "https://zitadel.com/blog/increase-spa-security-with-cloudflare-workers",
                },
                {
                  "type": "link",
                  "label": "Firezone",
                  "href": "https://www.firezone.dev/docs/authenticate/oidc/zitadel",
                },
                {
                  "type": "link",
                  "label": "Netbird",
                  "href": "https://docs.netbird.io/selfhosted/identity-providers",
                },
                {
                  "type": "link",
                  "label": "Nextcloud",
                  "href": "https://zitadel.com/blog/zitadel-as-sso-provider-for-selfhosting",
                },
                {
                  "type": "link",
                  "label": "Psono",
                  "href": "https://doc.psono.com/admin/configuration/oidc-zitadel.html",
                },
                {
                  "type": "link",
                  "label": "Zoho Desk",
                  "href": "https://help.zoho.com/portal/en/kb/desk/user-management-and-security/data-security/articles/setting-up-saml-single-signon-for-help-center#Zitadel_IDP",
                },
              ],
            },
            {
              "type": "category",
              "label": "Tools",
              "items": [
                {
                  "type": "link",
                  "label": "Argo CD",
                  "href": "https://argo-cd.readthedocs.io/en/latest/operator-manual/user-management/zitadel/",
                },
                "guides/integrate/tools/apache2",
                "guides/integrate/authenticated-mongodb-charts",
                "examples/identity-proxy/oauth2-proxy",
              ],
            },
          ],
        },
        {
          "type": "category",
          "label": "Build your own Login UI",
          "items": [
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
          "type": "category",
          "label": "Actions",
          "items": [
            {
              "type": "category",
              "label": "V1",
              "items": [
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
              "type": "category",
              "label": "V2",
              "items": [
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
          "type": "category",
          "label": "Migrate",
          "items": [
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
      "type": "category",
      "label": "Configure Identity & Policies",
      "items": [
        "guides/manage/user/reg-create-user",
        "guides/manage/terraform-provider",
        {
          "type": "category",
          "label": "Identity Providers",
          "items": [
            {
              "type": "category",
              "label": "Social Logins",
              "items": [
                "guides/integrate/identity-providers/google",
                "guides/integrate/identity-providers/apple",
                "guides/integrate/identity-providers/github",
                "guides/integrate/identity-providers/gitlab",
                "guides/integrate/identity-providers/linkedin-oauth",
              ],
            },
            {
              "type": "category",
              "label": "Enterprise (SAML & OIDC)",
              "items": [
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
              "type": "category",
              "label": "Legacy & Directory (LDAP)",
              "items": [
                "guides/integrate/identity-providers/ldap",
                "guides/integrate/identity-providers/openldap",
              ],
            },
            {
              "type": "category",
              "label": "Custom & Generic",
              "items": [
                "guides/integrate/identity-providers/generic-oidc",
                "guides/integrate/identity-providers/jwt_idp",
                "guides/integrate/identity-providers/mocksaml",
              ],
            },
            {
              "type": "category",
              "label": "Guides",
              "items": [
                "guides/integrate/identity-providers/migrate",
                "guides/integrate/identity-providers/additional-information",
              ]
            }
          ],
        },
        {
          "type": "category",
          "label": "Policies",
          "items": [
            "guides/manage/console/default-settings",
            "concepts/structure/policies",
          ],
        },
        {
          "type": "category",
          "label": "Roles & Permissions",
          "items": [
            "guides/manage/console/roles",
            "guides/integrate/retrieve-user-roles",
            "concepts/structure/granted_projects",
            "guides/manage/console/managers",
          ],
        },
        {
          "type": "category",
          "label": "Compliance & Security",
          "items": [
            "concepts/features/audit-trail",
            "guides/integrate/external-audit-log",
          ],
        },
      ],
    },
    {
      "type": "category",
      "label": "Test & Debug",
      "items": [
        {
          "type": "link",
          "label": "OIDC Playground",
          "href": "https://zitadel.com/playgrounds/oidc",
        },
      ],
    },
    {
      "type": "category",
      "label": "Deploy & Operate",
      "items": [
        "guides/manage/console/console-overview",
        {
          "type": "category",
          "label": "Customer Portal",
          "items": [
            "guides/manage/cloud/start",
            "guides/manage/cloud/instances",
            "guides/manage/cloud/settings",
            "guides/manage/cloud/billing",
            "guides/manage/cloud/users",
          ],
        },
        {
          "type": "category",
          "label": "Self-Hosted",
          "items": [
            "self-hosting/deploy/overview",
            "self-hosting/deploy/linux",
            "self-hosting/deploy/macos",
            "self-hosting/deploy/devcontainer",
            "self-hosting/deploy/compose",
            "self-hosting/deploy/kubernetes",
            {
              "type": "category",
              "label": "Manage",
              "items": [
                {
                  "type": "category",
                  "label": "Production & Operations",
                  "items": [
                    "self-hosting/manage/production",
                    "self-hosting/manage/productionchecklist",
                    "self-hosting/manage/usage_control",
                  ],
                },
                {
                  "type": "category",
                  "label": "Configuration",
                  "items": [
                    "self-hosting/manage/configure/configure",
                    "self-hosting/manage/custom-domain",
                    "self-hosting/manage/tls_modes",
                    "self-hosting/manage/http2",
                    "self-hosting/manage/login-client",
                  ],
                },
                {
                  "type": "category",
                  "label": "Reverse Proxy",
                  "items": [
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
                  "type": "category",
                  "label": "Observability",
                  "items": [
                    "self-hosting/manage/service_ping",
                    {
                      "type": "category",
                      "label": "Metrics",
                      "items": ["self-hosting/manage/metrics/prometheus"],
                    },
                  ],
                },
                {
                  "type": "category",
                  "label": "Tools",
                  "items": [
                    {
                      "type": "category",
                      "label": "Command Line Interface",
                      "items": ["self-hosting/manage/cli/mirror"],
                    },
                  ],
                },
              ],
            },
            "self-hosting/deploy/troubleshooting/troubleshooting",
          ],
        },
        {
          "type": "category",
          "label": "Scaling & Performance",
          "items": [
            "self-hosting/manage/updating_scaling",
            "self-hosting/manage/database/database",
            "self-hosting/manage/cache",
          ],
        },
      ],
    },
    {
      "type": "category",
      "label": "Architecture & Concepts",
      "items": [
        {
          "type": "category",
          "label": "System Architecture",
          "items": [
            "concepts/architecture/solution",
            "concepts/architecture/software",
            "concepts/architecture/secrets",
          ],
        },
        "concepts/principles",
        {
          "type": "category",
          "label": "Advanced Topics",
          "items": [
            {
              "type": "category",
              "label": "Event Store",
              "items": [
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
      "type": "category",
      "label": "Product, Releases & Support",
      "items": [
        {
          "type": "category",
          "label": "Product Features",
          "items": [
            "product/roadmap",
          ],
        },
        {
          "type": "category",
          "label": "Releases",
          "items": [
            {
              "type": "link",
              "label": "Changelog",
              "href": "https://zitadel.com/changelog",
            },
            "product/release-cycle",
          ],
        },
        {
          "type": "category",
          "label": "Support Resources",
          "items": [
            "support/troubleshooting",
            "support/technical_advisory",
            {
              "type": "link",
              "label": "Support States",
              "href": "https://help.zitadel.com/zitadel-support-states",
            },
            {
              "type": "link",
              "label": "Zitadel Release Cycle",
              "href": "https://help.zitadel.com/zitadel-software-release-cycle",
            },
            {
              "type": "link",
              "label": "Knowledge Base",
              "href": "https://help.zitadel.com",
            }
          ],
        },
        {
          "type": "link",
          "label": "Pricing",
          "href": "https://zitadel.com/pricing",
        },
        {
          "type": "link",
          "label": "Contact Us",
          "href": "https://zitadel.com/contact",
        },
        "guides/manage/cloud/support",
      ],
    },
]

def map_path(p):
    if not isinstance(p, str): return p
    # Relative to "guides/" folder in "content/docs/guides/meta.json"
    
    # Strip "guides/"
    if p.startswith("guides/"):
        return p[7:] # remove "guides/"
    
    # Prepend "../" for siblings
    if p.startswith("concepts/") or \
       p.startswith("apis/") or \
       p.startswith("examples/") or \
       p.startswith("sdk-examples/") or \
       p.startswith("self-hosting/") or \
       p.startswith("product/") or \
       p.startswith("support/"):
       return "../" + p
       
    return p

def map_item(item):
    if isinstance(item, str):
        return map_path(item)
    if item["type"] == "category":
        return {
            "title": item["label"],
            "pages": [map_item(i) for i in item.get("items", [])]
        }
    if item["type"] == "link":
        return {
            "title": item["label"],
            "url": item["href"]
        }
    if item["type"] == "doc":
        if "label" in item:
             return {
                 "title": item["label"],
                 "url": map_path(item["id"])
             }
        return map_path(item["id"])
    return None

meta_pages = [map_item(i) for i in guides_sidebar]

print(json.dumps({"root": True, "pages": meta_pages}, indent=2))
