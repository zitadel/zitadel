import { server, startIdentityProviderFlow } from "@/lib/zitadel";
import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const body = await request.json();
  if (body) {
    let { idpId, successUrl, failureUrl } = body;

    return startIdentityProviderFlow(server, {
      idpId,
      urls: {
        successUrl,
        failureUrl,
      },
    })
      .then((resp) => {
        return NextResponse.json(resp);
      })
      .catch((error) => {
        return NextResponse.json(error, { status: 500 });
      });
  } else {
    return NextResponse.json({}, { status: 400 });
  }
}
