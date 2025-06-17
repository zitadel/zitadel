FROM login-pnpm AS login-dev-base
RUN pnpm install --frozen-lockfile --prefer-offline --workspace-root --filter .

