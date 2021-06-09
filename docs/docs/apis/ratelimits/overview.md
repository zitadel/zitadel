---
title: Rate Limit Policy
---

:::caution
This is subject to change
:::

## Why do we rate limit

To protect our users from unavalabilty due to overloads we impose rate limits on certain API.
This helps us garuantee the perfomance and availabilty of ZITADEL Cloud to all it's users.

## How is the rate limit implemented

ZITADEL Clouds rate limit is built around a `IP` oriented model. Please be aware that we also utilize a service for DDoS mitigation.
So if you simply change your `IP` address and run the same request again and again you might be get blocked at some point.

If you are blocked you will recieve a `http status 429`.

:::tip
You should consider to implement [exponential backoff](https://en.wikipedia.org/wiki/Exponential_backoff) into your application to prevent a blocking loop.
:::

## What rate limits do apply

### Login, Register, Reset Limits

For the rate limits of the Login, Register and Reset features please visit [Login Rate Limits](accounts)

### API Rate Limits

For our API rate limits please check the [API Endpoint Rate Limits](api)

## Load Testing

If you would like to conduct load testing of ZITADEL Cloud or a managed instance, you MUST request to do so with a minimum of 2 weeks notice before the test by contacting us at support@zitadel.ch.
You MUST NOT conduct load testing without prior approval by us. Without prior approval and setup there is a high risk of beeing flaged by our (D)DOS solution als malicous traffic. This can have a severy impact on your service qualtiy.
