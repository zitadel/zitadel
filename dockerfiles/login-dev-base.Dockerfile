FROM login-pnpm AS login-dev-base
COPY \
  turbo.json \
  .npmrc \
  package.json \
  ./
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter .

