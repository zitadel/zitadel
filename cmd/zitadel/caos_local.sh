BASEDIR=$(dirname "$0")

gopass sync --store zitadel-secrets

# Tracing
gopass zitadel-secrets/zitadel/developer/default/zitadel-svc-account-eventstore-local | base64 -D > "$BASEDIR/local_svc-account-tracing.json"
export GOOGLE_APPLICATION_CREDENTIALS="$BASEDIR/local_svc-account-tracing.json"

export ZITADEL_TRACING_PROJECT_ID=caos-citadel-test
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
export ZITADEL_OIDC_KEYS_ID=OIDCKey_1
export ZITADEL_COOKIE_KEY=CookieKey_1

# Notifications
export DEBUG_MODE=TRUE
export TWILIO_SERVICE_SID=$(gopass zitadel-secrets/zitadel/dev/twilio-sid)
export TWILIO_TOKEN=$(gopass zitadel-secrets/zitadel/dev/twilio-auth-token)
export TWILIO_SENDER_NAME=CAOS AG
export SMTP_HOST=smtp.gmail.com:465
export SMTP_USER=zitadel@caos.ch
export SMTP_PASSWORD=$(gopass zitadel-secrets/zitadel/google/emailappkey)
export EMAIL_SENDER_ADDRESS=noreply@caos.ch
export EMAIL_SENDER_NAME=CAOS AG
export SMTP_TLS=TRUE
export CHAT_URL=$(gopass zitadel-secrets/zitadel/dev/google-chat-url | base64 -D)

#OIDC
export ZITADEL_ISSUER=http://localhost:50022
export ZITADEL_ACCOUNTS=http://localhost:50031
export ZITADEL_AUTHORIZE=http://localhost:50022
export ZITADEL_OAUTH=http://localhost:50022
export ZITADEL_CONSOLE=http://localhost:4200
export CAOS_OIDC_DEV=true
export ZITADEL_COOKIE_DOMAIN=localhost

#Console
export ZITADEL_CONSOLE_ENV_DIR=../../console/src/assets/