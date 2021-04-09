---
title: GitOps Mode on an existing Kubernetes cluster
---

I'd like to have a reproducible ZITADEL environment and a pull-based configuration management for safe and comfortable day-two operations

First of all, copy the files [orbiter.yml](examples/orbiter/gce/orbiter.yml) and [boom.yml](examples/boom/boom.yml) to the root of a new git Repository.

```bash
# Downloading the zitadelctl binary
curl -s https://api.github.com/repos/caos/zitadel/releases/tags/v0.118.2 | grep "browser_download_url.*zitadelctl-$(uname | awk '{print tolower($0)}')-amd64" | cut -d '"' -f 4 | sudo wget -i - -O /usr/local/bin/zitadelctl && sudo chmod +x /usr/local/bin/zitadelctl && sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl
sudo chmod +x /usr/local/bin/zitadelctl
sudo chown $(id -u):$(id -g) /usr/local/bin/zitadelctl

# Create an orb file at ${HOME}/.orb/config
MY_GIT_REPO="git@github.com:me/my-orb.git"
zitadelctl configure --repourl ${MY_GIT_REPO} --masterkey "$(openssl rand -base64 21)"

```
