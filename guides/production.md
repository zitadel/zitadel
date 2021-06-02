# Production Build

This can also be run locally!

```bash
DOCKER_BUILDKIT=1 docker build -f build/dockerfile . -t zitadel:local --build-arg ENV=prod
```
