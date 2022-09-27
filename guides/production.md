# Production Build

To create a production build to run locally, create a snapshot release with goreleaser:

```sh
goreleaser release --snapshot --rm-dist
```

This can be released to production (if you have credentials configured) using gorelease as well:

```sh
goreleaser release
```
