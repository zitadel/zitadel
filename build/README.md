# Build Readme

## Console

### Build Container

```shell
docker build -f build/console.dockerfile --platform linux/amd64 . -t zitadel-console
```

> The `--platform linux/amd64` is needed if you run an arm device (including apple silicon) since [grpc-node does only include an m1 arm binary but not linux](https://github.com/grpc/grpc-node/issues/1405)

### Run Container

```shell
docker run -p 8080:8080 zitadel-console
```

## Backend

## Docs