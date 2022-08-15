---
title: ZITADEL Cloud Rate Limits
---

| Path                 | Description                            | Throttling                           | One Minute Banning        |
|----------------------|----------------------------------------|--------------------------------------|----------------------------------------|
| /ui/login*           | Global Login, Register and Reset Limit | 10 requests per second over a minute | 15 requests per sencond over 3 minutes |
| *various paths* [^1] | ZITADEL API                            | 4 requests per second over a minute  | 8 requests per second over 3 minutes   |

[^1] API paths reqular expression:
```regex
/system/v[0-9]+/.*|/auth/v[0-9]+/.|/admin/v[0-9]+/.|/management/v[0-9]+/.*|zitadel\.system\.v[0-9]+\.SystemService/.*|zitadel\.admin\.v[0-9]+\.AdminService/.*|zitadel\.auth\.v[0-9]+\.AuthService/.*|zitadel\.management\.v[0-9]+\.ManagementService/.*
```