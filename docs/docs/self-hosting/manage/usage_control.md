---
title: Usage Control
sidebar_label: Usage Control
---

If you have a self-hosted ZITADEL environment, you can limit the usage of your [instances](/concepts/structure/instance).
For example, if you provide your customers [their own virtual instances](/concepts/structure/instance#multiple-virtual-instances) with access on their own domains, you can design a pricing model based on the usage of their instances.
The usage control features are currently limited to the instance level only.

## Limit Audit Trails

You can restrict the maximum age of events returned by the following APIs:

- [Events Search](/apis/resources/admin/admin-service-list-events), See also the [Event API guide](guides/integrate/event-api)
- [My User History](/apis/resources/auth/auth-service-list-my-user-changes)
- [A Users History](/apis/resources/mgmt/management-service-list-user-changes)
- [An Applications History](/apis/resources/mgmt/management-service-list-app-changes)
- [An Organizations History](/apis/resources/mgmt/management-service-list-org-changes)
- [A Projects History](/apis/resources/mgmt/management-service-list-project-changes)
- [A Project Grants History](/apis/resources/mgmt/management-service-list-project-grant-changes)

You can set a global default limit as well as a default limit [for new virtual instances](/concepts/structure/instance#multiple-virtual-instances) in the ZITADEL configuration.
The following snippets shows the defaults:

```yaml
# AuditLogRetention limits the number of events that can be queried via the events API by their age.
# A value of "0s" means that all events are available.
# If an audit log retention is set using an instance limit, it will overwrite the system default.
AuditLogRetention: 0s # ZITADEL_AUDITLOGRETENTION
DefaultInstance:
  Limits:
    # AuditLogRetention limits the number of events that can be queried via the events API by their age.
    # A value of "0s" means that all events are available.
    # If this value is set, it overwrites the system default unless it is not reset via the admin API.
    AuditLogRetention: # ZITADEL_DEFAULTINSTANCE_LIMITS_AUDITLOGRETENTION
```

You can also set a limit for [a specific virtual instance](/concepts/structure/instance#multiple-virtual-instances) using the [system API](/category/apis/resources/system/limits).

## Quotas

Quotas enables you to limit usage and/or register webhooks that trigger on configurable usage levels for certain units.
For example, you might want to report usage to an external billing tool and notify users when 80 percent of a quota is exhausted.

ZITADEL supports limiting authenticated requests and action run seconds with quotas.

For using the quotas feature you have to activate it in your ZITADEL configurations *Quotas* section.
The following snippets shows the defaults:

```yaml
Quotas:
  Access:
    # If enabled, authenticated requests are counted and potentially limited depending on the configured quota of the instance
    Enabled: false # ZITADEL_QUOTAS_ACCESS_ENABLED
    Debounce:
      MinFrequency: 0s # ZITADEL_QUOTAS_ACCESS_DEBOUNCE_MINFREQUENCY
      MaxBulkSize: 0 # ZITADEL_QUOTAS_ACCESS_DEBOUNCE_MAXBULKSIZE
    ExhaustedCookieKey: "zitadel.quota.exhausted" # ZITADEL_QUOTAS_ACCESS_EXHAUSTEDCOOKIEKEY
    ExhaustedCookieMaxAge: "300s" # ZITADEL_QUOTAS_ACCESS_EXHAUSTEDCOOKIEMAXAGE
  Execution:
    # If enabled, all action executions are counted and potentially limited depending on the configured quota of the instance
    Enabled: false # ZITADEL_QUOTAS_EXECUTION_DATABASE_ENABLED
    Debounce:
      MinFrequency: 0s # ZITADEL_QUOTAS_EXECUTION_DEBOUNCE_MINFREQUENCY
      MaxBulkSize: 0 # ZITADEL_QUOTAS_EXECUTION_DEBOUNCE_MAXBULKSIZE
```

Once you have activated the quotas feature, you can configure quotas [for your virtual instances](/concepts/structure/instance#multiple-virtual-instances) using the [system API](/category/apis/resources/system/quotas) or the *DefaultInstances.Quotas* section.
The following snippets shows the defaults:

```yaml
DefaultInstance:
  Quotas:
    # Items take a slice of quota configurations, whereas, for each unit type and instance, one or zero quotas may exist.
    # The following unit types are supported

    # "requests.all.authenticated"
    # The sum of all requests to the ZITADEL API with an authorization header,
    # excluding the following exceptions
    # - Calls to the System API
    # - Calls that cause internal server errors
    # - Failed authorizations
    # - Requests after the quota already exceeded

    # "actions.all.runs.seconds"
    # The sum of all actions run durations in seconds
    Items:
#      - Unit: "requests.all.authenticated"
#        # From defines the starting time from which the current quota period is calculated.
#        # This is relevant for querying the current usage.
#        From: "2023-01-01T00:00:00Z"
#        # ResetInterval defines the quota periods duration
#        ResetInterval: 720h # 30 days
#        # Amount defines the number of units for this quota
#        Amount: 25000
#        # Limit defines whether ZITADEL should block further usage when the configured amount is used
#        Limit: false
#        # Notifications are emitted by ZITADEL when certain quota percentages are reached
#        Notifications:
#            # Percent defines the relative amount of used units, after which a notification should be emitted.
#          - Percent: 100
#            # Repeat defines, whether a notification should be emitted each time when a multitude of the configured Percent is used.
#            Repeat: true
#            # CallURL is called when a relative amount of the quota is used.
#            CallURL: "https://httpbin.org/post"
```

### Exhausted Authenticated Requests

If a quota is configured to limit requests and the quotas amount is exhausted, all further requests are blocked except requests to the System API.
Also, a cookie is set, to make it easier to block further traffic before it reaches your ZITADEL runtime.

### Exhausted Action Run Seconds

If a quota is configured to limit action run seconds and the quotas amount is exhausted, all further actions will fail immediately with a context timeout exceeded error.
The action that runs into the limit also fails with the context timeout exceeded error.

