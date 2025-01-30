import { newSystemToken } from "@zitadel/client/node";

export async function systemAPIToken({
  serviceRegion,
}: {
  serviceRegion: string;
}) {
  const REGIONS = ["eu1", "us1"].map((region) => {
    return {
      id: region,
      audience: process.env[region + "_AUDIENCE"],
      userID: process.env[region + "_SYSTEM_USER_ID"],
      token: Buffer.from(
        process.env[
          region.toUpperCase() + "_SYSTEM_USER_PRIVATE_KEY"
        ] as string,
        "base64",
      ).toString("utf-8"),
    };
  });

  const region = REGIONS.find((region) => region.id === serviceRegion);

  if (!region || !region.audience || !region.userID || !region.token) {
    throw new Error("Invalid region");
  }

  const token = newSystemToken({
    audience: region.audience,
    subject: region.userID,
    key: region.token,
  });

  return token;
}
