import { createCallback } from "@zitadel/nextjs/auth/session";
import { errorMessage, jsonResponse } from "../utils";

interface CreateCallbackBody {
  authRequestId?: string;
  sessionId?: string;
  sessionToken?: string;
}

export async function POST(request: Request) {
  const body = (await request.json()) as CreateCallbackBody;
  const authRequestId = body.authRequestId?.trim();
  const sessionId = body.sessionId?.trim();
  const sessionToken = body.sessionToken?.trim();

  if (!authRequestId || !sessionId || !sessionToken) {
    return jsonResponse(
      {
        error: "Missing required fields: authRequestId, sessionId, sessionToken",
      },
      400,
    );
  }

  try {
    const callbackUrl = await createCallback({
      authRequestId,
      sessionId,
      sessionToken,
    });

    return jsonResponse({
      callbackUrl,
    });
  } catch (error) {
    return jsonResponse(
      {
        error: "Failed to create callback",
        details: errorMessage(error),
      },
      400,
    );
  }
}
