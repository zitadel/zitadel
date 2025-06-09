FROM login-pnpm AS login-dev-base
COPY \
  turbo.json \
  .npmrc \
  package.json \
  ./
COPY apps/login/package.json ./apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter . --filter=apps/login
