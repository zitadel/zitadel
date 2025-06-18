FROM login-client AS login-test-unit
COPY apps/login/package.json ./apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./apps/login
COPY apps/login ./apps/login
RUN pnpm test:unit:standalone
