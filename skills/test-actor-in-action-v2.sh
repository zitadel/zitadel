#!/usr/bin/env bash
#
# test-actor-in-action-v2.sh
#
# End-to-end check that the impersonating actor is handed to Actions v2 execution
# targets. See skills/test-actor-in-action-v2.md for the full runbook.
#
# Safety: refuses to run against anything but a local Zitadel instance. It enables
# impersonation instance wide and creates users, a project and an OIDC client.
#
set -euo pipefail

URL="http://localhost:8080"
PAT_FILE="admin.pat"
RECEIVER="webhook.site"
KEEP=0
ASSUME_YES=0

usage() {
  cat <<'EOF'
Usage: skills/test-actor-in-action-v2.sh [options]

  --url URL    Zitadel instance to test (default: http://localhost:8080).
               Only localhost / 127.0.0.1 / [::1] are accepted.
  --pat FILE   File holding an admin personal access token (default: admin.pat).
  --local      Receive webhooks with a local HTTP sink instead of webhook.site.
               Use this when offline or rate limited. Requires python3.
  --keep       Do not clean up the resources this script creates.
  --yes        Skip the interactive confirmation.
  -h, --help   Show this help.
EOF
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --url) URL="${2:?--url needs a value}"; shift 2 ;;
    --pat) PAT_FILE="${2:?--pat needs a value}"; shift 2 ;;
    --local) RECEIVER="local"; shift ;;
    --keep) KEEP=1; shift ;;
    --yes|-y) ASSUME_YES=1; shift ;;
    -h|--help) usage; exit 0 ;;
    *) echo "unknown option: $1" >&2; usage >&2; exit 2 ;;
  esac
done

URL="${URL%/}"

# ---------------------------------------------------------------- output helpers

RED=''; GREEN=''; YELLOW=''; BOLD=''; RESET=''
if [[ -t 1 ]]; then
  RED=$'\033[31m'; GREEN=$'\033[32m'; YELLOW=$'\033[33m'; BOLD=$'\033[1m'; RESET=$'\033[0m'
fi

FAILURES=0
step()  { printf '\n%s==> %s%s\n' "$BOLD" "$*" "$RESET"; }
info()  { printf '    %s\n' "$*"; }
warn()  { printf '    %sWARN%s %s\n' "$YELLOW" "$RESET" "$*"; }
die()   { printf '\n%sERROR%s %s\n' "$RED" "$RESET" "$*" >&2; exit 1; }

pass()  { printf '    %sPASS%s %s\n' "$GREEN" "$RESET" "$*"; }
fail()  { printf '    %sFAIL%s %s\n' "$RED" "$RESET" "$*"; FAILURES=$((FAILURES + 1)); }
skip()  { printf '    %sSKIP%s %s\n' "$YELLOW" "$RESET" "$*"; }

check_eq() { # check_eq <description> <expected> <actual>
  if [[ "$2" == "$3" ]]; then pass "$1"; else fail "$1 (expected '$2', got '$3')"; fi
}

# ---------------------------------------------------------------- safety guards

# Only ever run this against a local test instance. It flips an instance wide
# security policy and creates users.
host="${URL#*://}"; host="${host%%/*}"; host="${host%%:*}"
case "$host" in
  localhost|127.0.0.1|'[::1]'|::1) ;;
  *) die "refusing to run against '$URL': only localhost, 127.0.0.1 and [::1] are allowed." ;;
esac

for dep in curl jq; do
  command -v "$dep" >/dev/null 2>&1 || die "missing required dependency: $dep"
done
if [[ "$RECEIVER" == "local" ]]; then
  command -v python3 >/dev/null 2>&1 || die "--local requires python3"
fi

[[ -f "$PAT_FILE" ]] || die "PAT file '$PAT_FILE' not found. Pass --pat FILE or create it (see the runbook)."

# Strip whitespace, quotes and backticks. Copying a PAT out of chat or markdown
# very easily drags a stray backtick along, which surfaces much later as an
# opaque "token contains an invalid number of segments".
ADMIN_PAT="$(tr -d ' \t\r\n`"'"'" < "$PAT_FILE")"
[[ -n "$ADMIN_PAT" ]] || die "PAT file '$PAT_FILE' is empty"

