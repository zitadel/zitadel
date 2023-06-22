FROM bufbuild/buf:1.21.0 as protos

RUN buf export https://github.com/envoyproxy/protoc-gen-validate.git --path validate --output /proto
RUN buf export https://github.com/grpc-ecosystem/grpc-gateway.git --path protoc-gen-openapiv2 --output /proto
RUN buf export https://github.com/googleapis/googleapis.git --path google/api/annotations.proto --path google/api/http.proto --path google/api/field_behavior.proto --output /proto
RUN buf export https://github.com/zitadel/zitadel.git --path ./proto/zitadel --output /proto

FROM scratch AS config

COPY mocked-services.cfg .
COPY initial-stubs initial-stubs
COPY --from=protos /proto .

FROM golang:1.20.5-alpine3.18 as grpc-mock

RUN go install github.com/eliobischof/grpc-mock/cmd/grpc-mock@01b09f60db1b501178af59bed03b2c22661df48c

COPY --from=config / .

ENTRYPOINT [ "sh", "-c", "grpc-mock -v 1 -proto $(tr '\n' ',' < ./mocked-services.cfg) -stub-dir ./initial-stubs" ]
