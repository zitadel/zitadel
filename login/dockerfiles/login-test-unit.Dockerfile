FROM login-client AS login-test-unit
# Copy package.json first for better dependency caching
COPY apps/login/package.json ./apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./apps/login
# Copy source code
COPY apps/login ./apps/login
# Run unit tests (equivalent to turbo test:unit:standalone)
RUN cd apps/login && pnpm test:unit
