FROM node:18 as npm-base

ENV PROTOC_ARCH x86_64
ARG PROTOC_VERSION=3.18.0
ARG PROTOC_ZIP=protoc-${PROTOC_VERSION}-linux-${PROTOC_ARCH}.zip
ARG GRPC_WEB_VERSION=1.3.0
ARG GATEWAY_VERSION=2.6.0
ARG VALIDATOR_VERSION=0.6.2
ARG GRPC_WEB=protoc-gen-grpc-web-${GRPC_WEB_VERSION}-linux-x86_64

WORKDIR /console

RUN apt install curl -y \
    && curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/$PROTOC_ZIP \
    && unzip -o $PROTOC_ZIP -d /usr/local bin/protoc \
    && unzip -o $PROTOC_ZIP -d /proto include/* \
    && rm -f $PROTOC_ZIP \
    && curl -OL https://github.com/grpc/grpc-web/releases/download/${GRPC_WEB_VERSION}/${GRPC_WEB} \
    && mv ${GRPC_WEB} /usr/local/bin/protoc-gen-grpc-web \
    && chmod +x /usr/local/bin/protoc-gen-grpc-web \
    && curl https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/v${VALIDATOR_VERSION}/validate/validate.proto --create-dirs -o include/validate/validate.proto  \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v${GATEWAY_VERSION}/protoc-gen-openapiv2/options/annotations.proto --create-dirs -o include/protoc-gen-openapiv2/options/annotations.proto \
    && curl https://raw.githubusercontent.com/grpc-ecosystem/grpc-gateway/v${GATEWAY_VERSION}/protoc-gen-openapiv2/options/openapiv2.proto --create-dirs -o include/protoc-gen-openapiv2/options/openapiv2.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto --create-dirs -o include/google/api/annotations.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto --create-dirs -o include/google/api/http.proto \
    && curl https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/field_behavior.proto --create-dirs -o include/google/api/field_behavior.proto

COPY proto/ include/.

COPY console/package.json console/package-lock.json ./
RUN npm ci

RUN ls -la node_modules/google-proto-files/

RUN mkdir -p src/app/proto/generated
RUN protoc -I=include -I=node_modules/google-proto-files --js_out=import_style=commonjs,binary:src/app/proto/generated --grpc-web_out=import_style=typescript,mode=grpcweb:src/app/proto/generated $(find include -iname "*.proto")

COPY console .


#######################
## copy for local dev
#######################
FROM scratch as npm-copy
COPY --from=npm-base /console/src/app/proto/generated /console/src/app/proto/generated

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as angular-build

#RUN npm run lint
RUN npm run prodbuild

#######################
## Only Copy Assets
#######################
FROM scratch as angular-export
COPY --from=angular-build /console/dist/console .
