echo $(uname -m)

if [ "$(uname)" = "Darwin" ]; then
  curl -L -o protoc-gen-connect-openapi.tar.gz https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_darwin_all.tar.gz
else
  ARCH=$(uname -m)
  case $ARCH in
    x86_64)
      ARCH="amd64"
      ;;
    aarch64|arm64)
      ARCH="arm64"
      ;;
    *)
      echo "Unsupported architecture: $ARCH"
      exit 1
      ;;
  esac
  curl -L -o protoc-gen-connect-openapi.tar.gz https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_linux_${ARCH}.tar.gz
fi
tar -xvf protoc-gen-connect-openapi.tar.gz protoc-gen-connect-openapi