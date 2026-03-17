#!/bin/sh
set -e

# Create zitadel system user if it doesn't exist
if ! id -u zitadel >/dev/null 2>&1; then
    useradd --system --no-create-home --shell /usr/sbin/nologin zitadel
fi

# Set ownership
chown -R zitadel:zitadel /etc/zitadel
chown -R zitadel:zitadel /opt/zitadel

# Reload systemd
systemctl daemon-reload

echo ""
echo "ZITADEL installed successfully."
echo ""
echo "Next steps:"
echo "  1. Configure /etc/zitadel/defaults.yaml"
echo "  2. Set ZITADEL_MASTERKEY in /etc/systemd/system/zitadel.service"
echo "  3. Configure a reverse proxy (caddy/nginx/traefik) for path-based routing"
echo "  4. systemctl enable --now zitadel zitadel-login"
echo ""
