FROM mcr.microsoft.com/devcontainers/typescript-node:20-bookworm

ENV SHELL=/bin/bash \
    DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    LC_ALL=C.UTF-8 \
    CI=1 \
    PNPM_HOME=/home/node/.local/share/pnpm \
    PATH=/home/node/.local/share/pnpm:$PATH


RUN apt-get update && \
    apt-get --no-install-recommends install -y \
    libgtk2.0-0 libgtk-3-0 libgbm-dev libnotify-dev libnss3 libxss1 libasound2 libxtst6 xauth xvfb && \
    apt-get clean && \
    corepack enable && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack prepare pnpm@9.1.2 --activate
