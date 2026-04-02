# ZITADEL Cloud

Cloud administration and multi-instance console for ZITADEL Cloud.

This app composes shared console pages with multi-instance routing and adds cloud-specific features.

## Architecture

```
apps/cloud/
├── app/
│   ├── instances/           # Instance list, create
│   ├── instances/[id]/      # Per-instance console (re-uses console pages)
│   │   ├── users/
│   │   ├── organizations/
│   │   ├── projects/
│   │   └── ...
│   ├── billing/             # Cloud-only
│   ├── usage/               # Cloud-only
│   └── debug/               # Preview/dev only — test instance config
```

## How It Relates to Console

- `apps/console` — deploys standalone for self-hosted customers (single instance)
- `apps/cloud` — deploys for ZITADEL Cloud (multi-instance, wraps console pages)

Both share page components. The difference is routing and how the API target is resolved.

## Running

```bash
pnpm nx dev cloud
```

## License

Source Available (see LICENSE file in this directory)
