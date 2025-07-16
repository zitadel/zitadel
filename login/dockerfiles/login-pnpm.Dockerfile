ARG NODE_VERSION=20
FROM node:${NODE_VERSION}-bookworm AS login-pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable && COREPACK_ENABLE_DOWNLOAD_PROMPT=0 corepack prepare pnpm@9.1.2 --activate && \
    apt-get update && apt-get install -y --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /build
COPY turbo.json .npmrc package.json pnpm-lock.yaml pnpm-workspace.yaml ./
ENTRYPOINT ["pnpm"]
