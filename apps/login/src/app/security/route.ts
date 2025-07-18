import { createServiceForHost } from "@/lib/service";
import { getServiceUrlFromHeaders } from "@/lib/service-url";
import { Client } from "@zitadel/client";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { headers } from "next/headers";
import { NextResponse } from "next/server";

export async function GET() {
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const settingsService: Client<typeof SettingsService> =
    await createServiceForHost(SettingsService, serviceUrl);

  const settings = await settingsService
    .getSecuritySettings({})
    .then((resp) => (resp.settings ? resp.settings : undefined));

  const response = NextResponse.json({ settings }, { status: 200 });

  // Add Cache-Control header to cache the response for up to 1 hour
  response.headers.set(
    "Cache-Control",
    "public, max-age=3600, stale-while-revalidate=86400",
  );

  return response;
}
