# Verify the token actor reaches Actions v2 execution targets

> One of the runbooks in [`skills/`](README.md). Read [skills/README.md](README.md) first for the
> ground rules that apply to all of them: local instances only, ask before changing anything,
> clean up afterwards.

Impersonated tokens carry an *actor*: the user who performed the token exchange. Zitadel
exposes that actor to Actions, so a target can tell "this request is Bob acting as Alice"
apart from "this request is Alice".

This runbook proves the whole chain end to end against a local instance:

| Surface | What is checked |
| --- | --- |
| Token exchange | the exchanged access token and id token carry the `act` claim |
| `function/preaccesstoken` | the webhook payload carries `actor` when a JWT is minted |
| `function/preuserinfo` | the webhook payload carries `actor`, both when the id token is created and when `/oidc/v1/userinfo` is called |
| Negative control | a non impersonated token produces a payload with **no** `actor` key at all |

Relevant code: `ContextInfo.Actor` in `internal/api/oidc/userinfo.go`, populated from
`token.actor` (userinfo endpoint) and `session.Actor` (`internal/api/oidc/token.go`).
Actions v1 gets the same data as `ctx.v1.actor` via `internal/actions/object/token.go`.

## Prerequisites

- A Zitadel instance **already running locally** on the branch you want to test.
- An admin personal access token in a file, `admin.pat` in the repository root by default.
  It must hold exactly the token: no quotes, no backticks, no extra lines.
- `curl` and `jq` on `PATH`.
- Outbound HTTPS to `webhook.site`, or see [Receiving webhooks locally](#receiving-webhooks-locally).
- The `oidcTokenExchange` instance feature enabled. The script checks this and prints the
  command to enable it if it is off.

Everything else is created and removed by the script. No pre-existing impersonator, end user,
project or client is required.

## Safety

**Local test instances only.** The script refuses any `--url` whose host is not `localhost`,
`127.0.0.1` or `[::1]`, because it:

- enables impersonation instance wide (`enableImpersonation` in the security policy),
- creates a machine user holding `IAM_END_USER_IMPERSONATOR`,
- creates a human user, a project and an OIDC client.

It asks for confirmation before touching anything. If you are driving this with an AI
assistant, have it show you the plan and get your explicit go-ahead before it passes `--yes`.

By default everything created is deleted again on the way out, including on failure, and the
original `enableImpersonation` value is restored.

## Run it

```bash
./skills/test-actor-in-action-v2.sh
```

Options:

| Flag | Meaning |
| --- | --- |
| `--url URL` | instance to test, default `http://localhost:8080` |
| `--pat FILE` | admin PAT file, default `admin.pat` |
| `--local` | use a local HTTP sink instead of webhook.site, see below |
| `--keep` | keep everything that was created, for debugging |
| `--yes` | skip the confirmation prompt |

A successful run ends with:

```
==> Results
    PASS access token sub is the impersonated user
    PASS access token act.sub is the impersonator
    PASS access token act.iss is the issuer
    PASS id token sub is the impersonated user
    PASS id token act.sub is the impersonator
    PASS userinfo sub is the impersonated user
    PASS function/preaccesstoken payload carries the actor
    PASS both function/preuserinfo payloads carry the actor
    PASS non-impersonated payload omits the actor key entirely

All assertions passed.
```

The exit status is non-zero if any assertion fails, and the received payloads are dumped so you
can see what actually arrived.

## What the payload looks like

Note the keys are snake_case, they come straight from `domain.TokenActor`:

```json
{
  "function": "function/preuserinfo",
  "userinfo": { "sub": "385694975679005038" },
  "actor": {
    "user_id": "385694975578341742",
    "issuer": "http://localhost:8080"
  }
}
```

`actor` is `omitempty`: for a normal, non impersonated token the key is absent entirely rather
than present and null. Existing targets therefore see no change in their payload.

For a delegation chain, `actor.actor` holds the previous actor.

## Receiving webhooks locally

`--local` starts a throwaway HTTP sink on `127.0.0.1` and points the target at it. This is only
usable if the instance was started with its SSRF guard relaxed: `HTTPClient.DenyList` in
`cmd/defaults.yaml` blocks `localhost` and the loopback and private CIDR ranges by default, and
the target is rejected with `Errors.Target.DeniedURL` otherwise.

To use it, start the instance with the deny list cleared, for example:

```bash
ZITADEL_HTTPCLIENT_DENYLIST= zitadel start-from-init --masterkeyFromEnv
```

Never do that on anything but a local test instance. Without that, use the default
webhook.site receiver; the script creates an anonymous token, uses it, reads it back through
the webhook.site API and deletes it again.

## Troubleshooting

| Symptom | Cause |
| --- | --- |
| `the PAT ... was rejected` | the PAT file has extra characters. A trailing backtick picked up from a copy-paste is the classic one, and only shows up much later as `token contains an invalid number of segments`. |
| `unsupported_grant_type` | the `oidcTokenExchange` instance feature is off. |
| `Errors.TokenExchange.Impersonation.PolicyDisabled` | `enableImpersonation` is off in the security policy. |
| `actor_token invalid` | the impersonator token is malformed or expired. |
| `Errors.Target.DeniedURL` | the endpoint is in `HTTPClient.DenyList`, see above. |
| `could not create a webhook.site token` | webhook.site rate limit, ten per minute for anonymous use. Wait a minute or use `--local`. |
| no payloads arrive | the instance cannot reach the endpoint. Check outbound network access. |

The token endpoint hides the underlying error by default. To see it, enable the
`debugOidcParentError` instance feature, reproduce, then turn it back off:

```bash
curl -X PUT http://localhost:8080/v2beta/features/instance \
  -H "Authorization: Bearer $(cat admin.pat)" -H 'Content-Type: application/json' \
  -d '{"debugOidcParentError": true}'
```

## Doing it by hand

If the script is broken or you want to poke at a single step, these are the two calls that
matter. Everything else is setup.

Impersonate by user id and ask for a JWT:

```bash
curl -X POST http://localhost:8080/oauth/v2/token \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -u "${CLIENT_ID}:${CLIENT_SECRET}" \
  -d 'grant_type=urn:ietf:params:oauth:grant-type:token-exchange' \
  -d "subject_token=${END_USER_ID}" \
  -d 'subject_token_type=urn:zitadel:params:oauth:token-type:user_id' \
  -d "actor_token=${IMPERSONATOR_PAT}" \
  -d 'actor_token_type=urn:ietf:params:oauth:token-type:access_token' \
  -d 'requested_token_type=urn:ietf:params:oauth:token-type:jwt' \
  -d 'scope=openid profile email'
```

Then use the returned token on the userinfo endpoint:

```bash
curl http://localhost:8080/oidc/v1/userinfo -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

See [the token exchange guide](../apps/docs/content/guides/integrate/token-exchange.mdx) for the
other impersonation flavours, including impersonation with a real subject token and refreshing
an impersonated token.
