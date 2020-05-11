BASEDIR=$(dirname "$0")

# Tracing
gopass zitadel-secrets/zitadel/developer/default/zitadel-svc-account-eventstore-local | base64 -D > "$BASEDIR/local_svc-account-tracing.json"
export GOOGLE_APPLICATION_CREDENTIALS="$BASEDIR/local_svc-account-tracing.json"

export ZITADEL_TRACING_PROJECT_ID=caos-zitadel-test
export ZITADEL_TRACING_FRACTION=0.1

# Log
export ZITADEL_LOG_LEVEL=debug

# Cockroach
export ZITADEL_EVENTSTORE_HOST=localhost
export ZITADEL_EVENTSTORE_PORT=26257

# Keys
gopass zitadel-secrets/zitadel/developer/default/keys.yaml > "$BASEDIR/local_keys.yaml"
export ZITADEL_KEY_PATH="$BASEDIR/local_keys.yaml"

export ZITADEL_USER_VERIFICATION_KEY=UserVerificationKey_1
export ZITADEL_OTP_VERIFICATION_KEY=OTPVerificationKey_1
