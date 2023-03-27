---
title: Production Setup
---

As soon as you successfully deployed ZITADEL as a proof of concept using one of our [deployment guides](/docs/self-hosting/deploy/overview),
you are ready to configure ZITADEL for production usage.


## TL;DR
We created a [Production Checklist](./productionchecklist.md) as a technical step by step guide.

## High Availability

We recommend running ZITADEL highly available using an orchestrator that schedules ZITADEL on multiple servers,
like [Kubernetes](/docs/self-hosting/deploy/kubernetes).
For keeping startup times fast when scaling ZITADEL,
you should also consider using separate jobs with `zitadel init` and `zitadel setup`,
so your workload containers just have to execute `zitadel start`.
Read more about separating the init and setup phases on the [Updating and Scaling page](/docs/self-hosting/manage/updating_scaling).

## Configuration

Read [on the configure page](/docs/self-hosting/manage/configure) about the available options you have to configure ZITADEL.

## Networking

- To make ZITADEL available at the domain of your choice, [you need to configure the ExternalDomain property](/docs/self-hosting/manage/custom-domain).
- To enable and restrict access to **HTTPS**, head over to [the description of your TLS options](/docs/self-hosting/manage/tls_modes).
- If you want to front ZITADEL with a reverse proxy, web application firewall or content delivery network, make sure to support **[HTTP/2](/docs/self-hosting/manage/http2)**.
- You can also refer to some **[example reverse proxy configurations](/docs/self-hosting/manage/reverseproxy/reverse_proxy)**.
- The ZITADEL Console web GUI uses many gRPC-Web stubs. This results in a fairly big JavaScript bundle. You might want to compress it using [Gzip](https://www.gnu.org/software/gzip/) or [Brotli](https://github.com/google/brotli).
- Serving and caching the assets using a content delivery network could improve network latencies and shield your ZITADEL runtime.

## Monitoring

By default, [**metrics**](/apis/observability/metrics) are exposed at /debug/metrics in OpenTelemetry (otel) format.

Also, you can enable **tracing** in the ZITADEL configuration.

```yaml
Tracing:
  # Choose one in "otel", "google", "log" and "none"
  Type: google
  Fraction: 1
  MetricPrefix: zitadel
```

## Database

### Prefer CockroachDB

ZITADEL supports [CockroachDB](https://www.cockroachlabs.com/) and [PostgreSQL](https://www.postgresql.org/).
We highly recommend using CockroachDB,
as horizontal scaling is much easier than with PostgreSQL.
Also, if you are concerned about multi-regional data locality,
[the way to go is with CockroachDB](https://www.cockroachlabs.com/docs/stable/multiregion-overview.html).

### Configure ZITADEL

Depending on your environment, you maybe would want to tweak some settings about how ZITADEL interacts with the database in the database section of your ZITADEL configuration. Read more about your [database configuration options](/docs/self-hosting/manage/database).

```yaml
Database:
  cockroach:
    Host: localhost
    Port: 26257
    Database: zitadel
    //highlight-start
    MaxOpenConns: 20
    MaxConnLifetime: 30m
    MaxConnIdleTime: 30m
    //highlight-end
    Options: ""
```

You also might want to configure how [projections](/concepts/eventstore/implementation#projections) are computed. These are the default values:

```yaml
# The Projections section defines the behaviour for the scheduled and synchronous events projections.
Projections:
  # Time interval between scheduled projections
  RequeueEvery: 60s
  # Time between retried database statements resulting from projected events
  RetryFailedAfter: 1s
  # Retried execution number of database statements resulting from projected events
  MaxFailureCount: 5
  # Number of concurrent projection routines
  ConcurrentInstances: 1
  # Limit of returned events per query
  BulkLimit: 200
  # If HandleInactiveInstances this is false, only instances are projected,
  # for which at least a projection relevant event exists withing the timeframe
  # from twice the RequeueEvery time in the past until the projections current time
  HandleInactiveInstances: false
  # In the Customizations section, all settings from above can be overwritten for each specific projection
  Customizations:
    Projects:
      BulkLimit: 2000
    # The Notifications projection is used for sending emails and SMS to users
    Notifications:
      # As notification projections don't result in database statements, retries don't have an effect
      MaxFailureCount: 0
    # The NotificationsQuotas projection is used for calling quota webhooks
    NotificationsQuotas:
      # Delivery guarantee requirements are probably higher for quota webhooks
      HandleInactiveInstances: true
      # As quota notification projections don't result in database statements, retries don't have an effect
      MaxFailureCount: 0
```

### Manage your Data

When designing your backup strategy,
it is worth knowing that
[ZITADEL is event sourced](/docs/concepts/eventstore/overview).
That means, ZITADEL itself is able to recompute its
whole state from the records in the table eventstore.events.
The timestamp of your last record in the events table
defines up to which point in time ZITADEL can restore its state.

The ZITADEL binary itself is stateless,
so there is no need for a special backup job.

Generally, for maintaining your database management system in production,
please refer to the corresponding docs
[for CockroachDB](https://www.cockroachlabs.com/docs/stable/recommended-production-settings.html)
or [for PostgreSQL](https://www.postgresql.org/docs/current/admin.html).


## Data Initialization

- You can configure instance defaults in the DefaultInstance section.
  If you plan to eventually create [multiple virtual instances](/concepts/structure/instance#multiple-virtual-instances), these defaults take effect.
  Also, these configurations apply to the first instance, that ZITADEL automatically creates for you.
  Especially the following properties are of special interest for your production setup.

```yaml
DefaultInstance:
  OIDCSettings:
    AccessTokenLifetime: 12h
    IdTokenLifetime: 12h
    RefreshTokenIdleExpiration: 720h #30d
    RefreshTokenExpiration: 2160h #90d
  # this configuration sets the default email configuration
  SMTPConfiguration:
    # configuration of the host
    SMTP:
      #for example smtp.mailtrap.io:2525
      Host:
      User:
      Password:
    TLS:
    # if the host of the sender is different from ExternalDomain set DefaultInstance.DomainPolicy.SMTPSenderAddressMatchesInstanceDomain to false
    From:
    FromName:
```

- If you don't want to use the DefaultInstance configuration for the first instance that ZITADEL automatically creates for you during the [setup phase](/self-hosting/manage/configure#database-initialization), you can provide a FirstInstance YAML section using the --steps argument.
- Learn how to configure ZITADEL via the [Console user interface](/guides/manage/console/overview).
- Probably, you also want to [apply your custom branding](/guides/manage/customize/branding), [hook into certain events](/guides/manage/customize/behavior), [customize texts](/guides/manage/customize/texts) or [add metadata to your users](/guides/manage/customize/user-metadata).
- If you want to automatically create ZITADEL resources, you can use the [ZITADEL Terraform Provider](/guides/manage/terraform/basics).

## Quotas

If you host ZITADEL as a service,
you might want to [limit usage and/or execute tasks on certain usage units and levels](/self-hosting/manage/quotas).
