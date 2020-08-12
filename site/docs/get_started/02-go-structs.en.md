---
title: Go structs
---

### Go structures

You can reference go struct tables from our go struct generator.
Provide a `doc_assets` folder with all generated files in it.
Make sure that the `.md` file consists of no other than the table itself and metadata which defines name and description of the struct


| Attribute                   | Description                                                                     | Default | Collection  |
| --------------------------- | ------------------------------------------------------------------------------- | ------- | ----------  |
| boomVersion                 | Version of BOOM which should be reconciled                                      |         |             |
| forceApply                  | Relative folder path where the currentstate is written to                       |         |             |
| currentStatePath            | Flag if --force should be used by apply of resources                            |         |             |
| preApply                    | Spec for the yaml-files applied before applications   |         |             |
| postApply                   | Spec for the yaml-files applied after applicatio      |         |             |
| prometheus-operator         | Spec for the Prometheus-Operator ,               |         |             |
| logging-operator            | Spec for the Banzaicloud Logging-Operator ,         |         |             |
| prometheus-node-exporter    | Spec for the Prometheus-Node-Exporter ,        |         |             |
| prometheus-systemd-exporter | Spec for the Prometheus-Systemd-Exporter , |         |             |
| grafana                     | Spec for the Grafana , [                          |         |             |
| ambassador                  | Spec for the Ambassador ,                                |         |             |
| kube-state-metrics          | Spec for the Kube-State-Metrics ,                 |         |             |
| argocd                      | Spec for the Argo-CD ,                              |         |             |
| prometheus                  | Spec for the Prometheus instance ,                    |         |             |
| loki                        | Spec for the Loki instance ,             |         |             |

#### References

To reference a table ...