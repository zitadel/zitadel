---
title: Rate Limit Policy
custom_edit_url: null
---

Last updated on February 24, 2025

This policy is an annex to the [Terms of Service](../terms-of-service) and clarifies your obligations while using our Services, specifically how we will use rate limiting to enforce certain aspects of our [Acceptable Use Policy](acceptable-use-policy).

## Why do we rate limit

To ensure the availability of our Services and to avoid slow or failed requests by our Customers, due to overloads, we impose rate limits on certain API. These limits helps us guarantee the performance and availability of ZITADEL Cloud.

## How is the rate limit implemented

ZITADEL Clouds rate limit is built around a `IP` oriented model.
Please be aware that we also utilize a service for DDoS mitigation.
So if you simply change your `IP` address and run the same request again and again you might get blocked at some point.

If you are blocked you will receive a `http status 429`.

:::tip Implement exponential backoff
You should consider to implement [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) into your application to prevent a blocking loop.
:::

:::info Raising limits
We understand that there are certain scenarios where your users access ZITADEL from shared IP Addresses.
For example if you use a corporate proxy or Network Address Translation NAT.
Please [get in touch](https://zitadel.com/contact) with us to discuss your requirements, and we'll find a solution.
:::

## What rate limits do apply

For ZITADEL Cloud, we have dedicated rate limits for the user interfaces (login, register, console,...) and the APIs.

Rate limits are implemented with the following rules:

| Path                 | Description                            | Rate Limiting                        | One Minute Banning                    |
|----------------------|----------------------------------------|--------------------------------------|---------------------------------------|
| /ui/\*               | Global Login, Register and Reset Limit | 10 requests per second over a minute | 15 requests per second over 3 minutes |
| All other paths      | All gRPC-, REST and OAuth APIs         | 50 requests per second over a minute | 50 requests per second over 3 minutes |

## Load Testing

If you would like to conduct load testing of ZITADEL Cloud or a managed instance, you MUST request to do so with a minimum of 2 weeks notice before the test by contacting us at [support@zitadel.com](mailto:support@zitadel.com).
You MUST NOT conduct load testing without prior approval by us. Without prior approval and setup there is a high risk of being flagged by our DDoS solution as malicious traffic. This can have a severe impact on your service quality or result in termination of your agreement.
