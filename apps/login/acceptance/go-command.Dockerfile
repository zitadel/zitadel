ARG LOGIN_TEST_ACCEPTANCE_GOLANG_TAG="golang:1.24-alpine"

FROM ${LOGIN_TEST_ACCEPTANCE_GOLANG_TAG}
RUN apk add curl jq
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /go-command .
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s \
  CMD curl -f http://localhost:${PORT}/healthy || exit 1
ENTRYPOINT [ "/go-command" ]
