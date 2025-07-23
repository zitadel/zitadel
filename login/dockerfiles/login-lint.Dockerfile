FROM login-dev-base AS login-lint
# Copy linting configuration files first for better caching
COPY .prettierrc .prettierignore ./
COPY apps/login/package.json apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter apps/login
# Copy source code
COPY . .
# Run linting and formatting (equivalent to turbo lint)
RUN cd apps/login && pnpm lint && pnpm exec prettier --check .
