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

### Instances

This limit defines how many separate instances you can create. 
An instance is your own private space to manage all your organizations, users, and applications. 
Using multiple instances is the ideal way to create completely separate environments (e.g., for Development and Production).

### Organizations

In Zitadel, each "Organization" represents a distinct tenant, perfect for your B2B customers or for creating separate environments.
This limit defines the total number of organizations you can create to manage your users and resources in a multi-tenancy setup.

### Users per Organization

This is the total number of users that can be added to any single organization.
The count includes both human and machine accounts.

### Administrator Users

This limit applies to users assigned an administrative role at the instance level, giving them broad oversight and management capabilities.
It counts all instance-wide administrators, from users with full, unrestricted access (IAM_OWNER) to those with view-only administrative permissions (IAM_OWNER_VIEWER), or any other role.
You are billed for the total number of administrator users on your account at the time the invoice is created.

### SCIM Connections

Automate user provisioning by connecting your identity sources to Zitadel. 
This limit defines how many different SCIM clients you can configure to send user data to your instance, enabling automatic creation, updates, and deactivation of users.

### Identity Providers (Social Login, Enterprise Login, etc.)

We count all configured external identity providers, with the one exception of LDAP connections. 
This includes social logins (like Google, GitHub, and Apple) as well as generic providers using protocols like OpenID Connect, OAuth, and SAML.

### LDAP Identity Providers

We count each configured LDAP connection separately from other identity providers. 
Your bill is based on the total number of these LDAP providers connected to your account at the time the invoice is created.

### Actions

Your subscription plan determines the number of "Actions" you can configure. For billing, we count the total number of actions you have configured at the time your invoice is created.

### Policies

Policies are rules that enforce security behaviors, like password complexity or requiring multi-factor authentication (MFA).

Your bill is not based on the number of policies you create, but on how many times they are actively enforced/activated.
We calculate your usage by summing up every policy that is enabled within each of your organizations at the end of the billing cycle.

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
