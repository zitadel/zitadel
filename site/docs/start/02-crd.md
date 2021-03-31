---
title: CRD Mode on an existing Kubernetes cluster
---

I'd like to see an automatically operated ZITADEL instance running on my own [Kubernetes](https://kubernetes.io/) cluster

```bash
# Downloading the zitadelctl binary
curl -s https://api.github.com/repos/caos/zitadel/releases/tags/v0.118.2 | grep "browser_download_url.*zitadelctl-$(uname | awk '{print tolower($0)}')-amd64" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/zitadelctl && sudo chmod +x /usr/local/bin/zitadelctl && sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl

# Deploying the operators to the current-context of your ~/.kube/config file
zitadelctl takeoff

# Downloading the configuration templates
wget https://raw.githubusercontent.com/caos/zitadel/crd-mode-docs/site/docs/start/templates/crd/database.yml
wget https://raw.githubusercontent.com/caos/zitadel/crd-mode-docs/site/docs/start/templates/crd/zitadel.yml

# Before applying, adjust the values in ./database.yml and zitadel.yml using your favorite text editor to match your environment.
# Especially the values for the domain, cluster DNS and storage class are important
kubectl apply --filename ./database.yml,./zitadel.yml

# Write the minimal Secrets
wget https://raw.githubusercontent.com/caos/zitadel/crd-mode-docs/site/docs/start/templates/example_keys && zitadelctl writesecret zitadel.keys.existing --file ./example_keys

# Enjoy watching the zitadel pods becoming ready
watch "kubectl --namespace caos-zitadel get pods"
```

ZITADEL needs [gRPC-Web](https://grpc.io/docs/platforms/web/basics/) for client-server communication, which the widely spread [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/) doesn't support out-of-the-box but Ambassador does. If you don't have an [Ambassador](https://www.getambassador.io/) running, we recommend you run it with our operator [BOOM](https://github.com/caos/orbos/blob/v3.1.4/docs/boom/boom.md).

```bash
# Downloading the orbctl binary
curl -s https://api.github.com/repos/caos/orbos/releases/tags/v3.1.4 | grep "browser_download_url.*orbctl.$(uname).$(uname -m)" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/orbctl

# Deploying the operator to the current-context of your ~/.kube/config file
orbctl takeoff

# Downloading the configuration template
wget https://raw.githubusercontent.com/caos/zitadel/crd-mode-docs/site/docs/start/templates/boom.yml

# Before applying, adjust the values in ./boom.yml using your favorite text editor to match your environment.
# Especially the value for proxyProtocol is of special interest
kubectl apply --filename ./boom.yml

# Enjoy watching the ambassador pod becoming ready
watch "kubectl --namespace caos-system get pods"
```
