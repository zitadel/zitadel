CGO_ENABLED=0 go build \
  -a \
  -installsuffix cgo \
  -ldflags "-X main.Version=$(git rev-parse --abbrev-ref HEAD | sed -e 's/heads\///')" \
  -o zitadelctl \
  ./cmd/zitadelctl/main.go
  
