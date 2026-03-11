import { createSession, getSession } from "@zitadel/nextjs/auth/session";
import { errorMessage, jsonResponse } from "../utils";

interface CreateSessionBody {
  loginName?: string;
  password?: string;
}

export async function POST(request: Request) {
  const body = (await request.json()) as CreateSessionBody;
  const loginName = body.loginName?.trim();
  const password = body.password;

  if (!loginName || !password) {
    return jsonResponse(
      {
        error: "Missing required fields: loginName and password",
      },
      400,
    );
  }

  try {
    const session = await createSession({
      checks: {
        user: {
          search: {
            case: "loginName",
            value: loginName,
          },
        },
        password: {
          password,
        },
      },
    } as Parameters<typeof createSession>[0]);

    return jsonResponse({
      sessionId: session.sessionId,
      sessionToken: session.sessionToken,
      details: session.details,
      challenges: session.challenges,
    });
  } catch (error) {
    return jsonResponse(
      {
        error: "Failed to create session",
        details: errorMessage(error),
      },
      400,
    );
  }
}

export async function GET(request: Request) {
  const url = new URL(request.url);
  const sessionId = url.searchParams.get("sessionId")?.trim();
  const sessionToken = url.searchParams.get("sessionToken")?.trim() || undefined;

  if (!sessionId) {
    return jsonResponse(
      {
        error: "Missing required query parameter: sessionId",
      },
      400,
    );
  }

  try {
    const response = await getSession({ sessionId, sessionToken });
    return jsonResponse({
      session: response.session,
    });
  } catch (error) {
    return jsonResponse(
      {
        error: "Failed to load session",
        details: errorMessage(error),
      },
      400,
    );
  }
}
