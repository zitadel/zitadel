# @zitadel/examples/angular

Angular SPA example using the same shared `@zitadel/examples/spa-bff` server contract as the React example.

## What this demonstrates

- Angular SPA on `http://localhost:4200`
- Login redirect, callback, and token exchange handled server-side by `examples/spa-bff`
- Browser receives only minimal session state from `/session`
- Protected user info is fetched through `/api/userinfo`

## Run

```sh
cp examples/spa-bff/.env.example examples/spa-bff/.env.local
# set SPA_ORIGIN=http://localhost:4200 and your ZITADEL values

pnpm --filter @zitadel/examples/spa-bff dev
pnpm --filter @zitadel/examples/angular dev
```

Then open `http://localhost:4200`.
