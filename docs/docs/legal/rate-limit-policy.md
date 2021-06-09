---
title: Rate Limit Policy
custom_edit_url: null
---
## Introduction

This policy is an annex to the [Terms of Service](terms-of-service) and clarifies your obligations while using our Services, specifically how we will use rate limiting to enforce certain aspects of our [Acceptable Use Policy](acceptable-use-policy).

## Why do we rate limit

To ensure the availability of our Services and to avoid slow or failed requests by our Customers, due to overloads, we impose rate limits on certain API. These limits helps us guarantee the performance and availability of ZITADEL Cloud.

## How is the rate limit implemented

ZITADEL Clouds rate limit is built around a `IP` oriented model. Please be aware that we also utilize a service for DDoS mitigation.
So if you simply change your `IP` address and run the same request again and again you might be get blocked at some point.

If you are blocked you will receive a `http status 429`.

:::tip
You should consider to implement [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) into your application to prevent a blocking loop.
:::

## What rate limits do apply

### Login, Register, Reset Limits

For the rate limits of the Login, Register and Reset features please visit [Login Rate Limits](/docs/apis/ratelimits/accounts)

### API Rate Limits

For our API rate limits please check the [API Endpoint Rate Limits](/docs/apis/ratelimits/api)

## Load Testing

If you would like to conduct load testing of ZITADEL Cloud or a managed instance, you MUST request to do so with a minimum of 2 weeks notice before the test by contacting us at support@zitadel.ch.  
You MUST NOT conduct load testing without prior approval by us. Without prior approval and setup there is a high risk of being flagged by our DDoS solution als malicious traffic. This can have a severe impact on your service quality or result in termination of your agreement.
