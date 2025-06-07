FROM login-dev-base AS login-dev-dependencies

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
COPY apps/login/cypress/package.json ./apps/login/cypress/

RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile

ENTRYPOINT ["pnpm"]
