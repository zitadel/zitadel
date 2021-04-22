---
title: GitOps Mode on an existing Kubernetes cluster
---

:::tip What I need
I'd like to have a reproducible ZITADEL environment and a pull-based configuration management for safe and comfortable day-two operations.
:::

First, copy the template files [database.yml](./templates/gitops/database.yml) and [zitadel.yml](./templates/gitops/zitadel.yml) to the root of a new git Repository. Then adjust the values in database.yml and zitadel.yml to match your environment. Especially the values for the domain, cluster DNS, storage class, email and Twilio are important.  

Now open a terminal and execute the following commands.

```bash
# Download the zitadelctl binary
curl -s https://api.github.com/repos/caos/zitadel/releases/tags/v1.0.0 | grep "browser_download_url.*zitadelctl-$(uname | awk '{print tolower($0)}')-amd64" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/zitadelctl && sudo chmod +x /usr/local/bin/zitadelctl && sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl
sudo chmod +x /usr/local/bin/zitadelctl
sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl

# Create an orb file at ${HOME}/.orb/config
MY_GIT_REPO="git@github.com:me/my-orb.git"
zitadelctl --gitops configure --repourl ${MY_GIT_REPO} --masterkey "$(openssl rand -base64 21)"

# Write the Twiilio sender ID and auth token so that ZITADEL is able to send your users SMS.
TWILIO_SID=<My Twilio Sender ID>
TWILIO_AUTH_TOKEN=<My Twilio auth token>
zitadelctl --gitops writesecret zitadel.twiliosid.encrypted --value $SID
zitadelctl --gitops writesecret zitadel.twilioauthtoken.encrypted --value $TWILIO_AUTH_TOKEN

# Write your email relays app key so that ZITADEL is able to verify your users email addresses
EMAIL_APP_KEY=<My email relays app key>
zitadelctl --gitops writesecret zitadel.twilioauthtoken.encrypted --value $EMAIL_APP_KEY

# Deploy the operators to the current-context of your ~/.kube/config file
zitadelctl --gitops takeoff

# Enjoy watching the zitadel pods becoming ready
watch "kubectl --namespace caos-zitadel get pods"
```

ZITADEL needs [gRPC-Web](https://grpc.io/docs/platforms/web/basics/) for client-server communication, which the widely spread [NGINX Ingress Controller](https://kubernetes.github.io/ingress-nginx/) doesn't support out-of-the-box but Ambassador does. If you don't have an [Ambassador](https://www.getambassador.io/) running, we recommend you run it with our operator [BOOM](https://github.com/caos/orbos/blob/v4.0.0/docs/boom/boom.md). Do so by adding the template [boom.yml](./templates/boom.yml) to the root of your Repository and execute the following commands. 

```bash
# Download the orbctl binary
curl -s https://api.github.com/repos/caos/orbos/releases/tags/v4.0.0 | grep "browser_download_url.*orbctl.$(uname).$(uname -m)" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/orbctl
sudo chmod +x /usr/local/bin/orbctl
sudo chown $(id -u):$(id -g) /usr/local/bin/orbctl

# Deploy the operator to the current-context of your ~/.kube/config file
orbctl --gitops takeoff

# Enjoy watching the ambassador pod becoming ready
watch "kubectl --namespace caos-system get pods"
```

Congratulations, you can accept traffic at four new ZITADEL [subdomains](/docs/apis/domains) now.
