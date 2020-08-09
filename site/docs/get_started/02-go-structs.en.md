---
title: Go structs
---

### Go structures

You can reference go struct tables from our go struct generator.
Provide a `doc_assets` folder with all generated files in it.
Make sure that the `.md` file consists of no other than the table itself and metadata which defines name and description of the struct

Take a look at the following example

```md
    --- 
    title: ToolsetSpec
    description: BOOM reconciles itself if a boomVersion is defined, if no boomVersion is defined there is no reconciling.
    ---

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
```

which produces the following table:


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

#### References

To reference a table ...