FROM node:20-alpine AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN corepack enable

RUN apk add --no-cache libc6-compat bash git

WORKDIR /app

COPY \
  turbo.json \
  .npmrc \
  package.json \
  pnpm-lock.yaml \
  pnpm-workspace.yaml \
  ./

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
COPY packages/zitadel-eslint-config/package.json ./packages/zitadel-eslint-config/
COPY packages/zitadel-prettier-config/package.json ./packages/zitadel-prettier-config/
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
COPY packages/zitadel-tailwind-config/package.json ./packages/zitadel-tailwind-config/
COPY packages/zitadel-tsconfig/package.json ./packages/zitadel-tsconfig/
COPY apps/login/package.json ./apps/login/

RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile

COPY . .

ENTRYPOINT ["pnpm"]
