ARG GO_VERSION=1.17

#######################
##  run initialize setup step for ZITADEL tests
#######################
FROM golang:${GO_VERSION}
WORKDIR /go/src/github.com/zitadel/zitadel

# copy setup
COPY . .
RUN e2e/e2e-setup-compose.sh
