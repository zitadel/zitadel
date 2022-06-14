---
title: Near Production Example
---

With this configuration, you create an almost production-ready docker-compose environment for ZITADEL.

The stack consists of three long-running containers:
- A secure [CockroachDB](https://www.cockroachlabs.com/docs/stable/)
- A secure ZITADEL container configured for a custom domain
- A [Traefik](https://doc.traefik.io/traefik/) reverse proxy with upstream HTTP/2 enabled, issuing a self-signed TLS certificate

```bash
# Download the docker compose example configuration. For example:
wget https://docs.zitadel.com/docs/guides/installation/near-production-example/docker-compose.yaml

# Download and adjust the example configuration file containing standard configuration
wget https://docs.zitadel.com/docs/guides/installation/near-production-example/example-zitadel-config.yaml

# Download and adjust the example configuration file containing secret configuration
wget https://docs.zitadel.com/docs/guides/installation/near-production-example/example-zitadel-secrets.yaml

# Download and adjust the example configuration file containing database initialization configuration
wget https://docs.zitadel.com/docs/guides/installation/near-production-example/example-zitadel-init-steps.yaml

# A single ZITADEL instance always needs the same 32 characters long masterkey
# If you haven't done so already, you can generate a new one.
# For example:
export ZITADEL_MASTERKEY="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 32)"

# Run the database and application containers
docker compose up --detach
```

Make `127.0.0.1` available at `my.domain`. For example, this can be achived with an entry `127.0.1.1 my.domain` in the `/etc/hosts` file.

Open your browser at https://my.domain/ui/console/. You can safely proceed, if your browser warns you about the insecure self-signed TLS certificate.
With the configuration from the example files, you can log in with the following credentials:
- **username**: *root@<span></span>my-org.my.domain*
- **password**: *RootPassword1!*
