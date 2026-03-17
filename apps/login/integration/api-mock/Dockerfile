FROM bufbuild/buf:1.64.0 AS dependencies
RUN buf export https://github.com/envoyproxy/protoc-gen-validate.git --path validate --output /proto && \
    buf export https://github.com/grpc-ecosystem/grpc-gateway.git --path protoc-gen-openapiv2 --output /proto && \
    buf export https://github.com/googleapis/googleapis.git --path google/api/annotations.proto --path google/api/http.proto --path google/api/field_behavior.proto --output /proto

    
FROM golang:1.25.8-alpine3.22 AS mock-zitadel

RUN go install github.com/conblem/grpc-mock/cmd/grpc-mock@f3c4cdd486f5afe7effa3a8aae5d26f45aa684c9

COPY mocked-services.cfg .
COPY initial-stubs initial-stubs
COPY --from=dependencies /proto/ ./
COPY --from=zitadel-protos /zitadel/ ./zitadel/

ENTRYPOINT [ "sh", "-c", "grpc-mock -v 1 -proto $(tr '\n' ',' < ./mocked-services.cfg) -stub-dir ./initial-stubs -mock-addr :22222" ]
