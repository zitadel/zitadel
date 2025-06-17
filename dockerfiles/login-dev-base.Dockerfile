FROM login-pnpm AS login-dev-base
RUN --mount=type=cache,target=${PNPM_HOME} \
    pnpm install --frozen-lockfile --prefer-offline --workspace-root --filter .