if [[ $ASSUME_YES -ne 1 ]]; then
  cat <<EOF

${BOLD}This will modify the Zitadel instance at ${URL}${RESET}

  - enable impersonation in the instance security policy
  - create a machine user (impersonator) with a PAT and IAM_END_USER_IMPERSONATOR
  - create a human user to impersonate
  - create a project and an OIDC client with the token exchange grant
  - create an Actions v2 target and two executions
  - receive webhooks via: ${RECEIVER}
EOF
  if [[ $KEEP -eq 1 ]]; then
    echo "  - keep everything afterwards (--keep)"
  else
    echo "  - delete all of the above afterwards"
  fi
  echo
  printf 'Continue? [y/N] '
  read -r reply < /dev/tty 2>/dev/null || reply=""
  case "$reply" in [yY]|[yY][eE][sS]) ;; *) die "aborted by user; nothing was created." ;; esac
fi

# ---------------------------------------------------------------- api helpers

RUN_ID="$$-$(date +%s)"

api() { # api <method> <path> [json-body]
  local method="$1" path="$2" body="${3:-}"
  local args=(-sS -X "$method" "$URL$path"
              -H "Authorization: Bearer $ADMIN_PAT"
              -H 'Content-Type: application/json'
              -H 'Accept: application/json')
  [[ -n "$body" ]] && args+=(--data "$body")
  curl "${args[@]}"
}

api_ok() { # api_ok <method> <path> [json-body] -- dies on a Zitadel error envelope
  local out
  out="$(api "$@")"
  if jq -e 'type == "object" and has("code") and has("message")' >/dev/null 2>&1 <<<"$out"; then
    die "$1 $2 failed: $(jq -c . <<<"$out")"
  fi
  printf '%s' "$out"
}

