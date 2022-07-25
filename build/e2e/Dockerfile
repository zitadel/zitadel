#######################
##  run initialize setup step for ZITADEL tests
#######################
FROM zitadel-test-base

# Invalidate caches when only test code changes
COPY ./docs/docs/guides/installation/run/docker-compose.yaml ./docs/docs/guides/installation/run/docker-compose.yaml
COPY ./cmd/e2e-setup ./cmd/e2e-setup
COPY ./e2e ./e2e

# copy setup
CMD e2e/e2e-setup.sh
