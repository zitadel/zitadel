# protoc-gen-authoption

Proto options to annotate auth methods in protos

## Generate protos/templates
protos: `go generate authoption/generate.go`  
templates/install: `go generate generate.go`

## Usage
```
// proto file
import "authoption/options.proto";

service MyService {

    rpc Hello(Hello) returns (google.protobuf.Empty) {
        option (google.api.http) = {
        get: "/hello"
        };

        option (caos.zitadel.utils.v1.auth_option) = {
            zitadel_permission: "hello.read"
            zitadel_check_param: "id"
        };
    }

    message Hello {
        string id = 1;
    }
}
```
Caos Auth Option is used for granting groups
On each zitadel role is specified which auth methods are allowed to call

Get protoc-get-authoption: ``go get github.com/zitadel/zitadel/internal/protoc/protoc-gen-authoption``

Protc-Flag: ``--authoption_out=.``