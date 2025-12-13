---
title: ZITADEL Ready and Health Endpoints
---




ZITADEL exposes a `Ready`- and `Healthy` endpoint to allow external systems like load balancers, orchestration systems, uptime probes and others to check the status.

## Ready

The `Ready` endpoint is located on the path `/debug/ready` and allows systems to probe if a ZITADEL process is ready to serve and accept traffic.
This endpoint is useful for operations like [zero downtime upgrade](/docs/learn/architecture/solution#zero-downtime-updates) since it allows systems like Kubernetes to verify that ZITADEL is working on something (e.g. database schema migration) but is not yet ready to accept traffic.

<Callout type="info" title="In Kubernetes this is called the `readinessProbe`.">

</Callout>

## Healthy

The `Health` endpoint is located on the path `/debug/healthz` and allows systems to probe if a ZITADEL process is still alive.
This helps system like kubernetes or a load balancer to observe if the process is still alive to accept traffic.

<Callout type="info" title="In Kubernetes this is called the `livenessProbe`.">

</Callout>
