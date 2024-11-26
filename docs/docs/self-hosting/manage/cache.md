---
title: Caches
sidebar_label: Caches
---

ZITADEL supports the use of a caches to speed up the lookup of frequently needed objects. As opposed to HTTP caches which might reside between ZITADEL and end-user applications, the cache build into ZITADEL uses active invalidation when an object gets updated. Another difference is that HTTP caches only cache the result of a complete request and the built-in cache stores objects needed for the internal business logic. For example, each request made to ZITADEL needs to retrieve and set [instance](/docs/concepts/structure/instance) information in middleware.

:::info
Caches is currently an [experimental beta](/docs/support/software-release-cycles-support#beta) feature.
:::

## Configuration

The `Caches` configuration entry defines *connectors* which can be used by several objects. It is possible to mix *connectors* with different objects based on operational needs.

```yaml
Caches:
  Connectors:
    SomeConnector:
      Enabled: true
      SomeOption: foo
  SomeObject:
    # Connector must be enabled above.
    # When connector is empty, this cache will be disabled.
    Connector: "SomeConnector"
    MaxAge: 1h
    LastUsage: 10m
    # Log enables cache-specific logging. Default to error log to stderr when omitted.
    Log:
      Level: error
```

For a full configuration reference, please see the [runtime configuration file](/docs/self-hosting/manage/configure#runtime-configuration-file) section's `defaults.yaml`.

## Connectors

ZITADEL supports a number of *connectors*. Connectors integrate a cache with a storage backend. Users can combine connectors with the type of object cache depending on their operational and performance requirements.
When no connector is specified for an object cache, then no caching is performed. This is the current default.

### Auto prune

Some connectors take an `AutoPrune` option. This is provided for caches which don't have built-in expiry and cleanup routines. The auto pruner is a routine launched by ZITADEL and scans and removes outdated objects in the cache. Pruning can take a cost as they typically involve some kind of scan. However, using a long interval can cause higher storage utilization.

```yaml
Caches:
  Connectors:
    Memory:
    Enabled: true
    # AutoPrune removes invalidated or expired object from the cache.
    AutoPrune:
      Interval: 1m
      TimeOut: 5s
```

### Redis cache

Redis is supported in simple mode. Cluster and Sentinel are not yet supported. There is also a circuit-breaker provided which prevents a single point of failure, should the single Redis instance become unavailable.

Benefits:

- Centralized cache with single source of truth
- Consistent invalidation
- Very fast when network latency is kept to a minimum
- Built-in object expiry, no pruner required

Drawbacks:

- Increased operational overhead: need to run a Redis instance as part of your infrastructure.
- When running multiple servers of ZITADEL in different regions, network roundtrip time might impact performance, neutralizing the benefit of a cache.

#### Circuit breaker

A [circuit breaker](https://learn.microsoft.com/en-us/previous-versions/msp-n-p/dn589784(v=pandp.10)?redirectedfrom=MSDN) is provided for the Redis connector, to prevent a single point of failure in the case persistent errors. When the circuit breaker opens, the cache is temporary disabled and ignored. ZITADEL will continue to operate using queries to the database.

```yaml
Caches:
  Connectors:
    Redis:
      Enabled: true
      Addr: localhost:6379
      # Many other options...
      CircuitBreaker:
        # Interval when the counters are reset to 0.
        # 0 interval never resets the counters until the CB is opened.
        Interval: 0
        # Amount of consecutive failures permitted
        MaxConsecutiveFailures: 5
        # The ratio of failed requests out of total requests
        MaxFailureRatio: 0.1
        # Timeout after opening of the CB, until the state is set to half-open.
        Timeout: 60s
        # The allowed amount of requests that are allowed to pass when the CB is half-open.
        MaxRetryRequests: 1
```

### PostgreSQL cache

PostgreSQL can be used to store objects in unlogged tables. [Unlogged tables](https://www.postgresql.org/docs/current/sql-createtable.html#SQL-CREATETABLE-UNLOGGED) do not write to the WAL log and are therefore faster than regular tables. If the PostgreSQL server crashes, the data from those tables are lost. ZITADEL always creates the cache schema in the `zitadel` database during [setup](./updating_scaling#the-setup-phase). This connector requires a [pruner](#auto-prune) routine.

Benefits:

- Centralized cache with single source of truth
- No operational overhead. Reuses the query connection pool and the existing `zitadel` database.
- Consistent invalidation
- Faster than regular queries which often contain `JOIN` clauses.

Drawbacks:

- Slowest of the available caching options
- Might put additional strain on the database server, limiting horizontal scalability
- CockroachDB does not support unlogged tables. When this connector is enabled against CockroachDB, it does work but little to no performance benefit is to be expected.

### Local memory cache

ZITADEL is capable of caching object in local application memory, using hash-maps. Each ZITADEL server manages its own copy of the cache. This connector requires a [pruner](#auto-prune) routine.

Benefits:

- Fastest of the available caching options
- No operational overhead

Drawbacks:

- Inconsistent invalidation. An object validated in one ZITADEL server will not get invalidated in other servers.
- There's no single source of truth. Different servers may operate on a different version of an object
- Data is duplicated in each server, consuming more total memory inside a deployment.
 
The drawbacks restricts its usefulness in distributed deployments. However simple installations running a single server can benefit greatly from this type of cache. For example test, development or home deployments.
If inconsistency is acceptable for short periods of time, one can choose to use this type of cache in distributed deployments with short max age configuration. 

**For example**: A ZITADEL deployment with 2 servers is serving 1000 req/sec total. The installation only has one instance[^1]. There is only a small amount of data cached (a few kB) so duplication is not a problem in this case. It is acceptable for [instance level setting](/docs/guides/manage/console/default-settings) to be out-dated for a short amount of time. When the memory cache is enabled for the instance objects, with a max age of 1 second, the instance only needs to be obtained from the database 2 times per second (once for each server). Saving 998 of redundant queries. Once an instance level setting is changed, it takes up to 1 second for all the servers to get the new state.

## Objects

The following section describes the type of objects ZITADEL can currently cache. Objects are actively invalidated at the cache backend when one of their properties is changed. Each object cache defines:

- `Connector`: Selects the used [connector](#connectors) back-end. Must be activated first.
- `MaxAge`: the amount of time that an object is considered valid. When this age is passed the object is ignored (cache miss) and possibly cleaned up by the [pruner](#auto-prune) or other built-in garbage collection.
- `LastUsage`: defines usage based lifetime. Each time an object is used, its usage timestamp is updated. Popular objects remain cached, while unused objects are cleaned up. This option can be used to indirectly limit the size of the cache.
- `Log`: allows specific log settings for the cache. This can be used to debug a certain cache without having to change the global log level.

```yaml
Caches:
  SomeObject:
    # Connector must be enabled above.
    # When connector is empty, this cache will be disabled.
    Connector: ""
    MaxAge: 1h
    LastUsage: 10m
    # Log enables cache-specific logging. Default to error log to stderr when omitted.
    Log:
      Level: error
      AddSource: true
      Formatter:
        Format: text
```

### Instance

All HTTP and gRPC requests sent to ZITADEL receive an instance context. The instance is usually resolved by the domain from the request. In some cases, like the [system service](/docs/apis/resources/system/system-service), the instance can be resolved by its ID. An instance object contains many of the [default settings](/docs/guides/manage/console/default-settings):

- Instance [features](/docs/guides/manage/console/default-settings#features)
- Instance domains: generated and [custom](/docs/guides/manage/cloud/instances#add-custom-domain)
- [Trusted domains](/docs/apis/resources/admin/admin-service-add-instance-trusted-domain)
- Security settings ([IFrame policy](/docs/guides/solution-scenarios/configurations#embedding-zitadel-in-an-iframe))
- Limits[^2]
- [Allowed languages](/docs/guides/manage/console/default-settings#languages)

These settings typically change infrequently in production. ***Every*** request made to ZITADEL needs to query for the instance. This is a typical case of set once, get many times where a cache can provide a significant optimization.

### Milestones

Milestones are used to track the administrator's progress in setting up their instance. Milestones are used to render *your next steps* in the [console](/docs/guides/manage/console/overview) landing page.
Milestones are reached upon the first time a certain action is performed. For example the first application created or the first human login. In order to push a "reached" event only once, ZITADEL must keep track of the current state of milestones by an eventstore query every time an eligible action is performed. This can cause an unwanted overhead on production servers, therefore they are cached.

As an extra optimization, once all milestones are reached by the instance, an in-memory flag is set and the milestone state is never queried again from the database nor cache.
For single instance setups which fulfilled all milestone (*your next steps* in console) it is not needed to enable this cache. We mainly use it for ZITADEL cloud where there are many instances with *incomplete* milestones.

### Organization

Most resources like users, project and applications are part of an [organization](/docs/concepts/structure/organizations). Therefore many parts of the ZITADEL logic search for an organization by ID or by their primary domain.
Organization objects are quite small and receive infrequent updates after they are created:

- Change of organization name
- Deactivation / Reactivation
- Change of primary domain
- Removal

## Examples

Currently caches are in beta and disabled by default. However, if you want to give caching a try, the following sections contains some suggested configurations for different setups.

The following configuration is recommended for single instance setups with a single ZITADEL server:

```yaml
Caches:
  Memory:
    Enabled: true
  Instance:
    Connector: "memory"
    MaxAge: 1h
  Organization:
    Connector: "memory"
    MaxAge: 1h
```

The following configuration is recommended for single instance setups with high traffic on multiple servers, where Redis is not available:

```yaml
Caches:
  Memory:
    Enabled: true
  Postgres:
    Enabled: true
  Instance:
    Connector: "memory"
    MaxAge: 1s
  Milestones:
    Connector: "postgres"
    MaxAge: 1h
    LastUsage: 10m
  Organization:
    Connector: "memory"
    MaxAge: 1s
```

When running many instances on multiple servers:

```yaml
Caches:
  Connectors:
    Redis:
      Enabled: true
      # Other connection options
      
  Instance:
    Connector: "redis"
    MaxAge: 1h
    LastUsage: 10m
  Milestones:
    Connector: "redis"
    MaxAge: 1h
    LastUsage: 10m
  Organization:
    Connector: "redis"
    MaxAge: 1h
    LastUsage: 10m
```
----

[^1]: Many deployments of ZITADEL have only one or few [instances](/docs/concepts/structure/instance). Multiple instances are mostly used for ZITADEL cloud, where each customer gets at least one instance.

[^2]: Limits are imposed by the system API, usually when customers exceed their subscription in ZITADEL cloud.