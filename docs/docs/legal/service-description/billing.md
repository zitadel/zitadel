---
title: Pricing and billing of ZITADEL services
sidebar_label: Billing
custom_edit_url: null
--- 

Last updated on November 15, 2023

This annex of the [Framework Agreement](../terms-of-service) describes the pricing and billing of our Services.

## Pricing

You can find pricing information on our [website](https://zitadel.com/pricing).

### Enterprise pricing

Customer and ZITADEL may agree on an individual per-customer pricing via an Enterprise Agreement.

## Billing Metrics

### Monthly amount

Monthly amount means the available usage per measure for one billing period.
The amount is reset to zero with the start of a new billing period.

### Daily active user (DAU)

Daily Active Users (DAU) are counted as users who authenticate or refresh their token during the given day.
To calculate the monthly amount we take the sum of DAU over a given month.
Included are users that either login with local accounts or users that login with an external identity provider.
Service users that authenticate or access the management API are counted against Daily Active Users.

### Active external identity providers

To calculate the monthly amount we take the sum of activated external identity providers over all instances on each day and calculate the average over a given month, rounded up to the next integer.
Excluded are configured identity providers that are not activated.

### Action minutes

Action minutes mean execution time, rounded up to 1 second, of custom code execution via a customer defined Action.

### Management API requests

Management API requests means any request to the following API endpoints requiring a valid authorization header.
Excluded are requests with a server error, public endpoints, health endpoints, and endpoints to load UI assets.

Management endpoints:

- /zitadel.*
- /v2alpha*
- /v2beta*
- /admin*
- /management*
- /system*

### Admin users

Admin users means users within the customer portal that can manage a customer's account including billing, instances, analytics and additional services.

### Audit trail history (events)

Audit trail history (events) means past events that can be retrieved via API or GUI.
Typically all changes to any object in within ZITADEL are saved as events and can be used for audit trail and analytics purposes.
The number of past events that can be retrieved may be limited by your subscription.

### Access and runtime logs (logs)

Access and runtime logs (logs) means logs that are available about your instance.
Logs may contain information about success or failure reasons for API requests and Action executions, output from Actions, rate limit violations, and system health.
You might be able to retrieve logs only for a limited period of time based on your subscription.

### Custom domains

Custom domains mean domains that you can configure to reach your ZITADEL instance.
By default ZITADEL creates a custom domain for you when creating new instances.
Besides the included amount each additional custom domain is charged.

## Payment cycle

If not agreed otherwise, the payment frequency is monthly.
Your invoice will contain both pre-paid items for the current billing period and usage-based charges from the last billing period.
