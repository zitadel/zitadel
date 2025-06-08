FROM golang:1.20.5-alpine3.18

RUN go install github.com/eliobischof/grpc-mock/cmd/grpc-mock@01b09f60db1b501178af59bed03b2c22661df48c

COPY mocked-services.cfg .
COPY initial-stubs initial-stubs
COPY --from=protos . .

ENTRYPOINT [ "sh", "-c", "grpc-mock -v 1 -proto $(tr '\n' ',' < ./mocked-services.cfg) -stub-dir ./initial-stubs -mock-addr :22222" ]
