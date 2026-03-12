import { createServiceForHost } from "@/lib/service";
import { Client } from "@zitadel/client";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { NextResponse } from "next/server";

export async function GET() {
  if (!process.env.ZITADEL_API_URL) {
    return new NextResponse("ZITADEL_API_URL is not set", {
      status: 503,
      headers: { "Content-Type": "text/plain" },
    });
  }

  try {
    const settingsService: Client<typeof SettingsService> = await createServiceForHost(SettingsService, {
      baseUrl: process.env.ZITADEL_API_URL,
    });
    await settingsService.getGeneralSettings({});
    return new NextResponse("OK", {
      status: 200,
      headers: { "Content-Type": "text/plain" },
    });
  } catch (e) {
    return new NextResponse(e instanceof Error ? e.message : String(e), {
      status: 503,
      headers: { "Content-Type": "text/plain" },
    });
  }
}
