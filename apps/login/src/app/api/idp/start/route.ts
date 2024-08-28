import { startIdentityProviderFlow } from "@/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";
import { toJson } from "@zitadel/client";
import { StartIdentityProviderIntentResponseSchema } from "@zitadel/proto/zitadel/user/v2/user_service_pb";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    let { idpId, successUrl, failureUrl } = body;

    return startIdentityProviderFlow({
      idpId,
      urls: {
        successUrl,
        failureUrl,
      },
    })
      .then((resp) => {
        return NextResponse.json(
          toJson(StartIdentityProviderIntentResponseSchema, resp),
        );
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.json({}, { status: 400 });
  }
}
