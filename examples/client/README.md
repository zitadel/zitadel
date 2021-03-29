# Example

This guide describes how to generate clients to interact with ZITADEL.

ZITADEL decided to not check in generated files after v0.104.5.

As the go-sdk is not ready yet we recommend to to build the client in your own project.

## Requirements

 - docker

## Generate client stub

### PROJECT_PATH

The PROJECT_PATH argument is needed for replacing imports in the generated files.
The path MUST represent the folder where the generated ZITADEL packages will reside in.

This replacement is needed for the message proto.

### TAG_NAME

It's recommended to clone a specific tag.

For example: TAG_NAME=v0.118.3

`DOCKER_BUILDKIT=1 docker build --target zitadel-copy -t zitadel:example --build-arg PROJECT_PATH=github.com/caos/zitadel/examples/client --build-arg TAG_NAME=master -f Dockerfile . -o .`
