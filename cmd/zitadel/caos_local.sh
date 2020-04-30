# Tracing
gopass citadel-secrets/citadel/developer/default/citadel-svc-account-eventstore-local | base64 -D > local_svc-account-tracing.json
export GOOGLE_APPLICATION_CREDENTIALS="$(pwd)/local_svc-account-tracing.json"

export ZITADEL_TRACING_PROJECT_ID=caos-citadel-test
export ZITADEL_TRACING_FRACTION=0.1

# Log
export ZITADEL_LOG_LEVEL=debug

# Cockroach
export ZITADEL_EVENTSTORE_HOST=localhost
export ZITADEL_EVENTSTORE_PORT=26257

# Keys
gopass citadel-secrets/citadel/developer/default/keys.yaml > local_keys.yaml
export ZITADEL_KEY_PATH="$(pwd)/local_keys.yaml"

export ZITADEL_USER_VERIFICATION_KEY=UserVerificationKey_1
export ZITADEL_OTP_VERIFICATION_KEY=OTPVerificationKey_1