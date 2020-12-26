#! /bin/sh

set -eux

GEN_PATH=src/app/proto/generated

echo "Create folders"
mkdir -p $GEN_PATH

echo "Generate grpc"

protoc \
  -I=proto/include \
  -I=node_modules/google-proto-files \
  --js_out=import_style=commonjs,binary:$GEN_PATH \
  --grpc-web_out=import_style=commonjs+dts,mode=grpcweb:$GEN_PATH \
  proto/include/zitadel/*.proto

echo "Generate annotations js file (compatibility)"

# mkdir -p $GEN_PATH/proto
# touch $GEN_PATH/proto/message_pb.js
# echo "export {}" > $GEN_PATH/proto/message_pb.d.ts

# mkdir -p $GEN_PATH/google/api/
# touch $GEN_PATH/google/api/annotations_pb.js
# echo "export {}" > $GEN_PATH/google/api/annotations_pb.d.ts

# mkdir -p $GEN_PATH/validate
# touch $GEN_PATH/validate/validate_pb.js
# echo "export {}" > $GEN_PATH/validate/validate_pb.d.ts

# mkdir -p $GEN_PATH/protoc-gen-swagger/options
# touch $GEN_PATH/protoc-gen-swagger/options/annotations_pb.js
# echo "export {}" > $GEN_PATH/protoc-gen-swagger/options/annotations_pb.d.ts

# mkdir -p $GEN_PATH/authoption
# touch $GEN_PATH/authoption/options_pb.js
# echo "export {}" > $GEN_PATH/authoption/options_pb.d.ts