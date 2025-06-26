---
title: Pricing and billing of ZITADEL services
sidebar_label: Billing
custom_edit_url: null
--- 

Last updated on June 30, 2025

This annex of the [Framework Agreement](../terms-of-service) describes the pricing and billing of our Services.

## Pricing

You can find pricing information on our [website](https://zitadel.com/pricing).

### Enterprise pricing

Customer and ZITADEL may agree on an individual per-customer pricing via an Enterprise Agreement.

## Billing Metrics

### Monthly amount

Monthly amount means the available usage per measure for one billing period. 
This amount is based on the highest number recorded during that billing period. 
The amount is reset to zero with the start of a new billing period.

### Custom domains

Custom domains mean domains that you can configure to reach your ZITADEL instance.
By default ZITADEL creates a custom domain for you when creating new instances.
Besides the included amount each additional custom domain is charged.

### API Requests

API requests means any request to any API endpoints requiring a valid authorization header.
Excluded are requests with a server error, public endpoints, health endpoints, and endpoints to load UI assets.

### Log Drain (Instance Logs)

Access and runtime logs (logs) means logs that are available about your instance.
Logs may contain information about success or failure reasons for API requests and Action executions, output from Actions, rate limit violations, and system health.
The total volume of logs you can retrieve is determined by the GB allowance included in your subscription.

### B2C Organizations

We categorize an organization as Business-to-Consumer (B2C) when it contains more than 10,000 users.

### M2M Organizations

A Machine-to-Machine (M2M) Organization is one that includes more than five machine users.

### SCIM Connections

Automate user provisioning by connecting your identity sources to Zitadel.
This metric tracks the number of unique SCIM clients you've configured to send user data to your instance, enabling automatic user creation, updates, and deactivation.

### Identity Providers (Social Login, Enterprise Login, etc.)

We count all configured external identity providers, with the one exception of LDAP connections. This includes social logins (like Google, GitHub, and Apple) as well as generic providers using protocols like OpenID Connect, OAuth, and SAML.

For billing purposes, your invoice reflects the peak number of Identity providers connected to your account throughout the month.

### LDAP Identity Providers

We count each configured LDAP connection separately from other identity providers.
Your bill is based on the highest number of LDAP providers connected to your account at any point during the billing cycle.

### Actions

The number of 'Actions' you can configure is set by your subscription plan. We bill based on the maximum number of actions you had configured at any time within the billing period.

### Policies

Policies are rules that enforce security behaviors, like password complexity or requiring multi-factor authentication (MFA).

Your bill is not based on the number of policies you create, but on how many times they are actively enforced/activated.
Your usage is determined by the highest count of enabled policies across all your organizations recorded at any point during the billing cycle.
List of policies:
- Force MFA / Factors
- Password Complexity
- Password Expiry
- User Lockout

### Audit trail history / Audit Logs (events)

Audit trail history (events) means past events that can be retrieved via API or GUI.
Typically all changes to any object in within ZITADEL are saved as events and can be used for audit trail and analytics purposes.
The number of past events that can be retrieved may be limited by your subscription.

## Payment cycle

If not agreed otherwise, the payment frequency is monthly.
Your invoice will contain both pre-paid items for the current billing period and usage-based charges from the last billing period.
