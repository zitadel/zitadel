import "server-only";

import { getUserByID, ServiceConfig } from "@/lib/zitadel";

export async function getUserAvatarUrl({
  serviceConfig,
  userId,
}: {
  serviceConfig: ServiceConfig;
  userId?: string;
}): Promise<string | undefined> {
  if (!userId) return undefined;

  try {
    const { user } = await getUserByID({ serviceConfig, userId });
    return user?.type.case === "human" ? user.type.value.profile?.avatarUrl : undefined;
  } catch {
    // A missing or unavailable profile must not prevent authentication.
    return undefined;
  }
}
