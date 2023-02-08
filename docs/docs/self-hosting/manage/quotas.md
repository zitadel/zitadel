---
title: Usage Quotas
---

Quotas is an enterprise feature that is relevant if you want to host ZITADEL as a service.
It enables you to limit usage and/or register webhooks that trigger on configurable usage levels for certain units.
For example, you might want to report usage to an external billing tool and notify users when 80 percent of a quota is exhausted.
Quotas are currently supported [for the instance level only](/concepts/structure/instance).
Please refer to the [system API docs](/apis/proto/system#addquota) for detailed explanations about how to use the quotas feature.

ZITADEL supports limiting authenticated requests and action run seconds

## Authenticated Requests

For using the quotas feature for authenticated requests you have to enable the database logstore for access logs in your ZITADEL configurations LogStore section:

```ts reference
https://github.com/zitadel/zitadel/blob/add-quotas/cmd/defaults.yaml#L323-L336
```

If a quota is configured to limit requests and the quotas amount is exhausted, all further requests are blocked except requests to the System API.
Also, a cookie is set, to make it easier to block further traffic before it reaches your ZITADEL runtime.


## Action Run Seconds

For using the quotas feature for action run seconds you have to enable the database logstore for execution logs in your ZITADEL configurations LogStore section:

```ts reference
https://github.com/zitadel/zitadel/blob/add-quotas/cmd/defaults.yaml#L345-L357
```

If a quota is configured to limit action run seconds and the quotas amount is exhausted, all further actions will fail immediately with a context timeout exceeded error.
The action that runs into the limit also fails with the context timeout exceeded error.

