import { createZitadelApiClient } from "@zitadel/nextjs/api";

type SignupPasswordBody = {
  email?: unknown;
  password?: unknown;
  givenName?: unknown;
  familyName?: unknown;
  username?: unknown;
  organizationId?: unknown;
};

function toOptionalString(value: unknown) {
  if (typeof value !== "string") {
    return undefined;
  }

  const trimmed = value.trim();
  return trimmed.length > 0 ? trimmed : undefined;
}

export async function POST(request: Request) {
  let body: SignupPasswordBody;
  try {
    body = (await request.json()) as SignupPasswordBody;
  } catch {
    return Response.json(
      {
        ok: false,
        endpoint: "/api/demo/signup/password",
        method: "password",
        error: { message: "Request body must be valid JSON." },
      },
      { status: 400 },
    );
  }

  const email = toOptionalString(body.email);
  const password = typeof body.password === "string" ? body.password : undefined;
  const givenName = toOptionalString(body.givenName);
  const familyName = toOptionalString(body.familyName);
  const username = toOptionalString(body.username);
  const organizationId = toOptionalString(body.organizationId);

  const requestSummary = {
    email,
    username,
    givenName,
    familyName,
    organizationId,
    hasPassword: Boolean(password),
  };

  const missingFields = [
    email ? null : "email",
    password ? null : "password",
    givenName ? null : "givenName",
    familyName ? null : "familyName",
  ].filter((field): field is string => Boolean(field));

  if (missingFields.length > 0) {
    return Response.json(
      {
        ok: false,
        endpoint: "/api/demo/signup/password",
        method: "password",
        request: requestSummary,
        error: { message: `Missing required field(s): ${missingFields.join(", ")}` },
      },
      { status: 400 },
    );
  }

  try {
    const api = await createZitadelApiClient();
    const result = await api.userService.addHumanUser({
      email: {
        email,
        verification: {
          case: "isVerified",
          value: false,
        },
      },
      username: username ?? email,
      profile: {
        givenName,
        familyName,
      },
      passwordType: {
        case: "password",
        value: {
          password,
          changeRequired: false,
        },
      },
      ...(organizationId
        ? {
            organization: {
              org: {
                case: "orgId",
                value: organizationId,
              },
            },
          }
        : {}),
    });

    return Response.json(
      {
        ok: true,
        endpoint: "/api/demo/signup/password",
        method: "password",
        request: requestSummary,
        result: {
          userId: result.userId,
          details: result.details,
          emailCode: result.emailCode,
          phoneCode: result.phoneCode,
        },
      },
      { status: 201 },
    );
  } catch (error) {
    const message = error instanceof Error ? error.message : "Unexpected signup API error.";
    const code =
      typeof error === "object" && error && "code" in error
        ? String((error as { code: unknown }).code)
        : undefined;
    const rawMessage =
      typeof error === "object" && error && "rawMessage" in error
        ? String((error as { rawMessage: unknown }).rawMessage)
        : undefined;
    const status = code === "6" ? 409 : 500;

    return Response.json(
      {
        ok: false,
        endpoint: "/api/demo/signup/password",
        method: "password",
        request: requestSummary,
        error: {
          message,
          code,
          rawMessage,
        },
      },
      { status },
    );
  }
}
