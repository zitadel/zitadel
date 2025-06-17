FROM node:20-bookworm AS login-pnpm
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
RUN apt-get update && apt-get install -y --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /build

RUN  --mount=type=cache,target=${PNPM_HOME} \
  pnpm config set store-dir ${PNPM_HOME}

COPY \
  turbo.json \
  .npmrc \
  package.json \
  pnpm-lock.yaml \
  pnpm-workspace.yaml \
  ./

ENTRYPOINT ["pnpm"]
