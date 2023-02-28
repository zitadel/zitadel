---
title: ZITADEL Cloud Rate Limits
---

Rate limits are implemented according to our [rate limit policy](/legal/rate-limit-policy.md) with the following rules:

| Path                     | Description                            | Rate Limiting                        | One Minute Banning                     |
|--------------------------|----------------------------------------|--------------------------------------|----------------------------------------|
| /ui/login*               | Global Login, Register and Reset Limit | 10 requests per second over a minute | 15 requests per sencond over 3 minutes |
| All other paths | All gRPC- and REST APIs as well as the ZITADEL Customer Portal | 10 requests per second over a minute       | 10 requests per second over 3 minutes   |
