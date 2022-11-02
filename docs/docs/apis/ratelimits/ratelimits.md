---
title: ZITADEL Cloud Rate Limits
---

Rate limits are implemented according to our
[rate limit policy](/docs/legal/rate-limit-policy.md) with the following rules:

| Path                     | Description                                                                                         | Throttling                           | One Minute Banning                     |
| ------------------------ | --------------------------------------------------------------------------------------------------- | ------------------------------------ | -------------------------------------- |
| /ui/login\*              | Global Login, Register and Reset Limit                                                              | 10 requests per second over a minute | 15 requests per sencond over 3 minutes |
| _Various API paths_ [^1] | All other gRPC- and REST APIs<br/> - Management API<br/>- Admin API<br/>- Auth API<br/>- System API | 4 requests per second over a minute  | 8 requests per second over 3 minutes   |

[^1] API paths:

<details>
    <summary>Open to see the reqular expression</summary>
    <pre>
/system/v[0-9]+/.*|/auth/v[0-9]+/.|/admin/v[0-9]+/.|/management/v[0-9]+/.*|zitadel\.system\.v[0-9]+\.SystemService/.*|zitadel\.admin\.v[0-9]+\.AdminService/.*|zitadel\.auth\.v[0-9]+\.AuthService/.*|zitadel\.management\.v[0-9]+\.ManagementService/.*
    </pre>
</details>
