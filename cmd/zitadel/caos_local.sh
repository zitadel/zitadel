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

# Notifications
export TWILIO_SERVICE_SID=AC983b9e61d074abae7d3be7b56ecbc697
export TWILIO_TOKEN=e1748c2a4514473c336e08f78a81d3d4
export TWILIO_SENDER_NAME=CAOS AG
export SMTP_HOST=smtp.gmail.com:465
export SMTP_USER=zitadel@caos.ch
export SMTP_PASSWORD=yvubsweihyamlcwo
export EMAIL_SENDER_ADDRESS=noreply@caos.ch
export EMAIL_SENDER_NAME=CAOS AG
export SMTP_TLS=TRUE
export CHAT_URL=https://chat.googleapis.com/v1/spaces/AAAAdl31ujI/messages?key\=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token\=cDv8MZv4MPeeCgbFo66U5bDJ9R73UI9A_pdd1sQHdPQ%3D
