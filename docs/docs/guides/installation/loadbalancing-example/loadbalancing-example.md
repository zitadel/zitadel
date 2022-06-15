---
title: Load Balancing Example
---

With this example configuration, you create a near production environment for ZITADEL with [Docker Compose](https://docs.docker.com/compose/).

The stack consists of three long-running containers:
- A [Traefik](https://doc.traefik.io/traefik/) reverse proxy with upstream HTTP/2 enabled, issuing a self-signed TLS certificate
- A secure ZITADEL container configured for a custom domain
- A secure [CockroachDB](https://www.cockroachlabs.com/docs/stable/)

You will need to download the following files:
- [docker-compose.yaml](./docker-compose.yaml)
- [example-zitadel-config.yaml](./example-zitadel-config.yaml)
- [example-zitadel-secrets.yaml](./example-zitadel-secrets.yaml)
- [example-zitadel-init-steps.yaml](./example-zitadel-init-steps.yaml)

The setup is tested against Docker version 20.10.17 and Docker Compose version v2.2.3

```bash
# Download the docker compose example configuration. For example:
wget https://docs.zitadel.com/docs/guides/installation/loadbalancing-example/docker-compose.yaml

# Download and adjust the example configuration file containing standard configuration
wget https://docs.zitadel.com/docs/guides/installation/loadbalancing-example/example-zitadel-config.yaml

# Download and adjust the example configuration file containing secret configuration
wget https://docs.zitadel.com/docs/guides/installation/loadbalancing-example/example-zitadel-secrets.yaml

# Download and adjust the example configuration file containing database initialization configuration
wget https://docs.zitadel.com/docs/guides/installation/loadbalancing-example/example-zitadel-init-steps.yaml

# A single ZITADEL instance always needs the same 32 characters long masterkey
# If you haven't done so already, you can generate a new one.
# For example:
export ZITADEL_MASTERKEY="$(tr -dc A-Za-z0-9 </dev/urandom | head -c 32)"

# Run the database and application containers
docker compose up --detach
```

Make `127.0.0.1` available at `my.domain`. For example, this can be achived with an entry `127.0.1.1 my.domain` in the `/etc/hosts` file.

Open your favorite internet browser at [https://my.domain/ui/console/](https://my.domain/ui/console/).
You can safely proceed, if your browser warns you about the insecure self-signed TLS certificate.
This is the IAM admin users login according to your configuration in the [example-zitadel-init-steps.yaml](./example-zitadel-init-steps.yaml):
- **username**: *root@<span></span>my-org.my.domain*
- **password**: *RootPassword1!*
