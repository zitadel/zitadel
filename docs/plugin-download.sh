if [ "$(uname)" = "Darwin" ]; then
  wget -O protoc-gen-connect-openapi.tar.gz https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_darwin_all.tar.gz
else
  wget -O protoc-gen-connect-openapi.tar.gz https://github.com/sudorandom/protoc-gen-connect-openapi/releases/download/v0.18.0/protoc-gen-connect-openapi_0.18.0_linux_$(uname -m).tar.gz
fi
tar -xzvf protoc-gen-connect-openapi.tar.gz