jwt_payload() { # jwt_payload <jwt> -- prints the decoded payload as JSON
  local seg pad
  seg="${1#*.}"; seg="${seg%%.*}"
  pad=$(( (4 - ${#seg} % 4) % 4 ))
  while (( pad-- > 0 )); do seg="${seg}="; done
  printf '%s' "$seg" | tr '_-' '/+' | base64 -d 2>/dev/null
}

# ---------------------------------------------------------------- cleanup

CLEAN_EXECUTIONS=0
TARGET_ID=""
PROJECT_ID=""
IMPERSONATOR_ID=""
ENDUSER_ID=""
WEBHOOK_TOKEN=""
SINK_PID=""
SINK_DIR=""
IMPERSONATION_WAS=""

cleanup() {
  local code=$?
  set +e

  if [[ -n "$SINK_PID" ]]; then kill "$SINK_PID" >/dev/null 2>&1; wait "$SINK_PID" 2>/dev/null; fi

  if [[ $KEEP -eq 1 ]]; then
    step "Leaving resources in place (--keep)"
    info "project:      ${PROJECT_ID:-–}"
    info "impersonator: ${IMPERSONATOR_ID:-–}"
    info "end user:     ${ENDUSER_ID:-–}"
    info "target:       ${TARGET_ID:-–}"
    [[ -n "$WEBHOOK_TOKEN" ]] && info "webhook.site: https://webhook.site/#!/view/$WEBHOOK_TOKEN"
    [[ -n "$SINK_DIR" ]] && info "sink log:     $SINK_DIR/payloads.jsonl"
    exit "$code"
  fi

  step "Cleaning up"
  if [[ $CLEAN_EXECUTIONS -eq 1 ]]; then
    # There is no DeleteExecution endpoint; setting an empty target list removes it.
    for fn in preuserinfo preaccesstoken; do
      api PUT /v2/actions/executions "{\"condition\":{\"function\":{\"name\":\"$fn\"}},\"targets\":[]}" >/dev/null
    done
    info "removed executions"
  fi
  [[ -n "$TARGET_ID"       ]] && api DELETE "/v2/actions/targets/$TARGET_ID"        >/dev/null && info "removed target"
  [[ -n "$PROJECT_ID"      ]] && api DELETE "/management/v1/projects/$PROJECT_ID"   >/dev/null && info "removed project and client"
  [[ -n "$IMPERSONATOR_ID" ]] && api DELETE "/v2/users/$IMPERSONATOR_ID"            >/dev/null && info "removed impersonator"
  [[ -n "$ENDUSER_ID"      ]] && api DELETE "/v2/users/$ENDUSER_ID"                 >/dev/null && info "removed end user"
  if [[ -n "$IMPERSONATION_WAS" ]]; then
    api PUT /admin/v1/policies/security "{\"enableImpersonation\":$IMPERSONATION_WAS}" >/dev/null
    info "restored enableImpersonation=$IMPERSONATION_WAS"
  fi
  if [[ -n "$WEBHOOK_TOKEN" ]]; then
    curl -sS -X DELETE "https://webhook.site/token/$WEBHOOK_TOKEN" >/dev/null 2>&1
    info "removed webhook.site token"
  fi
  [[ -n "$SINK_DIR" ]] && rm -rf "$SINK_DIR"

  exit "$code"
}
trap cleanup EXIT

# ---------------------------------------------------------------- preflight

step "Preflight"

whoami_json="$(api GET /auth/v1/users/me)"
if ! jq -e '.user.id' >/dev/null 2>&1 <<<"$whoami_json"; then
  die "the PAT in '$PAT_FILE' was rejected by $URL: $(jq -c . <<<"$whoami_json" 2>/dev/null || printf '%s' "$whoami_json")
Check that the file holds exactly the token, with no surrounding quotes, backticks or extra lines."
fi
info "authenticated as $(jq -r '.user.userName' <<<"$whoami_json")"

ISSUER="$(curl -sS "$URL/.well-known/openid-configuration" | jq -r '.issuer')"
[[ -n "$ISSUER" && "$ISSUER" != "null" ]] || die "could not read the issuer from $URL/.well-known/openid-configuration"
info "issuer $ISSUER"

features="$(api GET /v2beta/features/instance)"
if [[ "$(jq -r '.oidcTokenExchange.enabled // false' <<<"$features")" != "true" ]]; then
  die "the oidcTokenExchange instance feature is disabled. Enable it with:
  curl -X PUT $URL/v2beta/features/instance \\
    -H \"Authorization: Bearer \\\$(cat $PAT_FILE)\" -H 'Content-Type: application/json' \\
    -d '{\"oidcTokenExchange\": true}'"
fi
info "oidcTokenExchange feature enabled"

# ---------------------------------------------------------------- receiver

step "Setting up the webhook receiver ($RECEIVER)"

if [[ "$RECEIVER" == "local" ]]; then
  SINK_DIR="$(mktemp -d)"
  python3 -u -c '
import http.server, sys, threading
out = sys.argv[1]
lock = threading.Lock()

class Handler(http.server.BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def do_POST(self):
        length = int(self.headers.get("Content-Length") or 0)
        body = self.rfile.read(length)
        with lock, open(out, "ab") as fh:
            fh.write(body.replace(b"\n", b"") + b"\n")
        self.send_response(200)
        self.send_header("Content-Length", "0")
        self.end_headers()

    def log_message(self, *args):
        pass

server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), Handler)
print(server.server_address[1], flush=True)
server.serve_forever()
' "$SINK_DIR/payloads.jsonl" > "$SINK_DIR/port" &
  SINK_PID=$!
  for _ in $(seq 1 50); do
    [[ -s "$SINK_DIR/port" ]] && break
    sleep 0.1
  done
  SINK_PORT="$(cat "$SINK_DIR/port")"
  [[ -n "$SINK_PORT" ]] || die "the local sink failed to start"
  ENDPOINT="http://127.0.0.1:$SINK_PORT"
  receiver_payloads() { [[ -f "$SINK_DIR/payloads.jsonl" ]] && cat "$SINK_DIR/payloads.jsonl" || true; }
else
  token_json="$(curl -sS -X POST https://webhook.site/token -H 'Content-Type: application/json' -d '{}')"
  WEBHOOK_TOKEN="$(jq -r '.uuid // empty' <<<"$token_json")"
  [[ -n "$WEBHOOK_TOKEN" ]] || die "could not create a webhook.site token (rate limited?): $token_json
Retry in a minute or run with --local."
  ENDPOINT="https://webhook.site/$WEBHOOK_TOKEN"
  receiver_payloads() {
    curl -sS "https://webhook.site/token/$WEBHOOK_TOKEN/requests?sorting=oldest&per_page=50" \
      | jq -r '.data[]?.content // empty'
  }
  info "inbox https://webhook.site/#!/view/$WEBHOOK_TOKEN"
fi
info "endpoint $ENDPOINT"

# ---------------------------------------------------------------- setup

step "Configuring the instance"

IMPERSONATION_WAS="$(api GET /admin/v1/policies/security | jq -r '.policy.enableImpersonation // false')"
api_ok PUT /admin/v1/policies/security '{"enableImpersonation":true}' >/dev/null
info "enableImpersonation true (was $IMPERSONATION_WAS)"

IMPERSONATOR_NAME="actor-test-impersonator-$RUN_ID"
IMPERSONATOR_ID="$(api_ok POST /management/v1/users/machine \
  "{\"userName\":\"$IMPERSONATOR_NAME\",\"name\":\"actor test impersonator\",\"description\":\"created by skills/test-actor-in-action-v2.sh\",\"accessTokenType\":\"ACCESS_TOKEN_TYPE_BEARER\"}" \
  | jq -r '.userId')"
info "impersonator $IMPERSONATOR_ID ($IMPERSONATOR_NAME)"

IMPERSONATOR_PAT="$(api_ok POST "/management/v1/users/$IMPERSONATOR_ID/pats" '{}' | jq -r '.token')"
[[ -n "$IMPERSONATOR_PAT" && "$IMPERSONATOR_PAT" != "null" ]] || die "could not create a PAT for the impersonator"

# IAM_END_USER_IMPERSONATOR carries the "impersonation" permission that
# CreateOIDCSession checks (internal/command/oidc_session.go).
api_ok POST /admin/v1/members \
  "{\"userId\":\"$IMPERSONATOR_ID\",\"roles\":[\"IAM_END_USER_IMPERSONATOR\"]}" >/dev/null
info "granted IAM_END_USER_IMPERSONATOR"

ENDUSER_NAME="actor-test-enduser-$RUN_ID"
ENDUSER_ID="$(api_ok POST /v2/users/human "$(jq -n --arg u "$ENDUSER_NAME" '{
  username: $u,
  profile: {givenName: "Actor", familyName: "Test"},
  email: {email: ($u + "@example.com"), isVerified: true}
}')" | jq -r '.userId')"
info "end user $ENDUSER_ID ($ENDUSER_NAME)"

# The client must not live in the ZITADEL project: impersonated tokens are
# rejected on the Zitadel API itself.
PROJECT_ID="$(api_ok POST /management/v1/projects "{\"name\":\"actor-test-$RUN_ID\"}" | jq -r '.id')"
app_json="$(api_ok POST "/management/v1/projects/$PROJECT_ID/apps/oidc" '{
  "name": "actor-test-client",
  "redirectUris": ["https://example.com/callback"],
  "responseTypes": ["OIDC_RESPONSE_TYPE_CODE"],
  "grantTypes": ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE", "OIDC_GRANT_TYPE_TOKEN_EXCHANGE"],
  "appType": "OIDC_APP_TYPE_WEB",
  "authMethodType": "OIDC_AUTH_METHOD_TYPE_BASIC",
  "accessTokenType": "OIDC_TOKEN_TYPE_BEARER",
  "devMode": true
}')"
CLIENT_ID="$(jq -r '.clientId' <<<"$app_json")"
CLIENT_SECRET="$(jq -r '.clientSecret' <<<"$app_json")"
info "project $PROJECT_ID, client $CLIENT_ID"

target_json="$(api POST /v2/actions/targets "$(jq -n --arg e "$ENDPOINT" --arg n "actor-test-$RUN_ID" '{
  name: $n, restWebhook: {interruptOnError: false}, timeout: "10s", endpoint: $e
}')")"
if [[ "$(jq -r '.message // empty' <<<"$target_json")" == *DeniedURL* ]]; then
  if [[ "$RECEIVER" == "local" ]]; then
    die "the instance refuses $ENDPOINT as a target.

Zitadel's SSRF guard denies loopback and private ranges by default
(HTTPClient.DenyList in cmd/defaults.yaml). To use --local, restart the
instance with the deny list relaxed, for example:

  ZITADEL_HTTPCLIENT_DENYLIST= zitadel start-from-init ...

Otherwise drop --local and use the webhook.site receiver."
  fi
  die "the instance refuses $ENDPOINT as a target: $(jq -c . <<<"$target_json")"
fi
TARGET_ID="$(jq -r '.id // empty' <<<"$target_json")"
[[ -n "$TARGET_ID" ]] || die "could not create the target: $(jq -c . <<<"$target_json")"
info "target $TARGET_ID"

for fn in preuserinfo preaccesstoken; do
  api_ok PUT /v2/actions/executions \
    "{\"condition\":{\"function\":{\"name\":\"$fn\"}},\"targets\":[\"$TARGET_ID\"]}" >/dev/null
done
CLEAN_EXECUTIONS=1
info "executions function/preuserinfo and function/preaccesstoken -> target"

# ---------------------------------------------------------------- exercise

step "Impersonating $ENDUSER_ID via token exchange"

exchange="$(curl -sS -X POST "$URL/oauth/v2/token" \
  -H 'Content-Type: application/x-www-form-urlencoded' -H 'Accept: application/json' \
  -u "$CLIENT_ID:$CLIENT_SECRET" \
  --data-urlencode 'grant_type=urn:ietf:params:oauth:grant-type:token-exchange' \
  --data-urlencode "subject_token=$ENDUSER_ID" \
  --data-urlencode 'subject_token_type=urn:zitadel:params:oauth:token-type:user_id' \
  --data-urlencode "actor_token=$IMPERSONATOR_PAT" \
  --data-urlencode 'actor_token_type=urn:ietf:params:oauth:token-type:access_token' \
  --data-urlencode 'requested_token_type=urn:ietf:params:oauth:token-type:jwt' \
  --data-urlencode 'scope=openid profile email')"

if [[ "$(jq -r '.error // empty' <<<"$exchange")" != "" ]]; then
  die "token exchange failed: $(jq -c . <<<"$exchange")
See the troubleshooting section of skills/test-actor-in-action-v2.md."
fi

ACCESS_TOKEN="$(jq -r '.access_token' <<<"$exchange")"
ID_TOKEN="$(jq -r '.id_token // empty' <<<"$exchange")"
at_claims="$(jwt_payload "$ACCESS_TOKEN")"
info "issued_token_type $(jq -r '.issued_token_type' <<<"$exchange")"

step "Calling the userinfo endpoint with the impersonated token"
userinfo="$(curl -sS "$URL/oidc/v1/userinfo" -H "Authorization: Bearer $ACCESS_TOKEN")"
info "sub $(jq -r '.sub // "none"' <<<"$userinfo")"

step "Negative control: userinfo without impersonation"
CONTROL_RAN=0
control_secret="$(api POST "/v2/users/$IMPERSONATOR_ID/secret" '{}' | jq -r '.clientSecret // empty')"
if [[ -n "$control_secret" ]]; then
  # client_credentials resolves the client by login name, not by user id
  # (see clientCredentialsAuth in internal/api/oidc/client_credentials.go).
  control_login="$(api GET "/v2/users/$IMPERSONATOR_ID" | jq -r '.user.loginNames[0] // empty')"
  control_token="$(curl -sS -X POST "$URL/oauth/v2/token" \
    -H 'Content-Type: application/x-www-form-urlencoded' -H 'Accept: application/json' \
    -u "$control_login:$control_secret" \
    --data-urlencode 'grant_type=client_credentials' \
    --data-urlencode "scope=openid urn:zitadel:iam:org:project:id:$PROJECT_ID:aud" \
    | jq -r '.access_token // empty')"
  if [[ -n "$control_token" ]]; then
    control_userinfo="$(curl -sS "$URL/oidc/v1/userinfo" -H "Authorization: Bearer $control_token")"
    if [[ "$(jq -r '.sub // empty' <<<"$control_userinfo")" == "$IMPERSONATOR_ID" ]]; then
      CONTROL_RAN=1
      info "control userinfo returned the impersonator's own sub"
    else
      warn "control userinfo call did not return the expected subject: $control_userinfo"
    fi
  else
    warn "could not obtain a client_credentials token for the control"
  fi
else
  warn "could not create a machine secret for the control"
fi

# ---------------------------------------------------------------- collect

step "Waiting for webhook payloads"

# Count payloads per function that carry the expected actor. The JSON keys come
# from domain.TokenActor, so they are snake_case: user_id / issuer.
count_with_actor() { # count_with_actor <function> <payloads>
  jq -s --arg fn "$1" --arg uid "$IMPERSONATOR_ID" --arg iss "$ISSUER" \
    '[.[] | select(.function == $fn and .actor.user_id == $uid and .actor.issuer == $iss)] | length' \
    <<<"$2"
}

count_without_actor() { # count_without_actor <payloads>
  jq -s --arg uid "$IMPERSONATOR_ID" \
    '[.[] | select(.function == "function/preuserinfo" and (has("actor") | not) and .userinfo.sub == $uid)] | length' \
    <<<"$1"
}

# Poll for the payloads we actually assert on rather than for a raw count: the
# negative control emits payloads too, so a total would be satisfied too early.
payloads=""
for i in $(seq 1 30); do
  payloads="$(receiver_payloads)"
  if [[ -n "$payloads" ]]; then
    [[ "$(count_with_actor function/preaccesstoken "$payloads")" -ge 1 &&
       "$(count_with_actor function/preuserinfo "$payloads")" -ge 2 &&
       ( $CONTROL_RAN -eq 0 || "$(count_without_actor "$payloads")" -ge 1 ) ]] && break
  fi
  sleep 1
done
count="$(grep -c . <<<"$payloads" || true)"
info "received $count payload(s) after ${i}s"

if [[ "$count" -eq 0 ]]; then
  die "no payloads reached the receiver. Is the instance able to reach $ENDPOINT?"
fi

# ---------------------------------------------------------------- assertions

step "Results"

check_eq "access token sub is the impersonated user" \
  "$ENDUSER_ID" "$(jq -r '.sub // "none"' <<<"$at_claims")"
check_eq "access token act.sub is the impersonator" \
  "$IMPERSONATOR_ID" "$(jq -r '.act.sub // "none"' <<<"$at_claims")"
check_eq "access token act.iss is the issuer" \
  "$ISSUER" "$(jq -r '.act.iss // "none"' <<<"$at_claims")"

if [[ -n "$ID_TOKEN" ]]; then
  id_claims="$(jwt_payload "$ID_TOKEN")"
  check_eq "id token sub is the impersonated user" \
    "$ENDUSER_ID" "$(jq -r '.sub // "none"' <<<"$id_claims")"
  check_eq "id token act.sub is the impersonator" \
    "$IMPERSONATOR_ID" "$(jq -r '.act.sub // "none"' <<<"$id_claims")"
else
  fail "the token exchange response carried no id_token"
fi

check_eq "userinfo sub is the impersonated user" \
  "$ENDUSER_ID" "$(jq -r '.sub // "none"' <<<"$userinfo")"

check_eq "function/preaccesstoken payload carries the actor" \
  "true" "$([[ "$(count_with_actor function/preaccesstoken "$payloads")" -ge 1 ]] && echo true || echo false)"
check_eq "both function/preuserinfo payloads carry the actor" \
  "true" "$([[ "$(count_with_actor function/preuserinfo "$payloads")" -ge 2 ]] && echo true || echo false)"

if [[ $CONTROL_RAN -eq 1 ]]; then
  check_eq "non-impersonated payload omits the actor key entirely" \
    "true" "$([[ "$(count_without_actor "$payloads")" -ge 1 ]] && echo true || echo false)"
else
  skip "negative control (could not mint a non-impersonated token)"
fi

if [[ $FAILURES -gt 0 ]]; then
  step "Payloads received"
  jq -c '{function, sub: .userinfo.sub, actor}' <<<"$payloads" || printf '%s\n' "$payloads"
  printf '\n%sFAILED%s: %d assertion(s)\n' "$RED" "$RESET" "$FAILURES"
  exit 1
fi

printf '\n%sAll assertions passed.%s\n' "$GREEN" "$RESET"
