import { createZitadelApiClient } from "@zitadel/nextjs/api";

interface RegisterOrganizationBody {
  name?: string;
  organizationId?: string;
}

function normalizeGrpcCode(code: unknown): string | undefined {
  if (typeof code === "number") {
    return String(code);
  }

  if (typeof code === "string" && code.trim()) {
    return code.trim();
  }

  return undefined;
}

function mapGrpcCodeToHttpStatus(code: string | undefined): number {
  switch (code?.toUpperCase()) {
    case "3":
    case "INVALID_ARGUMENT":
      return 400;
    case "5":
    case "NOT_FOUND":
      return 404;
    case "6":
    case "ALREADY_EXISTS":
      return 409;
    case "7":
    case "PERMISSION_DENIED":
      return 403;
    case "16":
    case "UNAUTHENTICATED":
      return 401;
    case "14":
    case "UNAVAILABLE":
      return 503;
    default:
      return 500;
  }
}

function extractErrorMessage(error: unknown): string {
  if (error && typeof error === "object") {
    const maybeRawMessage = (error as { rawMessage?: unknown }).rawMessage;
    if (typeof maybeRawMessage === "string" && maybeRawMessage) {
      return maybeRawMessage;
    }
  }

  if (error instanceof Error && error.message) {
    return error.message;
  }

  return "Unknown error";
}

function createErrorResponse(error: unknown): { status: number; body: Record<string, unknown> } {
  const code =
    error && typeof error === "object"
      ? normalizeGrpcCode((error as { code?: unknown }).code)
      : undefined;
  const status = mapGrpcCodeToHttpStatus(code);

  return {
    status,
    body: {
      error: "Failed to register organization",
      message: extractErrorMessage(error),
      code: code ?? null,
      permissionHint:
        status === 403
          ? "Missing org.create permission. Use a token/session with org.create on the instance."
          : null,
    },
  };
}

export async function POST(request: Request) {
  let body: RegisterOrganizationBody;
  try {
    body = (await request.json()) as RegisterOrganizationBody;
  } catch {
    return Response.json({ error: "Invalid JSON body" }, { status: 400 });
  }

  const name = body.name?.trim();
  if (!name) {
    return Response.json({ error: "Organization name is required" }, { status: 400 });
  }

  const organizationId = body.organizationId?.trim();

  try {
    const api = await createZitadelApiClient();
    const response = await api.organizationService.addOrganization({
      name,
      ...(organizationId ? { organizationId } : {}),
    });

    return Response.json(
      {
        organizationId: response.organizationId,
        createdAdmins: response.createdAdmins.map((admin) => ({
          userId: admin.userId,
          emailCode: admin.emailCode ?? null,
          phoneCode: admin.phoneCode ?? null,
        })),
      },
      { status: 201 },
    );
  } catch (error) {
    const { status, body: errorBody } = createErrorResponse(error);
    return Response.json(errorBody, { status });
  }
}
