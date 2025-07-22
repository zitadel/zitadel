FROM login-dev-base AS login-lint
COPY .prettierrc .prettierignore ./
COPY apps/login/package.json apps/login/
RUN  --mount=type=cache,id=pnpm,target=/pnpm/store \
     pnpm install --frozen-lockfile --workspace-root --filter apps/login
COPY . .
RUN pnpm lint && pnpm format
