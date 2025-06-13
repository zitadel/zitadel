FROM login-client AS login-standalone-builder
COPY --from=login-dev-base /build/apps/login apps/login
