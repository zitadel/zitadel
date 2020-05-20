BASEDIR=$(dirname "$0")

# Tracing
gopass citadel-secrets/citadel/developer/default/citadel-svc-account-eventstore-local | base64 -D > "$BASEDIR/local_svc-account-tracing.json"
export GOOGLE_APPLICATION_CREDENTIALS="$BASEDIR/local_svc-account-tracing.json"

export ZITADEL_TRACING_PROJECT_ID=caos-citadel-test
export ZITADEL_TRACING_FRACTION=0.1

# Log
export ZITADEL_LOG_LEVEL=debug

# Cockroach
export ZITADEL_EVENTSTORE_HOST=localhost
export ZITADEL_EVENTSTORE_PORT=26257

# Keys
gopass citadel-secrets/citadel/developer/default/keys.yaml > "$BASEDIR/local_keys.yaml"
export ZITADEL_KEY_PATH="$BASEDIR/local_keys.yaml"

export ZITADEL_USER_VERIFICATION_KEY=UserVerificationKey_1
export ZITADEL_OTP_VERIFICATION_KEY=OTPVerificationKey_1
export ZITADEL_OIDC_KEYS_ID=OIDCKey_1

#OIDC
export ZITADEL_ISSUER=http://localhost:50022
export ZITADEL_ACCOUNTS=http://localhost:50031
export ZITADEL_AUTHORIZE=http://localhost:50022
export ZITADEL_OAUTH=http://localhost:50022
export ZITADEL_CONSOLE=http://localhost:4200
export CAOS_OIDC_DEV=true