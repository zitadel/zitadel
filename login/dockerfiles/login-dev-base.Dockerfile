FROM login-pnpm AS login-dev-base
# Install all workspace dependencies with caching
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --prefer-offline --workspace-root --filter .

