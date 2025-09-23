FROM mcr.microsoft.com/devcontainers/typescript-node:20-bookworm

ENV SHELL=/bin/bash \
    DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8 \
    PNPM_HOME=/home/node/.local/share/pnpm \
    PATH=/home/node/.local/share/pnpm:$PATH

RUN apt-get update && \
    apt-get --no-install-recommends install -y \
    # Cypress dependencies
    libgtk2.0-0 libgtk-3-0 libgbm-dev libnotify-dev libnss3 libxss1 libasound2 libxtst6 xauth xvfb && \
    apt-get clean

COPY package.json /tmp/package.json
WORKDIR /tmp

RUN corepack enable

RUN pnpm --version

RUN chown -R node:node ${PNPM_HOME}

WORKDIR /workspaces/zitadel

USER node
