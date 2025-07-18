FROM bufbuild/buf:1.54.0 AS proto-files
RUN buf export https://github.com/envoyproxy/protoc-gen-validate.git --path validate --output /proto-files && \
    buf export https://github.com/grpc-ecosystem/grpc-gateway.git --path protoc-gen-openapiv2 --output /proto-files && \
    buf export https://github.com/googleapis/googleapis.git --path google/api/annotations.proto --path google/api/http.proto --path google/api/field_behavior.proto --output /proto-files && \
    buf export https://github.com/zitadel/zitadel.git --path ./proto/zitadel --output /proto-files

FROM scratch
COPY --from=proto-files /proto-files /
