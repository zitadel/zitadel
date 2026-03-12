import { createServiceForHost } from "@/lib/service";
import { createLogger } from "@/lib/logger";
import { Client } from "@zitadel/client";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { NextResponse } from "next/server";

const logger = createLogger("readiness");

export async function GET() {
  if (!process.env.ZITADEL_API_URL) {
    return new NextResponse("Service unavailable", {
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
    logger.error("Readiness check failed", { error: e });
    return new NextResponse("Service unavailable", {
      status: 503,
      headers: { "Content-Type": "text/plain" },
    });
  }
}
