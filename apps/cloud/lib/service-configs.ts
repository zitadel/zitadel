/**
 * Service configuration definitions for the debug UI.
 * This is a plain module (not "use server") so both client and server can import it.
 */

export interface ServiceConfig {
  title: string
  description: string
  vars: { key: string; label: string; type: "url" | "secret" | "text" | "json"; placeholder?: string }[]
}

export const serviceConfigs: ServiceConfig[] = [
  {
    title: "Login Instance",
    description: "The ZITADEL instance used for authentication",
    vars: [
      { key: "ZITADEL_API_URL", label: "API URL", type: "url", placeholder: "https://auth.rootd.ch" },
      { key: "ZITADEL_SERVICE_USER_TOKEN", label: "Service User Token", type: "secret" },
      { key: "EMAIL_VERIFICATION", label: "Email Verification", type: "text", placeholder: "false" },
    ],
  },
  {
    title: "Test Instances",
    description: "ZITADEL instances for local development and testing",
    vars: [
      { key: "ZITADEL_INSTANCES", label: "Instances (JSON array)", type: "json", placeholder: '[{"name":"Local","url":"http://localhost:8080","pat":"..."}]' },
    ],
  },
  {
    title: "Stripe",
    description: "Payment processing",
    vars: [
      { key: "STRIPE_SECRET_KEY", label: "Secret Key", type: "secret", placeholder: "sk_test_..." },
      { key: "NEXT_PUBLIC_STRIPE_PUBLIC_KEY", label: "Public Key", type: "text", placeholder: "pk_test_..." },
      { key: "STRIPE_WEBHOOK_SECRET", label: "Webhook Secret", type: "secret", placeholder: "whsec_..." },
    ],
  },
  {
    title: "HubSpot CRM",
    description: "Customer relationship management",
    vars: [
      { key: "HUBSPOT_ACCESS_TOKEN", label: "Access Token", type: "secret" },
      { key: "NEXT_PUBLIC_HUBSPOT_ID", label: "Portal ID (public)", type: "text" },
    ],
  },
  {
    title: "Sanity CMS",
    description: "Content management for docs and marketing",
    vars: [
      { key: "NEXT_PUBLIC_SANITY_PROJECT_ID", label: "Project ID", type: "text" },
      { key: "NEXT_PUBLIC_SANITY_DATASET", label: "Dataset", type: "text", placeholder: "production" },
      { key: "NEXT_PUBLIC_SANITY_API_TOKEN", label: "API Token", type: "secret" },
    ],
  },
]

/** Shape of a test instance for local dev */
export interface TestInstance {
  name: string
  url: string
  pat: string
}
