FROM login-pnpm AS login-dev-base
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter .

