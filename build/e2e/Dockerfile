ARG GO_VERSION=1.17

#######################
##  run initialize setup step for ZITADEL tests
#######################
FROM golang:${GO_VERSION}
RUN mkdir -p /cmd/e2e-setup
WORKDIR /

# copy setup
COPY . .
RUN e2e/e2e-setup-compose.sh
