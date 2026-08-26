#!/usr/bin/env bash
# Full publication sweep: all 11 targets that v4 has pages for.
# Each run is converted + published immediately and its raw CSV dropped.
# Clears the .hold gate first. ~35min per target, ~6.5h total.
cd "$(dirname "$0")"
rm -f .hold
for t in add_session human_password_login introspect machine_client_credentials_login \
         machine_jwt_profile_grant machine_pat_login manipulate_user oidc_session \
         otp_session password_session user_info; do
  echo "=== QUEUE start ${t} $(date -u '+%Y-%m-%d %H:%M:%S UTC') ==="
  if ./run-bench.sh "$t" 600 1800s; then
    ./postprocess.sh "$t" v4.17.1 || echo "=== QUEUE postprocess FAILED ${t} ==="
  else
    echo "=== QUEUE run FAILED ${t} ==="
  fi
  echo "=== QUEUE done ${t} $(date -u '+%Y-%m-%d %H:%M:%S UTC') ==="
done
echo "=== QUEUE ALL 11 DONE ==="
