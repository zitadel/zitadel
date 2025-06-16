FROM login-dev-base AS login-lint
COPY .prettierrc .prettierignore ./
COPY packages/zitadel-tsconfig packages/zitadel-tsconfig
COPY packages/zitadel-prettier-config packages/zitadel-prettier-config
COPY packages/zitadel-eslint-config packages/zitadel-eslint-config
COPY apps/login/package.json apps/login/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-login
COPY apps/login apps/login
