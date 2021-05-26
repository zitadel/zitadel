---
title: CRD Mode on an existing Kubernetes cluster
---

:::tip What I need
I'd like to see an automatically operated ZITADEL instance running on my own [Kubernetes](https://kubernetes.io/) cluster
:::

First, download the template configuration files [database.yml](./templates/crd/database.yml) and [zitadel.yml](./templates/crd/zitadel.yml). Then adjust the values in database.yml and zitadel.yml to match your environment. Especially the values for the domain, cluster DNS, storage class, email and Twilio are important.  

```bash
# Download the zitadelctl binary
curl -s https://api.github.com/repos/caos/zitadel/releases/latest | grep "browser_download_url.*zitadelctl-$(uname | awk '{print tolower($0)}')-amd64" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/zitadelctl && sudo chmod +x /usr/local/bin/zitadelctl && sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl
sudo chmod +x /usr/local/bin/zitadelctl
sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl

# Deploy the operators to the current-context of your ~/.kube/config file
zitadelctl takeoff

# As soon as the configuration is applied, the operators start their work
kubectl apply --filename ./database.yml,./zitadel.yml

# Write the encryption keys
wget https://raw.githubusercontent.com/caos/zitadel/main/site/docs/start/templates/example_keys && zitadelctl writesecret zitadel.keys.existing --file ./example_keys

# Write the Twiilio sender ID and auth token so that ZITADEL is able to send your users SMS.
TWILIO_SID=<My Twilio Sender ID>
TWILIO_AUTH_TOKEN=<My Twilio auth token>
zitadelctl writesecret zitadel.twiliosid.existing --value $SID
zitadelctl writesecret zitadel.twilioauthtoken.existing --value $TWILIO_AUTH_TOKEN

# Write your email relays app key so that ZITADEL is able to verify your users email addresses
EMAIL_APP_KEY=<My email relays app key>
zitadelctl writesecret zitadel.emailappkey.existing --value $EMAIL_APP_KEY

# Enjoy watching the zitadel pods becoming ready
watch "kubectl --namespace caos-zitadel get pods"
```

ZITADEL needs [gRPC-Web](https://grpc.io/docs/platforms/web/basics/) for client-server communication, which the widely spread [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/) doesn't support out-of-the-box but Ambassador does. If you don't have an [Ambassador](https://www.getambassador.io/) running, we recommend you run it with our operator [BOOM](https://github.com/caos/orbos/blob/v4.0.0/docs/boom/boom.md).

Download the template configuration file [boom.yml](./templates/boom.yml). Then adjust the values in boom.yml to match your environment.  

```bash
# Download the orbctl binary
curl -s https://api.github.com/repos/caos/orbos/releases/latest | grep "browser_download_url.*orbctl.$(uname).$(uname -m)" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/orbctl
sudo chmod +x /usr/local/bin/orbctl
sudo chown $(id -u):$(id -g) /usr/local/bin/orbctl

# Deploy the operator to the current-context of your ~/.kube/config file
orbctl takeoff

# As soon as the configuration is applied, BOOM starts its work
kubectl apply --filename ./boom.yml

# Enjoy watching the ambassador pod becoming ready
watch "kubectl --namespace caos-system get pods"
```

Congratulations, you can accept traffic at four new ZITADEL [subdomains](/docs/apis/domains) now.
