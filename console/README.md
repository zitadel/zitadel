# Console

This folder includes all the code for ZITADEL Console, the management frontend for ZITADEL.

## Requirements

### Build gRPC client

You need to generate the gRPC client for console. This can be done by using `yarn generate`.
This will generate the necessary files with the help of a `bufbuild` docker image.

If you have `buf` installed locally you can invoke `buf generate --include-imports` from the repository root to generate the files.

### Create a local env file

--- max can you insert quickly somehting here

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory. Use the `--prod` flag for a production build.

## Running end-to-end tests

Please refer to the [contributing guide](../CONTRIBUTING.md#console)

## Container Image

If you just want to start console locally without installing node you can fallback to our container image.
Execute the following commands from the repository root to build and start a local version of ZITADEL 

```shell
docker build -f console/Dockerfile . -t zitadel-console
```

```shell
docker run -p 8080:8080 zitadel-console
```
