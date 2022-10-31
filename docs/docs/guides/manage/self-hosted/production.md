---
title: Production Checklist
---

As soon as you successfully deployed ZITADEL as a proof of concept using one of our [deployment guides](/docs/guides/deploy/overview),
you are ready to configure ZITADEL for production usage.

## High Availability

We recommend running ZITADEL higly available using an orchestrator that schedules ZITADEL on multiple nodes, like [Kubernetes](/docs/guides/deploy/kubernetes).

## Binary Configuration

Read [on the configure page](/docs/guides/manage/self-hosted/configure) about the available options you have to configure the ZITADEL binary.

## Networking

- For enabling access via **HTTPS** and restricting access to **HTTPS**, head to [the description of your TLS options](/docs/guides/manage/self-hosted/tls_modes).
- If you want to front ZITADEL with a reverse proxy, web application firewall or content delivery network, make sure to support **[HTTP/2](/docs/guides/manage/self-hosted/http2)**.
- You can also refer to some **[example reverse proxy configurations](/docs/guides/manage/self-hosted/reverseproxy/reverse_proxy)**.

## Monitoring

By default, **metrics** are exposed at /debug/metrics in OpenTelemetry (otel) format.
They can easily be consumed by an [otel collector](https://opentelemetry.io/docs/collector/) and be sent to or scraped by various backends,
for example [Prometheus](https://github.com/open-telemetry/opentelemetry-collector-contrib/tree/main/exporter/prometheusexporter).

Also, you can enable **tracing** in the ZITADEL configuration.

```yaml
Tracing:
  # Choose one in "otel", "google", "log" and "none"
  Type: google
  Fraction: 1
  MetricPrefix: zitadel
```

## Database

Depending on your environment, you maybe would want to tweak some settings about how ZITADEL interacts with the database in the Database section of your ZITADEL configuration. Read more about your [database configuration options](/docs/guides/manage/self-hosted/database).

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

You also might want to configure how projections are computed. These are the default values:

```yaml
Projections:
  RequeueEvery: 60s
  RetryFailedAfter: 1s
  MaxFailureCount: 5
  ConcurrentInstances: 1
  BulkLimit: 200
  MaxIterators: 1
  Customizations:
    projects:
      BulkLimit: 2000
```

## Data Initialization

- You can configure instance defaults in the DefaultInstance section.
  If you plan to eventually create [multiple virtual instances](/docs/concepts/structure/instance#multiple-virtual-instances), these defaults take effect, too.
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

- If you don't want to use the DefaultInstance configuration for the first instance that ZITADEL automatically creates for you during the [startup phase](/docs/guides/manage/self-hosted/configure#database-initialization), you can provide a FirstInstance YAML section using the --steps argument.
- Learn how to configure ZITADEL via the [Console user interface](/docs/guides/manage/console/overview).
- Probably, you also want [apply your custom branding](/docs/guides/manage/customize/branding), [hook into certain events](/docs/guides/manage/customize/behavior), [customize texts](/docs/guides/manage/customize/texts) or [add metadata to your users](/docs/guides/manage/customize/user-metadata)
- If you want to automatically setup ZITADEL resources, you can use the [ZITADEL Terraform Provider](/docs/guides/manage/terraform/basics)
