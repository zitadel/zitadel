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

## What rate limits do apply

Please check to corresponding page.
