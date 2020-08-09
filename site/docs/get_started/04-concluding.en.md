---
title: Concluding
---

#### Overview

The `manifest.json`, and `favicon.ico` are mandatory. If its all set and done, reensure all of this listed files are commited and pushed.

```bash
├ docs
  ├ manifest.json
  ├ static
  │ ├ favicon.ico
  └ index.svelte
```


## Structure 
 

| Attribute                   | Description                                                                     | Default | Collection  |
| --------------------------- | ------------------------------------------------------------------------------- | ------- | ----------  |
| boomVersion                 | Version of BOOM which should be reconciled                                      |         |             |
| forceApply                  | Relative folder path where the currentstate is written to                       |         |             |
| currentStatePath            | Flag if --force should be used by apply of resources                            |         |             |
| preApply                    | Spec for the yaml-files applied before applications , [here](PreApply.md)       |         |             |
| postApply                   | Spec for the yaml-files applied after applications , [here](PostApply.md)       |         |             |
| prometheus-operator         | Spec for the Prometheus-Operator , [here](PrometheusOperator.md)                |         |             |
| logging-operator            | Spec for the Banzaicloud Logging-Operator , [here](LoggingOperator.md)          |         |             |
| prometheus-node-exporter    | Spec for the Prometheus-Node-Exporter , [here](PrometheusNodeExporter.md)       |         |             |
| prometheus-systemd-exporter | Spec for the Prometheus-Systemd-Exporter , [here](PrometheusSystemdExporter.md) |         |             |
| grafana                     | Spec for the Grafana , [here](grafana/Grafana.md)                               |         |             |
| ambassador                  | Spec for the Ambassador , [here](Ambassador.md)                                 |         |             |
| kube-state-metrics          | Spec for the Kube-State-Metrics , [here](KubeStateMetrics.md)                   |         |             |
| argocd                      | Spec for the Argo-CD , [here](argocd/Argocd.md)                                 |         |             |
| prometheus                  | Spec for the Prometheus instance , [here](Prometheus.md)                        |         |             |
| loki                        | Spec for the Loki instance , [here](Loki.md)                                    |         |             |