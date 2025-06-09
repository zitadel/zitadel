FROM login-dev-base AS login-lint
COPY packages/zitadel-tsconfig packages/zitadel-tsconfig
COPY packages/zitadel-prettier-config packages/zitadel-prettier-config
COPY packages/zitadel-eslint-config packages/zitadel-eslint-config
COPY apps/login apps/login
