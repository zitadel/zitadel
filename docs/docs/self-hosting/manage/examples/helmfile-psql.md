---
title: Helmfile + PostgreSQL
---

This example uses [`helmfile`](https://helmfile.readthedocs.io/) with [`helm-secrets`](https://github.com/jkroepke/helm-secrets) in combination with our [Kubernetes Deploy Guides](/docs/self-hosting/deploy/kubernetes) and the [PostgreSQL Database](/docs/self-hosting/manage/database#postgres) to provision a production-eligible ZITADEL instance. You can make this example work without `helmfile` by mirroring the `zitadel` configration to your declarative setup.

The secrets are safely retrieved from Azure KeyVault in this example but you can also use other secret stores supported by `helm-secrets` or provide the secrets directly within the `values`.

:::caution
Even though this examples strives to be complete, make sure to read our [Production Guide](/docs/self-hosting/manage/production) before you decide to use it as reference.
:::

```helmfile
environments:
  default:
    values:
      - zitadel:
        hostName: id.your-domain.tld
        mainKey: ref+azurekeyvault://vault-name/zitadel-main-key
        database:
          host: ref+azurekeyvault://vault-name/host
          port: ref+azurekeyvault://vault-name/port
          username: ref+azurekeyvault://vault-name/username
          password: ref+azurekeyvault://vault-name/password
          zitadelPassword: ref+azurekeyvault://vault-name/zitadel-zitadel-password

releases:
  - name: zitadel
    chart: zitadel/zitadel
    namespace: iam
    version: 5.0.0
    createNamespace: true
    wait: false
    values:
      - zitadel:
          # The certificate is public in this case so you can store it in VCS; if not you should source it from secrets as well
          dbSslRootCrt: |
            -----BEGIN CERTIFICATE-----
            MIIDrzCCAp[OMITTED FOR BREVITY]Mbp1ZWVbd4=
            -----END CERTIFICATE-----
          dbSslRootCrtSecret: null
          dbSslClientCrtSecret: null
          masterkey: {{ .Values.zitadel.mainKey | fetchSecretValue | quote }}
          configmapConfig:
            ExternalDomain: {{ .Values.zitadel.hostName }}
            ExternalPort: 443
            ExternalSecure: true
            LogStore:
              Access:
                Stdout:
                  Enabled: true
            TLS:
              Enabled: false # The example delegates this to the upstream ingress controller
            Database:
              postgres:
                Host: {{ .Values.zitadel.database.host | fetchSecretValue | quote }}
                Port: {{ .Values.zitadel.database.port | fetchSecretValue | quote }}
                Database: zitadel
                MaxOpenConns: 50
                MaxConnLifetime: 1h
                MaxConnIdleTime: 5m
                Options:
                User:
                  Username: z1t4d3l # A bit more tricky to guess but could also use a secret
                  Password: {{ .Values.zitadel.database.zitadelPassword | fetchSecretValue | quote }}
                  SSL:
                    Mode: verify-full
                    RootCert: /.secrets/ca.crt
                    Cert:
                    Key:
                Admin:
                  Username: {{ .Values.zitadel.database.username | fetchSecretValue | quote }}
                  Password: {{ .Values.zitadel.database.password | fetchSecretValue | quote }}
                  SSL:
                    Mode: verify-full
                    RootCert: /.secrets/ca.crt
                    Cert:
                    Key:
            ...
```