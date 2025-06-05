# BUILD STAGE
FROM node:20-alpine

WORKDIR /app

RUN apk add --no-cache libc6-compat bash git
RUN corepack enable && corepack prepare pnpm@latest --activate

# Copy remote turbo.json config for pruning
COPY turbo.json ./
COPY .npmrc ./

# pnpm store + turbo build cache
RUN mkdir -p .pnpm-store .next

# Copy just lockfile & manifests for better cache-hit
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
COPY packages/zitadel-eslint-config/package.json ./packages/zitadel-eslint-config/
COPY packages/zitadel-prettier-config/package.json ./packages/zitadel-prettier-config/
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
COPY packages/zitadel-tailwind-config/package.json ./packages/zitadel-tailwind-config/
COPY packages/zitadel-tsconfig/package.json ./packages/zitadel-tsconfig/
COPY apps/login/package.json ./apps/login/

RUN --mount=type=cache,target=/app/.pnpm-store \
    pnpm install --frozen-lockfile --store-dir .pnpm-store

# Full source
COPY . .
